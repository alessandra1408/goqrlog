package httpclient

import (
	"net/http"

	"github.com/alessandra1408/goqrlog/internal/config"
)

func NewHTTPClient(cfg *config.Config) *http.Client {

	var tr http.RoundTripper = &http.Transport{
		MaxIdleConnsPerHost: 2000,
	}

	return &http.Client{
		Timeout:   cfg.App.HTTPTimeout,
		Transport: tr,
	}
}
