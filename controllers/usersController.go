package controllers

import (
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
	"github.com/natha/kyoyu-backend/initializers"
	"github.com/natha/kyoyu-backend/models"
	"golang.org/x/crypto/bcrypt"
)

func SignUp(c *gin.Context) {
	// Get Email/Password from Body
	var body struct {
		Email    string
		Password string
		Role     models.UserRole
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request",
		})

		return
	}

	// Hash the Password
	hash, err := bcrypt.GenerateFromPassword([]byte(body.Password), 10)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})

		return
	}

	// Create the User
	user := models.User{Email: body.Email, Password: string(hash), Role: body.Role}
	result := initializers.DB.Create(&user)

	if result.Error != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": result.Error.Error(),
		})

		return
	}

	// Respond
	c.JSON(http.StatusOK, gin.H{
		"message": "User created successfully",
	})
}

func Login(c *gin.Context) {
	// Get the Email and Password from Body
	var body struct {
		Email    string
		Password string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request",
		})

		return
	}

	// Look up requested User
	var user models.User
	initializers.DB.Find(&user, "Email = ?", body.Email)

	if user.UserID == uuid.Nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})

		return
	}

	// Compare the Password sent with the hashed Password
	err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(body.Password))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Invalid Email or Password",
		})

		return
	}

	// Generate JWT Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.UserID,
		"exp": time.Now().Add(time.Hour * 12).Unix(),
	})

	// Sign the JWT
	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET")))

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to create token",
		})

		return
	}

	// Set the cookie
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    tokenString,
		Expires:  time.Now().Add(time.Hour * 12),
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)

	type UserWithoutPassword struct {
		UserID uuid.UUID `json:"userId"`
		Email  string   `json:"email"`
		Role   models.UserRole `json:"role"`
		Token  string  `json:"token"`
	}

	userWithoutPassword := UserWithoutPassword{
		UserID: user.UserID,
		Email:  user.Email,
		Role:   user.Role,
		Token:  tokenString,
	}

	c.JSON(http.StatusOK, gin.H{"message": "Login successful", "data": userWithoutPassword})
}

func Logout(c *gin.Context) {
	// Delete the "Authorization" cookie
	cookie := http.Cookie{
		Name:     "Authorization",
		Value:    "",
		Expires:  time.Now().Add(-time.Hour), // Expire the cookie
		HttpOnly: true,
		SameSite: http.SameSiteLaxMode,
		Secure:   false,
		Path:     "/",
	}
	http.SetCookie(c.Writer, &cookie)

	c.JSON(http.StatusOK, gin.H{
		"message": "Logout successful",
	})
}

func UserShow(c *gin.Context) {
	user, _ := c.Get("user")

	// Load user with associated posts
	var userWithPosts models.User
	result := initializers.DB.Preload("Posts").First(&userWithPosts, "user_id = ?", user.(models.User).UserID)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error retrieving user and posts",
		})
		return
	}

	// Create a user response without password
	userWithoutPassword := models.User{
		UserID: userWithPosts.UserID,
		Email:  userWithPosts.Email,
		Role:   userWithPosts.Role,
		Posts:  userWithPosts.Posts,
		CreatedAt: userWithPosts.CreatedAt,
		UpdatedAt: userWithPosts.UpdatedAt,
	}

	c.JSON(http.StatusOK, gin.H{
		"data": userWithoutPassword,
	})
}

func UserIndex(c *gin.Context) {
	users := []models.User{}

	// Load users with associated posts
	result := initializers.DB.Preload("Posts").Find(&users)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Error retrieving users and posts",
		})
		return
	}

	// Create user responses without passwords
	usersWithoutPasswords := make([]models.User, len(users))
	for i, user := range users {
		userWithoutPassword := models.User{
			UserID: user.UserID,
			Email:  user.Email,
			Role:   user.Role,
			Posts:  user.Posts,
			CreatedAt: user.CreatedAt,
			UpdatedAt: user.UpdatedAt,
		}
		usersWithoutPasswords[i] = userWithoutPassword
	}

	c.JSON(http.StatusOK, gin.H{
		"data": usersWithoutPasswords,
	})
}
