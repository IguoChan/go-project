package esx

import (
	"errors"
	"log"

	es7 "github.com/olivere/elastic/v7"
	"github.com/sirupsen/logrus"
)

type Client struct {
	*es7.Client
	opt *Options
}

type Options struct {
	Addrs    string
	Username string
	Password string

	// logger
	LogLogger    *log.Logger
	LogrusLogger *logrus.Logger
}

func NewClient(opt *Options) (*Client, error) {
	if opt == nil {
		return nil, errors.New("options is nil")
	}

	// logger
	var logger es7.Logger
	if opt.LogLogger != nil {
		logger = opt.LogLogger
	} else if opt.LogrusLogger != nil {
		logger = opt.LogrusLogger
	} else {
		logger = logrus.StandardLogger()
	}

	// client
	logger.Printf("test %+v", opt.Addrs)
	es, err := es7.NewClient(
		es7.SetURL(opt.Addrs),
		es7.SetBasicAuth(opt.Username, opt.Password),
		es7.SetErrorLog(logger),
		es7.SetInfoLog(logger),
		es7.SetTraceLog(logger),
		es7.SetSniff(false),
	)
	if err != nil {
		return nil, err
	}

	return &Client{
		Client: es,
		opt:    opt,
	}, nil
}
