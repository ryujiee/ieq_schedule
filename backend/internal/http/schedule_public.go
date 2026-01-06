package http

import (
	"net/http"
	"strconv"
	"time"

	"ieq/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func publicListScheduleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		fromStr := c.Query("from")
		toStr := c.Query("to")
		if fromStr == "" || toStr == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "from and to are required (YYYY-MM-DD)"})
			return
		}

		from, err := parseDateOnly(fromStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid from date"})
			return
		}
		to, err := parseDateOnly(toStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid to date"})
			return
		}

		var items []models.ScheduleAssignment
		if err := db.
			Preload("TeamFunction").
			Preload("Member").
			Where("date >= ? AND date <= ?", from, to).
			Order("date asc").
			Order("team_function_id asc").
			Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, items)
	}
}

func publicNextScheduleHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		limit := 10
		if s := c.Query("limit"); s != "" {
			if v, err := strconv.Atoi(s); err == nil && v > 0 && v <= 100 {
				limit = v
			}
		}

		// "hoje" como date-only
		now := time.Now()
		today := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, time.UTC)

		var items []models.ScheduleAssignment
		if err := db.
			Preload("TeamFunction").
			Preload("Member").
			Where("date >= ?", today).
			Order("date asc").
			Order("team_function_id asc").
			Limit(limit).
			Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"from":  today.Format("2006-01-02"),
			"limit": limit,
			"items": items,
		})
	}
}
