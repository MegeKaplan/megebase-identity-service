package handlers

import (
	"github.com/MegeKaplan/megebase-identity-service/services"
	"github.com/MegeKaplan/megebase-identity-service/utils"
	"github.com/MegeKaplan/megebase-identity-service/utils/response"
	"github.com/gin-gonic/gin"
)

type userHandler struct {
	userService services.UserService
}

func NewUserHandler(userService services.UserService) *userHandler {
	return &userHandler{userService: userService}
}

func (h *userHandler) GetMe() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID, exists := c.Get("userID")
		if !exists {
			utils.JSONError(c, response.ErrUnauthorized, "asd")
			return
		}

		user, err := h.userService.GetUserByID(userID.(string))
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}

		utils.JSONSuccess(c, response.UserFetched, user)
	}
}

func (h *userHandler) GetUserByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		user, err := h.userService.GetUserByID(id)
		if err != nil {
			utils.JSONError(c, err, err.Details)
			return
		}
		utils.JSONSuccess(c, response.UserFetched, user)
	}
}
