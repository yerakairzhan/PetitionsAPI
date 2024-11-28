package handlers

import (
	"petitionsGO/db"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	Queries      *db.Queries
	JWTSecretKey string
}

func (h *Handler) InitRoutes() *gin.Engine {
	router := gin.New()
	router.POST("/register", h.register)
	router.POST("/login", h.login)
	router.GET("/list/petition", h.listPetition)
	router.POST("/add/petition", h.createPetition)
	router.POST("/vote/:id/petition", h.votePetition)
	router.POST("/delete/:id/petition", h.deletePetition)

	return router
}
