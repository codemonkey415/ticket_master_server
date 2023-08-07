package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octocat0415/database"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

var eventCollection *mongo.Collection = database.OpenCollection(database.Client, "events")

func GetAllValidEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var result []bson.M

		var requestBody struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		// matchStage := bson.D{
		// 	bson.E{
		// 		Key: "$match",
		// 		Value: bson.M{
		// 			// "seats": bson.M{"$gt": []interface{}{}},
		// 			"is_active": 1,
		// 			"event_date": bson.M{
		// 				"$gte": requestBody.StartDate,
		// 				"$lte": requestBody.EndDate,
		// 			},
		// 		},
		// 	},
		// }

		// matchStage := bson.D{
		// 	bson.E{
		// 		Key: "$match",
		// 		Value: bson.M{
		// 			"$and": []bson.M{
		// 				{"event_date": bson.M{"$gte": requestBody.StartDate}},
		// 				{"event_date": bson.M{"$lte": requestBody.EndDate}},
		// 			},
		// 		},
		// 	},
		// }

		// lookupStage := bson.D{
		// 	bson.E{
		// 		Key: "$lookup",
		// 		Value: bson.M{
		// 			"from":         "seats",
		// 			"localField":   "event_id",
		// 			"foreignField": "event_id",
		// 			"as":           "seats",
		// 		},
		// 	},
		// }

		// projectStage := bson.D{
		// 	bson.E{
		// 		Key: "$project",
		// 		Value: bson.M{
		// 			"label": "$event.venue",
		// 		},
		// 	},
		// }
		groupStage := bson.D{
			bson.E{
				Key:   "$group",
				Value: bson.D{{Key: "_id", Value: "$event_id"}},
			},
		}

		// pipeline := mongo.Pipeline{matchStage /*lookupStage,*/ /*, projectStage*/}
		// pipeline := mongo.Pipeline{matchStage, lookupStage}
		pipeline := mongo.Pipeline{groupStage}
		results, err := eventCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
			return
		}
		defer results.Close(ctx)
		for results.Next(ctx) {
			var single bson.M
			if err = results.Decode(&single); err != nil {
				log.Fatal(err)
			}
			result = append(result, single)
		}

		c.JSON(http.StatusOK, gin.H{"result": result})
	}
}

func GetVenueByDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var result []bson.M

		var requestBody struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		matchStage := bson.D{
			bson.E{
				Key: "$match",
				Value: bson.M{
					"event_date": bson.M{
						"$gte": requestBody.StartDate,
						"$lte": requestBody.EndDate,
					},
				},
			},
		}

		projectStage := bson.D{
			bson.E{
				Key: "$project",
				Value: bson.M{
					"venue": "$venue",
				},
			},
		}

		groupStage := bson.D{
			bson.E{
				Key:   "$group",
				Value: bson.D{{Key: "_id", Value: "$venue"}},
			},
		}

		pipeline := mongo.Pipeline{matchStage, projectStage, groupStage}
		results, err := eventCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
			return
		}
		defer results.Close(ctx)
		for results.Next(ctx) {
			var single bson.M
			if err = results.Decode(&single); err != nil {
				log.Fatal(err)
			}
			result = append(result, single)
		}

		c.JSON(http.StatusOK, result)
	}
}

func GetEventByVenue() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var result []bson.M

		var requestBody struct {
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
			Venue     string    `json:"venue"`
		}

		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error1": err.Error()})
			return
		}

		matchStage := bson.D{
			bson.E{
				Key: "$match",
				Value: bson.M{
					"event_date": bson.M{
						"$gte": requestBody.StartDate,
						"$lte": requestBody.EndDate,
					},
					"venue": requestBody.Venue,
				},
			},
		}

		projectStage := bson.D{
			bson.E{
				Key: "$project",
				Value: bson.M{
					"event_date": "$event_date",
					"event_id":   "$event_id",
					"name":       "$name",
				},
			},
		}

		// groupStage := bson.D{
		// 	bson.E{
		// 		Key:   "$group",
		// 		Value: bson.D{{Key: "_id", Value: "$event_id"}},
		// 	},
		// }

		pipeline := mongo.Pipeline{matchStage, projectStage}
		results, err := eventCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error2": err.Error()})
			return
		}
		defer results.Close(ctx)
		for results.Next(ctx) {
			var single bson.M
			if err = results.Decode(&single); err != nil {
				log.Fatal(err)
			}
			result = append(result, single)
		}

		c.JSON(http.StatusOK, result)
	}
}
