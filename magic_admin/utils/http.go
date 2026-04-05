package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"
)

func Get[T any](url string) (*T, error) {
	return get[T](url)
}

func get[T any](url string) (*T, error) {
	obj := new(T)
	resp, err := http.Get(url)
	if err != nil {
		return obj, err
	}
	defer resp.Body.Close() // 确保在函数结束时关闭响应体
	// 检查响应状态码
	if resp.StatusCode != http.StatusOK {
		return obj, fmt.Errorf("error:%d", resp.StatusCode)
	}
	// 读取响应体
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(body, &obj)
	return obj, err
}

func POST[T any](url string, data []byte) (*T, error) {
	return post[T](url, data)
}

func post[T any](uri string, data []byte) (*T, error) {
	obj := new(T)
	// 创建一个新的 POST 请求
	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(data))
	if err != nil {
		return obj, err
	}
	// 设置请求头
	req.Header.Set("Content-Type", "application/json")
	// 发送请求
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return obj, err

	}
	defer resp.Body.Close()
	// 读取响应
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return obj, err
	}
	err = json.Unmarshal(body, &obj)
	return obj, err
}
