package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"
	"time"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
	"golang.org/x/exp/maps"
)

type Todo struct {
	Id string `json:"id"`
	Title string `json:"title"`
	Difficulty string `json:"difficulty"`
	LastUpdated int64 `json:"lastUpdated"`
	Status int32 `json:"status"`
}

func SetupTodoRoutes(r *gin.Engine, client *db.Client) {
	completeTodos(r, client)
	updateTodo(r, client)
	createTodo(r, client)
	deleteTodo(r, client)
	processExpiredTodos(r, client)
}

func completeTodos(r *gin.Engine, client *db.Client) {
	r.POST("todo/complete", func (c* gin.Context) {
		id := c.PostForm("id")
		status, err := strconv.ParseInt(c.PostForm("status"), 10, 32); if err != nil {
			log.Fatal(err)
		}
		lastUpdated := time.Now().UnixMilli()

		ref := client.NewRef("todos/" + id)
		newTodo := map[string]interface{}{"id":id,"lastUpdated":lastUpdated,"status":status}
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
		status, parseErr := strconv.ParseInt(c.PostForm("status"), 10, 32); if parseErr != nil {
			log.Fatal(parseErr)
		}

		ref := client.NewRef("todos/" + id)
		newTodo := map[string]interface{}{
			"id":id,
			"title":title,
			"difficulty":difficulty,
			"lastUpdated": -1,
			"status":status,
		}
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

func processExpiredTodos(r *gin.Engine, client *db.Client) {
	/**
		For every todo:
			if its status is COMPLETE:
				- if the lastUpdated timestamp is outside of the boundaries of the current day (12AM PST), mark the todo as INCOMPLETE and update its timestamp
				- otherwise, do nothing since the todo was completed on the current day

			if its status is INCOMPLETE:
				- if the lastUpdated timestamp is outside of the boundaries of the current day, then we need to deduct rankExp for it (add its penalty to a total and do a batch update on the current rank and exp total)
				- otherwise, do nothing since the user still has the current day to complete the todo


		For rank penalty processing
			- Grab the current rank and expTotal
			- Subtract the penalties from the expTotal. If we go below 0, then we demote the rank
				- Example: Curr rank: Bronze II, Rank Penalty is -20, curr exp is 10. Final result: Bronze I, curr exp: 90
				- Example: Curr rank: Bronze III, Rank Penalty is -120, curr exp is 10, final result: Bronze I, curr exp: 90
	*/
	r.POST("todo/expire", func (c *gin.Context) {
		loc, err := time.LoadLocation("America/Los_Angeles")
		if err != nil {
			log.Fatal(err)
		}
		now := time.Now().In(loc)
		bod := beginningOfDay(now)

		todoRef := client.NewRef("todos")
		var todos map[string]Todo
		if getTodosErr := todoRef.Get(context.TODO(), &todos); getTodosErr != nil {
			log.Fatalln("error fetching TODOs from firebase DB: ", getTodosErr)
		}
		todoArr := maps.Values(todos)
		todosToUpdate := make(map[string]interface{})
		var penalties []int
		var currTimestamp = time.Now().UnixMilli()

		for _, todo := range todoArr {
			if (todo.LastUpdated < bod.UnixMilli()) {
				// If todo is incomplete, add a penalty
				if (todo.Status == 0) {
					penalties = append(penalties, -10)
				} else {
					// Otherwise, change its status to incomplete
					todo.Status = 0
				}
				// Change its last updated so that we don't update this todo again
				todo.LastUpdated = currTimestamp
				todosToUpdate[todo.Id] = todo
			}
		}

		var prog Progression
		progRef := client.NewRef("prog")
		if fetchProgErr := progRef.Get(context.TODO(), &prog); fetchProgErr != nil {
			log.Fatal(fetchProgErr)
		}

		currExpTotal := int(prog.ExpTotal)
		currRank := int(prog.Rank)
		for _, penalty := range penalties {
			currExpTotal += penalty
			if (currExpTotal < 0) {
				if (currRank == 0) {
					currExpTotal = 0
				} else {
					currRank--
					currExpTotal += 100
				}
			}
		}
		progPayload := map[string]interface{}{"rank":currRank,"expTotal":currExpTotal}
		if updateErr := progRef.Update(context.TODO(), progPayload); updateErr != nil {
			log.Fatal(updateErr)
		}

		if len(todosToUpdate) != 0 {
			if updateTodoErr := todoRef.Update(context.TODO(), todosToUpdate); updateTodoErr != nil {
				log.Fatal(updateTodoErr)
			}
		}

		c.JSON(http.StatusOK, gin.H{
			"progression": map[string]interface{}{
				"old": prog,
				"new": progPayload,
			},
			"todos": todoArr,
		})
	})
}

func beginningOfDay(t time.Time) time.Time {
	year, month, day := t.Date()
	return time.Date(year, month, day, 0, 0, 0, 0, t.Location())
}