package http_runner

import (
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
)

var DefaultHeaders = map[string]string{
	"accept":          "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"accept-language": "en-US,en;q=0.9",
	"cache-control":   "max-age=0",
}

type IHttpRunner interface {
	GetJson(requestData IJsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	GetHtml(requestData IHtmlRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	GetFile(requestData IFileRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	PostJson(requestData IJsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	PutJson(requestData IJsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	PostForm(requestData IFormRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
}

type IBaseRequest interface {
	Url() string

	IsHeadersSet() bool
	SetHeaders(headers map[string]string)
	Headers() map[string]string

	IsRetryOptionSet() bool
	SetRetryOption(count int)
	RetryOption() int

	IsTimeoutOptionSet() bool
	SetTimeoutOption(timeout time.Duration)
	TimeoutOption() time.Duration

	IsFollowRedirectOptionSet() bool
	SetFollowRedirectOption(follow bool)
	FollowRedirectOption() bool
}

//

type IJsonRequestData interface {
	IBaseRequest

	IsValueSet() bool
	SetValue(bytes []byte)
	Value() []byte
}

type JsonRequestData struct {
	url            string
	value          *[]byte
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (j *JsonRequestData) Url() string {
	return j.url
}

func (j *JsonRequestData) IsHeadersSet() bool {
	return j.headers != nil
}
func (j *JsonRequestData) SetHeaders(headers map[string]string) {
	j.headers = &headers
}
func (j *JsonRequestData) Headers() map[string]string {
	return *j.headers
}

func (j *JsonRequestData) IsValueSet() bool {
	return j.value != nil
}
func (j *JsonRequestData) SetValue(bytes []byte) {
	j.value = &bytes
}
func (j *JsonRequestData) Value() []byte {
	return *j.value
}

func (j *JsonRequestData) IsRetryOptionSet() bool {
	return j.retryCount != nil
}
func (j *JsonRequestData) SetRetryOption(count int) {
	j.retryCount = &count
}
func (j *JsonRequestData) RetryOption() int {
	return *j.retryCount
}

func (j *JsonRequestData) IsTimeoutOptionSet() bool {
	return j.timeout != nil
}
func (j *JsonRequestData) SetTimeoutOption(timeout time.Duration) {
	j.timeout = &timeout
}
func (j *JsonRequestData) TimeoutOption() time.Duration {
	return *j.timeout
}

func (j *JsonRequestData) IsFollowRedirectOptionSet() bool {
	return j.followRedirect != nil
}
func (j *JsonRequestData) SetFollowRedirectOption(follow bool) {
	j.followRedirect = &follow
}
func (j *JsonRequestData) FollowRedirectOption() bool {
	return *j.followRedirect
}

func NewJsonRequestData(url string) IJsonRequestData {
	return &JsonRequestData{url: url}
}

//

type IHtmlRequestData interface {
	IBaseRequest

	IsValueSet() bool
	SetValue(bytes []byte)
	Value() []byte
}

type HtmlRequestData struct {
	url            string
	value          *[]byte
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (h *HtmlRequestData) Url() string {
	return h.url
}

func (h *HtmlRequestData) IsHeadersSet() bool {
	return h.headers != nil
}
func (h *HtmlRequestData) SetHeaders(headers map[string]string) {
	h.headers = &headers
}
func (h *HtmlRequestData) Headers() map[string]string {
	return *h.headers
}

func (h *HtmlRequestData) IsValueSet() bool {
	return h.value != nil
}
func (h *HtmlRequestData) SetValue(bytes []byte) {
	h.value = &bytes
}
func (h *HtmlRequestData) Value() []byte {
	return *h.value
}

func (h *HtmlRequestData) IsRetryOptionSet() bool {
	return h.retryCount != nil
}
func (h *HtmlRequestData) SetRetryOption(count int) {
	h.retryCount = &count
}
func (h *HtmlRequestData) RetryOption() int {
	return *h.retryCount
}

func (h *HtmlRequestData) IsTimeoutOptionSet() bool {
	return h.timeout != nil
}
func (h *HtmlRequestData) SetTimeoutOption(timeout time.Duration) {
	h.timeout = &timeout
}
func (h *HtmlRequestData) TimeoutOption() time.Duration {
	return *h.timeout
}

func (h *HtmlRequestData) IsFollowRedirectOptionSet() bool {
	return h.followRedirect != nil
}
func (h *HtmlRequestData) SetFollowRedirectOption(follow bool) {
	h.followRedirect = &follow
}
func (h *HtmlRequestData) FollowRedirectOption() bool {
	return *h.followRedirect
}

func NewHtmlRequestData(url string) IHtmlRequestData {
	return &HtmlRequestData{url: url}
}

//

type IFormRequestData interface {
	IBaseRequest

	IsValuesSet() bool
	SetValues(values map[string]string)
	Values() map[string]string
}

type FormRequestData struct {
	url            string
	values         *map[string]string
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (f *FormRequestData) Url() string {
	return f.url
}

func (f *FormRequestData) IsHeadersSet() bool {
	return f.headers != nil
}
func (f *FormRequestData) SetHeaders(headers map[string]string) {
	f.headers = &headers
}
func (f *FormRequestData) Headers() map[string]string {
	return *f.headers
}

func (f *FormRequestData) IsValuesSet() bool {
	return f.values != nil
}
func (f *FormRequestData) SetValues(values map[string]string) {
	f.values = &values
}
func (f *FormRequestData) Values() map[string]string {
	return *f.values
}

func (f *FormRequestData) IsRetryOptionSet() bool {
	return f.retryCount != nil
}
func (f *FormRequestData) SetRetryOption(count int) {
	f.retryCount = &count
}
func (f *FormRequestData) RetryOption() int {
	return *f.retryCount
}

func (f *FormRequestData) IsTimeoutOptionSet() bool {
	return f.timeout != nil
}
func (f *FormRequestData) SetTimeoutOption(timeout time.Duration) {
	f.timeout = &timeout
}
func (f *FormRequestData) TimeoutOption() time.Duration {
	return *f.timeout
}

func (f *FormRequestData) IsFollowRedirectOptionSet() bool {
	return f.followRedirect != nil
}
func (f *FormRequestData) SetFollowRedirectOption(follow bool) {
	f.followRedirect = &follow
}
func (f *FormRequestData) FollowRedirectOption() bool {
	return *f.followRedirect
}

func NewFormRequestData(url string) IFormRequestData {
	return &FormRequestData{url: url}
}

//

type IFileRequestData interface {
	IBaseRequest

	FilePath() string
}

type FileRequestData struct {
	url            string
	filePath       string
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (j *FileRequestData) Url() string {
	return j.url
}

func (j *FileRequestData) FilePath() string {
	return j.filePath
}

func (j *FileRequestData) IsHeadersSet() bool {
	return j.headers != nil
}
func (j *FileRequestData) SetHeaders(headers map[string]string) {
	j.headers = &headers
}
func (j *FileRequestData) Headers() map[string]string {
	return *j.headers
}

func (j *FileRequestData) IsRetryOptionSet() bool {
	return j.retryCount != nil
}
func (j *FileRequestData) SetRetryOption(count int) {
	j.retryCount = &count
}
func (j *FileRequestData) RetryOption() int {
	return *j.retryCount
}

func (j *FileRequestData) IsTimeoutOptionSet() bool {
	return j.timeout != nil
}
func (j *FileRequestData) SetTimeoutOption(timeout time.Duration) {
	j.timeout = &timeout
}
func (j *FileRequestData) TimeoutOption() time.Duration {
	return *j.timeout
}

func (j *FileRequestData) IsFollowRedirectOptionSet() bool {
	return j.followRedirect != nil
}
func (j *FileRequestData) SetFollowRedirectOption(follow bool) {
	j.followRedirect = &follow
}
func (j *FileRequestData) FollowRedirectOption() bool {
	return *j.followRedirect
}

func NewFileRequestData(url, filePath string) IFileRequestData {
	return &FileRequestData{url: url, filePath: filePath}
}

func integrateCookies(requestData IBaseRequest, request *resty.Request, cookieJar []*http.Cookie) error {
	parsedUrl, err := url.Parse(requestData.Url())
	if err != nil {
		return err
	}

	cookieJarMap := make(map[string]*http.Cookie)

	for _, cookie := range cookieJar {
		if strings.HasSuffix(parsedUrl.Host, strings.TrimPrefix(cookie.Domain, ".")) {
			cookieComplexKey := cookie.Domain + cookie.Name + cookie.Path
			if _, present := cookieJarMap[cookieComplexKey]; !present {
				cookieJarMap[cookieComplexKey] = cookie
			}
		}
	}

	for _, cookie := range cookieJarMap {
		request.SetCookie(cookie)
	}

	return err
}
