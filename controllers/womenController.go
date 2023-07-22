package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octocat0415/database"
	"github.com/octocat0415/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var womenCollection *mongo.Collection = database.OpenCollection(database.Client, "women")

func getAllWomen() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var women []models.Women

		defer cancel()
		if err := c.BindJSON(&women); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		results, err := womenCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var woman models.Women
			if err = results.Decode(&woman); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			women = append(women, woman)
		}

		c.JSON(http.StatusOK, women)
	}
}
