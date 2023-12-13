package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gunsirawich/todo/auth"
	"github.com/gunsirawich/todo/todo"
	"github.com/joho/godotenv"
	"golang.org/x/time/rate"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

func main() {

	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}

	defer os.Remove("/tmp/live")

	err = godotenv.Load("local.env")
	if err != nil {
		log.Println("please consider environment variables: %s", err)
	}

	//dsn := "host=localhost user=postgres password=mysecretpassword dbname=mytest port=5432 sslmode=disable"
	db, err := gorm.Open(postgres.Open(os.Getenv("DB_CONN")), &gorm.Config{})
	if err != nil {
		panic("failed to connect database")
	}

	db.AutoMigrate(&todo.Todo{})

	r := gin.Default()

	r.GET("/healthz", func(c *gin.Context) {
		c.Status(200)
	})

	userhandler := auth.NewUserHandler(db)
	r.POST("/register", userhandler.Register)
	r.POST("/login", userhandler.Login)
	//r.GET("/tokenz", auth.AccessToken(os.Getenv("Sign")))

	protected := r.Group("", auth.Protect([]byte(os.Getenv("SIGN"))))

	handler := todo.NewTodoHandler(db)
	protected.POST("/todos", handler.NewTask)
	protected.GET("/todos", handler.List)
	protected.PUT("/todos/", handler.Update)
	protected.DELETE("/todos", handler.Remove)

	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer stop()

	s := &http.Server{
		Addr:           ":" + os.Getenv("PORT"),
		Handler:        r,
		ReadTimeout:    10 * time.Second,
		WriteTimeout:   10 * time.Second,
		MaxHeaderBytes: 1 << 20,
	}

	go func() {
		if err := s.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("listen: %s\n", err)
		}
	}()

	<-ctx.Done()
	stop()
	fmt.Println("shutting down gracefully, press Ctrl+C again to force")

	timeoutCtx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	if err := s.Shutdown(timeoutCtx); err != nil {
		fmt.Println(err)
	}
}

var limiter = rate.NewLimiter(5, 5)

func limitedHandler(c *gin.Context) {
	if !limiter.Allow() {
		c.AbortWithStatus(http.StatusTooManyRequests)
		return
	}
	c.JSON(200, gin.H{
		"message": "pong",
	})
}
