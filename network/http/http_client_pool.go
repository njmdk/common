package http

import (
	"encoding/json"
	"errors"
	"sort"
	"sync"
	"sync/atomic"
	"time"

	"github.com/njmdk/common/logger"
	"github.com/njmdk/common/utils"
)

type errorCountDomain struct {
	*Client
	ErrorCount int64
}

type ErrorCountDomainSort []*errorCountDomain

func (this_ ErrorCountDomainSort) Len() int {
	return len(this_)
}

func (this_ ErrorCountDomainSort) Less(i, j int) bool {
	ic := atomic.LoadInt64(&this_[i].ErrorCount)
	jc := atomic.LoadInt64(&this_[j].ErrorCount)
	return ic < jc
}
func (this_ ErrorCountDomainSort) Swap(i, j int) {
	this_[i], this_[j] = this_[j], this_[i]
}

type DomainError struct {
	Domain string `json:"domain"`
	Err    error  `json:"err"`
}

type DomainErrors []*DomainError

func (this_ DomainErrors) Error() string {
	d, _ := json.Marshal(this_)
	return string(d)
}

type ClientPool struct {
	pool *sync.Map
	log  *logger.Logger
}

func NewPool(urls []string, log *logger.Logger, debug bool) (*ClientPool, error) {
	if len(urls) == 0 {
		return nil, errors.New("invalid urls")
	}
	pool := &sync.Map{}
	for _, v := range urls {
		c, err := NewHTTPClientWithTimeout(v, nil, log, false, time.Second*5)
		if err != nil {
			return nil, err
		}
		if !debug {
			c.ShowLog = false
		}
		pool.Store(v, &errorCountDomain{Client: c})
	}
	return &ClientPool{
		pool: pool,
		log:  log,
	}, nil
}

func (this_ *ClientPool) GetJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, bodyData interface{}, out interface{}) error {
	var cs ErrorCountDomainSort
	this_.pool.Range(func(key, value interface{}) bool {
		cs = append(cs, value.(*errorCountDomain))
		return true
	})
	sort.Sort(cs)
	var errs DomainErrors
	for _, v := range cs {
		err := v.GETJsonResponse(url, queryParams, header, bodyData, out)
		if err == nil {
			return nil
		}
		errs = append(errs, &DomainError{Domain: v.GetHost(), Err: err})
		cv, ok := this_.pool.Load(v.GetHost())
		if ok {
			c, ok := cv.(errorCountDomain)
			if ok {
				atomic.AddInt64(&c.ErrorCount, 1)
			}
		}
	}
	return errs
}

func (this_ *ClientPool) PostFormJsonResponse(url string, queryParams *utils.URLValues, header *utils.HTTPHeader, values *utils.URLValues, out interface{}) error {
	var cs ErrorCountDomainSort
	this_.pool.Range(func(key, value interface{}) bool {
		cs = append(cs, value.(*errorCountDomain))
		return true
	})
	sort.Sort(cs)
	var errs DomainErrors
	for _, v := range cs {
		respData, err := v.POSTForm(url, queryParams, header, values)
		if err == nil {
			return json.Unmarshal(respData, out)
		}
		errs = append(errs, &DomainError{Domain: v.GetHost(), Err: err})
		cv, ok := this_.pool.Load(v.GetHost())
		if ok {
			c, ok := cv.(errorCountDomain)
			if ok {
				atomic.AddInt64(&c.ErrorCount, 1)
			}
		}
	}
	return errs
}
