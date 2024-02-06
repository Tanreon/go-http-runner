package http_runner

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nadoo/glider/rule"

	log "github.com/sirupsen/logrus"
)

type ProxyHttpRunner struct {
	defHeaders map[string]string
	client     *resty.Client
}

func NewAdvancedProxyHttpRunner(dialer *rule.Proxy, retryCount int, timeout time.Duration, headers map[string]string) (IHttpRunner, error) {
	// CREATE TRANSPORT FOR HTTP
	transport := http.Transport{
		DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, err error) {
			return dialer.NextDialer(addr).Dial(network, addr)
		},
	}
	// CREATE TRANSPORT FOR HTTP

	// CREATE A RESTY CLIENT WITH PROXY
	client := resty.New()
	client.SetTransport(&transport)
	client.SetRetryCount(retryCount)
	client.SetTimeout(timeout) // 30 * time.Second
	client.SetDisableWarn(true)

	if !log.IsLevelEnabled(log.TraceLevel) {
		restyLogger := log.New()
		restyLogger.SetOutput(ioutil.Discard)

		client.SetLogger(restyLogger)
	}

	//// Using raw func into resty.SetRedirectPolicy
	//client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
	//	// Implement your logic here
	//
	//	// return nil for continue redirect otherwise return error to stop/prevent redirect
	//	return nil
	//}))

	//client.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//client.Header.Add("accept-encoding", "gzip, deflate, br")
	//client.Header.Add("accept-language", "en-US,en;q=0.9")
	//client.Header.Add("cache-control", "max-age=0")
	//client.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")

	runner := &ProxyHttpRunner{
		defHeaders: headers,
		client:     client,
	}
	// CREATE A RESTY CLIENT WITH PROXY

	return runner, nil
}

func NewProxyHttpRunner(dialer *rule.Proxy) (IHttpRunner, error) {
	return NewAdvancedProxyHttpRunner(dialer, 3, time.Second*30, DefaultHeaders)
}

func (p *ProxyHttpRunner) GetJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestOptions.IsRetryOptionSet() {
		p.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := p.client.R()

	if len(p.defHeaders) > 0 {
		for key, value := range p.defHeaders {
			request.SetHeaderVerbatim(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.SetHeaderVerbatim(key, value)
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

func (p *ProxyHttpRunner) GetHtml(requestOptions IHtmlRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestOptions.IsRetryOptionSet() {
		p.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := p.client.R()

	if len(p.defHeaders) > 0 {
		for key, value := range p.defHeaders {
			request.SetHeaderVerbatim(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.SetHeaderVerbatim(key, value)
		}
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Get(requestOptions.Url())
}

func (p *ProxyHttpRunner) GetFile(requestOptions IFileRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestOptions.IsRetryOptionSet() {
		p.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := p.client.R().SetOutput(requestOptions.FilePath())

	if len(p.defHeaders) > 0 {
		for key, value := range p.defHeaders {
			request.SetHeaderVerbatim(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.SetHeaderVerbatim(key, value)
		}
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Get(requestOptions.Url())
}

func (p *ProxyHttpRunner) PostJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestOptions.IsRetryOptionSet() {
		p.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := p.client.R()

	if len(p.defHeaders) > 0 {
		for key, value := range p.defHeaders {
			request.SetHeaderVerbatim(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.SetHeaderVerbatim(key, value)
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

func (p *ProxyHttpRunner) PutJson(requestOptions IJsonRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestOptions.IsRetryOptionSet() {
		p.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := p.client.R()

	if len(p.defHeaders) > 0 {
		for key, value := range p.defHeaders {
			request.SetHeaderVerbatim(key, value)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.SetHeaderVerbatim(key, value)
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

func (p *ProxyHttpRunner) PostForm(requestOptions IFormRequestOptions, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestOptions.IsFollowRedirectOptionSet() {
		if !requestOptions.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestOptions.IsRetryOptionSet() {
		p.client.SetRetryCount(requestOptions.RetryOption())
	}
	if requestOptions.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestOptions.TimeoutOption())
	}

	request := p.client.R()

	if len(p.defHeaders) > 0 {
		for key, value := range p.defHeaders {
			request.SetHeaderVerbatim(key, value)
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
			request.SetHeaderVerbatim(key, value)
		}
	}

	if len(request.Header.Get("Content-Type")) <= 0 {
		request.Header.Set("Content-Type", "x-www-form-urlencoded")
	}

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Post(requestOptions.Url())
}
