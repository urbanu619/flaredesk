package core

import (
	"crypto/ecdsa"
	"encoding/base64"
	"fmt"
	"github.com/btcsuite/btcd/btcutil/hdkeychain"
	"github.com/btcsuite/btcd/chaincfg"
	"github.com/demdxx/gocast"
	"github.com/ethereum/go-ethereum/accounts"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/tyler-smith/go-bip39"
	"go_server/model/system"
	"strings"
	"time"
)

type Sign struct {
	SysName       string         `mapstructure:"sys-name" json:"sys-name" yaml:"sys-name"` // 本系统名
	SysSignPrefix string         `mapstructure:"sys-sign-prefix" json:"sys-sign-prefix" yaml:"sys-sign-prefix"`
	SysAddress    string         `mapstructure:"sys-address" json:"sys-address" yaml:"sys-address"` // 本系统地址
	SysKey        string         `mapstructure:"sys-key" json:"sys-key" yaml:"sys-key"`             // 本系统密钥 -- 使用本系统密码本加密
	Internals     []*InternalSys `mapstructure:"internals" json:"internals" yaml:"internals"`
}

type InternalSys struct {
	SysName    string `mapstructure:"sys-name" json:"sys-name" yaml:"sys-name"`
	SysUrl     string `mapstructure:"sys-url" json:"sys-url" yaml:"sys-url"`                // 系统地址
	SignPrefix string `mapstructure:"sign-prefix" json:"sign-prefix" yaml:"sign-prefix"`    // 签名前缀
	SignExpSec int64  `mapstructure:"sign-exp-sec" json:"sign-exp-sec" yaml:"sign-exp-sec"` // 签名超时-S
	SysAddress string `mapstructure:"sys-address" json:"sys-address" yaml:"sys-address"`    // 管理系统地址
}

var signConfer *Sign

func signConf() *Sign {
	if signConfer == nil {
		// 初始化
		rows := make([]*system.SysSignConfig, 0)
		err := MainDb().
			Model(&system.SysSignConfig{}).
			Where("1 = 1").
			Find(&rows).Error
		if err != nil {
			panic(err)
		}
		signConfer = new(Sign)
		for _, row := range rows {
			if row.IsSystemSign {
				signConfer.SysSignPrefix = row.SignName
				signConfer.SysKey = row.SignPriKey
				signConfer.SysAddress = row.SignAddress
				signConfer.SysName = row.SignName
				continue
			}
			if signConfer.Internals == nil {
				signConfer.Internals = make([]*InternalSys, 0)
			}
			signConfer.Internals = append(signConfer.Internals, &InternalSys{
				SysName:    row.SignName,
				SysAddress: row.SignAddress,
				SysUrl:     row.SysUrl,
				SignPrefix: row.SignName,
				SignExpSec: row.SignExpSec,
			})
		}
	}
	return signConfer
}

// 构建本系统签名信息

func BuildSignMessage() (string, error) {
	sysPrivateKey, err := Base64Decode(signConf().SysKey)
	if err != nil {
		return "", err
	}
	msg := fmt.Sprintf("%s%d", signConf().SysSignPrefix, time.Now().Unix())
	sign, err := signMessage(msg, sysPrivateKey)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", msg, sign), nil
}

type Signer struct {
	SignPrefix string `mapstructure:"sign-prefix" json:"sign-prefix" yaml:"sign-prefix"`    // 签名前缀
	SignExpSec int64  `mapstructure:"sign-exp-sec" json:"sign-exp-sec" yaml:"sign-exp-sec"` // 签名超时-S
	SysAddress string `mapstructure:"sys-address" json:"sys-address" yaml:"sys-address"`    // 管理系统地址
}

// 通过系统名获取外部系统签名器 用于签名验证

var signerMap = make(map[string]*Signer)

func SignerBySysName(sysName string) (*Signer, error) {
	v, ok := signerMap[sysName]
	if ok {
		return v, nil
	}
	for _, item := range signConf().Internals {
		if item.SysName == sysName {
			signerMap[sysName] = &Signer{
				SignPrefix: item.SignPrefix,
				SignExpSec: item.SignExpSec,
				SysAddress: item.SysAddress,
			}
			return signerMap[sysName], nil
		}
	}
	return nil, fmt.Errorf("sys name not found")
}

// 本系统地址 -- 暴露给其他系统使用

func SysAddress() (string, error) {
	sysPrivateKey, err := Base64Decode(signConf().SysKey)
	if err != nil {
		return "", err
	}
	return ethPriKey2HexAddress(sysPrivateKey)
}

func SysSignPrefix() string {
	return signConf().SysSignPrefix
}

const (
	base64Table = "7PQX12RVW3YZaDEFGbcdeIJjkKLMNO56fghiABCHlSTUmnopqrxyz04stuvw89+/"
)

var coder = base64.NewEncoding(base64Table)

func Base64Encode(src string) string { //编码
	return coder.EncodeToString([]byte(src))
}

func Base64Decode(src string) (string, error) { //解码
	bts, err := coder.DecodeString(src)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s", bts), nil
}

func signMessage(message, hexPri string) (string, error) {
	priKey, err := priHexKeyToECDSA(hexPri)
	if err != nil {
		return "", err
	}
	return signMessageWithEcdsa(message, priKey)
}

// signMessageWithEcdsa 使用私钥签名消息
func signMessageWithEcdsa(message string, privateKey *ecdsa.PrivateKey) (string, error) {
	// 计算消息的Keccak256哈希
	hash := crypto.Keccak256Hash([]byte(message))

	// 签名
	signature, err := crypto.Sign(hash.Bytes(), privateKey)
	if err != nil {
		return "", err
	}

	// 添加恢复ID(以太坊要求)
	signature[64] += 27

	return hexutil.Encode(signature), nil
}

func (s *Signer) VerifySignature(message string) error {
	msg, signature, ok := s.checkMessage(message)
	if !ok {
		return fmt.Errorf("非法签名")
	}
	return verifySignature(msg, signature, common.HexToAddress(s.SysAddress))
}

// 消息格式
// 1 去除消息头部信息
// 2 获取前10位数字
// 3 对比时间 不得相差3秒

func (s *Signer) checkMessage(message string) (string, string, bool) {

	if !strings.HasPrefix(message, s.SignPrefix) {
		fmt.Println(fmt.Sprintf("签名未包含前缀:%s", s.SignPrefix))
		return "", "", false
	}

	removePrefix := strings.TrimPrefix(message, s.SignPrefix)
	timestamp := removePrefix[:10]
	notUnixTime := time.Now().Unix()
	if notUnixTime > gocast.ToInt64(timestamp)+s.SignExpSec || gocast.ToInt64(timestamp) > notUnixTime {
		return "", "", false
	}
	signature := removePrefix[10:]
	return s.SignPrefix + removePrefix[:10], signature, true
}

// verifySignature 验证签名
func verifySignature(message, signature string, address common.Address) error {
	// 解码签名
	sig, err := hexutil.Decode(signature)
	if err != nil {
		return err
	}

	// 恢复ID必须为27或28
	if sig[64] != 27 && sig[64] != 28 {
		return fmt.Errorf("无效的恢复ID")
	}

	// 调整恢复ID为0或1(以太坊库要求)
	sig[64] -= 27

	// 计算消息哈希
	hash := crypto.Keccak256Hash([]byte(message))

	// 从签名恢复公钥
	publicKey, err := crypto.SigToPub(hash.Bytes(), sig)
	if err != nil {
		return err
	}

	// 从公钥生成地址
	recoveredAddr := crypto.PubkeyToAddress(*publicKey)

	// 验证地址是否匹配
	if recoveredAddr != address {
		return fmt.Errorf("签名验证失败: 地址不匹配")
	}

	return nil
}

// 生成助记词 使用标准：BIP39(i -128 /256)
// 生成“128的熵”
//1. 生成一个长度为 128 ~ 256 位的随机序列（后面称为熵）；
//- 熵的取值长度需要为 32 的整数倍的值。所以取值可能为：【128, 160, 192, 224, 256】。
//2. 获取熵的校验和；
//- 取熵哈希后的前 n 位作为校验和（n = 熵长度/32）。
//- 校验和的理论值为【4，5，6，7，8】。
//3. 生成新的序列；
//- 新的序列方式为 = 熵 + 校验和
//- 新序列的长度理论值可能为：【128+4，160+5，192+6，224+7，256+8】即【132，165，198，231，264】
//4. 获得 m 个 11 位二进制数。
//- 将新生成的序列，按照 11 位进行平分。
//- m 的理论可能值为：【12，15，18，21，24】
//- 为什么是 11 位？助记词库总共有 2048 个单词，2^11 正好是 2048。所以用 11 位二进制数刚好可以将所有的助记词定位。
//5. 获得 m 个助记词；
//- 根据第四步得到的 m 个 11位二进制数，去助记词表定位每一个助记词。
//- 得到助记词个数的理论值为：【12，15，18，21，24】
// 波场使用的熵长度为256 128生成12位单词助记词

func GenerateMnemonicWords(i int) string {
	entropy, err := bip39.NewEntropy(i)
	if err != nil {
		return ""
	}
	mnemonic, _ := bip39.NewMnemonic(entropy)
	return mnemonic
}

func MnemonicToEthInfo(mnemonic string) (ethHexAddress string, subPriKey string, err error) {
	// 验证助记词有效性
	if !bip39.IsMnemonicValid(mnemonic) {
		return "", "", fmt.Errorf("invalid mnemonic phrase")
	}

	// 生成种子
	seed := bip39.NewSeed(mnemonic, "") // 第二个参数是密码，可以为空

	// 创建主密钥
	master, err := hdkeychain.NewMaster(seed, &chaincfg.MainNetParams)
	if err != nil {
		return "", "", fmt.Errorf("failed to create master key: %v", err)
	}

	// 解析派生路径
	dp, err := accounts.ParseDerivationPath("m/44'/60'/0'/0/0")
	if err != nil {
		return "", "", fmt.Errorf("failed to parse derivation path: %v", err)
	}

	// 派生密钥
	for _, n := range dp {
		var child *hdkeychain.ExtendedKey
		if master.IsAffectedByIssue172() {
			child, err = master.Derive(n)
		} else {
			child, err = master.DeriveNonStandard(n)
		}
		if err != nil {
			return "", "", fmt.Errorf("failed to derive key: %v", err)
		}
		master = child
	}

	// 获取ECDSA私钥
	privateKey, err := master.ECPrivKey()
	if err != nil {
		return "", "", fmt.Errorf("failed to get EC private key: %v", err)
	}

	privateKeyECDSA := privateKey.ToECDSA()
	publicKey := privateKeyECDSA.Public().(*ecdsa.PublicKey)

	// 生成地址
	addr := crypto.PubkeyToAddress(*publicKey)
	ethHexAddress = addr.Hex()

	// 编码私钥
	privateKeyBytes := crypto.FromECDSA(privateKeyECDSA)
	subPriKey = hexutil.Encode(privateKeyBytes)

	return ethHexAddress, subPriKey, nil
}

func GenerateSysInfo() (string, string, error) {
	mn := GenerateMnemonicWords(128)
	addr, pri, err := MnemonicToEthInfo(mn)
	if err != nil {
		return "", "", err
	}
	pri = Base64Encode(pri)
	return addr, pri, nil
}

func EthAddressCheck(addr string) bool {
	return common.IsHexAddress(addr)
}

func priHexKeyToECDSA(hexPriKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(hexPriKey, "0x"))
}

func ethPriKey2HexAddress(priKeyHash string) (string, error) {
	priKey, err := priHexKeyToECDSA(priKeyHash)
	if err != nil {
		return "", err
	}
	pubKey := priKey.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pubKey)
	return addr.Hex(), nil
}

func EthPriKey2HexAddress(priKeyHash string) (string, error) {
	return ethPriKey2HexAddress(priKeyHash)
}
