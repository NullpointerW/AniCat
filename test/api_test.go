package test

import (
	"bytes"
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"testing"
)

func TestApi(t *testing.T) {
	c := &http.Client{}
	undecode := `2TM6ac9b293c0f550168509242317768790196`
	data := []byte(undecode)

	// 计算md5散列值
	hash := md5.Sum(data)

	// 将散列值转换为16进制字符串
	md5str := hex.EncodeToString(hash[:])
	fmt.Println(md5str)
	values := url.Values{}
	values.Set("gatewayCode", "TM6ac9b293c0f550")
	values.Set("breakerIndex", "2")
	values.Set("rts", "1685092703")

	r, _ := http.NewRequest("POST", "http://admin.tmiotb.com:8888/api/breaker/info", bytes.NewBufferString(values.Encode()))
	r.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Set("appid", "17768790196")
	r.Header.Set("sign", md5str)
	

	resp, err := c.Do(r)
	if err != nil {
		t.Error(err)
		t.FailNow()
	}
	b, _ := io.ReadAll(resp.Body)
	fmt.Println(string(b))
}
