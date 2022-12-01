package http_runner

import (
	"io"
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
	GetJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error)
	GetHtml(requestOptions IHtmlRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error)
	GetFile(requestOptions IFileRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error)
	PostJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error)
	PutJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error)
	PostForm(requestOptions IFormRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error)
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

type IJsonRequestOptions interface {
	IBaseRequest

	IsValueSet() bool
	SetValue(bytes []byte)
	Value() []byte
}

type JsonRequestOptions struct {
	url            string
	value          *[]byte
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (j *JsonRequestOptions) Url() string {
	return j.url
}

func (j *JsonRequestOptions) IsHeadersSet() bool {
	return j.headers != nil
}
func (j *JsonRequestOptions) SetHeaders(headers map[string]string) {
	j.headers = &headers
}
func (j *JsonRequestOptions) Headers() map[string]string {
	return *j.headers
}

func (j *JsonRequestOptions) IsValueSet() bool {
	return j.value != nil
}
func (j *JsonRequestOptions) SetValue(bytes []byte) {
	j.value = &bytes
}
func (j *JsonRequestOptions) Value() []byte {
	return *j.value
}

func (j *JsonRequestOptions) IsRetryOptionSet() bool {
	return j.retryCount != nil
}
func (j *JsonRequestOptions) SetRetryOption(count int) {
	j.retryCount = &count
}
func (j *JsonRequestOptions) RetryOption() int {
	return *j.retryCount
}

func (j *JsonRequestOptions) IsTimeoutOptionSet() bool {
	return j.timeout != nil
}
func (j *JsonRequestOptions) SetTimeoutOption(timeout time.Duration) {
	j.timeout = &timeout
}
func (j *JsonRequestOptions) TimeoutOption() time.Duration {
	return *j.timeout
}

func (j *JsonRequestOptions) IsFollowRedirectOptionSet() bool {
	return j.followRedirect != nil
}
func (j *JsonRequestOptions) SetFollowRedirectOption(follow bool) {
	j.followRedirect = &follow
}
func (j *JsonRequestOptions) FollowRedirectOption() bool {
	return *j.followRedirect
}

func NewJsonRequestOptions(url string) IJsonRequestOptions {
	return &JsonRequestOptions{url: url}
}

//

type IHtmlRequestOptions interface {
	IBaseRequest

	IsValueSet() bool
	SetValue(bytes []byte)
	Value() []byte
}

type HtmlRequestOptions struct {
	url            string
	value          *[]byte
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (h *HtmlRequestOptions) Url() string {
	return h.url
}

func (h *HtmlRequestOptions) IsHeadersSet() bool {
	return h.headers != nil
}
func (h *HtmlRequestOptions) SetHeaders(headers map[string]string) {
	h.headers = &headers
}
func (h *HtmlRequestOptions) Headers() map[string]string {
	return *h.headers
}

func (h *HtmlRequestOptions) IsValueSet() bool {
	return h.value != nil
}
func (h *HtmlRequestOptions) SetValue(bytes []byte) {
	h.value = &bytes
}
func (h *HtmlRequestOptions) Value() []byte {
	return *h.value
}

func (h *HtmlRequestOptions) IsRetryOptionSet() bool {
	return h.retryCount != nil
}
func (h *HtmlRequestOptions) SetRetryOption(count int) {
	h.retryCount = &count
}
func (h *HtmlRequestOptions) RetryOption() int {
	return *h.retryCount
}

func (h *HtmlRequestOptions) IsTimeoutOptionSet() bool {
	return h.timeout != nil
}
func (h *HtmlRequestOptions) SetTimeoutOption(timeout time.Duration) {
	h.timeout = &timeout
}
func (h *HtmlRequestOptions) TimeoutOption() time.Duration {
	return *h.timeout
}

func (h *HtmlRequestOptions) IsFollowRedirectOptionSet() bool {
	return h.followRedirect != nil
}
func (h *HtmlRequestOptions) SetFollowRedirectOption(follow bool) {
	h.followRedirect = &follow
}
func (h *HtmlRequestOptions) FollowRedirectOption() bool {
	return *h.followRedirect
}

func NewHtmlRequestOptions(url string) IHtmlRequestOptions {
	return &HtmlRequestOptions{url: url}
}

//

type FileInfo struct {
	fileName string
	reader   io.Reader
}

type IFormRequestOptions interface {
	IBaseRequest

	IsValuesSet() bool
	SetValues(values map[string]string)
	Values() map[string]string

	IsFilesSet() bool
	SetFiles(files map[string]FileInfo)
	Files() map[string]FileInfo
}

type FormRequestOptions struct {
	url            string
	values         *map[string]string
	files          *map[string]FileInfo
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (f *FormRequestOptions) Url() string {
	return f.url
}

func (f *FormRequestOptions) IsHeadersSet() bool {
	return f.headers != nil
}
func (f *FormRequestOptions) SetHeaders(headers map[string]string) {
	f.headers = &headers
}
func (f *FormRequestOptions) Headers() map[string]string {
	return *f.headers
}

func (f *FormRequestOptions) IsValuesSet() bool {
	return f.values != nil
}
func (f *FormRequestOptions) SetValues(values map[string]string) {
	f.values = &values
}
func (f *FormRequestOptions) Values() map[string]string {
	return *f.values
}

func (f *FormRequestOptions) IsFilesSet() bool {
	return f.files != nil
}
func (f *FormRequestOptions) SetFiles(files map[string]FileInfo) {
	f.files = &files
}
func (f *FormRequestOptions) Files() map[string]FileInfo {
	return *f.files
}

func (f *FormRequestOptions) IsRetryOptionSet() bool {
	return f.retryCount != nil
}
func (f *FormRequestOptions) SetRetryOption(count int) {
	f.retryCount = &count
}
func (f *FormRequestOptions) RetryOption() int {
	return *f.retryCount
}

func (f *FormRequestOptions) IsTimeoutOptionSet() bool {
	return f.timeout != nil
}
func (f *FormRequestOptions) SetTimeoutOption(timeout time.Duration) {
	f.timeout = &timeout
}
func (f *FormRequestOptions) TimeoutOption() time.Duration {
	return *f.timeout
}

func (f *FormRequestOptions) IsFollowRedirectOptionSet() bool {
	return f.followRedirect != nil
}
func (f *FormRequestOptions) SetFollowRedirectOption(follow bool) {
	f.followRedirect = &follow
}
func (f *FormRequestOptions) FollowRedirectOption() bool {
	return *f.followRedirect
}

func NewFormRequestOptions(url string) IFormRequestOptions {
	return &FormRequestOptions{url: url}
}

//

type IFileRequestOptions interface {
	IBaseRequest

	FilePath() string
}

type FileRequestOptions struct {
	url            string
	filePath       string
	headers        *map[string]string
	retryCount     *int
	timeout        *time.Duration
	followRedirect *bool
}

func (j *FileRequestOptions) Url() string {
	return j.url
}

func (j *FileRequestOptions) FilePath() string {
	return j.filePath
}

func (j *FileRequestOptions) IsHeadersSet() bool {
	return j.headers != nil
}
func (j *FileRequestOptions) SetHeaders(headers map[string]string) {
	j.headers = &headers
}
func (j *FileRequestOptions) Headers() map[string]string {
	return *j.headers
}

func (j *FileRequestOptions) IsRetryOptionSet() bool {
	return j.retryCount != nil
}
func (j *FileRequestOptions) SetRetryOption(count int) {
	j.retryCount = &count
}
func (j *FileRequestOptions) RetryOption() int {
	return *j.retryCount
}

func (j *FileRequestOptions) IsTimeoutOptionSet() bool {
	return j.timeout != nil
}
func (j *FileRequestOptions) SetTimeoutOption(timeout time.Duration) {
	j.timeout = &timeout
}
func (j *FileRequestOptions) TimeoutOption() time.Duration {
	return *j.timeout
}

func (j *FileRequestOptions) IsFollowRedirectOptionSet() bool {
	return j.followRedirect != nil
}
func (j *FileRequestOptions) SetFollowRedirectOption(follow bool) {
	j.followRedirect = &follow
}
func (j *FileRequestOptions) FollowRedirectOption() bool {
	return *j.followRedirect
}

func NewFileRequestOptions(url, filePath string) IFileRequestOptions {
	return &FileRequestOptions{url: url, filePath: filePath}
}

func integrateCookies(requestOptions IBaseRequest, request *resty.Request, cookieJar []*http.Cookie) error {
	parsedUrl, err := url.Parse(requestOptions.Url())
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
