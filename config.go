package main

import (
	"flag"
	"os"

	"github.com/BurntSushi/toml"
)

type FileRoute struct {
	Route   string
	Folder  string
	Headers map[string]string
}

type Config struct {
	Server struct {
		HTTPPort int
		SSL      struct {
			Enabled          bool
			PrivateKey       string
			CertificateChain string
		}
	}
	Static []FileRoute
	Proxy  struct {
		Folder  string
		Headers map[string]string
	}
}

// NewConfig parse config from config.toml
func NewConfig() Config {
	c := Config{}
	file := flag.String("conf", "config.toml", "config file")
	flag.Parse()
	logger.Infof("config file %s", *file)
	_, err := toml.DecodeFile(*file, &c)
	if err != nil {
		logger.Fatalf("read config from %s: %s", *file, err.Error())
	}
	for _, route := range c.Static {
		info, err := os.Stat(route.Folder)
		if err != nil {
			logger.Fatalf("get stat %s errors: %s", route.Folder, err.Error())
		}
		if !info.IsDir() {
			logger.Fatalf("path %s is not a folder", route.Folder)
		}
		if route.Headers == nil {
			route.Headers = map[string]string{}
		}
		if _, ok := route.Headers["Cache-Control"]; !ok {
			route.Headers["Cache-Control"] = "public,max-age=360,no-cache"
		}
	}

	if err := os.MkdirAll(c.Proxy.Folder, os.ModePerm); err != nil {
		logger.Fatalf("mkdir %s errors: %s", c.Proxy.Folder, err.Error())
	}
	if c.Proxy.Headers == nil {
		c.Proxy.Headers = map[string]string{}
	}
	if _, ok := c.Proxy.Headers["Cache-Control"]; !ok {
		c.Proxy.Headers["Cache-Control"] = "public,max-age=360,no-cache"
	}

	logger.Infof("config %+v", c)
	return c
}
