package http_runner

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty/v2"
	"github.com/nadoo/glider/rule"
	"github.com/weppos/publicsuffix-go/publicsuffix"

	log "github.com/sirupsen/logrus"
)

type ProxyHttpRunner struct {
	client *resty.Client
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
	for key, value := range headers {
		client.Header.Add(key, value)
	}
	// CREATE A RESTY CLIENT WITH PROXY

	return &ProxyHttpRunner{client}, nil
}

func NewProxyHttpRunner(dialer *rule.Proxy) (IHttpRunner, error) {
	return NewAdvancedProxyHttpRunner(dialer, 3, time.Second*30, defaultHeaders)
}

func (p *ProxyHttpRunner) GetJson(requestData IJsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestData.IsFollowRedirectOptionSet() {
		if !requestData.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestData.IsRetryOptionSet() {
		p.client.SetRetryCount(requestData.RetryOption())
	}
	if requestData.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestData.TimeoutOption())
	}

	request := p.client.R()

	if requestData.IsHeadersSet() {
		for key, value := range requestData.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "application/json")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.Url())
		if err != nil {
			return nil, err
		}
		domain, err := publicsuffix.Domain(parsedUrl.Host)
		if err != nil {
			return nil, err
		}

		for _, cookie := range cookieJar {
			if strings.HasSuffix(cookie.Domain, domain) {
				request.SetCookies([]*http.Cookie{cookie})
			}
		}
	}

	return request.Get(requestData.Url())
}

func (p *ProxyHttpRunner) GetHtml(requestData IHtmlRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestData.IsFollowRedirectOptionSet() {
		if !requestData.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestData.IsRetryOptionSet() {
		p.client.SetRetryCount(requestData.RetryOption())
	}
	if requestData.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestData.TimeoutOption())
	}

	request := p.client.R()

	if requestData.IsHeadersSet() {
		for key, value := range requestData.Headers() {
			request.Header.Add(key, value)
		}
	}

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.Url())
		if err != nil {
			return nil, err
		}
		domain, err := publicsuffix.Domain(parsedUrl.Host)
		if err != nil {
			return nil, err
		}

		for _, cookie := range cookieJar {
			if strings.HasSuffix(cookie.Domain, domain) {
				request.SetCookies([]*http.Cookie{cookie})
			}
		}
	}

	return request.Get(requestData.Url())
}

func (p *ProxyHttpRunner) GetFile(requestData IFileRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestData.IsFollowRedirectOptionSet() {
		if !requestData.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestData.IsRetryOptionSet() {
		p.client.SetRetryCount(requestData.RetryOption())
	}
	if requestData.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestData.TimeoutOption())
	}

	request := p.client.R().SetOutput(requestData.FilePath())

	if requestData.IsHeadersSet() {
		for key, value := range requestData.Headers() {
			request.Header.Add(key, value)
		}
	}

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.Url())
		if err != nil {
			return nil, err
		}
		domain, err := publicsuffix.Domain(parsedUrl.Host)
		if err != nil {
			return nil, err
		}

		for _, cookie := range cookieJar {
			if strings.HasSuffix(cookie.Domain, domain) {
				request.SetCookies([]*http.Cookie{cookie})
			}
		}
	}

	return request.Get(requestData.Url())
}

func (p *ProxyHttpRunner) PostJson(requestData IJsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestData.IsFollowRedirectOptionSet() {
		if !requestData.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestData.IsRetryOptionSet() {
		p.client.SetRetryCount(requestData.RetryOption())
	}
	if requestData.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestData.TimeoutOption())
	}

	request := p.client.R()

	if requestData.IsValueSet() {
		request.SetBody(requestData.Value)
	}

	if requestData.IsHeadersSet() {
		for key, value := range requestData.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "application/json")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.Url())
		if err != nil {
			return nil, err
		}
		domain, err := publicsuffix.Domain(parsedUrl.Host)
		if err != nil {
			return nil, err
		}

		for _, cookie := range cookieJar {
			if strings.HasSuffix(cookie.Domain, domain) {
				request.SetCookies([]*http.Cookie{cookie})
			}
		}
	}

	return request.Post(requestData.Url())
}

func (p *ProxyHttpRunner) PutJson(requestData IJsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestData.IsFollowRedirectOptionSet() {
		if !requestData.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	} else {
		p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
			return http.ErrUseLastResponse
		}))
	}

	if requestData.IsRetryOptionSet() {
		p.client.SetRetryCount(requestData.RetryOption())
	}
	if requestData.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestData.TimeoutOption())
	}

	request := p.client.R()

	if requestData.IsValueSet() {
		request.SetBody(requestData.Value)
	}

	if requestData.IsHeadersSet() {
		for key, value := range requestData.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "application/json")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.Url())
		if err != nil {
			return nil, err
		}
		domain, err := publicsuffix.Domain(parsedUrl.Host)
		if err != nil {
			return nil, err
		}

		for _, cookie := range cookieJar {
			if strings.HasSuffix(cookie.Domain, domain) {
				request.SetCookies([]*http.Cookie{cookie})
			}
		}
	}

	return request.Put(requestData.Url())
}

func (p *ProxyHttpRunner) PostForm(requestData IFormRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	if requestData.IsFollowRedirectOptionSet() {
		if !requestData.FollowRedirectOption() {
			p.client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error { // disable redirect
				return http.ErrUseLastResponse
			}))
		}
	}

	if requestData.IsRetryOptionSet() {
		p.client.SetRetryCount(requestData.RetryOption())
	}
	if requestData.IsTimeoutOptionSet() {
		p.client.SetTimeout(requestData.TimeoutOption())
	}

	request := p.client.R()

	if requestData.IsValuesSet() {
		request.SetFormData(requestData.Values())
	}

	if requestData.IsHeadersSet() {
		for key, value := range requestData.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "x-www-form-urlencoded")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.Url())
		if err != nil {
			return nil, err
		}
		domain, err := publicsuffix.Domain(parsedUrl.Host)
		if err != nil {
			return nil, err
		}

		for _, cookie := range cookieJar {
			if strings.HasSuffix(cookie.Domain, domain) {
				request.SetCookies([]*http.Cookie{cookie})
			}
		}
	}

	return request.Post(requestData.Url())
}
