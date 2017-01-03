package httpclient

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strings"
	"sync"
	"time"
)

/*
1. 增加DNSCache，默认缓存5s
2. 为dial read write等接口增加超时返回机制
*/
//http://stackoverflow.com/questions/16895294/how-to-set-timeout-for-http-get-requests-in-golang
var (
	DefaultClient    *http.Client  = NewTimeoutClient(10*time.Second, 20*time.Second)
	DnsCacheDuration time.Duration = 5 * time.Second
	dnsCache                       = &DnsCache{caches: make(map[string]DnsCacheItem)}
)

type TimeoutConn struct {
	net.Conn
	timeout time.Duration
}

type DnsCacheItem struct {
	IP        string
	CacheTime int64
}

type DnsCache struct {
	sync.RWMutex
	caches map[string]DnsCacheItem
}

// the cache will not remove without a trigger of http get
func (this *DnsCache) Get(addr string) string {
	if DnsCacheDuration <= 0 {
		return addr
	}
	host, port, err := net.SplitHostPort(addr)
	if err != nil {
		return addr
	}
	this.RLock()
	item, ok := this.caches[host]
	this.RUnlock()

	if !ok || time.Now().Unix()-item.CacheTime > int64(DnsCacheDuration/time.Second) {
		go func() {
			netAddr, err := net.ResolveTCPAddr("tcp", addr)
			if err == nil {
				this.Lock()
				this.caches[host] = DnsCacheItem{IP: netAddr.IP.String(), CacheTime: time.Now().Unix()}
				this.Unlock()
			}
		}()
	}
	if ok {
		return fmt.Sprintf("%s:%s", item.IP, port)
	} else {
		return addr
	}
}

func NewTimeoutConn(conn net.Conn, timeout time.Duration) *TimeoutConn {
	return &TimeoutConn{conn, timeout}
}

func (c *TimeoutConn) Read(b []byte) (n int, err error) {
	c.SetReadDeadline(time.Now().Add(c.timeout))
	return c.Conn.Read(b)
}

func (c *TimeoutConn) Write(b []byte) (n int, err error) {
	c.SetWriteDeadline(time.Now().Add(c.timeout))
	return c.Conn.Write(b)
}

func TimeoutDialer(cTimeout time.Duration, rwTimeout time.Duration) func(net, addr string) (c net.Conn, err error) {
	return func(netw, addr string) (net.Conn, error) {
		if DnsCacheDuration > 0 {
			addr = dnsCache.Get(addr)
		}

		d := net.Dialer{Timeout: cTimeout, DualStack: true}
		conn, err := d.Dial(netw, addr)
		if err != nil {
			return nil, err
		}
		return NewTimeoutConn(conn, rwTimeout), nil
	}
}

func NewTimeoutClient(connectTimeout time.Duration, readWriteTimeout time.Duration) *http.Client {
	return &http.Client{
		Transport: &http.Transport{
			Dial:               TimeoutDialer(connectTimeout, readWriteTimeout),
			Proxy:              http.ProxyFromEnvironment,
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
			DisableKeepAlives:  true,
		},
	}
}

func DoGet(urlStr string, params url.Values) ([]byte, error) {
	if params != nil {
		urlStr += "?" + params.Encode()
	}
	resp, err := DefaultClient.Get(urlStr)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, err
}

func DoPost(urlStr string, params url.Values) ([]byte, error) {
	var postReader io.Reader = nil
	if params != nil {
		postReader = strings.NewReader(params.Encode())
	}
	resp, err := DefaultClient.Post(urlStr, "application/x-www-form-urlencoded", postReader)
	if err != nil {
		return nil, err
	}
	body, err := ioutil.ReadAll(resp.Body)
	defer resp.Body.Close()
	if err != nil {
		return nil, err
	}
	return body, err
}

func DoMultiPartPost(urlStr string, params url.Values, files url.Values) ([]byte, error) {
	var err error
	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)
	for key, values := range files {
		fileName := values[0]
		if strings.HasPrefix(fileName, "http://") {
			resp, err := DefaultClient.Get(fileName)
			if err != nil {
				return nil, err
			}
			part, err := writer.CreateFormFile(key, "pic")
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, resp.Body)
			defer resp.Body.Close()
		} else {
			file, err := os.Open(fileName)
			if err != nil {
				return nil, err
			}
			defer file.Close()
			part, err := writer.CreateFormFile(key, filepath.Base(fileName))
			if err != nil {
				return nil, err
			}
			_, err = io.Copy(part, file)
		}
	}
	for key, values := range params {
		_ = writer.WriteField(key, values[0])
	}
	err = writer.Close()
	if err != nil {
		return nil, err
	}

	req, err := http.NewRequest("POST", urlStr, body)
	if err != nil {
		return nil, err
	}
	req.Header.Add("Content-Type", writer.FormDataContentType())
	resp, err := DefaultClient.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body2, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return body2, nil
}

func DoRequest(req *http.Request) (*http.Response, error) {
	return DefaultClient.Do(req)
}
