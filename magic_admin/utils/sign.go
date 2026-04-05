package utils

import (
	"crypto/ecdsa"
	"fmt"
	"github.com/demdxx/gocast"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ethereum/go-ethereum/crypto"
	"strings"
	"time"
)

const (
	sysPrefix = "AICGOLD"
	pri       = "f05ffe77f132b34608e5bbbe5a16527738e3b22dbd50a48c1988529bc7c71b36"
	addr      = "0x484557aC895B3eEa2cE5B970d4E696DB11D44949"
)

func SystemAddress() string {
	return ethPriKey2HexAddress(pri)
}

func BuildSignMessage() (string, error) {
	msg := fmt.Sprintf("%s%d", sysPrefix, time.Now().Unix())
	sign, err := SignMessage(msg)
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s%s", msg, sign), nil
}

func priHexKeyToECDSA(hexPriKey string) (*ecdsa.PrivateKey, error) {
	return crypto.HexToECDSA(strings.TrimPrefix(hexPriKey, "0x"))
}

func ethPriKey2HexAddress(priKeyHash string) string {
	priKey, err := priHexKeyToECDSA(priKeyHash)
	if err != nil {
		panic(err)
	}
	pubKey := priKey.Public().(*ecdsa.PublicKey)
	addr := crypto.PubkeyToAddress(*pubKey)
	return addr.Hex()
}

func EthAddressCheck(addr string) bool {
	return common.IsHexAddress(addr)
}

func SignMessage(message string) (string, error) {
	priKey, err := priHexKeyToECDSA(pri)
	if err != nil {
		return "", err
	}
	return signMessage(message, priKey)
}

// signMessage 使用私钥签名消息
func signMessage(message string, privateKey *ecdsa.PrivateKey) (string, error) {
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

// 消息格式
// 1 去除消息头部信息
// 2 获取前10位数字
// 3 对比时间 不得相差3秒
func checkMessage(message string) (string, string, bool) {
	if !strings.HasPrefix(message, sysPrefix) {
		return "", "", false
	}
	removePrefix := strings.TrimPrefix(message, sysPrefix)
	timestamp := removePrefix[:10]
	notUnixTime := time.Now().Unix()
	if notUnixTime > gocast.ToInt64(timestamp)+3 || gocast.ToInt64(timestamp) > notUnixTime {
		return "", "", false
	}
	signature := removePrefix[10:]
	return sysPrefix + removePrefix[:10], signature, true
}

func VerifySignature(message string) error {
	msg, signature, ok := checkMessage(message)
	if !ok {
		return fmt.Errorf("非法签名")
	}
	return verifySignature(msg, signature, common.HexToAddress(addr))
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
