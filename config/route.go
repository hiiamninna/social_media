package config

import (
	"encoding/json"
	"fmt"
	"net/http"
	"social_media/collections"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/adaptor"
	"github.com/gofiber/fiber/v2/middleware/recover"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func AppRoute() *fiber.App {

	context, err := NewContext()
	if err != nil {
		fmt.Println(time.Now().Format(TIME_FORMAT), "new context : "+err.Error())
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
	// Authentication & Authorization
	NewRoute(app, context, "/v1/user/register", "POST", false, context.CTL.User.Register)
	NewRoute(app, context, "/v1/user/login", "POST", false, context.CTL.User.Login)
	// Update Account
	NewRoute(app, context, "/v1/user", "PATCH", true, context.CTL.User.UpdateProfile)
	// Link Phone Number / Email
	NewRoute(app, context, "/v1/user/link", "POST", true, context.CTL.User.UpdateLinkEmail)
	NewRoute(app, context, "/v1/user/link/phone", "POST", true, context.CTL.User.UpdateLinkPhone)
	// Friends
	NewRoute(app, context, "/v1/friend", "GET", true, context.CTL.Friend.List)
	NewRoute(app, context, "/v1/friend", "POST", true, context.CTL.Friend.Create)
	NewRoute(app, context, "/v1/friend", "DELETE", true, context.CTL.Friend.Delete)
	// Image upload
	NewRoute(app, context, "/v1/image", "POST", true, context.CTL.Image.ImageUpload)
	// Post
	NewRoute(app, context, "/v1/post", "POST", true, context.CTL.Post.Create)
	NewRoute(app, context, "/v1/post", "GET", true, context.CTL.Post.List)
	// Comment
	NewRoute(app, context, "/v1/post/comment", "POST", true, context.CTL.Comment.Create)

	return app
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
		fmt.Println(time.Now().Format(TIME_FORMAT), "marshal : "+err.Error())
	}

	return value
}

func parseContextWithMatrics(path string, method string, f func(*fiber.Ctx) (int, string, interface{}, interface{}, error)) fiber.Handler {
	return func(c *fiber.Ctx) error {
		startTime := time.Now()

		code, message, resp, meta, err := f(c)

		duration := time.Since(startTime).Seconds()

		statusCode := fmt.Sprintf("%d", c.Response().StatusCode())

		c.Set("Content-Type", "application/json")

		if err != nil {
			fmt.Println(time.Now().Format(TIME_FORMAT), err)
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
