package utils

import (
	"errors"
	"fmt"
	"net/http"
	"sync/atomic"
	"time"
)

var (
	errorRequestPoolFull    = errors.New("ERROR_HTTP_REQUEST_POOL_FULL")
	errorRequestCallTimeout = errors.New("ERROR_HTTP_REQUEST_TIMEOUNT")
	errorRequestNil         = errors.New("ERROR_HTTP_REQUEST_NIL")
	defaultPoolNum          = 100
)

//HTTPData http请求和响应
type HTTPData struct {
	Request   *http.Request
	Response  *http.Response
	Err       error
	ExtraData interface{} //http 请求的自定义信息
	ended     chan bool
}

//NewHTTPData HTTPData constructor
func NewHTTPData(request *http.Request) *HTTPData {
	return &HTTPData{Request: request, Response: nil, Err: nil, ended: make(chan bool, 1)}
}

//NewHTTPDataWithExtra HTTPData constructor with extra data
func NewHTTPDataWithExtra(request *http.Request, extra interface{}) *HTTPData {
	return &HTTPData{Request: request, Response: nil, Err: nil, ended: make(chan bool, 1), ExtraData: extra}
}

//HTTPConnectionPool http连接池
type HTTPConnectionPool struct {
	name        string
	timeout     time.Duration  //超时时间
	poolNum     int            //连接池数目
	requestPool chan *HTTPData //连接池
	httpClient  *http.Client
	timeoutNum  int64 //超时请求次数
	poolFullNum int64 //连接池满次数
	totalNum    int64 //总请求次数
}

//NewHTTPConnectionPool http连接池构造函数
func NewHTTPConnectionPool(timeout time.Duration, poolNum int) *HTTPConnectionPool {
	if poolNum < 0 {
		poolNum = defaultPoolNum
	}
	pool := new(HTTPConnectionPool)
	pool.name = "default"
	pool.timeout = timeout
	pool.poolNum = poolNum
	pool.requestPool = make(chan *HTTPData, poolNum)
	pool.httpClient = &http.Client{
		Transport: &http.Transport{
			MaxIdleConnsPerHost: poolNum * 6,
		},
		Timeout: timeout,
	}
	go func() {
		for i := 0; i < poolNum; i++ {
			go pool.newWorker()
			time.Sleep(time.Millisecond * 5) //avoid connection  frequency limit
		}
	}()
	return pool
}

// SetName 设置连接池名字，方便统计
func (cp *HTTPConnectionPool) SetName(name string) {
	cp.name = name
}

func (cp *HTTPConnectionPool) newWorker() {
	for {

		request := <-cp.requestPool
		if request.Request != nil {
			request.Response, request.Err = cp.httpClient.Do(request.Request)
		} else {
			request.Response, request.Err = nil, errorRequestNil
		}
		request.ended <- true
	}
}

// Request http请求接口
func (cp *HTTPConnectionPool) Request(request *http.Request) (*http.Response, error) {
	atomic.AddInt64(&cp.totalNum, 1)
	httpData := NewHTTPData(request)
	select {
	case cp.requestPool <- httpData:
	default:
		atomic.AddInt64(&cp.poolFullNum, 1)
		return nil, errorRequestPoolFull
	}
	select {
	case <-httpData.ended:
		return httpData.Response, httpData.Err
	case <-time.After(cp.timeout):
		atomic.AddInt64(&cp.timeoutNum, 1)
		return nil, errorRequestCallTimeout
	}
}

// BatchRequest http批量请求接口
func (cp *HTTPConnectionPool) BatchRequest(httpDatas []*HTTPData) {
	for _, httpData := range httpDatas {
		atomic.AddInt64(&cp.totalNum, 1)
		select {
		case cp.requestPool <- httpData:
		default:
			atomic.AddInt64(&cp.poolFullNum, 1)
			httpData.Response = nil
			httpData.Err = errorRequestPoolFull
			httpData.ended <- true
		}
	}
	for _, httpData := range httpDatas {
		select {
		case <-httpData.ended:
		case <-time.After(cp.timeout):
			atomic.AddInt64(&cp.timeoutNum, 1)
			httpData.Response = nil
			httpData.Err = errorRequestCallTimeout
		}
	}
}

//Status 获取连接池状态并初始化状态
func (cp *HTTPConnectionPool) Status() string {
	totalPoolNum := cp.poolNum
	poolNum := len(cp.requestPool)
	totalNum := cp.totalNum
	poolFullNum := cp.poolFullNum
	timeoutNum := cp.timeoutNum
	atomic.StoreInt64(&cp.totalNum, 0)
	atomic.StoreInt64(&cp.poolFullNum, 0)
	atomic.StoreInt64(&cp.timeoutNum, 0)
	return fmt.Sprintf("HTTPConnectionPool Status: name=%s, totalPoolNum=%d, usedPoolNum=%d, totalNum=%d, poolFullNum=%d, timeoutNum=%d",
		cp.name, totalPoolNum, poolNum, totalNum, poolFullNum, timeoutNum)
}
