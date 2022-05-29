package httpx

import (
	"io/ioutil"
	"net/http"
	"testing"

	"github.com/sirupsen/logrus"
)

func TestClient_GET(t *testing.T) {
	c := NewClient(&Option{
		Mode:      Http,
		Ip:        "220.181.38.148",
		Port:      80,
		UriPrefix: "",
		HC:        nil,
		Logger:    logrus.StandardLogger(),
	})

	c.GET("", func(resp *http.Response) error {
		body, _ := ioutil.ReadAll(resp.Body)
		logrus.Info("body:", string(body))
		return nil
	})
}
