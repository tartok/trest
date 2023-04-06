package trest

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

var (
	Url         = "http://localhost:8081"
	AccessToken = ""
)

func do(method, path string, query map[string]interface{}, body interface{}, token string, contentType string) (string, []byte, error) {
	var r io.Reader
	if body != nil {
		switch s := body.(type) {
		case string:
			r = bytes.NewReader([]byte(s))
		case *bytes.Buffer:
			r = bytes.NewReader(s.Bytes())
		default:
			d, err := json.Marshal(body)
			if err != nil {
				return "", nil, err
			}
			r = bytes.NewReader(d)
		}
	}
	var q []byte
	if len(query) > 0 {
		for k, v := range query {
			q = append(q, []byte(fmt.Sprintf("&%s=%v", k, v))...)
		}
		q[0] = '?'
	}
	request, err := http.NewRequest(method, Url+path+string(q), r)
	if err != nil {
		return "", nil, err
	}
	if contentType == "" {
		contentType = "application/json; charset=utf-8"
	}
	request.Header.Set("Content-Type", contentType)
	if token != "" {
		request.Header.Set("authorization", fmt.Sprintf("Bearer %s", token))
	}
	resp, err := http.DefaultClient.Do(request)
	if err != nil {
		return "", nil, err
	}
	var b []byte
	if resp.StatusCode != http.StatusNotFound {
		b, _ = io.ReadAll(resp.Body)
	}
	return resp.Status, b, err
}
func Post(path string, query map[string]interface{}, body interface{}, contentType string) (string, []byte, error) {
	return do(http.MethodPost, path, query, body, AccessToken, contentType)
}
func Put(path string, query map[string]interface{}, body interface{}, contentType string) (string, []byte, error) {
	return do(http.MethodPut, path, query, body, AccessToken, contentType)
}
func Patch(path string, query map[string]interface{}, body interface{}, contentType string) (string, []byte, error) {
	return do(http.MethodPatch, path, query, body, AccessToken, contentType)
}
func Delete(path string, query map[string]interface{}, body interface{}, contentType string) (string, []byte, error) {
	return do(http.MethodDelete, path, query, body, AccessToken, contentType)
}
func Get(path string, query map[string]interface{}, body interface{}) (string, []byte, error) {
	return do(http.MethodGet, path, query, body, AccessToken, "")
}

// PrintResult
func PrintResult(status string, body []byte, err error) (string, []byte, error) {
	d := bytes.NewBuffer(nil)
	if len(body) > 0 {
		e := json.Indent(d, body, " ", " ")
		if e != nil {
			panic(e)
		}
	}
	fmt.Printf("%s %v\n%s", status, err, d.Bytes())
	return status, body, err
}
