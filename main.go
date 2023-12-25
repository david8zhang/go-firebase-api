package main

import (
	"context"
	"log"
	"os"

	firebase "firebase.google.com/go/v4"
	"github.com/david8zhang/go-firebase/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://habit-rank-default-rtdb.firebaseio.com",
	}

	app, err := firebase.NewApp(ctx, conf)
	if err != nil {
		log.Fatalln("error in initializing firebase app: ", err)
	}

	client, err := app.Database(ctx)
	if err != nil {
		log.Fatalln("error in creating firebase db client: ", err)
	}

	routes.CreateTodo(r, client)
	routes.GetTodos(r, client)
	routes.DeleteTodo(r, client)

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
		log.Printf("defaulting to port %s", port)
}
	r.Run(":" + port)
}