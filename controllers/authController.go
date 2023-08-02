package controllers

import (
	"context"
	"crypto/rand"
	"fmt"
	"log"
	"math/big"

	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"

	"github.com/octocat0415/database"

	helper "github.com/octocat0415/helpers"
	"github.com/octocat0415/models"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

var userCollection *mongo.Collection = database.OpenCollection(database.Client, "user")
var validate = validator.New()

func generateRandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

	var result string
	for i := 0; i < length; i++ {
		randomIndex, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		result += string(charset[randomIndex.Int64()])
	}

	return result, nil
}

// HashPassword is used to encrypt the password before it is stored in the DB
func HashPassword(password string) string {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	if err != nil {
		log.Panic(err)
	}

	return string(bytes)
}

// VerifyPassword checks the input password while verifying it with the passward in the DB.
func VerifyPassword(userPassword string, providedPassword string) (bool, string) {
	err := bcrypt.CompareHashAndPassword([]byte(providedPassword), []byte(userPassword))
	check := true
	msg := ""

	if err != nil {
		msg = fmt.Sprintf("Incorrect password")
		check = false
	}

	return check, msg
}

// CreateUser is the api used to tget a single user
func SignUp() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User

		defer cancel()
		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"bind error": err.Error()})
			return
		}

		validationErr := validate.Struct(user)
		if validationErr != nil {
			c.JSON(http.StatusBadRequest, gin.H{"validate error": validationErr.Error()})
			return
		}
		fmt.Println("user: ", user.Email)

		count, err := userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		password := HashPassword(*user.Password)
		user.Password = &password

		count, err = userCollection.CountDocuments(ctx, bson.M{"email": user.Email})
		defer cancel()
		if err != nil {
			log.Panic(err)
			c.JSON(http.StatusInternalServerError, gin.H{"error": "error occured while checking for the email"})
			return
		}

		if count > 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "this email already exists"})
			return
		}

		user.Created_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Updated_at, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Due_Date, _ = time.Parse(time.RFC3339, time.Now().Format(time.RFC3339))
		user.Is_Approved = false
		user.ID = primitive.NewObjectID()
		user.User_id = user.ID.Hex()
		user.Role = "user"
		token, refreshToken, _ := helper.GenerateAllTokens(*user.Email, *user.First_name, *user.Last_name, user.User_id)
		user.Token = &token
		user.Refresh_token = &refreshToken

		resultInsertionNumber, insertErr := userCollection.InsertOne(ctx, user)
		if insertErr != nil {
			msg := fmt.Sprintf("User item was not created")
			c.JSON(http.StatusInternalServerError, gin.H{"error": msg})
			return
		}
		defer cancel()

		fmt.Println(resultInsertionNumber)

		c.JSON(http.StatusOK, gin.H{"message": "Successfully registered!"})

	}
}

// Login is the api used to get a single user
func Login() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var user models.User
		var foundUser models.User

		if err := c.BindJSON(&user); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"email": user.Email}).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}

		passwordIsValid, msg := VerifyPassword(*user.Password, *foundUser.Password)
		defer cancel()
		if !passwordIsValid {
			c.JSON(http.StatusBadRequest, gin.H{"error": msg})
			return
		}

		if foundUser.Role != "owner" {
			if !foundUser.Is_Approved {
				c.JSON(http.StatusBadRequest, gin.H{"error": "User is not approved"})
				return
			}

			currentTime := time.Now()

			if foundUser.Due_Date.Before(currentTime) {
				parsedDueDateStr := foundUser.Due_Date.Format(time.RFC3339)
				parsedDueDate, err := time.Parse(time.RFC3339, parsedDueDateStr)
				formattedDate := parsedDueDate.Format("January 2, 2006")
				if err != nil {
					c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to parse due date"})
					return
				}

				c.JSON(http.StatusBadRequest, gin.H{"error": "This account has expired at " + formattedDate})
				return
			}
		}

		token, refreshToken, _ := helper.GenerateAllTokens(*foundUser.Email, *foundUser.First_name, *foundUser.Last_name, foundUser.User_id)

		updatedUser := helper.UpdateAllTokens(token, refreshToken, foundUser.User_id)

		reservations, err := GetSeatDataForUser(foundUser)
		if err != nil {
			// handle the error if necessary
		}

		data := map[string]interface{}{
			"user":         updatedUser,
			"reservations": reservations,
		}

		c.JSON(http.StatusOK, data)
	}
}

func ForgotPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		var requestBody struct {
			Email string `json:"email"`
		}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		randomstr, err := generateRandomString(10)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		hashedToken := HashPassword(randomstr)

		filter := bson.M{"email": requestBody.Email}
		update := bson.M{"$set": bson.M{"reset_password_token": hashedToken}}

		result, err := userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "User doesn't exist"})
			return
		}

		if result.ModifiedCount == 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		} else {

			_, err := helper.SendResetPasswordLink(requestBody.Email, hashedToken)

			if err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"msg": "Email was sent successfully"})
		}

	}
}

func ResetPassword() gin.HandlerFunc {
	return func(c *gin.Context) {
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()

		token := c.Param("token")

		var foundUser models.User

		var requestBody struct {
			Password string `jsong:"password"`
		}

		if err := c.BindJSON(&requestBody); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		filter := bson.M{"reset_password_token": token}
		// update := bson.M{"$set": bson.M{"reset_password_token": hashedToken}}

		err := userCollection.FindOne(ctx, filter).Decode(&foundUser)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid Token"})
			return
		}

		newPassword := HashPassword(requestBody.Password)

		update := bson.M{"$set": bson.M{"password": newPassword}}

		_, err = userCollection.UpdateOne(ctx, filter, update)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"message": "Password changed successfully"})
	}
}
