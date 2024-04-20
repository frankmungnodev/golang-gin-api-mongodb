package middlewares

import (
	"context"
	"fmt"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"

	"github.com/fatih/color"

	"franky/go-api-gin/models"
	"franky/go-api-gin/mongodb"
)

type ProtectRoute struct {
}

func (protectRoute *ProtectRoute) Protect(c *gin.Context) {
	tokenString := extractTokenFromHeader(c)
	if tokenString == "" {
		color.New(color.BgRed).Println("middlewares protect > token not found")
		respondUnauthorized(c)
		return
	}

	claims, err := validateToken(tokenString, "asdfasdf")
	if err != nil {
		color.New(color.BgRed).Println(err.Error())
		respondUnauthorized(c)
		return
	}

	userEmail, ok := claims["email"].(string)
	if !ok {
		color.New(color.BgRed).Println("middlewares protect > email not found from claimed token")
		respondUnauthorized(c)
		return
	}

	user, err := findUserByEmail(userEmail)
	if err != nil {
		color.New(color.BgRed).Println("middlewares protect > ", err.Error())
		respondUnauthorized(c)
		return
	}

	c.Set("user", user)
	c.Next()
}

func extractTokenFromHeader(c *gin.Context) string {
	parts := strings.Split(c.GetHeader("Authorization"), " ")
	if len(parts) != 2 || parts[0] != "Bearer" {
		return ""
	}
	return parts[1]
}

func validateToken(tokenString string, secretKey string) (jwt.MapClaims, error) {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte(secretKey), nil
	})
	if err != nil {
		return nil, err
	}

	if !token.Valid {
		return nil, fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, fmt.Errorf("invalid token format claim error")
	}

	return claims, nil
}

func findUserByEmail(userEmail string) (models.User, error) {
	var user models.User
	err := mongodb.UserCollection.FindOne(context.Background(), bson.M{"email": userEmail}).Decode(&user)
	return user, err
}

func respondUnauthorized(c *gin.Context) {
	errorResponse := &models.ErrorInfo{Message: "Please login to continue"}
	c.AbortWithStatusJSON(http.StatusUnauthorized, errorResponse)
}
