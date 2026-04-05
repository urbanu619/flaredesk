package core

import (
	"bytes"
	"context"
	"embed"
	"encoding/base64"
	"fmt"
	"github.com/afocus/captcha"
	"github.com/emirpasic/gods/maps/hashmap"
	"github.com/redis/go-redis/v9"
	"go_server/base/config"
	"go_server/utils"
	"image/color"
	"image/png"
	"io/fs"
	"strings"
	"time"
)

type CaptchaEngine struct {
}

//go:embed *.ttf
var fontFs embed.FS

var cap *captcha.Captcha
var capStore *hashmap.Map

const (
	CaptchaImgWidth    = 200
	CaptchaImgHeight   = 60
	CaptchaDisturbance = captcha.NORMAL // 干扰度
	CaptchaKeyLong     = 4
	CaptchaStrType     = captcha.NUM // 字符类型
	CaptchaOpenTimes   = 3           // 触发验证码次数
	CaptchaTimeOut     = 30          // 过期时间 S
)

func CapEngine() *CaptchaEngine {
	if cap == nil {
		cap = captcha.New()
		// 设置字体
		//_ = capCore.SetFont("./comic.ttf") // 这里改为绝对路径
		buf, err := fs.ReadFile(fontFs, "comic.ttf")
		if err != nil {
			panic(err)
		}
		err = cap.AddFontFromBytes(buf)
		if err != nil {
			panic(err)
		}
		// 设置验证码大小 128 64
		cap.SetSize(CaptchaImgWidth, CaptchaImgHeight)
		// 设置干扰强度
		cap.SetDisturbance(CaptchaDisturbance)
		// 设置前景色 可以多个 随机替换文字颜色 默认黑色
		// 第一位：4690D8
		// 第二位：F19A80
		// 第三位：9CD881
		// 第四位：CB7590
		cap.SetFrontColor(color.RGBA{70, 144, 216, 255},
			color.RGBA{156, 216, 129, 255},
			color.RGBA{203, 117, 144, 255},
			color.RGBA{241, 154, 128, 255},
		)
		// 设置背景色 可以多个 随机替换背景色 默认白色
		//capCore.SetBkgColor(color.RGBA{255, 0, 0, 255}, color.RGBA{0, 0, 255, 255}, color.RGBA{0, 153, 0, 255})
		//cap.SetBkgColor(color.RGBA{33, 33, 33, 255})
		cap.SetBkgColor(color.RGBA{255, 0, 0, 100}, color.RGBA{0, 0, 255, 100}, color.RGBA{0, 153, 0, 100})

	}
	if capStore == nil {
		capStore = hashmap.New()
	}
	return &CaptchaEngine{}
}

func (c *CaptchaEngine) Generate() (bs64 string, cid string) {
	img, code := c.getImg()
	cid = base64.StdEncoding.EncodeToString([]byte(utils.GenStrUuid()))
	// todo: set cid to store
	capStore.Put(cid, fmt.Sprintf("%s", code))
	w := bytes.NewBuffer(nil)
	_ = png.Encode(w, img)
	return fmt.Sprintf("data:image/png;base64,%s", base64.StdEncoding.EncodeToString(w.Bytes())), cid
}

func (c *CaptchaEngine) Verify(cid, answer string) (match bool) {
	if answer == config.VerifyCode {
		return true
	}
	vv, ok := capStore.Get(cid)
	if !ok {
		return ok
	}
	v, _ := vv.(string)
	vvLower := strings.ToLower(strings.TrimSpace(v))
	capStore.Remove(cid)
	return vvLower == strings.ToLower(strings.TrimSpace(answer))
}

func (c *CaptchaEngine) VerifyAndSendTick(cid, answer string) (string, bool) {
	if c.Verify(cid, answer) {
		ticket := base64.StdEncoding.EncodeToString([]byte(utils.GenStrUuid()))
		capStore.Put(ticket, true)
		return ticket, true
	}
	return "", false
}

func (c *CaptchaEngine) VerifyTick(ticket string) bool {
	_, ok := capStore.Get(ticket)
	return ok
}

func (c *CaptchaEngine) RemoveTick(ticket string) {
	capStore.Remove(ticket)
}

func (*CaptchaEngine) getImg() (*captcha.Image, string) {
	return cap.Create(CaptchaKeyLong, CaptchaStrType)
}

func (c *CaptchaEngine) capIpKey(ip string) string {
	return fmt.Sprintf("Cap:%s", ip)
}

func (c *CaptchaEngine) CapCheck(redisDb redis.UniversalClient, ip string) (int, bool) {
	vCmd := redisDb.Get(context.Background(), c.capIpKey(ip))
	t, err := vCmd.Int()
	if err != nil {
		return t, false
	}
	return t, t > CaptchaOpenTimes
}

func (c *CaptchaEngine) CapAdd(redisDb redis.UniversalClient, ip string) string {
	vCmd := redisDb.Get(context.Background(), c.capIpKey(ip))
	t, _ := vCmd.Int()
	sCmd := redisDb.Set(context.Background(), c.capIpKey(ip), t+1, time.Minute*time.Duration(CaptchaTimeOut))
	return sCmd.Val()
}

func (c *CaptchaEngine) CapClear(redisDb redis.UniversalClient, ip string) int64 {
	iCmd := redisDb.Del(context.Background(), c.capIpKey(ip))
	return iCmd.Val()
}
