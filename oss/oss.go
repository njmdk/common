package oss

import (
	"context"
	"crypto/tls"
	"errors"
	"fmt"
	"io"
	"net"
	"net/http"
	"time"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
)

type OSS struct {
	Endpoint        string `json:"endpoint"`
	AccessKeyID     string `json:"access_key_id"`
	AccessKeySecret string `json:"access_key_secret"`
	Bucket          string `json:"bucket"`
	ClientAddr      string `json:"client_addr"`
	client          *oss.Client
	bucket          *oss.Bucket
}

func NewOSS(endpoint, accessKeyID, accessKeySecret string) (*OSS, error) {
	c, err := oss.New(endpoint, accessKeyID, accessKeySecret)
	if err != nil {
		return nil, err
	}
	return &OSS{
		Endpoint:        endpoint,
		AccessKeyID:     accessKeyID,
		AccessKeySecret: accessKeySecret,
		client:          c,
	}, nil
}

func (this_ *OSS) New() error {
	client := &http.Client{
		Transport: &http.Transport{
			Proxy: nil,
			DialContext: func(ctx context.Context, network, addr string) (conn net.Conn, e error) {
				return net.DialTimeout(network, addr, time.Second*10)
			},
			MaxIdleConns:        10,
			MaxIdleConnsPerHost: 1000,
			IdleConnTimeout:     time.Minute * 5,
			TLSClientConfig:     &tls.Config{InsecureSkipVerify: true},
		},
		Timeout: time.Second * 10,
		CheckRedirect: func(req *http.Request, via []*http.Request) error {
			if len(via) >= 100 {
				return errors.New("stopped after 100 redirects")
			}
			return nil
		},
	}
	c, err := oss.New(this_.Endpoint, this_.AccessKeyID, this_.AccessKeySecret, oss.HTTPClient(client))
	if err != nil {
		return err
	}
	this_.client = c
	b, err := this_.client.Bucket(this_.Bucket)
	if err != nil {
		return err
	}
	this_.bucket = b
	return nil
}

func (this_ *OSS) IsObjectExists(uri string) (bool, error) {
	return this_.bucket.IsObjectExist(uri)
}

func (this_ *OSS) CheckAndPut(download string, uri string) error {
	ok, err := this_.IsObjectExists(uri)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	resp, err := this_.client.HTTPClient.Get(download)
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return errors.New(fmt.Sprintf("download get %d", resp.StatusCode))
	}
	err = this_.bucket.PutObject(uri, resp.Body)
	return err
}

func (this_ *OSS) Put(reader io.Reader, uri string) error {
	ok, err := this_.IsObjectExists(uri)
	if err != nil {
		return err
	}
	if ok {
		return nil
	}
	err = this_.bucket.PutObject(uri, reader)
	return err
}