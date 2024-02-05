package Proxy

import (
	"log"
	"net/http"
	"net/url"
	"time"
)

// Proxy 设置代理，并返回一个http客户端
func Proxy() *http.Client {
	proxyAddr := "http://127.0.0.1:7890"
	proxy, err := url.Parse(proxyAddr)
	if err != nil {
		log.Fatal(err)
	}

	netTransport := &http.Transport{
		Proxy:                 http.ProxyURL(proxy),
		MaxIdleConnsPerHost:   10,
		ResponseHeaderTimeout: time.Second * time.Duration(5),
	}

	httpClient := &http.Client{
		Timeout:   time.Second * 10,
		Transport: netTransport,
	}
	return httpClient
}
