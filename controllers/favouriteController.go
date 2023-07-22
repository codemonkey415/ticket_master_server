package controllers

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/octocat0415/database"
	"github.com/octocat0415/models"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var favouriteCollection *mongo.Collection = database.OpenCollection(database.Client, "favourite")

func getAllFavourite() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var favourite []models.Favourite

		defer cancel()
		if err := c.BindJSON(&favourite); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(favourite)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		results, err := favouriteCollection.Find(ctx, bson.M{})

		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleFavourite models.Favourite
			if err = results.Decode(&singleFavourite); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
				return
			}

			favourite = append(favourite, singleFavourite)
		}

		c.JSON(http.StatusOK, favourite)

	}
}

func postFavourite() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var favourite models.Favourite

		defer cancel()
		if err := c.BindJSON(&favourite); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		validationErr := validate.Struct(favourite)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}

		result, err := userCollection.InsertOne(ctx, favourite)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": validationErr.Error()})
			return
		}
		c.JSON(http.StatusOK, result)
	}
}

func deleteFavourite() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		favouriteId := c.Param("id")

		defer cancel()

		objId, _ := primitive.ObjectIDFromHex(favouriteId)

		result, err := favouriteCollection.DeleteOne(ctx, bson.M{"id": objId})
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if result.DeletedCount < 1 {
			c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, result)
	}
}
