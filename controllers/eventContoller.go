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

var eventCollection *mongo.Collection = database.OpenCollection(database.Client, "events")

func GetAllEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var events []models.Event

		pipeline := bson.A{
			bson.M{
				"$match": bson.M{
					"is_active": 1,
				},
			},
			bson.M{
				"$addFields": bson.M{
					"label": "$name",
				},
			},
			// bson.M{"$limit": 100},
		}

		results, err := eventCollection.Aggregate(ctx, pipeline)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleEvent models.Event
			if err = results.Decode(&singleEvent); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			events = append(events, singleEvent)
		}

		c.JSON(http.StatusOK, events)
	}
}
