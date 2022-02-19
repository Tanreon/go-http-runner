package http_runner

import (
	"net/http"

	"github.com/go-resty/resty"
)

type IHttpRunner interface {
	GetJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	GetHtml(requestData HtmlRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	GetFile(requestData FileRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	PostJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	PutJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	PostForm(requestData FormRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
}

type JsonRequestData struct {
	Url     string
	Value   []byte
	Headers map[string]string
}

type HtmlRequestData struct {
	Url     string
	Value   []byte
	Headers map[string]string
}

type FormRequestData struct {
	Url     string
	Headers map[string]string
	Values  map[string]string
}

type FileRequestData struct {
	Url      string
	Headers  map[string]string
	FilePath string
}
