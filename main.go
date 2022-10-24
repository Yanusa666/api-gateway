package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"sf-api-gateway/internal/config"
	"sf-api-gateway/internal/http_server"
	"sf-api-gateway/internal/http_server/handlers"
	"sf-api-gateway/internal/service/censor"
	"sf-api-gateway/internal/service/comments"
	"sf-api-gateway/internal/service/news"
	"sf-api-gateway/pkg/logger"
	"syscall"
	"time"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	cfg := config.NewConfig()

	lgr, err := logger.NewLogger(os.Stdout, cfg.LogLevel)
	if err != nil {
		log.Fatalln(err)
	}

	lgr = lgr.With().
		CallerWithSkipFrameCount(2).
		Str("app", "sf-api-gateway").
		Logger()

	shutdownCh := make(chan os.Signal, 1)
	signal.Notify(shutdownCh, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT)

	_censor := censor.NewCensor(cfg, lgr)
	_news := news.NewNews(cfg, lgr)
	_comments := comments.NewComments(cfg, lgr)

	handler := handlers.NewHandler(cfg, lgr, _censor, _news, _comments)
	httpServer, listenHTTPErr := http_server.NewServer(cfg, lgr, handler)

mainLoop:
	for {
		select {
		case <-ctx.Done():
			break mainLoop

		case err = <-listenHTTPErr:
			if err != nil {
				lgr.Error().Err(err).Msg("http server error")
				shutdownCh <- syscall.SIGTERM
			}

		case sig := <-shutdownCh:
			lgr.Info().Msgf("shutdown signal received: %s", sig.String())

			if err = httpServer.Shutdown(); err != nil {
				lgr.Error().Err(err).Msg("shutdown http server error")
			}

			_censor.Shutdown()
			_news.Shutdown()
			_comments.Shutdown()

			lgr.Info().Msg("server loop stopped")
			cancel()
			time.Sleep(1 * time.Second)
		}
	}
}
