package handlers_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"himakiwa/handlers/middleware"
	"io"
	"net/http"
	"net/http/cookiejar"
	"net/http/httptest"
	"net/url"
	"strings"
	"testing"
)

type HttpMock struct {
	client http.Client
	server *httptest.Server
}

func NewReqAuthHttpMock(handler http.HandlerFunc, jwt string) *HttpMock {
	h := middleware.AuthMiddleware(handler)
	server := httptest.NewServer(h)

	url, _ := url.Parse(server.URL)
	jar, _ := cookiejar.New(nil)
	jwtCookie := &http.Cookie{Name: "token", Value: jwt}
	jar.SetCookies(url, []*http.Cookie{jwtCookie})
	client := http.Client{Jar: jar}
	return &HttpMock{client, server}
}

func NewHttpMock(handler http.HandlerFunc) *HttpMock {
	server := httptest.NewServer(http.HandlerFunc(handler))
	client := http.Client{}
	return &HttpMock{client, server}
}

func (h *HttpMock) do(t *testing.T, method string, body io.Reader) Resp {
	req, err := http.NewRequest(method, h.server.URL, body)
	if err != nil {
		t.Error(err)
	}
	resp, err := h.client.Do(req)
	if err != nil {
		t.Error(err)
	}
	return NewResp(resp)
}

func (h *HttpMock) Get(t *testing.T) Resp {
	return h.do(t, "GET", nil)
}

func (h *HttpMock) Post(t *testing.T, body string) Resp {
	return h.do(t, "POST", bytes.NewBufferString(body))
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
