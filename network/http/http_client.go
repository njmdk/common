package http

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
	"net/url"
	"strings"
	"sync"
	"time"

	"go.uber.org/zap"

	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/network/ulimit"
	"github.com/njmdk/common/timer"
	"github.com/njmdk/common/utils"
)

var emptyJSONBytes = []byte("{}")

type Client struct {
	client     *http.Client
	url        *url.URL
	pool       *sync.Pool
	dealReq    DealReq
	logger     *logger.Logger
	isMd5Query bool
	ShowLog    bool
	host       string
	UseCookie  bool
	cookie     []*http.Cookie
}

type DealReq func(req *http.Request)

func NewHTTPClientWithTimeout(host string, dealReq DealReq, logger *logger.Logger, isMd5Query bool, timeout time.Duration) (*Client, error) {
	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	if err := ulimit.SetRLimit(); err != nil {
		return nil, err
	}

	return &Client{
		host:   host,
		logger: logger,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: nil,
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					return net.DialTimeout(network, addr, timeout)
				},
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 1000,
				IdleConnTimeout:     time.Minute * 5,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: timeout,
			CheckRedirect: func(req *http.Request, via []*http.Request) error {
				if len(via) >= 100 {
					return errors.New("stopped after 100 redirects")
				}
				return nil
			},
		},
		url: uri,
		pool: &sync.Pool{
			New: func() interface{} {
				return &url.URL{}
			},
		},
		dealReq:    dealReq,
		isMd5Query: isMd5Query,
		ShowLog:    true,
	}, nil
}

func NewHTTPClient(host string, dealReq DealReq, logger *logger.Logger, isMd5Query bool) (*Client, error) {
	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	if err := ulimit.SetRLimit(); err != nil {
		return nil, err
	}

	return &Client{
		host:   host,
		logger: logger,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: nil,
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					return net.DialTimeout(network, addr, time.Second*5)
				},
				MaxIdleConns:        10,
				MaxIdleConnsPerHost: 1000,
				IdleConnTimeout:     time.Minute * 5,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 5,
		},
		url: uri,
		pool: &sync.Pool{
			New: func() interface{} {
				return &url.URL{}
			},
		},
		dealReq:    dealReq,
		isMd5Query: isMd5Query,
		ShowLog:    true,
	}, nil
}

func NewHTTPClientShort(host string, dealReq DealReq, logger *logger.Logger, isMd5Query bool) (*Client, error) {
	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	if err := ulimit.SetRLimit(); err != nil {
		return nil, err
	}

	return &Client{
		host:   host,
		logger: logger,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: nil,
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					return net.DialTimeout(network, addr, time.Second*10)
				},
				DisableKeepAlives: true,
				TLSClientConfig:   &tls.Config{InsecureSkipVerify: true},
			},
			Timeout: time.Second * 10,
		},
		url: uri,
		pool: &sync.Pool{
			New: func() interface{} {
				return &url.URL{}
			},
		},
		dealReq:    dealReq,
		isMd5Query: isMd5Query,
		ShowLog:    true,
	}, nil
}

func NewHTTPProxyClient(host string, dealReq DealReq, logger *logger.Logger, isMd5Query bool, proxyStr string) (*Client, error) {
	uri, err := url.Parse(host)
	if err != nil {
		return nil, err
	}

	if err := ulimit.SetRLimit(); err != nil {
		return nil, err
	}
	proxy, err := url.Parse(proxyStr)
	if err != nil {
		return nil, err
	}
	return &Client{
		host:   host,
		logger: logger,
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxy),
				DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
					return net.DialTimeout(network, addr, time.Second*15)
				},
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 1000,
				TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
			},
		},
		url: uri,
		pool: &sync.Pool{
			New: func() interface{} {
				return &url.URL{}
			},
		},
		dealReq:    dealReq,
		isMd5Query: isMd5Query,
		ShowLog:    true,
	}, nil
}

func (this_ *Client) getURLStruct() *url.URL {
	u := this_.pool.Get().(*url.URL)
	u.Host = this_.url.Host
	u.Scheme = this_.url.Scheme
	u.Path = this_.url.Path
	u.RawQuery = ""

	return u
}

func (this_ *Client) SetUseCookie(b bool) {
	this_.UseCookie = b
}

func (this_ *Client) putURLStruct(v *url.URL) {
	this_.pool.Put(v)
}

func (this_ *Client) GetHost() string {
	return this_.host
}

func (this_ *Client) POSTForm1(uri string, queryParams *utils.URLValues, header *utils.HTTPHeader, values *utils.URLValues) (respData []byte, err error) {
	start := timer.Now()
	var urlInfo *url.URL

	defer func() {
		if err != nil {
			this_.logger.Error("http client do: error", zap.Error(err),
				zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
				zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", values.Encode()),
				zap.ByteString("respData", respData))
		} else {
			if this_.ShowLog {
				this_.logger.Info("http client do: success",
					zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
					zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", values.Encode()),
					zap.ByteString("respData", respData))
			}
		}
	}()

	urlInfo = this_.getURLStruct()
	defer this_.putURLStruct(urlInfo)

	urlInfo.Path = uri
	if queryParams.Len() > 0 {
		urlInfo.RawQuery = queryParams.Encode()
	}
	method := "POST"
	var req *http.Request
	if values == nil {
		req, err = http.NewRequest(method, urlInfo.String(), nil)
	} else {
		req, err = http.NewRequest(method, urlInfo.String(), strings.NewReader(values.Encode()))
	}
	if err != nil {
		return nil, err
	}

	if this_.dealReq != nil {
		this_.dealReq(req)
	}

	if header.Len() > 0 {
		for k, v := range header.Header {
			for _, v1 := range v {
				req.Header.Add(k, v1)
			}
		}
	}
	if this_.UseCookie {
		for _, v := range this_.cookie {
			req.AddCookie(v)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resp *http.Response

	resp, err = this_.client.Do(req)
	if err != nil {
		return nil, err
	}
	if this_.UseCookie {
		cookies := resp.Cookies()
		if len(cookies) > 0 {
			this_.cookie = cookies
		}
	}
	data, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	return data, err
}

func (this_ *Client) POSTForm(uri string, queryParams *utils.URLValues, header *utils.HTTPHeader, values *utils.URLValues) (respData []byte, err error) {
	start := timer.Now()
	var urlInfo *url.URL

	defer func() {
		if err != nil {
			this_.logger.Error("http client do: error", zap.Error(err),
				zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
				zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", values.Encode()),
				zap.ByteString("respData", respData))
		} else {
			if this_.ShowLog {
				this_.logger.Info("http client do: success",
					zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
					zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", values.Encode()),
					zap.ByteString("respData", respData))
			}
		}
	}()

	urlInfo = this_.getURLStruct()
	defer this_.putURLStruct(urlInfo)

	urlInfo.Path = uri
	if queryParams.Len() > 0 {
		if this_.isMd5Query {
			sign := utils.MD5String(queryParams.Encode())
			queryParams.Add("sign", sign)
		}
		urlInfo.RawQuery = queryParams.Encode()
	}
	method := "POST"
	var req *http.Request
	if values == nil {
		req, err = http.NewRequest(method, urlInfo.String(), nil)
	} else {
		req, err = http.NewRequest(method, urlInfo.String(), strings.NewReader(values.Encode()))
	}
	if err != nil {
		return nil, err
	}

	if this_.dealReq != nil {
		this_.dealReq(req)
	}

	if header.Len() > 0 {
		for k, v := range header.Header {
			for _, v1 := range v {
				req.Header.Add(k, v1)
			}
		}
	}
	if this_.UseCookie {
		for _, v := range this_.cookie {
			req.AddCookie(v)
		}
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	var resp *http.Response

	resp, err = this_.client.Do(req)
	if err != nil {
		return nil, err
	}
	if this_.UseCookie {
		cookies := resp.Cookies()
		if len(cookies) > 0 {
			this_.cookie = cookies
		}
	}
	data, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	return data, err
}

func (this_ *Client) DoRaw(method, uri string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData []byte) (respData []byte, err error) {
	start := timer.Now()

	var urlInfo *url.URL

	defer func() {
		if err != nil {
			this_.logger.Error("http client do: error", zap.Error(err), zap.String("method", method),
				zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
				zap.Any("queryParams", queryParams), zap.Any("header", header), zap.ByteString("bodyData", bodyData),
				zap.ByteString("respData", respData))
		} else {
			if this_.ShowLog {
				this_.logger.Info("http client do: success", zap.String("method", method),
					zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
					zap.Any("queryParams", queryParams), zap.Any("header", header), zap.ByteString("bodyData", bodyData),
					zap.ByteString("respData", respData))
			}
		}
	}()

	urlInfo = this_.getURLStruct()
	defer this_.putURLStruct(urlInfo)

	urlInfo.Path = uri
	if queryParams.Len() > 0 {
		if this_.isMd5Query {
			sign := utils.MD5String(queryParams.Encode())
			queryParams.Add("sign", sign)
		}
		urlInfo.RawQuery = queryParams.Encode()
	}

	var req *http.Request
	if len(bodyData) == 0 {
		req, err = http.NewRequest(method, urlInfo.String(), nil)
	} else {
		req, err = http.NewRequest(method, urlInfo.String(), bytes.NewBuffer(bodyData))
	}
	if err != nil {
		return nil, err
	}

	if this_.dealReq != nil {
		this_.dealReq(req)
	}

	if header.Len() > 0 {
		for k, v := range header.Header {
			for _, v1 := range v {
				req.Header.Add(k, v1)
			}
		}
	}
	if this_.UseCookie {
		for _, v := range this_.cookie {
			req.AddCookie(v)
		}
	}
	var resp *http.Response

	resp, err = this_.client.Do(req)
	if err != nil {
		return nil, err
	}
	if this_.UseCookie {
		cookies := resp.Cookies()
		if len(cookies) > 0 {
			this_.cookie = cookies
		}
	}
	data, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	return data, err
}

func (this_ *Client) GETWithStatus(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, expectedStatus int) ([]byte, error) {
	return this_.DoWithStatus("GET", url, queryParams, header, bodyData, expectedStatus)
}

type ErrorStatus struct {
	StatusCode int
	Status     string
}

func (this_ *ErrorStatus) Error() string {
	return fmt.Sprintf("http response Status:%s,StatusCode:%d", this_.Status, this_.StatusCode)
}

func (this_ *Client) DoWithStatus(method, uri string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, expectedStatus int) (respData []byte, err error) {
	start := timer.Now()

	var urlInfo *url.URL

	defer func() {
		if err != nil {
			this_.logger.Error("http client do: error", zap.Error(err), zap.String("method", method),
				zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
				zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData),
				zap.ByteString("respData", respData))
		} else {
			if this_.ShowLog {
				this_.logger.Info("http client do: success", zap.String("method", method),
					zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
					zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData),
					zap.ByteString("respData", respData))
			}
		}
	}()

	urlInfo = this_.getURLStruct()
	defer this_.putURLStruct(urlInfo)

	urlInfo.Path = uri
	if queryParams.Len() > 0 {
		if this_.isMd5Query {
			sign := utils.MD5String(queryParams.Encode())
			queryParams.Add("sign", sign)
		}
		urlInfo.RawQuery = queryParams.Encode()
	} else {
		var urlTemp *url.URL
		urlTemp, err = url.ParseRequestURI(uri)
		if err != nil {
			return
		}
		urlInfo.Path = urlTemp.Path
		urlInfo.RawQuery = urlTemp.RawQuery
	}

	var req *http.Request
	if bodyData == nil {
		req, err = http.NewRequest(method, urlInfo.String(), nil)
	} else {
		reqBodyData := emptyJSONBytes
		if bodyData != nil {
			reqBodyData, err = json.Marshal(bodyData)
			if err != nil {
				return nil, err
			}
		}
		req, err = http.NewRequest(method, urlInfo.String(), bytes.NewBuffer(reqBodyData))
	}
	if err != nil {
		return nil, err
	}

	if this_.dealReq != nil {
		this_.dealReq(req)
	}

	if header.Len() > 0 {
		for k, v := range header.Header {
			for _, v1 := range v {
				req.Header.Add(k, v1)
			}
		}
	}
	if this_.UseCookie {
		for _, v := range this_.cookie {
			req.AddCookie(v)
		}
	}

	var resp *http.Response

	resp, err = this_.client.Do(req)
	if err != nil {
		return nil, err
	}
	if this_.UseCookie {
		cookies := resp.Cookies()
		if len(cookies) > 0 {
			this_.cookie = cookies
		}
	}
	if resp.StatusCode != expectedStatus {
		return nil, &ErrorStatus{resp.StatusCode, resp.Status}
	}
	data, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	return data, err
}

func (this_ *Client) Do(method, uri string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}) (respData []byte, err error) {
	start := timer.Now()

	var urlInfo *url.URL

	defer func() {
		if err != nil {
			this_.logger.Error("http client do: error", zap.Error(err), zap.String("method", method),
				zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
				zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData),
				zap.ByteString("respData", respData))
		} else {
			if this_.ShowLog {
				this_.logger.Info("http client do: success", zap.String("method", method),
					zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
					zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData),
					zap.ByteString("respData", respData))
			}
		}
	}()

	urlInfo = this_.getURLStruct()
	defer this_.putURLStruct(urlInfo)

	urlInfo.Path = uri
	if queryParams.Len() > 0 {
		if this_.isMd5Query {
			sign := utils.MD5String(queryParams.Encode())
			queryParams.Add("sign", sign)
		}
		urlInfo.RawQuery = queryParams.Encode()
	} else {
		var urlTemp *url.URL
		urlTemp, err = url.ParseRequestURI(uri)
		if err != nil {
			return
		}
		urlInfo.Path = urlTemp.Path
		urlInfo.RawQuery = urlTemp.RawQuery
	}

	var req *http.Request
	if bodyData == nil {
		req, err = http.NewRequest(method, urlInfo.String(), nil)
	} else {
		reqBodyData := emptyJSONBytes
		if bodyData != nil {
			reqBodyData, err = json.Marshal(bodyData)
			if err != nil {
				return nil, err
			}
		}
		req, err = http.NewRequest(method, urlInfo.String(), bytes.NewBuffer(reqBodyData))
	}
	if err != nil {
		return nil, err
	}

	if this_.dealReq != nil {
		this_.dealReq(req)
	}

	if header.Len() > 0 {
		for k, v := range header.Header {
			for _, v1 := range v {
				req.Header.Add(k, v1)
			}
		}
	}
	if this_.UseCookie {
		for _, v := range this_.cookie {
			req.AddCookie(v)
		}
	}
	var resp *http.Response

	resp, err = this_.client.Do(req)
	if err != nil {
		return nil, err
	}
	if this_.UseCookie {
		cookies := resp.Cookies()
		if len(cookies) > 0 {
			this_.cookie = cookies
		}
	}
	data, err := ioutil.ReadAll(resp.Body)
	_ = resp.Body.Close()

	return data, err
}

func (this_ *Client) DoReturnResponse(method, uri string, queryParams *utils.URLValues, header *utils.HTTPHeader, cookies []*http.Cookie, bodyData interface{}) (resp *http.Response, err error) {
	start := timer.Now()

	var urlInfo *url.URL

	defer func() {
		if err != nil {
			this_.logger.Error("http client do: error", zap.Error(err), zap.String("method", method),
				zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
				zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData))
		} else {
			if this_.ShowLog {
				this_.logger.Info("http client do: success", zap.String("method", method),
					zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
					zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData))
			}
		}
	}()

	urlInfo = this_.getURLStruct()
	defer this_.putURLStruct(urlInfo)

	urlInfo.Path = uri
	if queryParams.Len() > 0 {
		if this_.isMd5Query {
			sign := utils.MD5String(queryParams.Encode())
			queryParams.Add("sign", sign)
		}
		urlInfo.RawQuery = queryParams.Encode()
	} else {
		var urlTemp *url.URL
		urlTemp, err = url.ParseRequestURI(uri)
		if err != nil {
			return
		}
		urlInfo.Path = urlTemp.Path
		urlInfo.RawQuery = urlTemp.RawQuery
	}

	var req *http.Request
	if bodyData == nil {
		req, err = http.NewRequest(method, urlInfo.String(), nil)
	} else {
		reqBodyData := emptyJSONBytes
		if bodyData != nil {
			reqBodyData, err = json.Marshal(bodyData)
			if err != nil {
				return nil, err
			}
		}
		req, err = http.NewRequest(method, urlInfo.String(), bytes.NewBuffer(reqBodyData))
	}
	if err != nil {
		return nil, err
	}

	if this_.dealReq != nil {
		this_.dealReq(req)
	}

	if header.Len() > 0 {
		for k, v := range header.Header {
			for _, v1 := range v {
				req.Header.Add(k, v1)
			}
		}
	}

	for _, v := range cookies {
		req.AddCookie(v)
	}

	resp, err = this_.client.Do(req)
	if err != nil {
		return nil, err
	}
	return resp, nil
}

func (this_ *Client) DoWithResult(method, uri string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) (err error) {
	start := timer.Now()

	var (
		urlInfo  *url.URL
		respData []byte
	)

	defer func() {
		if this_.logger != nil {
			if err != nil {
				this_.logger.Error("http client do: error", zap.Error(err), zap.String("method", method),
					zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
					zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData),
					zap.ByteString("respData", respData))
			} else {
				if this_.ShowLog {
					this_.logger.Info("http client do: success", zap.String("method", method),
						zap.String("HOST", this_.url.Host), zap.String("urlInfo", urlInfo.String()), zap.Duration("costs", timer.Now().Sub(start)),
						zap.Any("queryParams", queryParams), zap.Any("header", header), zap.Any("bodyData", bodyData),
						zap.ByteString("respData", respData), zap.Any("out", out))
				}
			}
		}
	}()

	urlInfo = this_.getURLStruct()
	defer this_.putURLStruct(urlInfo)

	urlInfo.Path = uri
	if queryParams.Len() > 0 {
		if this_.isMd5Query {
			sign := utils.MD5String(queryParams.Encode())
			queryParams.Add("sign", sign)
			queryParams.Del("app_secret")
		}
		urlInfo.RawQuery = queryParams.Encode()
	}

	var req *http.Request
	if bodyData == nil {
		req, err = http.NewRequest(method, urlInfo.String(), nil)
	} else {
		reqBodyData := emptyJSONBytes
		if bodyData != nil {
			reqBodyData, err = json.Marshal(bodyData)
			if err != nil {
				return err
			}
		}
		req, err = http.NewRequest(method, urlInfo.String(), bytes.NewBuffer(reqBodyData))
	}
	if err != nil {
		return err
	}

	if this_.dealReq != nil {
		this_.dealReq(req)
	}

	if header.Len() > 0 {
		for k, v := range header.Header {
			for _, v1 := range v {
				req.Header.Add(k, v1)
			}
		}
	}
	if this_.UseCookie {
		for _, v := range this_.cookie {
			req.AddCookie(v)
		}
	}
	var resp *http.Response

	resp, err = this_.client.Do(req)
	if err != nil {
		return err
	}
	if this_.UseCookie {
		cookies := resp.Cookies()
		if len(cookies) > 0 {
			this_.cookie = cookies
		}
	}
	if resp.StatusCode != http.StatusOK {
		ioutil.ReadAll(resp.Body)
		return fmt.Errorf("http resp status_code[%d],status[%s]", resp.StatusCode, resp.Status)
	}
	respData, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		_ = resp.Body.Close()
		return err
	}

	_ = resp.Body.Close()

	return json.Unmarshal(respData, out)
}

func (this_ *Client) GET(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}) ([]byte, error) {
	return this_.Do("GET", url, queryParams, header, bodyData)
}

func (this_ *Client) POST(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}) ([]byte, error) {
	return this_.Do("POST", url, queryParams, header, bodyData)
}

func (this_ *Client) HEAD(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}) ([]byte, error) {
	return this_.Do("HEAD", url, queryParams, header, bodyData)
}

func (this_ *Client) OPTIONS(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}) ([]byte, error) {
	return this_.Do("OPTIONS", url, queryParams, header, bodyData)
}

func (this_ *Client) PUT(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}) ([]byte, error) {
	return this_.Do("PUT", url, queryParams, header, bodyData)
}

func (this_ *Client) DELETE(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}) ([]byte, error) {
	return this_.Do("DELETE", url, queryParams, header, bodyData)
}

func (this_ *Client) GETJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) error {
	return this_.DoWithResult("GET", url, queryParams, header, bodyData, out)
}

func (this_ *Client) POSTJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) error {
	return this_.DoWithResult("POST", url, queryParams, header, bodyData, out)
}

func (this_ *Client) PUTJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) error {
	return this_.DoWithResult("PUT", url, queryParams, header, bodyData, out)
}

func (this_ *Client) DELETEJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) error {
	return this_.DoWithResult("DELETE", url, queryParams, header, bodyData, out)
}

func (this_ *Client) HEADJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) error {
	return this_.DoWithResult("HEAD", url, queryParams, header, bodyData, out)
}

func (this_ *Client) OPTIONSJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) error {
	return this_.DoWithResult("HEAD", url, queryParams, header, bodyData, out)
}
