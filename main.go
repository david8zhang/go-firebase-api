package main

import (
	"context"
	"log"

	firebase "firebase.google.com/go/v4"
	"github.com/david8zhang/go-firebase/routes"
	"github.com/gin-gonic/gin"
	"google.golang.org/api/option"
)

func main() {
	r := gin.Default()
	ctx := context.Background()
	conf := &firebase.Config{
		DatabaseURL: "https://habit-rank-default-rtdb.firebaseio.com",
	}

	opt := option.WithCredentialsFile("env/habit-rank-firebase-adminsdk-8xo36-576cec4f1c.json")
	app, err := firebase.NewApp(ctx, conf, opt)
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

	r.Run()
}