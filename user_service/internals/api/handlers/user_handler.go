package handlers

import (
	"net/http"
	"os"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"github.com/palashbhasme/ecommerce_microservices/common/middlewares"
	common "github.com/palashbhasme/ecommerce_microservices/common/models"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/mapper"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/request"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/api/dto/response"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/models"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/domain/repository"
	"github.com/palashbhasme/ecommerce_microservices/user_service/internals/services"
	"github.com/palashbhasme/ecommerce_microservices/user_service/utils"
	"go.uber.org/zap"
)

var validate = request.NewValidator()
var jwtKey = []byte(os.Getenv("JWT_SECRET"))

type UserHandler struct {
	Repo repository.UserRepository
	log  *zap.Logger
}

func InitializeRoutes(router *gin.Engine, log *zap.Logger, repo repository.UserRepository) {
	handler := &UserHandler{
		Repo: repo,
		log:  log,
	}
	authconfig := common.NewAuthConfig(os.Getenv("JWT_SECRET"))
	api := router.Group("/api")
	{
		userRoutes := api.Group("/users/v1")

		// Public Routes (No Middleware)
		userRoutes.POST("/", handler.CreateUser)     // Signup
		userRoutes.POST("/login", handler.LoginUser) // Login

		// Protected Routes (Require Authentication)
		protectedRoutes := userRoutes.Group("/")
		protectedRoutes.Use(middlewares.AuthMiddleware(*authconfig)) // Apply authentication middleware
		{
			protectedRoutes.GET("/:id", handler.GetUserById)
			protectedRoutes.PUT("/:id", handler.UpdateUser)

			//admin specific routes
			adminRoutes := protectedRoutes.Group("/")
			adminRoutes.Use(middlewares.AdminMiddleware()) //Apply admin middleware
			{
				adminRoutes.DELETE("/:id", handler.DeleteUser)
				adminRoutes.GET("/getall", handler.GetAllUsers)

			}
		}

		// Test Route (Public)
		userRoutes.GET("/test", func(c *gin.Context) {
			c.JSON(200, gin.H{"message": "alive"})
		})
	}
}

func (h *UserHandler) LoginUser(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required"`
		Password string `json:"password" binding:"required"`
	}

	// Validate request payload
	if err := c.ShouldBindJSON(&req); err != nil {
		h.log.Error("Invalid request data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	user, err := h.Repo.GetUserByEmail(req.Email)
	if err != nil {
		h.log.Error("Database error while checking existing user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not check for existing user",
		})
		return
	}
	if user == nil {
		h.log.Warn("User does not exists", zap.String("email", req.Email))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User with this email does not exists",
		})
		return
	}

	if err := utils.ComparePasswordHash(req.Password, user.Account.PasswordHash); err != nil {
		h.log.Error("password provided does not match", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Password is incorrect",
		})
		return
	}

	expirationTime := time.Now().Add(60 * time.Minute)
	claims := models.Claims{
		Role: user.Account.Role,
		StandardClaims: jwt.StandardClaims{
			Subject:   user.Email,
			ExpiresAt: expirationTime.Unix(),
		},
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(jwtKey)
	if err != nil {
		h.log.Error("error genearting token", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "could not generate token"})
		return
	}

	c.SetCookie("token", tokenString, int(expirationTime.Unix()), "/", "localhost", false, true)
	c.JSON(http.StatusOK, gin.H{"success": "user logged in"})

}
func (h *UserHandler) CreateUser(c *gin.Context) {
	var userRequest request.UserRequest

	// Validate request payload
	if err := c.ShouldBindJSON(&userRequest); err != nil {
		h.log.Error("Invalid request data", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request data",
			"error":   err.Error(),
		})
		return
	}

	// Validate struct fields
	if err := validate.Struct(userRequest); err != nil {
		h.log.Error("Validation failed", zap.Error(err))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Invalid request body",
		})
		return
	}

	// Check if the user already exists
	existingUser, err := h.Repo.GetUserByEmail(userRequest.Email)
	if err != nil {
		h.log.Error("Database error while checking existing user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Could not check for existing user",
		})
		return
	}

	if existingUser != nil {
		h.log.Warn("User already exists", zap.String("email", userRequest.Email))
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "User with this email already exists",
		})
		return
	}

	// Hash password
	hashedPassword, err := utils.GeneratePasswordHash(userRequest.Account.Password)
	if err != nil {
		h.log.Error("Failed to hash password", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to process user password",
		})
		return
	}

	// Map request to user model and set hashed password
	user := mapper.MapUserFromRequest(&userRequest)
	user.Account.PasswordHash = hashedPassword

	// Create user in database
	if err := h.Repo.CreateUser(&user); err != nil {
		h.log.Error("Failed to create user", zap.Error(err))
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to create user",
		})
		return
	}

	// Success response
	c.JSON(http.StatusCreated, gin.H{
		"message": "User created successfully",
	})
}

func (h *UserHandler) GetUserById(c *gin.Context) {
	id := c.Param("id")
	h.log.Info("Fetching user by ID", zap.String("id", id))

	objectID := models.MyObjectID(id)
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

	objectID := models.MyObjectID(id)

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
