package main

import (
	"context"
  "strconv"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis/v8"
	"math/rand"
	"net/http"
)

type App struct {
	router *gin.Engine
	db     *redis.Client
	ctx    context.Context
}

func (app *App) SetupRoutes() {
  app.router.GET("/s/:code", app.GoToLink())
	app.router.POST("/submit/link", app.CreateLink())
	app.router.GET("/", app.Index())
}

func (app *App) CreateLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		link := c.PostForm("link")
    code := strconv.Itoa(rand.Intn(1000))
    err := app.db.Set(app.ctx, code, link, 0).Err()
		if err != nil {
			c.String(http.StatusInternalServerError, "Cannot Save To DB")
		}
		c.String(http.StatusOK, "Code : %s", code)
	}
}

func (app *App) GoToLink() gin.HandlerFunc {
	return func(c *gin.Context) {
		code := c.Param("code")
    link, err := app.db.Get(app.ctx, code).Result()
		if err != nil {
			c.String(http.StatusBadRequest, "Link Not Found")
		}
		c.Redirect(http.StatusFound, link)
	}
}

func (app *App) Index() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{})
	}
}

func main() {
	app := App{}
	app.ctx = context.Background()
	app.router = gin.Default()
	app.router.LoadHTMLGlob("templates/*")
	app.db = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})
  app.SetupRoutes()
	app.router.Run(":5000")
}
