package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/app"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/db"
	"github.com/DannyTuanAnh/appointment-scheduler-application/internal/utils"
)

func main() {
	utils.LoadEnv()

	// 1. Initialize original context for the application
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	// 2. Initialize database connection
	if err := db.InitDB(); err != nil {
		log.Fatal("Failed to initialize database:", err)
		return
	}
	defer db.Close()

	// 3. Initialize application
	application := app.NewApplication(ctx)

	// 4. Run the application and capture any error message
	msg, err := application.Run(ctx)
	if err != nil {
		log.Fatalf("%s: %v", msg, err)
	}

	log.Println(msg)
}
