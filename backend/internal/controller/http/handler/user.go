package handler

import (
	"context"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"

	// "mado/internal"
	"mado/internal/auth"
	"mado/internal/auth/model"
	"mado/models"

	// "mado/internal/controller/http/httperr"
	"mado/internal/core/user"
)

type ECP struct {
	Ecp string `json:"ecp"    binding:"required,ecp"`
}

// UserService is a user service interface.
type UserService interface {
	// Create(ctx context.Context, dto user.CreateDTO) (user.User, error)
	Login(model.LoginRequirements) (*user.User, error)
	GetUser(ctx context.Context, userId int) (models.User, error)
	// GetAllRows()()
}

type userDeps struct {
	router *gin.RouterGroup

	userService UserService
}

type userHandler struct {
	userService UserService
}

func newUserHandler(deps userDeps) {
	handler := userHandler{
		userService: deps.userService,
	}

	usersGroup := deps.router.Group("/users")
	{
		usersGroup.GET("/", handler.getUser)
		usersGroup.POST("/", handler.createUser)    // api/users/
		usersGroup.POST("/login", handler.sendLink) // api/users/login
		usersGroup.POST("/confirm", handler.confirmCredentials)

	}

}

func (h userHandler) createUser(c *gin.Context) {
	fmt.Println("createUser")
	c.IndentedJSON(http.StatusOK, "User created")

}

func (h userHandler) getUser(c *gin.Context) {
	c.IndentedJSON(http.StatusOK, http.StatusText(http.StatusOK))
	fmt.Println("GetUser")
}

func (h userHandler) sendLink(c *gin.Context) {
	egovMobileLink, qrSigner, nonce := auth.PreparationStep()
	if egovMobileLink == nil || qrSigner == nil || nonce == nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"err": "egovMobileLink or qrSigner or nonce is nil"})
		return
	}
	requirements := model.LoginRequirements{QrSigner: qrSigner, Nonce: nonce}
	// go h.userService.Login(context.Background(), qrSigner, nonce)
	c.JSON(http.StatusOK, gin.H{"link": egovMobileLink, "requirements": requirements})
}

func (h userHandler) confirmCredentials(c *gin.Context) {
	var request *model.LoginRequirements
	if err := c.BindJSON(&request); err != nil {
		c.AbortWithStatus(http.StatusBadRequest)
		return
	}
	// request.Context = &gin.Context{}
	user, err := h.userService.Login(*request)
	if err != nil {
		fmt.Println(err)
		c.AbortWithStatus(http.StatusInternalServerError)
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": user})

}

// func (h userHandler) getAllRows(c *gin.Context) {
// 	rows := h.userService.GetAllRows()
// 	c.JSON(http.StatusOK, gin.H{"rows":rows})
// }
