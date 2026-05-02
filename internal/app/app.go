package app

import (
	"context"
	"log"
	"net/http"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/config"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/metrics"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/routes"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/validation"
	"github.com/gin-gonic/gin"
)

type ModelHTTP interface {
	Routes() routes.Routes
}

type Application struct {
	config  *config.Config
	route   *gin.Engine
	modules []ModelHTTP
}

func NewApplication(ctx context.Context) *Application {
	metrics.Init()

	cfg := config.NewConfigServer()

	// 1. Initialize the Gin router
	r := gin.Default()

	// 2. Initialize custom validator
	err := validation.InitValidator()
	if err != nil {
		log.Fatalf("Failed to initialize validator: %v", err)
		return nil
	}

	// 4. Initialize modules
	modules := []ModelHTTP{
		NewDealershipModule(),
		NewSkillModule(),
		NewServiceBayModule(),
		NewTechnicianModule(),
	}

	// 5. Register all routes from modules by calling the getModuleRoutes helper function to extract the routes from each module
	// and then passing them to the routes.RegisterRoutes function to register them with the Gin router
	routes.RegisterRoutes(ctx, r, getModuleRoutes(modules)...)

	return &Application{
		config:  cfg,
		route:   r,
		modules: modules,
	}
}

func (ac *Application) Run(ctx context.Context) (string, error) {
	// 1. Start server with shut down gracefully
	srv := &http.Server{
		Addr:    ":" + ac.config.Server.Port,
		Handler: ac.route,

		ReadTimeout:       ac.config.Server.ReadTimeout,
		ReadHeaderTimeout: ac.config.Server.ReadHeaderTimeout,
		WriteTimeout:      ac.config.Server.WriteTimeout,
		IdleTimeout:       ac.config.Server.IdleTimeout,

		MaxHeaderBytes: ac.config.Server.MaxHeaderBytes,
	}

	// 2. Create a channel to listen for server errors
	errChan := make(chan error, 1)

	// 3. Listen and serve in a goroutine
	go func() {
		log.Printf("Server is running on port %s...", srv.Addr)
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			errChan <- err
		}
	}()

	// 4. Wait for an error or a shutdown signal
	select {
	case err := <-errChan:
		return "Server error", err

	case <-ctx.Done():
		log.Println("Shutting down server...")
		shutdownCtx, cancel := context.WithTimeout(context.Background(), ac.config.Server.ShutdownTimeout)
		defer cancel()

		if err := srv.Shutdown(shutdownCtx); err != nil {
			return "Server forced to shutdown", err
		}
		return "Server exiting gracefully!", nil
	}

}

// getModuleRoutes is a helper function that takes a slice of Model interfaces
// and returns a slice of routes.Routes by calling the Routes() method on each module
func getModuleRoutes(models []ModelHTTP) []routes.Routes {
	routeList := make([]routes.Routes, len(models))
	for i, model := range models {
		routeList[i] = model.Routes()
	}
	return routeList
}
