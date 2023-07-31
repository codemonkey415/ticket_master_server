package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"

	"github.com/octocat0415/database"

	"github.com/octocat0415/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var seatCollection *mongo.Collection = database.OpenCollection(database.Client, "seats")

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/x-www-form-urlencoded")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		email := c.GetString("email")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var foundUser models.User

		err := userCollection.FindOne(ctx, bson.M{"email": email}).Decode(&foundUser)
		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}
		reservations, err := GetSeatDataForUser(foundUser)
		if err != nil {
			// handle the error if necessary
		}

		data := map[string]interface{}{
			"user":         foundUser,
			"reservations": reservations,
		}
		c.JSON(http.StatusOK, data)
	}
}

func GetSeatDataForUser(user models.User) ([]models.Ticket, error) {
	fmt.Println(user.Reservations)
	filter := bson.M{"_id": bson.M{"$in": user.Reservations}}

	var tickets []models.Ticket

	cur, err := seatCollection.Find(context.TODO(), filter)
	if err != nil {
		return nil, err
	}
	defer cur.Close(context.TODO())

	for cur.Next(context.TODO()) {
		var ticket models.Ticket
		err := cur.Decode(&ticket)
		if err != nil {
			return nil, err
		}
		tickets = append(tickets, ticket)
	}

	if err := cur.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

func GetUsers() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/x-www-form-urlencoded")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		pageSizeStr := c.Query("page_size")
		pageStr := c.Query("page")

		pageSize := 10
		page := 1

		if pageSizeStr != "" {
			pSize, err := strconv.Atoi(pageSizeStr)
			if err == nil && pSize > 0 {
				pageSize = pSize
			}
		}

		if pageStr != "" {
			p, err := strconv.Atoi(pageStr)
			if err == nil && p > 0 {
				page = p
			}
		}

		fmt.Println(pageSize)
		fmt.Println(page)

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var users []bson.M
		defer cancel()

		findOptions := options.Find()
		findOptions.SetLimit(int64(pageSize))
		findOptions.SetSkip(int64((page - 1) * pageSize))
		sort := bson.D{{Key: "role", Value: 1}}
		findOptions.SetSort(sort)

		filter := bson.M{"role": bson.M{"$ne": "owner"}}
		count, err := userCollection.CountDocuments(context.TODO(), filter)
		results, err := userCollection.Find(ctx, filter, findOptions)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}

		defer cancel()
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		defer results.Close(ctx)
		for results.Next(ctx) {
			var singleUser bson.M
			if err = results.Decode(&singleUser); err != nil {
				log.Fatal(err)
			}
			users = append(users, singleUser)
		}
		data := map[string]interface{}{
			"users": users,
			"count": count,
		}
		c.JSON(http.StatusOK, data)
	}
}

func ChangeStatus() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/x-www-form-urlencoded")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		user_id := c.Param("userid")
		parsedID, err := primitive.ObjectIDFromHex(user_id)

		fmt.Println(user_id)

		filter := bson.M{"_id": parsedID}

		var requestBody struct {
			IsApproved bool `json:"is_approved"`
		}
		err = c.ShouldBindJSON(&requestBody)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		update := bson.M{"$set": bson.M{"is_approved": requestBody.IsApproved}}
		result := userCollection.FindOneAndUpdate(ctx, filter, update)

		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}

		updatedUser := models.User{}
		err = result.Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedUser)

	}
}

func UpdateDueDate() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/x-www-form-urlencoded")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		month := c.Param("month")
		intMonth, err := strconv.Atoi(month)
		if err != nil {
			log.Fatal("Failed to convert month to int:", err)
		}

		userid := c.Param("userid")
		parsedID, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userid"})
			return
		}

		filter := bson.M{"_id": parsedID}

		var user models.User
		err = userCollection.FindOne(ctx, filter).Decode(&user)

		originalDueDate := user.Due_Date

		if err != nil {
			log.Fatal("Failed to decode user:", err)
		}

		parsedDueDate, err := time.Parse(time.RFC3339, originalDueDate.Format(time.RFC3339))
		if err != nil {
			log.Fatal("Failed to parse due date:", err)
		}

		updatedDueDate := parsedDueDate.AddDate(0, intMonth, 0)

		update := bson.D{{"$set", bson.D{{"due_date", updatedDueDate.Format(time.RFC3339)}}}}

		result := userCollection.FindOneAndUpdate(ctx, filter, update)
		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}
		updatedUser := models.User{}
		err = result.Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, updatedUser)
	}
}

func ChangeRole() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/x-www-form-urlencoded")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userid := c.Param("userid")
		parsedID, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userid"})
			return
		}

		filter := bson.M{"_id": parsedID}

		var requestBody struct {
			Role string `json:"role"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		update := bson.D{{"$set", bson.D{{"role", requestBody.Role}}}}

		result := userCollection.FindOneAndUpdate(ctx, filter, update)
		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}

		updatedUser := models.User{}
		err = result.Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedUser)
	}
}

func SaveReservations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/x-www-form-urlencoded")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userid := c.Param("userid")
		parsedID, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userid"})
			return
		}

		filter := bson.M{"_id": parsedID}

		var requestBody struct {
			Reservations []primitive.ObjectID `json:"reservations"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		update := bson.D{{"$addToSet", bson.D{{"reservations", bson.D{{"$each", requestBody.Reservations}}}}}}
		result := userCollection.FindOneAndUpdate(ctx, filter, update)

		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}
		updatedUser := models.User{}
		err = result.Decode(&updatedUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, updatedUser)
	}
}

func RemoveReservations() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Header("Content-Type", "application/x-www-form-urlencoded")
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "POST")
		c.Header("Access-Control-Allow-Headers", "Content-Type")

		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		userid := c.Param("userid")
		parsedID, err := primitive.ObjectIDFromHex(userid)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid userid"})
			return
		}

		filter := bson.M{"_id": parsedID}

		var requestBody struct {
			Reservations []primitive.ObjectID `json:"reservations"`
		}
		if err := c.ShouldBindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
			return
		}

		result := userCollection.FindOne(ctx, filter)
		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}

		user := models.User{}
		err = result.Decode(&user)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, reservationID := range requestBody.Reservations {
			for i, userReservationID := range user.Reservations {
				if userReservationID == reservationID {
					user.Reservations = append(user.Reservations[:i], user.Reservations[i+1:]...)
					break
				}
			}
		}

		update := bson.D{{"$set", bson.D{{"reservations", user.Reservations}}}}
		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update user"})
			return
		}

		c.JSON(http.StatusOK, user)
	}
}
