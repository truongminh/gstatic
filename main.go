package main

import (
	"fmt"
	"net/http"
	"os"

	"gstatic/clog"
)

func main() {
	c := NewConfig()
	static(c)
	proxy(c)
	listen(c)
}

func listen(c Config) {
	var server http.Handler
	var err error
	ssl := c.Server.SSL
	addr := fmt.Sprintf(":%d", c.Server.HTTPPort)
	if ssl.Enabled {
		if _, err := os.Stat(ssl.CertificateChain); err != nil {
			logger.Fatalf("missing certificate chain: %s", err.Error())
		}
		if _, err := os.Stat(ssl.PrivateKey); err != nil {
			logger.Fatalf("missing private key: %s", err.Error())
		}
		logger.Infof("listening on ssl %s", addr)
		err = http.ListenAndServeTLS(addr, ssl.CertificateChain, ssl.PrivateKey, server)
	} else {
		logger.Infof("listening on %s", addr)
		err = http.ListenAndServe(addr, server)
	}
	if err != nil {
		logger.Fatalf("listening error: %s", err.Error())
	}
}

var logger = clog.Prefix("main")
