package profile

import (
	"context"
	"fmt"
	"net/http"
	"net/http/pprof"

	"github.com/Meland-Inc/game-services/src/common/serviceLog"
	gkitLog "github.com/go-kit/kit/log"
)

type HTTPServerOpt func(*HTTPServer)

type HTTPServer struct {
	host     string
	port     int32
	stopChan chan chan struct{}
	logger   gkitLog.Logger
}

func WithHTTPServerLogger(rootLogger gkitLog.Logger) HTTPServerOpt {
	return func(s *HTTPServer) {
		s.logger = gkitLog.With(
			rootLogger,
			"component", "profile.http",
		)
	}
}

// 创建 HTTP 服务器实例
func NewHTTPServer(host string, port int32, opts ...HTTPServerOpt) *HTTPServer {
	s := &HTTPServer{
		stopChan: make(chan chan struct{}, 1),
	}

	for _, opt := range opts {
		opt(s)
	}
	return s
}

func (s *HTTPServer) Start() error {
	serviceLog.Info("========  pprof-http try ========")
	smux := http.NewServeMux()

	smux.HandleFunc("/debug/pprof/", pprof.Index)
	smux.HandleFunc("/debug/pprof/cmdline", pprof.Cmdline)
	smux.HandleFunc("/debug/pprof/profile", pprof.Profile)
	smux.HandleFunc("/debug/pprof/symbol", pprof.Symbol)
	smux.HandleFunc("/debug/pprof/trace", pprof.Trace)

	var httpHandler http.Handler
	{
		httpHandler = smux
	}

	httpServer := &http.Server{
		Addr:    fmt.Sprintf("%s:%d", s.host, s.port),
		Handler: httpHandler,
	}

	httpServerErrChan := make(chan error)
	go func(httpServerErrChan chan error, httpServer *http.Server) {
		httpServerErrChan <- httpServer.ListenAndServe()
	}(httpServerErrChan, httpServer)

	s.logger.Log(
		"message",
		fmt.Sprintf("http server started at %s", httpServer.Addr),
	)

	serviceLog.Info("========  pprof-http done ========")

	select {
	case c := <-s.stopChan:
		httpServer.Shutdown(context.Background())
		c <- struct{}{}
		return nil
	case err := <-httpServerErrChan:
		return err
	}
}
