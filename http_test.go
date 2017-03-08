package utils

import (
	"net/http"
	"testing"
	"time"
)

func Test_HTTPConnectionPool(t *testing.T) {
	pool := NewHTTPConnectionPool(time.Second, 1)
	if pool == nil {
		t.Fail()
	}
}

func Test_HTTPConnectionPoolNagativePoolNum(t *testing.T) {
	pool := NewHTTPConnectionPool(time.Second, -1)
	if pool == nil {
		t.Fail()
	}
}

func Test_HTTPRequestOK(t *testing.T) {
	pool := NewHTTPConnectionPool(2*time.Second, 1)
	if pool == nil {
		t.Fail()
	}
	request, err := http.NewRequest("GET", "http://www.mgtv.com", nil)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	response, err := pool.Request(request)
	if err != nil {
		t.Errorf("request err:%s", err.Error())
		t.Fail()
	}
	if response == nil {
		t.Fail()
	}
	t.Logf(pool.Status())
}

func Test_HTTPRequestNil(t *testing.T) {
	pool := NewHTTPConnectionPool(2*time.Second, 1)
	if pool == nil {
		t.Fail()
	}
	_, err := pool.Request(nil)
	if err == nil {
		t.Errorf("request err:%s", err.Error())
		t.Fail()
	}
	t.Logf(pool.Status())
}

func Test_HTTPRequestNoPool(t *testing.T) {
	pool := NewHTTPConnectionPool(2*time.Second, 0)
	if pool == nil {
		t.Fail()
	}
	request, err := http.NewRequest("GET", "http://www.mgtv.com", nil)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	response, err := pool.Request(request)
	if err == nil {
		t.Fail()
	}
	if response != nil {
		t.Fail()
	}
	t.Logf(pool.Status())
}

func Test_HTTPRequestTimeout(t *testing.T) {
	pool := NewHTTPConnectionPool(0*time.Second, 1)
	if pool == nil {
		t.Fail()
	}
	request, err := http.NewRequest("GET", "http://www.mgtv.com", nil)
	if err != nil {
		t.Error(err.Error())
		t.Fail()
		return
	}
	response, err := pool.Request(request)
	if err == nil {
		t.Fail()
	}
	if response != nil {
		t.Fail()
	}
	t.Logf(pool.Status())
}

func Test_HTTPBatchRequest(t *testing.T) {
	pool := NewHTTPConnectionPool(2*time.Second, 10)
	if pool == nil {
		t.Fail()
	}
	httpdatas := make([]*HTTPData, 0, 10)
	for i := 0; i < 10; i++ {
		request, err := http.NewRequest("GET", "http://www.mgtv.com", nil)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
		httpdatas = append(httpdatas, NewHTTPData(request))
	}
	pool.BatchRequest(httpdatas)
	for _, httpData := range httpdatas {
		if httpData.Err != nil {
			t.Error(httpData.Err.Error())
			t.Fail()
		}
	}
	t.Logf(pool.Status())
}

func Test_HTTPBatchRequestPoolFull(t *testing.T) {
	pool := NewHTTPConnectionPool(2*time.Second, 1)
	if pool == nil {
		t.Fail()
	}
	httpdatas := make([]*HTTPData, 0, 20)
	for i := 0; i < 10; i++ {
		request, err := http.NewRequest("GET", "http://www.mgtv.com", nil)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
		httpdatas = append(httpdatas, NewHTTPData(request))
	}
	pool.BatchRequest(httpdatas)
	errNum := 0
	for _, httpData := range httpdatas {
		if httpData.Err != nil {
			//t.Logf("err:%s", httpData.Err.Error())
			errNum++
		}
	}
	t.Logf("err num:%d", errNum)
	if errNum == 0 {
		t.Fail()
	}
	t.Logf(pool.Status())
}

func Test_HTTPBatchRequestTimeOut(t *testing.T) {
	pool := NewHTTPConnectionPool(1*time.Millisecond, 20)
	if pool == nil {
		t.Fail()
	}
	httpdatas := make([]*HTTPData, 0, 20)
	for i := 0; i < 10; i++ {
		request, err := http.NewRequest("GET", "http://www.mgtv.com", nil)
		if err != nil {
			t.Error(err.Error())
			t.Fail()
			return
		}
		httpdatas = append(httpdatas, NewHTTPData(request))
	}
	pool.BatchRequest(httpdatas)
	errNum := 0
	for _, httpData := range httpdatas {
		if httpData.Err != nil {
			//t.Logf("err:%s", httpData.Err.Error())
			errNum++
		}
	}
	t.Logf("err num:%d", errNum)
	if errNum == 0 {
		t.Fail()
	}
	t.Logf(pool.Status())
}
