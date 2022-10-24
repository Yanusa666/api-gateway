package comments

import (
	"crypto/tls"
	"github.com/rs/zerolog"
	"net/http"
	"net/url"
	"sf-api-gateway/internal/config"
	"time"
)

type Comments struct {
	cfg    *config.Config
	lgr    zerolog.Logger
	client *http.Client
	url    *url.URL
}

func NewComments(cfg *config.Config, lgr zerolog.Logger) *Comments {
	lgr = lgr.With().Str("service", "Comments").Logger()

	client := &http.Client{
		Timeout: 5 * time.Second,
		Transport: &http.Transport{
			TLSClientConfig:    &tls.Config{InsecureSkipVerify: true},
			DisableCompression: true,
		},
	}

	_url, err := url.Parse(cfg.Comments.URI)
	if err != nil {
		lgr.Fatal().Err(err).Msg("url.Parse failed")
	}

	return &Comments{
		cfg:    cfg,
		lgr:    lgr,
		client: client,
		url:    _url,
	}
}

func (c *Comments) Shutdown() {
	c.client.CloseIdleConnections()
}
