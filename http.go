package http_runner

import (
	"context"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/go-resty/resty"
	"github.com/nadoo/glider/rule"
	"github.com/weppos/publicsuffix-go/publicsuffix"

	log "github.com/sirupsen/logrus"
)

type IHttpRunner interface {
	getJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	getHtml(requestData HtmlRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	getFile(requestData FileRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	postJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	putJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
	postForm(requestData FormRequestData, cookieJar ...*http.Cookie) (*resty.Response, error)
}

type JsonRequestData struct {
	url     string
	value   []byte
	headers map[string]string
}

type HtmlRequestData struct {
	url     string
	value   []byte
	headers map[string]string
}

type FormRequestData struct {
	url     string
	headers map[string]string
	values  map[string]string
}

type FileRequestData struct {
	url      string
	headers  map[string]string
	filePath string
}

type DirectHttpRunner struct {
	client *resty.Client
}

func NewDirectHttpRunner(dialer *rule.Proxy) (IHttpRunner, error) {
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
	client.SetRetryCount(2)
	client.SetTimeout(15 * time.Second)
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

	//client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
	//	return http.ErrUseLastResponse
	//}))

	client.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//client.Header.Add("accept-encoding", "gzip, deflate, br")
	client.Header.Add("accept-language", "en-US,en;q=0.9")
	client.Header.Add("cache-control", "max-age=0")
	//client.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")
	// CREATE A RESTY CLIENT WITHOUT PROXY

	return &DirectHttpRunner{client}, nil
}

func (d *DirectHttpRunner) getJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := d.client.R()

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "application/json")

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Get(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DirectHttpRunner) getHtml(requestData HtmlRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := d.client.R()

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Get(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DirectHttpRunner) getFile(requestData FileRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := d.client.R().SetOutput(requestData.filePath)

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Get(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DirectHttpRunner) postJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := d.client.R()

	if requestData.value != nil {
		request.SetBody(requestData.value)
	}

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "application/json")

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Post(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DirectHttpRunner) putJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := d.client.R()

	if requestData.value != nil {
		request.SetBody(requestData.value)
	}

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "application/json")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Put(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (d *DirectHttpRunner) postForm(requestData FormRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := d.client.R()

	if len(requestData.values) > 0 {
		request.SetFormData(requestData.values)
	}

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Post(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

type ProxyHttpRunner struct {
	client *resty.Client
}

func NewProxyHttpRunner(dialer *rule.Proxy) (IHttpRunner, error) {
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
	client.SetRetryCount(3)
	client.SetTimeout(30 * time.Second)
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

	client.SetRedirectPolicy(resty.RedirectPolicyFunc(func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}))

	client.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//client.Header.Add("accept-encoding", "gzip, deflate, br")
	client.Header.Add("accept-language", "en-US,en;q=0.9")
	client.Header.Add("cache-control", "max-age=0")
	//client.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")
	// CREATE A RESTY CLIENT WITH PROXY

	return &ProxyHttpRunner{client}, nil
}

func (p *ProxyHttpRunner) getJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := p.client.R()

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "application/json")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Get(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *ProxyHttpRunner) getHtml(requestData HtmlRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := p.client.R()

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Get(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *ProxyHttpRunner) getFile(requestData FileRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := p.client.R().SetOutput(requestData.filePath)

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Get(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *ProxyHttpRunner) postJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := p.client.R()

	if requestData.value != nil {
		request.SetBody(requestData.value)
	}

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "application/json")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Post(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *ProxyHttpRunner) putJson(requestData JsonRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := p.client.R()

	if requestData.value != nil {
		request.SetBody(requestData.value)
	}

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "application/json")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Put(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}

func (p *ProxyHttpRunner) postForm(requestData FormRequestData, cookieJar ...*http.Cookie) (*resty.Response, error) {
	request := p.client.R()

	if len(requestData.values) > 0 {
		request.SetFormData(requestData.values)
	}

	for key, value := range requestData.headers {
		request.Header.Add(key, value)
	}

	request.Header.Add("Content-Type", "x-www-form-urlencoded")
	if len(cookieJar) > 0 {
		parsedUrl, err := url.Parse(requestData.url)
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

	response, err := request.Post(requestData.url)
	if err != nil {
		return nil, err
	}

	return response, nil
}
