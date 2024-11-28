package handlers

import (
	"fmt"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

func (h *Handler) generateJWT(userID int32, username string) (string, error) {
	expirationTime := time.Now().Add(24 * time.Hour)
	claims := jwt.MapClaims{
		"user_id":  userID,
		"username": username,
		"exp":      expirationTime.Unix(),
		"iat":      time.Now().Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(h.JWTSecretKey))
}

func verifyToken(c *gin.Context, tokenString string) error {
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		return []byte("secret"), nil
	})

	if err != nil {
		return err
	}

	if !token.Valid {
		return fmt.Errorf("invalid token")
	}

	c.Set("user", token)
	return nil
}

func (h *Handler) GetUserId(c *gin.Context) (string, error) {
	u, exists := c.Get("user")
	if !exists || u == nil {
		return "", fmt.Errorf("user not found in context")
	}
	token, ok := u.(*jwt.Token)
	if !ok {
		return "", fmt.Errorf("invalid token")
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", fmt.Errorf("invalid token claims")
	}

	fmt.Printf("Token Claims: %+v\n", claims)

	userIDInterface, exists := claims["user_id"]
	if !exists {
		return "", fmt.Errorf("user_id not found in token")
	}

	switch v := userIDInterface.(type) {
	case float64:
		return strconv.Itoa(int(v)), nil
	case string:
		return v, nil
	default:
		return "", fmt.Errorf("user_id has unexpected type")
	}
}

func getTokenFromHeader(c *gin.Context) (string, error) {
	tokenString := c.GetHeader("Authorization")
	if tokenString == "" || len(tokenString) <= len("Bearer ") {
		return "", fmt.Errorf("Missing or invalid authorization header")
	}
	return tokenString[len("Bearer "):], nil
}

func validateToken(c *gin.Context, tokenString string) error {
	if err := verifyToken(c, tokenString); err != nil {
		return fmt.Errorf("Invalid token: %w", err)
	}
	return nil
}

func getUserIDFromContext(h *Handler, c *gin.Context) (int32, error) {
	userID, err := h.GetUserId(c)
	if err != nil {
		return 0, fmt.Errorf("Failed to retrieve user_id: %w", err)
	}

	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		return 0, fmt.Errorf("user_id is not a valid integer: %w", err)
	}

	return int32(userIDInt), nil
}

func getPetitionIDFromParams(c *gin.Context) (int32, error) {
	petitionIDParam := c.Param("id")
	if petitionIDParam == "" {
		return 0, fmt.Errorf("petition ID is missing")
	}

	petitionID, err := strconv.Atoi(petitionIDParam)
	if err != nil {
		return 0, fmt.Errorf("Invalid petition ID: %w", err)
	}

	return int32(petitionID), nil
}
