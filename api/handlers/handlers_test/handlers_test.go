package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"time"
)

type HttpMock struct {
	server *httptest.Server
}

func NewHttpMock(handler http.HandlerFunc) *HttpMock {
	server := httptest.NewServer(http.HandlerFunc(handler))
	time.Sleep(1 * time.Second)
	return &HttpMock{server}
}

func (h *HttpMock) Post(t *testing.T, body string) Resp {
	resp, err := http.Post(h.server.URL, "application/json; charset=UTF-8", bytes.NewBufferString(body))
	if err != nil {
		t.Error(err)
	}
	return NewResp(resp)
}

func (h *HttpMock) Close() {
	h.server.Close()
}

type Resp struct {
	Resp       *http.Response
	StatusCode int
	Header     Header
}

func NewResp(resp *http.Response) Resp {
	c := resp.Header.Get("Set-Cookie")
	header := Header{cookies: strings.Split(c, ";")}
	statusCode := resp.StatusCode
	return Resp{resp, statusCode, header}
}

func (r *Resp) BodyJson(resstr interface{}) error {
	bytetoken, err := io.ReadAll(r.Resp.Body)
	if err != nil {
		return err
	}
	if err = json.Unmarshal(bytetoken, resstr); err != nil {
		return err
	}
	return nil
}

func (r *Resp) BodyClose() {
	r.Resp.Body.Close()
}

type Header struct {
	cookies []string
}

func (h *Header) GetCookie(key string) string {
	key = fmt.Sprintf("%s=", key)
	value := ""
	for _, c := range h.cookies {
		if strings.HasPrefix(strings.TrimSpace(c), key) {
			value = strings.TrimPrefix(strings.TrimSpace(c), key)
			break
		}
	}
	return value
}

func (h *Header) GetShortJwtCookie() string {
	return h.GetCookie("shortjwt")
}
