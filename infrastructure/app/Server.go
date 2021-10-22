package app

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
)

const (
	gracefulShutdownRequestGraceSeconds = 10
)

type ApplicationState struct {
	Handler    *gin.Engine
	HTTPServer *http.Server
	DB         *sql.DB
	Config     *ApplicationConfig
}

type ApplicationServer struct {
	State ApplicationState
}

func NewApplicationServer(userOptions *ApplicationState) *ApplicationServer {
	ConfigSetup("config", ".")

	http.DefaultClient.Timeout = 30 * time.Second

	state := userOptions
	if state == nil {
		state = &ApplicationState{}
	}

	// if strings.EqualFold(state.Config.ApplicationMode, "production") {
	// 	gin.SetMode(gin.ReleaseMode)
	// }

	state.Handler = gin.Default()

	if state.Config == nil {
		state.Config = GetConfig()
	}

	if state.HTTPServer == nil {
		state.HTTPServer = &http.Server{
			ReadTimeout:  10 * time.Second,
			WriteTimeout: 60 * time.Second,
			IdleTimeout:  60 * time.Second,
			Addr:         ":7000",
			Handler:      state.Handler,
		}
	}

	if state.DB == nil {
		db, err := sql.Open("sqlite3", state.Config.DatabaseConnection)
		if err != nil {
			log.Fatal(err)
		}

		state.DB = db
	}

	srv := ApplicationServer{
		State: ApplicationState{
			HTTPServer: state.HTTPServer,
			Handler:    state.Handler,
			// Config:     state.Config,
		},
	}

	srv.registerHandlers()

	return &srv
}

func (s *ApplicationServer) registerHandlers() {
	s.State.Handler.Handle(http.MethodGet, "/products")
	s.State.Handler.Handle(http.MethodPost, "/products")
	s.State.Handler.Handle(http.MethodPatch, "/products")
	s.State.Handler.Handle(http.MethodDelete, "/products")

	s.State.Handler.Handle(http.MethodGet, "/articles")
	s.State.Handler.Handle(http.MethodPost, "/articles")
	s.State.Handler.Handle(http.MethodPatch, "/articles")
	s.State.Handler.Handle(http.MethodDelete, "/articles")
}

// TerminationSignalWatcher will wait for interrupt signal to gracefully shutdown
// the server with a timeout of x seconds.
func TerminationSignalWatcher(srv *http.Server) {
	// Make a channel to receive operating system signals.
	quit := make(chan os.Signal, 1)

	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM, syscall.SIGABRT)
	<-quit
	log.Println("Shutting down server, because of received signal..")

	// The context is used to inform the server it has x seconds to finish
	// the request it is currently handling
	ctx, cancel := context.WithTimeout(
		context.Background(),
		gracefulShutdownRequestGraceSeconds*time.Second,
	)

	if err := srv.Shutdown(ctx); err != nil {
		log.Fatal("Server forced to shutdown:", err)
	}

	defer cancel()

	log.Println("Server exiting..")
}
