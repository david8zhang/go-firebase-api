package routes

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

func CreateTodo(r *gin.Engine, client *db.Client) {
	r.POST("todo", func (c *gin.Context) {
		id := c.PostForm("id")
		title := c.PostForm("title")
		difficulty := c.PostForm("difficulty")

		ref := client.NewRef("todos/" + id)
		newTodo := map[string]interface{}{"id":id, "title":title, "difficulty":difficulty}
		if err := ref.Set(context.TODO(), newTodo); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"created": id,
		})
	})
}