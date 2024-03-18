package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"social_media/collections"
	"social_media/controller"
	"social_media/library"
	"social_media/repository"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func main() {

	context, err := NewContext()
	if err != nil {
		fmt.Println("new context : %w", err.Error())
	}

	// set route
	app := fiber.New(fiber.Config{})

	// Use recover middleware to prevent crashes
	app.Use(recover.New())
	app.Get("/metrics", adaptor.HTTPHandler(promhttp.Handler()))

	// use this route to test
	NewRoute(app, context, "/", "GET", false, func(c *fiber.Ctx) (int, string, interface{}, interface{}, error) {
		return http.StatusOK, "running an api at port 8080", nil, collections.Meta{}, c.JSON(fiber.Map{})
	})
	// actual route
	NewRoute(app, context, "/v1/user/register", "POST", false, context.CTL.User.Register)
	NewRoute(app, context, "/v1/user/login", "POST", false, context.CTL.User.Login)

	NewRoute(app, context, "/v1/user", "PATCH", true, context.CTL.User.UpdateProfile)
	NewRoute(app, context, "/v1/user/link/email", "POST", true, context.CTL.User.UpdateLinkEmail)
	NewRoute(app, context, "/v1/user/link/phone", "POST", true, context.CTL.User.UpdateLinkPhone)

	//set channel to notify when app interrupted
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		_ = <-c
		fmt.Println("Gracefully shutting down...")
		_ = app.Shutdown()
	}()
	if err := app.Listen(":8080"); err != nil {
		fmt.Println("error on http listen : %w", err)
	}
}

func Response(message string, meta interface{}, data interface{}) []byte {

	response := struct {
		Message string      `json:"message"`
		Data    interface{} `json:"data"`
		Meta    interface{} `json:"meta"`
	}{
		Message: message,
		Data:    data,
		Meta:    meta,
	}

	value, err := json.Marshal(response)
	if err != nil {
		fmt.Println(err.Error())
	}

	return value
}

type Context struct {
	CFG library.Config
	CTL controller.Controller
	JWT library.JWT
	S3  library.S3
}

func NewContext() (Context, error) {
	// read config
	config, err := library.NewConfiguration()
	if err != nil {
		fmt.Println("new config : %w", err)
	}

	// set up jwt
	jwt := library.NewJWT(config.JWTSecret)

	// set up db
	db, err := library.NewDatabaseConnection(config.DB)
	if err != nil {
		fmt.Println("new db : %w", err)
	}

	s3, err := library.NewS3(config.S3Config)
	if err != nil {
		fmt.Println("new s3 : %w", err)
	}

	// set up repo and controller
	repo := repository.NewRepository(db)
	ctl := controller.NewController(repo, jwt, config.BcryptSalt, s3)

	return Context{
		CFG: config,
		JWT: jwt,
		CTL: ctl,
		S3:  s3,
	}, nil
}

// PROMETHEUS n GRAFANA
var (
	RequestHistogram = promauto.NewHistogramVec(prometheus.HistogramOpts{
		Name:    "http_request",
		Help:    "Histogram of the http request duration.",
		Buckets: prometheus.LinearBuckets(1, 1, 10), // Adjust bucket sizes as needed
	}, []string{"path", "method", "status"})
)

func NewRoute(app *fiber.App, ctx Context, path string, method string, useAuth bool, handler func(*fiber.Ctx) (int, string, interface{}, interface{}, error)) {
	if useAuth {
		app.Add(method, path, ctx.JWT.Authentication(), parseContextWithMatrics(path, method, handler))
	} else {
		app.Add(method, path, parseContextWithMatrics(path, method, handler))
	}
}

func parseContextWithMatrics(path string, method string, f func(*fiber.Ctx) (int, string, interface{}, interface{}, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		code, message, resp, meta, err := f(c)

		duration := time.Since(startTime).Seconds()

		statusCode := fmt.Sprintf("%d", c.Response().StatusCode())

		c.Set("Content-Type", "application/json")

		if err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:01:02 "), err)
			errBody := Response(message, meta, nil)
			c.Set("Content-Length", fmt.Sprintf("%d", len(errBody)))
			return c.Status(code).Send(errBody)
		}

		RequestHistogram.WithLabelValues(path, method, statusCode).Observe(duration)
		successBody := Response(message, meta, resp)
		c.Set("Content-Length", fmt.Sprintf("%d", len(successBody)))
		return c.Status(code).Send(successBody)
	}
}

func ParseContext(f func(*fiber.Ctx) (int, string, interface{}, interface{}, error)) func(*fiber.Ctx) error {
	return func(c *fiber.Ctx) error {
		code, message, resp, meta, err := f(c)
		c.Set("Content-Type", "application/json")

		if err != nil {
			fmt.Println(time.Now().Format("2006-01-02 15:01:02 "), err)
			errBody := Response(message, meta, nil)
			c.Set("Content-Length", fmt.Sprintf("%d", len(errBody)))
			return c.Status(code).Send(errBody)
		}

		successBody := Response(message, meta, resp)
		c.Set("Content-Length", fmt.Sprintf("%d", len(successBody)))
		return c.Status(code).Send(successBody)
	}
}
