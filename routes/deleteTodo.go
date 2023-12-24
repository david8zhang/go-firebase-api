package routes

import (
	"context"
	"log"
	"net/http"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

func DeleteTodo(r *gin.Engine, client *db.Client) {
	r.DELETE("todo", func (c *gin.Context) {
		todoId := c.Query("id")
		ref := client.NewRef("todos/" + todoId)
		if err := ref.Delete(context.TODO()); err != nil {
			log.Fatalln("error in deleting ref: ", err)
		}
		c.JSON(http.StatusOK, gin.H{
			"id": todoId,
		})
	})
}