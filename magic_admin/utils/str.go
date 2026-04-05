package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"github.com/shopspring/decimal"
	"golang.org/x/text/language"
	"golang.org/x/text/message"
	"math/rand"
	"regexp"
	"strings"
	"time"
)

// length 长度 生成一个由0-9 + 26个小写字母组成的指定长度字符串

func RandStr(length int) string {
	digitCharacter := []string{"1", "2", "3", "4", "5", "6", "7", "8", "9", "0", "a", "b", "c", "d", "e", "f", "g", "h", "i", "j", "k", "l", "m", "n", "o", "p", "q", "r", "s", "t", "u", "v", "w", "x", "y", "z"}
	ranChr := ""
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	//rand.Seed(time.Now().Unix())
	for j := 0; j < length; j++ {
		ranChr += digitCharacter[r.Intn(len(digitCharacter))]
	}
	return strings.ToUpper(ranChr)
}

func DecimalFmtWithText(d decimal.Decimal) string {
	p := message.NewPrinter(language.English)
	f, _ := d.Float64()
	return p.Sprintf("%.2f", f)
}

func Md5Encode(data string) string {
	h := md5.New()
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

func Md5ByteEncode(data []byte) string {
	h := md5.New()
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

func Md5JsonEncode(any interface{}) string {
	h := md5.New()
	data, _ := json.Marshal(any)
	h.Write(data)
	return hex.EncodeToString(h.Sum(nil))
}

// 安全码检查

func CheckSafePasswordLever6(ps string) error {
	if len(ps) != 6 {
		return fmt.Errorf("password len must be 6")
	}
	pattern := `^[0-9]{6}$` // 必须为数字
	reg := regexp.MustCompile(pattern)
	if !reg.MatchString(ps) {
		return fmt.Errorf("password must be nums")
	}
	return nil
}

// 单词变化为蛇形

func VarToSnakeCase(s string) string {
	// 在大写字母前添加下划线，并将所有非字母数字字符替换为下划线
	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")
	snake := matchFirstCap.ReplaceAllString(s, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	snake = strings.ToLower(snake)
	// 移除连续的下划线和非字母数字字符
	reg := regexp.MustCompile("[^a-z0-9]+")
	snake = reg.ReplaceAllString(snake, "_")
	snake = strings.Trim(snake, "_")
	return snake
}

func WordsToSnakeCase(s string) string {
	strsOne := strings.Split(s, ",")
	handOne := make([]string, 0)
	for _, item := range strsOne {
		newStrs := make([]string, 0)
		strs := strings.Split(item, " ")
		for _, str := range strs {
			newStrs = append(newStrs, VarToSnakeCase(str))
		}
		handOne = append(handOne, strings.Join(newStrs, " "))
	}
	return strings.Join(handOne, ",")
}
