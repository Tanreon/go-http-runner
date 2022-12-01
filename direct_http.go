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

type DirectHttpRunner struct {
	client *resty.Client
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

	//client.Header.Add("accept", "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9")
	//client.Header.Add("accept-encoding", "gzip, deflate, br")
	//client.Header.Add("accept-language", "en-US,en;q=0.9")
	//client.Header.Add("cache-control", "max-age=0")
	//client.Header.Add("user-agent", "Mozilla/5.0 (Windows NT 10.0; WOW64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/81.0.4044.113 Safari/537.36")
	for key, value := range headers {
		client.Header.Add(key, value)
	}
	// CREATE A RESTY CLIENT WITHOUT PROXY

	return &DirectHttpRunner{client}, nil
}

func NewDirectHttpRunner(dialer *rule.Proxy) (IHttpRunner, error) {
	return NewAdvancedDirectHttpRunner(dialer, 2, time.Second*15, DefaultHeaders)
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

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "application/json")

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

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Add(key, value)
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

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Add(key, value)
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

	if requestOptions.IsValueSet() {
		request.SetBody(requestOptions.Value())
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "application/json")

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

	if requestOptions.IsValueSet() {
		request.SetBody(requestOptions.Value())
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "application/json")

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

	if requestOptions.IsValuesSet() {
		request.SetFormData(requestOptions.Values())
	}

	if requestOptions.IsFilesSet() {
		for key, value := range requestOptions.Files() {
			request.SetFileReader(key, value.fileName, value.reader)
		}
	}

	if requestOptions.IsHeadersSet() {
		for key, value := range requestOptions.Headers() {
			request.Header.Add(key, value)
		}
	}

	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	if len(cookieJar) > 0 {
		if err := integrateCookies(requestOptions, request, cookieJar); err != nil {
			return nil, err
		}
	}

	return request.Post(requestOptions.Url())
}
