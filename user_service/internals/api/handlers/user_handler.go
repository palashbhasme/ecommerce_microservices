package handlers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/mapper"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/request"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/response"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/repository"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/services"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.uber.org/zap"
)

var validate = request.NewValidator()

type UserHandler struct {
	Repo repository.UserRepository
	log  *zap.Logger
}

func InitializeRoutes(router *gin.Engine, log *zap.Logger, repo repository.UserRepository) {
	handler := &UserHandler{
		Repo: repo,
		log:  log,
	}
	api := router.Group("/api")
	{
		userRoutes := api.Group("/users/v1")
		{
			userRoutes.POST("/", handler.CreateUser)
			userRoutes.GET("/:id", handler.GetUserById)
			userRoutes.PUT("/:id", handler.UpdateUser)
			userRoutes.DELETE("/:id", handler.DeleteUser)
			userRoutes.GET("/getAll", handler.GetAllUsers)
			userRoutes.GET("/test", func(c *gin.Context) {
				c.JSON(200, gin.H{"message": "alive"})
			})
		}
	}
}

func (h *UserHandler) CreateUser(c *gin.Context) {
	var userRequest request.UserRequest

	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.log.Error("Invalid request data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	if err := validate.Struct(userRequest); err != nil {
		h.log.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
		})
		return
	}

	user := mapper.MapUserFromRequest(&userRequest)

	err := h.Repo.CreateUser(&user)
	if err != nil {
		h.log.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")
	h.log.Info("Fetching user by ID", zap.String("id", id))

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.log.Error("Invalid user ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
		})
		return
	}
	user, err := h.Repo.GetUserById(objectID)
	if err != nil {
		h.log.Error("Failed to fetch user", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user",
		})
		return
	}

	userResponse := mapper.MapUserToResponse(*user)

	c.JSON(http.StatusOK, gin.H{
		"message": "User fetched successfully",
		"user":    userResponse,
	})

}

func (h *UserHandler) UpdateUser(c *gin.Context) {

	id := c.Param("id")
	h.log.Info("Updating user by ID", zap.String("id", id))

	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		h.log.Error("Invalid user ID", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid user ID",
			"error":   err.Error(),
		})
		return
	}

	existingUser, err := h.Repo.GetUserById(objectID)
	if err != nil {
		h.log.Error("Failed to fetch user", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch user",
			"error":   err.Error(),
		})
		return
	}

	var userRequest request.UpdateUserRequest
	if err = c.ShouldBindJSON(&userRequest); err != nil {
		h.log.Error("Invalid Request Data", zap.String("id", id))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	if err := validate.Struct(userRequest); err != nil {
		h.log.Error("Failed to validate user", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "invalid request body",
			"error":   err.Error(),
		})
		return
	}
	user, err := services.UpdateUser(existingUser, userRequest)
	if err != nil {
		h.log.Error("Failed to update user data", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "user not found",
			"error":   err.Error(),
		})
		return
	}

	err = h.Repo.UpdateUser(objectID, user)
	if err != nil {
		h.log.Error("Failed to update user", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to update user",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated successfully",
	})

}

func (h *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	h.log.Info("Deleting user by ID", zap.String("id", id))

	err := h.Repo.DeleteUser(id)
	if err != nil {
		h.log.Error("Failed to delete user", zap.String("id", id), zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to delete user",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "User deleted successfully",
	})
}

func (h *UserHandler) GetAllUsers(c *gin.Context) {
	h.log.Info("Fetching all users")

	users, err := h.Repo.GetAllUsers()
	if err != nil {
		h.log.Error("Failed to fetch users", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to fetch users",
		})
		return
	}

	userResponses := make([]response.UserResponse, 0)
	for _, user := range users {
		userResponses = append(userResponses, mapper.MapUserToResponse(*user))
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Users fetched successfully",
		"users":   userResponses,
	})
}
