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

var allProductCollection *mongo.Collection = database.OpenCollection(database.Client, "allProduct")

func getAllProducts() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var products []models.AllProduct

		defer cancel()
		if err := c.BindJSON(&products); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		results, err := allProductCollection.Find(ctx, bson.M{})
		defer cancel()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleProduct models.AllProduct
			if err = results.Decode(&singleProduct); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}

			products = append(products, singleProduct)
		}

		c.JSON(http.StatusOK, products)
	}
}
