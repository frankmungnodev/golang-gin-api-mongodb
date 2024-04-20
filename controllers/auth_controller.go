package controllers

import (
	"context"
	"franky/go-api-gin/envs"
	"franky/go-api-gin/models"
	"franky/go-api-gin/mongodb"
	"franky/go-api-gin/utils"
	"net/http"
	"os"
	"time"

	"github.com/fatih/color"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/golang-jwt/jwt/v4"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"golang.org/x/crypto/bcrypt"
)

type AuthController struct{}

type AuthRequestBody struct {
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=6"`
}

func (ac *AuthController) Register(c *gin.Context) {

	// Extract request body
	var register AuthRequestBody
	if err := c.ShouldBindJSON(&register); err != nil {
		color.New(color.FgRed).Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: "Invalid registration data"})
		return
	}

	// Validate request body
	if err := validator.New().Struct(register); err != nil {
		errors := err.(validator.ValidationErrors)
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: utils.GetValidationMessage(errors)})
		return
	}

	userCollection := mongodb.DB.Collection("users")

	// Check existing user with email
	var existsUser models.User
	err := userCollection.FindOne(context.Background(), bson.M{"email": register.Email}).Decode(&existsUser)
	if err == nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: "User with this email already exists"})
		return
	}

	// If user not registered yet! hash password
	hashedPassword, error := hashPassword(register.Password)
	if error != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.ErrorInfo{Message: "Cannot encrypt password"})
		return
	}

	// Create new user
	now := time.Now()
	result, err := userCollection.InsertOne(
		context.Background(),
		&models.User{Email: register.Email, Password: hashedPassword, CreatedAt: now, UpdatedAt: now},
	)

	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: "User with this email already exists"})
		} else {
			c.AbortWithStatusJSON(http.StatusInternalServerError, &models.ErrorInfo{Message: "Cannot create user"})
		}
		return
	}

	var user models.User
	err = userCollection.FindOne(context.Background(), bson.M{"_id": result.InsertedID}).Decode(&user)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.ErrorInfo{Message: "Failed to retrieve created user"})
		return
	}

	sendTokenResponse(c, &user)
}

func (ac *AuthController) Login(c *gin.Context) {
	// Extract request body
	var login AuthRequestBody
	if err := c.ShouldBindJSON(&login); err != nil {
		color.New(color.FgRed).Println(err.Error())
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: "Invalid login data"})
		return
	}

	// Validate request body
	if err := validator.New().Struct(login); err != nil {
		errors := err.(validator.ValidationErrors)
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: utils.GetValidationMessage(errors)})
		return
	}

	// Check existing user with email
	var existsUser models.User
	err := mongodb.UserCollection.FindOne(context.Background(), bson.M{"email": login.Email}).Decode(&existsUser)
	if err != nil {
		color.New(color.FgRed).Println("login error when find user with email: ", err)
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: "Invalid credentials"})
		return
	}

	isPasswordMatch := checkPasswordHash(login.Password, existsUser.Password)
	if !isPasswordMatch {
		color.New(color.FgRed).Println("login error password not match")
		c.AbortWithStatusJSON(http.StatusBadRequest, &models.ErrorInfo{Message: "Invalid credentials"})
		return
	}

	sendTokenResponse(c, &existsUser)
}

func (ac *AuthController) GetMe(c *gin.Context) {
	user := c.MustGet("user").(models.User)
	c.JSON(http.StatusOK, gin.H{"data": user})
}

func (ac *AuthController) Logout(c *gin.Context) {
	cookie := http.Cookie{
		Name:     "token",
		Value:    "",
		Expires:  time.Now(),
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(
		http.StatusOK,
		gin.H{
			"message": "Logout success",
		},
	)
}

// Private functions
func sendTokenResponse(c *gin.Context, user *models.User) {
	expiration := time.Now().Add(time.Duration(envs.JWT_EXPIRE_IN_DAYS) * 24 * time.Hour)
	tokenString, err := generateJWT(user.Email, expiration)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, &models.ErrorInfo{Message: "Failed to generate JWT token"})
		return
	}
	cookie := http.Cookie{
		Name:     "token",
		Value:    tokenString,
		Expires:  expiration,
		HttpOnly: true,
		Secure:   false,
		SameSite: http.SameSiteStrictMode,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)
	c.JSON(
		http.StatusCreated,
		gin.H{
			"token": tokenString,
			"data":  user,
		},
	)
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

// GenerateJWT generates a JWT token with the user's email as a claim
func generateJWT(email string, expiration time.Time) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["email"] = email
	claims["exp"] = expiration.Unix() // Token expires in 24 hours

	secretKey := []byte(os.Getenv("JWT_SECRET"))
	return token.SignedString(secretKey)
}
