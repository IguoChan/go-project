package httpx

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"net/http"
	"time"

	"github.com/IguoChan/go-project/pkg/util"

	"github.com/sirupsen/logrus"
)

type HttpMode string

const (
	Http  HttpMode = "http"
	Https HttpMode = "https"
)

type Client struct {
	host string
	hc   *http.Client
	*Option
}

type Option struct {
	Mode      HttpMode
	Ip        string
	Port      int
	UriPrefix string
	HC        *http.Client
	Logger    *logrus.Logger
}

func NewClient(opt *Option) *Client {
	hc := opt.HC
	if hc == nil {
		hc = &http.Client{
			Transport: &http.Transport{
				DialContext: (&net.Dialer{
					Timeout:   10 * time.Second,
					KeepAlive: 30 * time.Second,
					DualStack: true,
				}).DialContext,
				ForceAttemptHTTP2:     true,
				MaxConnsPerHost:       50,
				MaxIdleConnsPerHost:   50,
				IdleConnTimeout:       90 * time.Second,
				TLSHandshakeTimeout:   10 * time.Second,
				ExpectContinueTimeout: 1 * time.Second,
			},
			Timeout: time.Minute,
		}
	}
	if opt.Logger != nil {
		opt.Logger.SetFormatter(&logrus.TextFormatter{ForceColors: true, FullTimestamp: true})
	}
	return &Client{
		host:   fmt.Sprintf("%s://%s:%d%s", util.SetIfEmpty(string(opt.Mode), string(Http)), opt.Ip, opt.Port, opt.UriPrefix),
		Option: opt,
		hc:     hc,
	}
}

func (c *Client) GET(uri string, handler func(resp *http.Response) error) error {
	return c.BaseRequest(context.Background(), "GET", uri, nil, handler)
}

func (c *Client) POST(uri string, data interface{}, handler func(resp *http.Response) error) error {
	return c.BaseRequest(context.Background(), "GET", uri, data, handler)
}

func (c *Client) BaseRequest(ctx context.Context, method, uri string, data interface{}, handler func(resp *http.Response) error) error {
	// new request
	url := c.host + uri
	body := io.Reader(nil)
	reqData := ""
	if data != nil {
		d, err := json.Marshal(data)
		if err != nil {
			return err
		}
		reqData = string(d)
		body = bytes.NewBuffer(d)
	}
	req, err := http.NewRequest(method, url, body)
	if err != nil {
		return err
	}

	// do request
	start := time.Now()
	c.printStart(method, url, start)
	resp, err := c.hc.Do(req.WithContext(ctx))
	if err != nil {
		return err
	}
	c.printEnd(method, url, reqData, resp, time.Now(), time.Since(start))
	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		bs, _ := ioutil.ReadAll(resp.Body)
		return HttpError{
			StatusCode: resp.StatusCode,
			Resp:       string(bs),
		}
	}

	// handle resp
	if err = handler(resp); err != nil {
		return err
	}

	return nil
}

func (c *Client) JsonDecode(resp *http.Response, res interface{}) error {
	return json.NewDecoder(resp.Body).Decode(res)
}

func (c *Client) printStart(method, url string, time time.Time) {
	if c.Logger == nil {
		return
	}

	c.Logger.Info(fmt.Sprintf("[HTTP] %v | %s %s %s | %#v start\n",
		time.Format("2006/01/02 - 15:04:05.000"),
		green, method, reset,
		url,
	))
}

func (c *Client) printEnd(method, url, req string, resp *http.Response, time time.Time, spendTime time.Duration) {
	if c.Logger == nil {
		return
	}

	body, _ := ioutil.ReadAll(resp.Body)
	c.Logger.Info(fmt.Sprintf("[HTTP] %v | %s %s %s | %#v | %s %d %s | %v finish\n\trequest: %s\n\tresponse: %s\n",
		time.Format("2006/01/02 - 15:04:05.000"),
		green, method, reset,
		url,
		green, resp.StatusCode, reset,
		spendTime,
		req,
		string(body),
	))
	resp.Body = ioutil.NopCloser(bytes.NewBuffer(body))
}
