package routes

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

type Todo struct {
	Id string
	Title string
	Difficulty string
	LastUpdated int64
}

func GetTodos(r *gin.Engine, client *db.Client) {
	r.GET("todo", func (c *gin.Context) {

		ref := client.NewRef("todos")
		var todos map[string]Todo
		if err := ref.Get(context.TODO(), &todos); err != nil {
			log.Fatalln("error in reading from firebase DB: ", err)
		}

		c.JSON(http.StatusOK, gin.H{
			"todos": maps.Values(todos),
		})
	})
}