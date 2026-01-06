package http

import (
	"net/http"
	"strconv"
	"strings"

	"ieq/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type functionCreateReq struct {
	Name   string `json:"name"`
	Active *bool  `json:"active"` // opcional
}

type functionUpdateReq struct {
	Name   *string `json:"name"`
	Active *bool   `json:"active"`
}

func listFunctionsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.TeamFunction
		if err := db.Order("name asc").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func createFunctionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req functionCreateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		name := strings.TrimSpace(req.Name)
		if name == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
			return
		}

		active := true
		if req.Active != nil {
			active = *req.Active
		}

		item := models.TeamFunction{
			Name:   name,
			Active: active,
		}

		if err := db.Create(&item).Error; err != nil {
			// provável conflito de unique
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not create (maybe duplicated name)"})
			return
		}

		c.JSON(http.StatusCreated, item)
	}
}

func updateFunctionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var req functionUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		var item models.TeamFunction
		if err := db.First(&item, uint(id)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		if req.Name != nil {
			name := strings.TrimSpace(*req.Name)
			if name == "" {
				c.JSON(http.StatusBadRequest, gin.H{"error": "name cannot be empty"})
				return
			}
			item.Name = name
		}

		if req.Active != nil {
			item.Active = *req.Active
		}

		if err := db.Save(&item).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not update (maybe duplicated name)"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func deleteFunctionHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var item models.TeamFunction
		if err := db.First(&item, uint(id)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		if err := db.Delete(&item).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not delete"})
			return
		}

		c.Status(http.StatusNoContent)
	}
}
