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
	Rank string `json:"rank"`
}

func SetupProgressionRoutes(r *gin.Engine, client *db.Client) {
	updateExp(r, client)
	updateRank(r, client)
	getProgression(r, client)
}

func updateExp(r *gin.Engine, client *db.Client) {
	r.POST("prog/exp", func(c *gin.Context) {
		expTotalStr := c.PostForm("expTotal")
		expTotalVal, err := strconv.ParseInt(expTotalStr, 10, 32); if err != nil {
			log.Fatal(err)
		}
		ref := client.NewRef("prog/")
		expPayload := map[string]interface{}{"expTotal":expTotalVal}
		if updateErr := ref.Update(context.TODO(), expPayload); updateErr != nil {
			log.Fatal(updateErr)
		}
		c.JSON(http.StatusCreated, gin.H{
			"newExpValue": expTotalVal,
		})
	})
}

func updateRank(r *gin.Engine, client *db.Client) {
	r.POST("prog/rank", func(c *gin.Context) {
		newRank := c.PostForm("newRank")
		ref := client.NewRef("prog/")
		rankPayload := map[string]interface{}{"rank":newRank}
		if updateErr := ref.Update(context.TODO(), rankPayload); updateErr != nil {
			log.Fatal(updateErr)
		}
		c.JSON(http.StatusCreated, gin.H{
			"newRank": newRank,
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
		if (prog.Rank == "") {
			prog.Rank = "Bronze I"
		}
		c.JSON(http.StatusOK, prog)
	})
}