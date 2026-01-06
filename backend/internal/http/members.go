package http

import (
	"net/http"
	"strconv"
	"strings"

	"ieq/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type memberCreateReq struct {
	Name        string `json:"name"`
	Phone       string `json:"phone"`
	Active      *bool  `json:"active"`
	FunctionIDs []uint `json:"functionIds"`
}

type memberUpdateReq struct {
	Name        *string `json:"name"`
	Phone       *string `json:"phone"`
	Active      *bool   `json:"active"`
	FunctionIDs *[]uint `json:"functionIds"` // ponteiro: só atualiza se vier no body
}

func listMembersHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var items []models.Member
		if err := db.Preload("Functions").Order("name asc").Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}
		c.JSON(http.StatusOK, items)
	}
}

func getMemberHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var item models.Member
		if err := db.Preload("Functions").First(&item, uint(id)).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "not found"})
			return
		}

		c.JSON(http.StatusOK, item)
	}
}

func createMemberHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req memberCreateReq
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

		phone := strings.TrimSpace(req.Phone)

		var functions []models.TeamFunction
		if len(req.FunctionIDs) > 0 {
			if err := db.Where("id IN ?", req.FunctionIDs).Find(&functions).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
				return
			}
			if len(functions) != len(uniqueUints(req.FunctionIDs)) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "one or more functionIds are invalid"})
				return
			}
		}

		item := models.Member{
			Name:      name,
			Phone:     phone,
			Active:    active,
			Functions: functions,
		}

		if err := db.Create(&item).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not create"})
			return
		}

		// garantir preload no retorno
		if err := db.Preload("Functions").First(&item, item.ID).Error; err != nil {
			c.JSON(http.StatusCreated, item)
			return
		}
		c.JSON(http.StatusCreated, item)
	}
}

func updateMemberHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var req memberUpdateReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		var item models.Member
		if err := db.Preload("Functions").First(&item, uint(id)).Error; err != nil {
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

		if req.Phone != nil {
			item.Phone = strings.TrimSpace(*req.Phone)
		}

		if req.Active != nil {
			item.Active = *req.Active
		}

		// Atualiza campos básicos primeiro
		if err := db.Save(&item).Error; err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "could not update"})
			return
		}

		// Se veio functionIds, substitui vínculo (Replace)
		if req.FunctionIDs != nil {
			var functions []models.TeamFunction
			if len(*req.FunctionIDs) > 0 {
				if err := db.Where("id IN ?", *req.FunctionIDs).Find(&functions).Error; err != nil {
					c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
					return
				}
				if len(functions) != len(uniqueUints(*req.FunctionIDs)) {
					c.JSON(http.StatusBadRequest, gin.H{"error": "one or more functionIds are invalid"})
					return
				}
			}

			if err := db.Model(&item).Association("Functions").Replace(functions); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "could not update functions"})
				return
			}
		}

		if err := db.Preload("Functions").First(&item, item.ID).Error; err != nil {
			c.JSON(http.StatusOK, item)
			return
		}
		c.JSON(http.StatusOK, item)
	}
}

func deleteMemberHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil || id <= 0 {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid id"})
			return
		}

		var item models.Member
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

func uniqueUints(in []uint) []uint {
	seen := map[uint]struct{}{}
	out := make([]uint, 0, len(in))
	for _, v := range in {
		if _, ok := seen[v]; ok {
			continue
		}
		seen[v] = struct{}{}
		out = append(out, v)
	}
	return out
}
