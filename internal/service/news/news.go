package news

import (
	"crypto/tls"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
	"sf-api-gateway/internal/config"
	"time"
)

type News struct {
	cfg    *config.Config
	lgr    zerolog.Logger
	client *http.Client
	url    *url.URL
}

func NewNews(cfg *config.Config, lgr zerolog.Logger) *News {
	lgr = lgr.With().Str("service", "News").Logger()

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		},
	}

	_url, err := url.Parse(cfg.News.URI)
	if err != nil {
		lgr.Fatal().Err(err).Msg("url.Parse failed")
	}

	return &News{
		cfg:    cfg,
		lgr:    lgr,
		client: client,
		url:    _url,
	}
}

func (n *News) Shutdown() {
	n.client.CloseIdleConnections()
}
