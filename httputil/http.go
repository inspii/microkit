package httputil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

const contentTypeJSON = "application/json"

type StatusError interface {
	Status() int
}

type Error struct {
	status int
	body   []byte
}

func NewError(code int, body []byte) error {
	return &Error{
		status: code,
		body:   body,
	}
}

func (e Error) Status() int {
	return e.status
}

func (e Error) Error() string {
	return fmt.Sprintf("code=%d, body=%s", e.status, e.body)
}

func GetStatus(err error) int {
	if httpErr, ok := err.(StatusError); ok {
		return httpErr.Status()
	}

	if err == nil {
		return http.StatusOK
	} else {
		return http.StatusInternalServerError
	}
}

// RequestJSON 发起HTTP JSON请求
//
// method 请求方法
// url 请求地址
// timeout 超时时间
// requestBody 请求数据。会序列化为JSON字节数组
// responseBody 响应数据。会对响应数据进行JSON反序列化，并将结果存入responseBody
func RequestJSON(method string, url string, timeout time.Duration, requestBody interface{}, responseBody interface{}) error {
	var requestBodyBin *bytes.Reader
	if requestBody != nil {
		buf, _ := json.Marshal(requestBody)
		requestBodyBin = bytes.NewReader(buf)
	} else {
		requestBodyBin = bytes.NewReader(make([]byte, 0))
	}

	req, err := http.NewRequest(method, url, requestBodyBin)
	if err != nil {
		return err
	}
	req.Header.Add("Content-Type", contentTypeJSON)
	req.Header.Add("Accept", contentTypeJSON)

	client := http.DefaultClient
	client.Timeout = timeout
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	status := resp.StatusCode
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	if status < 200 || status >= 300 {
		return NewError(status, body)
	}
	if body == nil || responseBody == nil {
		return nil
	}

	contentType := resp.Header.Get("Content-Type")
	if !strings.Contains(contentType, contentTypeJSON) && len(body) > 0 {
		return err
	}

	return json.Unmarshal(body, responseBody)
}
