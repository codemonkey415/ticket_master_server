package controllers

import (
	"context"
	"fmt"
	"log"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/octocat0415/database"

	"github.com/octocat0415/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var ticketCollection *mongo.Collection = database.OpenCollection(database.Client, "seats")
var zeroDecimal primitive.Decimal128

type Seat struct {
	EventId     string `json:"event_id"`
	MinPrice    int32  `json:"min_price"`
	MaxPrice    int32  `json:"max_price"`
	RowName     string `json:"row_name"`
	SectionName string `json:"section_name"`
}

func GetAllTickets() gin.HandlerFunc {
	return func(c *gin.Context) {
		var seat Seat
		if err := c.BindJSON(&seat); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		eventId := seat.EventId
		row_name := seat.RowName
		section_name := seat.SectionName
		min_price := seat.MinPrice
		max_price := seat.MaxPrice

		var ctx, cancel = context.WithTimeout(context.Background(), 200*time.Second)
		var ticket []models.Ticket

		row_name_regex := bson.M{"$regex": row_name, "$options": "i"}
		section_name_regex := bson.M{"$regex": section_name, "$options": "i"}
		// event_id_regex := bson.M{"$regex": eventId, "$options": "i"}

		row_name_match := bson.M{"row_name": row_name_regex}
		section_name_match := bson.M{"section_name": section_name_regex}
		// event_id_match := bson.M{"event_id": event_id_regex}

		pipeline := bson.A{
			bson.M{
				"$addFields": bson.M{
					"price_decimal": bson.M{
						"$cond": bson.A{
							bson.M{"$eq": bson.A{"$price", ""}},
							zeroDecimal,
							bson.M{"$toDecimal": "$price"},
						},
					},
				},
			},
			bson.M{
				"$match": bson.M{
					"$and": bson.A{
						bson.M{"price_decimal": bson.M{"$gte": min_price}},
						bson.M{"price_decimal": bson.M{"$lte": max_price}},
					},
				},
			},
			bson.M{
				"$match": bson.M{
					"event_id": eventId,
					// "section_name": section_name,
					// "row_name":     row_name,
				},
			},
			// bson.M{"$match": row_name_match},
			// bson.M{"$match": section_name_match},
			// bson.M{"$match": event_id_match},
		}

		if section_name == "" {
			fmt.Println("1")
			pipeline = append(pipeline, bson.M{
				"$match": section_name_match,
			})
		} else {
			fmt.Println("2")
			pipeline = append(pipeline, bson.M{
				"$match": bson.M{
					"section_name": section_name,
				},
			})
		}

		if row_name == "" {
			pipeline = append(pipeline, bson.M{
				"$match": row_name_match,
			})
		} else {
			pipeline = append(pipeline, bson.M{
				"$match": bson.M{
					"row_name": row_name,
				},
			})
		}

		results, err := ticketCollection.Aggregate(ctx, pipeline)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleTicket models.Ticket
			if err = results.Decode(&singleTicket); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			ticket = append(ticket, singleTicket)
		}
		c.JSON(http.StatusOK, ticket)
	}
}

func GetEvents() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second) // Set a shorter timeout duration
		defer cancel()
		var result []bson.M

		var requestBody struct {
			Venue     string    `json:"venue"`
			StartDate time.Time `json:"start_date"`
			EndDate   time.Time `json:"end_date"`
		}
		err := c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		groupStage := bson.D{
			bson.E{
				Key:   "$group",
				Value: bson.D{{Key: "_id", Value: "$event_id"}},
			},
		}

		lookupStage := bson.D{
			bson.E{
				Key: "$lookup",
				Value: bson.M{
					"from":         "events",
					"localField":   "_id",
					"foreignField": "event_id",
					"as":           "event",
				},
			},
		}

		unwindStage := bson.D{{Key: "$unwind", Value: "$event"}}

		projectStage := bson.D{
			bson.E{
				Key: "$project",
				Value: bson.M{
					"event_date": "$event.event_date",
					"label":      "$event.name",
					"venue":      "$event.venue",
				},
			},
		}

		matchStage := bson.D{
			bson.E{
				Key: "$match",
				Value: bson.M{
					"event_date": bson.M{
						"$gte": requestBody.StartDate,
						"$lte": requestBody.EndDate,
					},
					"venue": bson.M{
						"$eq": requestBody.Venue,
					},
				},
			},
		}

		pipeline := mongo.Pipeline{groupStage, lookupStage, unwindStage, projectStage, matchStage}
		results, err := ticketCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func GetVenues() gin.HandlerFunc {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		groupStage := bson.D{
			bson.E{
				Key:   "$group",
				Value: bson.D{{Key: "_id", Value: "$event_id"}},
			},
		}

		pipeline := mongo.Pipeline{groupStage /*, lookupStage, unwindStage, projectStage, matchStage*/}
		results, err := ticketCollection.Aggregate(ctx, pipeline)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func GetSectionNames() gin.HandlerFunc {
	return func(c *gin.Context) {

		event_id := c.Param("eventid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var result []bson.M

		pipeline := bson.A{
			bson.M{
				"$match": bson.M{
					"event_id": event_id,
				},
			},
			bson.M{
				"$group": bson.M{
					"_id":   "$section_name",
					"label": bson.M{"$first": "$section_name"},
				},
			},
		}

		results, err := ticketCollection.Aggregate(ctx, pipeline)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func GetRows() gin.HandlerFunc {
	return func(c *gin.Context) {

		event_id := c.Param("eventid")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var result []bson.M

		pipeline := bson.A{
			bson.M{
				"$match": bson.M{
					"event_id": event_id,
				},
			},
			bson.M{
				"$group": bson.M{
					"_id": "$row_name",
				},
			},
		}

		results, err := ticketCollection.Aggregate(ctx, pipeline)

		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func CountSeats() gin.HandlerFunc {
	return func(c *gin.Context) {

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var result []bson.M

		groupStage := bson.D{
			bson.E{
				Key: "$group",
				Value: bson.D{
					{Key: "_id", Value: "$event_id"},
				},
			},
		}

		pipeline := mongo.Pipeline{groupStage}
		results, err := ticketCollection.Aggregate(ctx, pipeline)

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
