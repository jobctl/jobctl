package server

import (
	"context"
	"errors"
	"io/fs"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/dagu-dev/dagu/internal/config"
	"github.com/dagu-dev/dagu/internal/frontend/gen/restapi"
	"github.com/dagu-dev/dagu/internal/logger"
	"github.com/go-openapi/loads"
	flags "github.com/jessevdk/go-flags"

	"github.com/dagu-dev/dagu/internal/frontend/gen/restapi/operations"
	// nolint
	pkgmiddleware "github.com/dagu-dev/dagu/internal/frontend/middleware"

	"github.com/go-chi/chi/v5"
)

type Server struct {
	funcsConfig funcsConfig
	host        string
	port        int
	basicAuth   *BasicAuth
	authToken   *AuthToken
	tls         *config.TLS
	logger      logger.Logger
	server      *restapi.Server
	handlers    []Handler
	assets      fs.FS
}

type NewServerArgs struct {
	Host      string
	Port      int
	BasicAuth *BasicAuth
	AuthToken *AuthToken
	TLS       *config.TLS
	Logger    logger.Logger
	Handlers  []Handler
	AssetsFS  fs.FS

	// Configuration for the frontend
	NavbarColor string
	NavbarTitle string
	APIBaseURL  string
}

type BasicAuth struct {
	Username string
	Password string
}

type AuthToken struct {
	Token string
}

type Handler interface {
	Configure(api *operations.DaguAPI)
}

func New(params NewServerArgs) *Server {
	return &Server{
		host:      params.Host,
		port:      params.Port,
		basicAuth: params.BasicAuth,
		authToken: params.AuthToken,
		tls:       params.TLS,
		logger:    params.Logger,
		handlers:  params.Handlers,
		assets:    params.AssetsFS,
		funcsConfig: funcsConfig{
			NavbarColor: params.NavbarColor,
			NavbarTitle: params.NavbarTitle,
			APIBaseURL:  params.APIBaseURL,
		},
	}
}

func (svr *Server) Shutdown() {
	if svr.server == nil {
		return
	}
	err := svr.server.Shutdown()
	if err != nil {
		svr.logger.Warn("Server shutdown", "error", err)
	}
}

func (svr *Server) Serve(ctx context.Context) (err error) {
	middlewareOptions := &pkgmiddleware.Options{
		Handler: svr.defaultRoutes(chi.NewRouter()),
	}
	if svr.authToken != nil {
		middlewareOptions.AuthToken = &pkgmiddleware.AuthToken{
			Token: svr.authToken.Token,
		}
	}
	if svr.basicAuth != nil {
		middlewareOptions.AuthBasic = &pkgmiddleware.AuthBasic{
			Username: svr.basicAuth.Username,
			Password: svr.basicAuth.Password,
		}
	}
	pkgmiddleware.Setup(middlewareOptions)

	swaggerSpec, err := loads.Analyzed(restapi.SwaggerJSON, "")
	if err != nil {
		svr.logger.Error("failed to load API spec", "error", err)
		return err
	}
	api := operations.NewDaguAPI(swaggerSpec)
	for _, h := range svr.handlers {
		h.Configure(api)
	}

	svr.server = restapi.NewServer(api)
	defer svr.Shutdown()

	svr.server.Host = svr.host
	svr.server.Port = svr.port
	svr.server.ConfigureAPI()

	// Server run context
	serverCtx, serverStopCtx := context.WithCancel(ctx)

	// Listen for syscall signals for process to interrupt/quit
	sig := make(chan os.Signal, 1)
	signal.Notify(
		sig, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT,
	)
	go func() {
		<-sig

		// Trigger graceful shutdown
		err := svr.server.Shutdown()
		if err != nil {
			svr.logger.Error("server shutdown error", "error", err)
		}
		serverStopCtx()
	}()

	if svr.tls != nil {
		svr.server.TLSCertificate = flags.Filename(svr.tls.CertFile)
		svr.server.TLSCertificateKey = flags.Filename(svr.tls.KeyFile)
		svr.server.EnabledListeners = []string{"https"}
		svr.server.TLSHost = svr.host
		svr.server.TLSPort = svr.port
	}

	// Run the server
	err = svr.server.Serve()
	if err != nil && !errors.Is(err, http.ErrServerClosed) {
		svr.logger.Error("server error", "error", err)
	}

	// Wait for server context to be stopped
	<-serverCtx.Done()

	svr.logger.Info("server closed")

	return nil
}
