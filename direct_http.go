package http_runner

import (
	"context"
	NetworkRunner "github.com/Tanreon/go-network-runner"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nadoo/glider/rule"

	log "github.com/sirupsen/logrus"
)

type DirectHttpRunner struct {
	defHeaders map[string]string
	client     *resty.Client
}

func NewAdvancedDirectHttpRunner(dialer *rule.Proxy, retryCount int, timeout time.Duration, headers map[string]string) (IHttpRunner, error) {
	// CREATE TRANSPORT FOR HTTP
	transport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			return dialer.NextDialer(addr).Dial(network, addr)
		},
	}
	// CREATE TRANSPORT FOR HTTP

	// CREATE A RESTY CLIENT WITHOUT PROXY
	client := resty.New()
	client.SetTransport(&transport)
	client.SetRetryCount(retryCount)
	client.SetTimeout(timeout) // time.Second * 15
	client.SetDisableWarn(true)

	if !log.IsLevelEnabled(log.TraceLevel) {
		restyLogger := log.New()
		restyLogger.SetOutput(io.Discard)

		client.SetLogger(restyLogger)
	}

	//// Using raw func into resty.SetRedirectPolicy
	//client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
	//	// Implement your logic here
	//
	//	// return nil for continue redirect otherwise return error to stop/prevent redirect
	//	return nil
	//}))

	//client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
	//	return http.ErrUseLastResponse
	//}))

	//client.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//client.Header.Add("accept-encoding", "gzip, deflate, br")
	//client.Header.Add("accept-language", "en-US,en;q=0.9")
	//client.Header.Add("cache-control", "max-age=0")
	//client.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")

	runner := &DirectHttpRunner{
		defHeaders: headers,
		client:     client,
	}
	// CREATE A RESTY CLIENT WITHOUT PROXY

	return runner, nil
}

func NewDirectHttpRunner(dialer *rule.Proxy) (IHttpRunner, error) {
	return NewAdvancedDirectHttpRunner(dialer, 2, time.Second*15, DefaultHeaders)
}

func NewDefaultDirectHttpRunner() (IHttpRunner, error) {
	directDialer, err := NetworkRunner.NewDirectDialer()
	if err != nil {
		return nil, err
	}

	return NewDirectHttpRunner(directDialer)
}

func (d *DirectHttpRunner) GetJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestOptions.IsRetryOptionSet() {
		d.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		d.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := d.client.R()

	if len(d.defHeaders) > 0 {
		for key, value := range d.defHeaders {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Set(key, value)
		}
	}

	if len(request.Header.Get("Content-Type")) <= 0 {
		request.Header.Set("Content-Type", "application/json")
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Get(requestOptions.Url())
}

func (d *DirectHttpRunner) GetHtml(requestOptions IHtmlRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestOptions.IsRetryOptionSet() {
		d.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		d.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := d.client.R()

	if len(d.defHeaders) > 0 {
		for key, value := range d.defHeaders {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Set(key, value)
		}
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Get(requestOptions.Url())
}

func (d *DirectHttpRunner) GetFile(requestOptions IFileRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestOptions.IsRetryOptionSet() {
		d.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		d.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := d.client.R().SetOutput(requestOptions.FilePath())

	if len(d.defHeaders) > 0 {
		for key, value := range d.defHeaders {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Set(key, value)
		}
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Get(requestOptions.Url())
}

func (d *DirectHttpRunner) PostJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestOptions.IsRetryOptionSet() {
		d.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		d.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := d.client.R()

	if len(d.defHeaders) > 0 {
		for key, value := range d.defHeaders {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsValueSet() {
		request.SetBody(requestOptions.Value())
	}

	if len(request.Header.Get("Content-Type")) <= 0 {
		request.Header.Set("Content-Type", "application/json")
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Post(requestOptions.Url())
}

func (d *DirectHttpRunner) PutJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestOptions.IsRetryOptionSet() {
		d.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		d.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := d.client.R()

	if len(d.defHeaders) > 0 {
		for key, value := range d.defHeaders {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsValueSet() {
		request.SetBody(requestOptions.Value())
	}

	if len(request.Header.Get("Content-Type")) <= 0 {
		request.Header.Set("Content-Type", "application/json")
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Put(requestOptions.Url())
}

func (d *DirectHttpRunner) PostForm(requestOptions IFormRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			d.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestOptions.IsRetryOptionSet() {
		d.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		d.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := d.client.R()

	if len(d.defHeaders) > 0 {
		for key, value := range d.defHeaders {
			request.Header.Set(key, value)
		}
	}

	if requestOptions.IsFilesSet() {
		for key, value := range requestOptions.Files() {
			request.SetFileReader(key, value.fileName, value.reader)
		}
	}

	if requestOptions.IsValuesSet() {
		request.SetFormData(requestOptions.Values())
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Set(key, value)
		}
	}

	if len(request.Header.Get("Content-Type")) <= 0 {
		request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Post(requestOptions.Url())
}
