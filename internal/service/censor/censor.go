package censor

import (
	"crypto/tls"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
	"sf-api-gateway/internal/config"
	"time"
)

type Censor struct {
	cfg    *config.Config
	lgr    zerolog.Logger
	client *http.Client
	url    *url.URL
}

func NewCensor(cfg *config.Config, lgr zerolog.Logger) *Censor {
	lgr = lgr.With().Str("service", "Censor").Logger()

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		},
	}

	_url, err := url.Parse(cfg.Censor.URI)
	if err != nil {
		lgr.Fatal().Err(err).Msg("url.Parse failed")
	}

	return &Censor{
		cfg:    cfg,
		lgr:    lgr,
		client: client,
		url:    _url,
	}
}

func (c *Censor) Shutdown() {
	c.client.CloseIdleConnections()
}
