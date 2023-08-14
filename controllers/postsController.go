package controllers

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/natha/kyoyu-backend/initializers"
	"github.com/natha/kyoyu-backend/models"
)

func PostShow(c *gin.Context) {
	user, _ := c.Get("user")

	posts := []models.Post{}
	initializers.DB.Where("author_id = ?", user.(models.User).UserID).Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}

func PostCreate(c *gin.Context) {
	user, _ := c.Get("user")

	var body struct {
		Title string
		Slug  string
		Body  string
	}

	if c.Bind(&body) != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": "Failed to read request",
		})

		return
	}

	post := models.Post{
		Title:    body.Title,
		Slug:     body.Slug,
		Body:     body.Body,
		AuthorID: user.(models.User).UserID,
		AuthorEmail:    user.(models.User).Email,
	}

	if err := initializers.DB.Create(&post).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": "Failed to create post",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"data": post,
	})

}

func PostIndex(c *gin.Context) {
	posts := []models.Post{}

	initializers.DB.Find(&posts)

	c.JSON(http.StatusOK, gin.H{
		"data": posts,
	})
}
