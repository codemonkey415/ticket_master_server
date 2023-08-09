package controllers

import (
	"context"
	"fmt"
	"log"
	"strconv"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	jsoniter "github.com/json-iterator/go"

	"github.com/octocat0415/database"
	helper "github.com/octocat0415/helpers"

	"github.com/octocat0415/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var seatCollection *mongo.Collection = database.OpenCollection(database.Client, "seats")

func GetUser() gin.HandlerFunc {
	return func(c *gin.Context) {
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
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)

		defer cancel()

		user_id := c.Param("userid")
		parsedID, err := primitive.ObjectIDFromHex(user_id)

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

		update := bson.D{{Key: "$set", Value: bson.D{{Key: "role", Value: requestBody.Role}}}}

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

		update := bson.D{
			{Key: "$addToSet", Value: bson.D{
				{Key: "reservations", Value: bson.D{
					{Key: "$each", Value: requestBody.Reservations},
				}},
			}},
		}
		result := userCollection.FindOneAndUpdate(ctx, filter, update)

		if result.Err() != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
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

func RunNotifySeats(ch chan<- bool, email string) {
	go func() {
		for {
			/////////////////////////////////////////////////////
			ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) // Reduced timeout duration for better performance
			defer cancel()

			matchStage := bson.D{
				bson.E{
					Key: "$match",
					Value: bson.D{
						{Key: "email", Value: email},
					},
				},
			}

			lookupStage := bson.D{
				bson.E{
					Key: "$lookup",
					Value: bson.M{
						"from":         "seats",
						"localField":   "reservations",
						"foreignField": "_id",
						"as":           "reservations_detail",
					},
				},
			}

			pipeline := mongo.Pipeline{matchStage, lookupStage}
			results, err := userCollection.Aggregate(ctx, pipeline)
			if err != nil {
				log.Fatal(err)
			}
			defer results.Close(ctx)

			var result bson.M
			if results.Next(ctx) {
				if err = results.Decode(&result); err != nil {
					log.Fatal(err)
				}
			}

			reservationsJSON, err := jsoniter.Marshal(result["reservations_detail"])
			var availableReservations []models.Ticket

			if err != nil {
				log.Fatal(err)
			}

			var reservations []models.Ticket
			err = jsoniter.Unmarshal(reservationsJSON, &reservations)
			if err != nil {
				log.Fatal(err)
			}

			for _, reservation := range reservations {
				IsAvailable := reservation.IsAvailable
				if IsAvailable == 1 {
					availableReservations = append(availableReservations, reservation)
				}
			}
			var readableReservations string

			if len(availableReservations) > 0 {
				for _, ticket := range availableReservations {
					readableReservations += fmt.Sprintf("Event ID: %s\nRow Name: %s\nSeat Group: %s\nSeat Name: %s\nPrice: %s\n\n", ticket.EventId, ticket.RowName, ticket.SeatGroup, ticket.SeatName, ticket.Price)
				}
				helper.SendMail(email, readableReservations)
			}
			////////////////////////////////////////////////////
			time.Sleep(1 * time.Minute)
		}
	}()
}

func StartNotifySeats() gin.HandlerFunc {
	return func(c *gin.Context) {

		email := c.GetString("email")

		// Create a channel for communication
		ch := make(chan bool)

		// Start the goroutine function
		RunNotifySeats(ch, email)

		// Return a response to the front end
		c.JSON(http.StatusOK, gin.H{"message": "NotifySeats goroutine started"})
	}
}
