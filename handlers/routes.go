package handlers

import (
	"context"
	"database/sql"
	"net/http"
	"petitionsGO/db"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// Define CustomClaims struct
type CustomClaims struct {
	UserID   int32  `json:"user_id"`
	Username string `json:"username"`
	jwt.StandardClaims
}

// router.POST("/register", h.register)
// router.POST("/login", h.login)
// router.POST("/add/petition", h.createPetition)
// router.GET("/list/petition", h.listPetition)
// router.POST("/vote/:id/petition", h.votePetition)
// router.DELETE("/delete/:id/petition", h.deletePetition)

// R E G I S T E R
func (h *Handler) register(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	_, err := h.Queries.GetUserByUsername(context.Background(), input.Username)
	if err == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "user already exists"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to hash password"})
		return
	}

	// Save to the database
	err = h.Queries.CreateUser(context.Background(), db.CreateUserParams{
		Username:     input.Username,
		PasswordHash: string(hashedPassword),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to register user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "user registered successfully"})
}

// L O G I N
func (h *Handler) login(c *gin.Context) {
	var input struct {
		Username string `json:"username" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	user, err := h.Queries.GetUserByUsername(context.Background(), input.Username)
	if err == sql.ErrNoRows {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	} else if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch user"})
		return
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(input.Password))
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "invalid username or password"})
		return
	}

	// Generate JWT
	token, err := h.generateJWT(user.ID, user.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "login successful", "token": token})
}

// L I S T   P E T I T I O N S
func (h *Handler) listPetition(c *gin.Context) {
	sortField := c.DefaultQuery("sort_field", "created_at")
	sortOrder := c.DefaultQuery("sort_order", "desc")
	page := c.DefaultQuery("page", "1")
	limit := c.DefaultQuery("limit", "10")

	pageInt, err := strconv.Atoi(page)
	if err != nil || pageInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid page parameter"})
		return
	}

	limitInt, err := strconv.Atoi(limit)
	if err != nil || limitInt < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid limit parameter"})
		return
	}

	offset := (pageInt - 1) * limitInt

	if sortField != "created_at" && sortField != "number_votes" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort_field"})
		return
	}

	if sortOrder != "asc" && sortOrder != "desc" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid sort_order"})
		return
	}

	var petitions []db.Petition
	if sortField == "created_at" && sortOrder == "asc" {
		petitions, err = h.Queries.ListPetitionsByCreatedAtAsc(context.Background())
	} else if sortField == "created_at" && sortOrder == "desc" {
		petitions, err = h.Queries.ListPetitionsByCreatedAtDesc(context.Background())
	} else if sortField == "number_votes" && sortOrder == "asc" {
		petitions, err = h.Queries.ListPetitionsByVotesAsc(context.Background())
	} else if sortField == "number_votes" && sortOrder == "desc" {
		petitions, err = h.Queries.ListPetitionsByVotesDesc(context.Background())
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to retrieve petitions"})
		return
	}

	start := offset
	end := offset + limitInt
	if start > len(petitions) {
		c.JSON(http.StatusOK, gin.H{"petitions": []db.Petition{}})
		return
	}
	if end > len(petitions) {
		end = len(petitions)
	}

	c.JSON(http.StatusOK, gin.H{"petitions": petitions[start:end]})
}

// C R E A T E   P E T I T I O N
func (h *Handler) createPetition(c *gin.Context) {
	c.Header("Content-Type", "application/json")
	tokenString, err := getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := validateToken(c, tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	var input struct {
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	userID, err := h.GetUserId(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	// Convert userID to int32
	userIDInt, err := strconv.Atoi(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "user_id is not a valid integer"})
		return
	}

	// Save the petition to the database
	petition, err := h.Queries.CreatePetition(context.Background(), db.CreatePetitionParams{
		Title:       input.Title,
		Description: input.Description,
		UserID:      int32(userIDInt),
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create petition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Petition created successfully",
		"petition": petition,
	})
}

// V O T E    F O R     S O M E
func (h *Handler) votePetition(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	tokenString, err := getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := validateToken(c, tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	petitionID, err := getPetitionIDFromParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID
	userID, err := getUserIDFromContext(h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	hasVoted, err := h.Queries.HasUserVoted(context.Background(), db.HasUserVotedParams{
		UserID:     userID,
		PetitionID: petitionID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to check user vote"})
		return
	}
	if hasVoted {
		c.JSON(http.StatusForbidden, gin.H{"error": "User has already voted for this petition"})
		return
	}

	// Record vote
	err = h.Queries.RecordVote(context.Background(), db.RecordVoteParams{
		UserID:     userID,
		PetitionID: petitionID,
	})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to record vote"})
		return
	}

	err = h.Queries.IncrementVoteCount(context.Background(), petitionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update vote count"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Vote successfully recorded"})
}

func (h *Handler) deletePetition(c *gin.Context) {
	c.Header("Content-Type", "application/json")

	tokenString, err := getTokenFromHeader(c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	if err := validateToken(c, tokenString); err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	petitionID, err := getPetitionIDFromParams(c)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Get user ID
	userID, err := getUserIDFromContext(h, c)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}
	petition, err := h.Queries.GetPetitionByID(context.Background(), petitionID)
	if err != nil {
		if err == sql.ErrNoRows {
			c.JSON(http.StatusNotFound, gin.H{"error": "Petition not found"})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve petition"})
		}
		return
	}

	if petition.UserID != userID {
		c.JSON(http.StatusForbidden, gin.H{"error": "You are not authorized to delete this petition"})
		return
	}

	err = h.Queries.DeletePetition(context.Background(), petitionID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete petition"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Petition deleted successfully"})
}
