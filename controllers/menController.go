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

var menCollection *mongo.Collection = database.OpenCollection(database.Client, "men")

func getAllMen() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var men []models.Men

		defer cancel()
		if err := c.BindJSON(&men); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		results, err := menCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var man models.Men
			if err = results.Decode(&man); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			men = append(men, man)
		}

		c.JSON(http.StatusOK, men)
	}
}
