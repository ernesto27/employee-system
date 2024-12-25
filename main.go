package main

import (
	"context"
	"employees-system/db"
	"employees-system/internal/s3"
	"employees-system/router"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		panic(err)
	}

	// Database
	host := os.Getenv("DB_HOST")
	user := os.Getenv("DB_USERNAME")
	password := os.Getenv("DB_PASSWORD")
	port := os.Getenv("DB_PORT")
	database := os.Getenv("DB_DATABASE")

	dbInstance, err := db.NewMysql(host, user, password, port, database)
	if err != nil {
		panic(err)
	}

	// AWS S3
	awsRegion := os.Getenv("AWS_REGION")
	awsS3Bucket := os.Getenv("AWS_S3_BUCKET")
	awsAccesKey := os.Getenv("AWS_ACCESS_KEY")
	awsSecret := os.Getenv("AWS_SECRET")
	myS3 := s3.NewS3(awsRegion, awsAccesKey, awsSecret, awsS3Bucket)

	url := os.Getenv("URL")

	r := router.GetRouter(dbInstance.Db, myS3, url)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: r,
	}

	go func() {
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			fmt.Printf("listen: %s\n", err)
		}
	}()

	fmt.Println("Starting the server on port 8080")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	<-quit

	fmt.Println("Shutting down server...")

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		fmt.Printf("Server forced to shutdown: %s\n", err)
	}

	fmt.Println("Server exiting")
}
