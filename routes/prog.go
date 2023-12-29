package routes

import (
	"context"
	"log"
	"net/http"
	"strconv"

	"firebase.google.com/go/v4/db"
	"github.com/gin-gonic/gin"
)

type Progression struct {
	ExpTotal int32 `json:"expTotal"`
	Rank int32 `json:"rank"`
}

func SetupProgressionRoutes(r *gin.Engine, client *db.Client) {
	updateProg(r, client)
	getProgression(r, client)
}

func updateProg(r *gin.Engine, client *db.Client) {
	r.POST("prog", func(c *gin.Context) {
		expTotalStr := c.PostForm("expTotal")
		expTotalVal, err := strconv.ParseInt(expTotalStr, 10, 32); if err != nil {
			log.Fatal(err)
		}

		newRankStr := c.PostForm("newRank")
		newRank, err := strconv.ParseInt(newRankStr, 10, 32); if err != nil {
			log.Fatal(err)
		}

		ref := client.NewRef("prog/")
		progPayload := map[string]interface{}{"expTotal":expTotalVal, "newRank":newRank}
		if updateErr := ref.Update(context.TODO(), progPayload); updateErr != nil {
			log.Fatal(updateErr)
		}
		c.JSON(http.StatusCreated, gin.H{
			"newExpValue": expTotalVal,
			"newRank":newRank,
		})
	})
}

func getProgression(r *gin.Engine, client *db.Client) {
	r.GET("prog", func(c *gin.Context) {
		ref := client.NewRef("prog/")
		var prog Progression
		if err := ref.Get(context.TODO(), &prog); err != nil {
			log.Fatal(err)
		}
		c.JSON(http.StatusOK, prog)
	})
}