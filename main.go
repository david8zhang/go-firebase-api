package main

import (
	"context"
	"errors"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/david8zhang/go-firebase/routes"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {

			c.Header("Access-Control-Allow-Origin", "*")
			c.Header("Access-Control-Allow-Credentials", "true")
			c.Header("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
			c.Header("Access-Control-Allow-Methods", "POST,HEAD,PATCH, OPTIONS, GET, PUT")

			if c.Request.Method == "OPTIONS" {
					c.AbortWithStatus(204)
					return
			}

			c.Next()
	}
}

func main() {
	r := gin.Default()
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://habit-rank-default-rtdb.firebaseio.com",
	}
	
	var app *firebase.App
	var err error
	if _, err := os.Stat("habit-rank-firebase-adminsdk-8xo36-576cec4f1c.json"); errors.Is(err, os.ErrNotExist) {
		app, err = firebase.NewApp(ctx, conf)
		if err != nil {
			log.Fatalln("error in initializing firebase app: ", err)
		}
	} else {
		opt := option.WithCredentialsFile("habit-rank-firebase-adminsdk-8xo36-576cec4f1c.json")
		app, err = firebase.NewApp(ctx, conf, opt)
		if err != nil {
			log.Fatalln("error in initializing firebase app: ", err)
		}
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase db client: ", err)
	}

	r.Use(CORSMiddleware())
	routes.SetupRoutes(r, client)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
}
	r.Run(":" + port)
}