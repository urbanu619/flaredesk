package main

import (
	"go_server/base/core"
	"testing"
)

func TestCreateSysAddress(t *testing.T) {
	mn := core.GenerateMnemonicWords(128)
	addr, pri, err := core.MnemonicToEthInfo(mn)
	if err != nil {
		t.Log(err.Error())
		return
	}
	t.Log("系统地址", addr)
	t.Log("私钥", pri)
	t.Log("转码私钥", core.Base64Encode(pri))
}

// 获取当前系统签名信息

func TestGetSysAddress(t *testing.T) {
	core.Log.Info("系统签名前缀:", core.SysSignPrefix())
	sysAddr, err := core.SysAddress()
	if err != nil {
		core.Log.Info(err)
		return
	}
	core.Log.Info("系统地址:", sysAddr)
	t.Log(core.Base64Decode(""))
}
