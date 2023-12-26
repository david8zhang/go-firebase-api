package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

type Todo struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Difficulty string `json:"difficulty"`
	LastUpdated int64 `json:"lastUpdated"`
}

func SetupTodoRoutes(r *gin.Engine, client *db.Client) {
	getTodos(r, client)
	completeTodos(r, client)
	updateTodo(r, client)
	createTodo(r, client)
	deleteTodo(r, client)
}

func getTodos(r *gin.Engine, client *db.Client) {
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

func completeTodos(r *gin.Engine, client *db.Client) {
	r.POST("todo/complete", func (c* gin.Context) {
		id := c.PostForm("id")
		lastUpdatedStr := c.PostForm("lastUpdated")
		lastUpdatedVal, err := strconv.ParseInt(lastUpdatedStr, 10, 64)
		if err != nil {
			log.Fatal(err)
		}

		ref := client.NewRef("todos/" + id)
		newTodo := map[string]interface{}{"id":id,"lastUpdated":lastUpdatedVal}
		if updateErr := ref.Update(context.TODO(), newTodo); updateErr != nil {
			log.Fatal(updateErr)
		}
		c.JSON(http.StatusCreated, gin.H{
			"updated": id,
		})
	})
}

func updateTodo(r *gin.Engine, client *db.Client) {
	r.POST("todo/update", func (c *gin.Context) {
		id := c.PostForm("id")
		title := c.PostForm("title")
		difficulty := c.PostForm("difficulty")

		ref := client.NewRef("todos/" + id)
		newTodo := map[string]interface{}{"id":id, "title":title, "difficulty":difficulty}
		if err := ref.Update(context.TODO(), newTodo); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"updated": id,
		})
	})
}

func createTodo(r *gin.Engine, client *db.Client) {
	r.POST("todo/new", func (c *gin.Context) {
		id := c.PostForm("id")
		title := c.PostForm("title")
		difficulty := c.PostForm("difficulty")

		ref := client.NewRef("todos/" + id)
		newTodo := map[string]interface{}{"id":id, "title":title, "difficulty":difficulty,"lastUpdated": -1}
		if err := ref.Set(context.TODO(), newTodo); err != nil {
			log.Fatal(err)
		}

		c.JSON(http.StatusCreated, gin.H{
			"created": id,
		})
	})
}

func deleteTodo(r *gin.Engine, client *db.Client) {
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