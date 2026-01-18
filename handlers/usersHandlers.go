package handlers

import (
	"e-commerce-API/config"
	"e-commerce-API/models"
	"e-commerce-API/repositories"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type AuthHandler struct {
	authRepo *repositories.AuthRepository
}

func NewAuthHandler(authRepo *repositories.AuthRepository) *AuthHandler {
	return &AuthHandler{authRepo: authRepo}
}

func (handler *AuthHandler) RegisterUser(c *gin.Context) {
	var request models.RegisterUser

	err := c.BindJSON(&request)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not bind json object"))
		return
	}

	passwordHash, err := bcrypt.GenerateFromPassword([]byte(request.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("Could not hash password"))
		return
	}

	user := models.User{
		Name:         request.Name,
		Email:        request.Email,
		PasswordHash: string(passwordHash),
	}

	id, err := handler.authRepo.Create(c, user)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("could not create user"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"id": id})
}

func (handler *AuthHandler) SignIn(c *gin.Context) {
	var request models.SignInRequest
	err := c.BindJSON(&request)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError("Invalid parameters"))
		return
	}
	user, err := handler.authRepo.FindByEmail(c, request.Email)

	if err != nil {
		c.JSON(http.StatusBadRequest, models.NewApiError(err.Error()))
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(request.Password))

	if err != nil {
		c.JSON(http.StatusUnauthorized, models.NewApiError("Invalid credials"))
		return
	}
	claims := jwt.RegisteredClaims{
		Subject:   strconv.Itoa(user.ID),
		ExpiresAt: jwt.NewNumericDate(time.Now().Add(config.Config.JwtExpiresIn)),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString([]byte(config.Config.JwtSecretKey))
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.NewApiError("could not sign JWT"))
		return
	}
	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func (handler *AuthHandler) SignOut(c *gin.Context) {
	c.Status(http.StatusOK)
}
