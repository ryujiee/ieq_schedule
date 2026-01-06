package http

import (
	"errors"
	"net/http"
	"time"

	"ieq/backend/internal/models"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type daySaveReq struct {
	Items []daySaveItem `json:"items"`
}

type daySaveItem struct {
	FunctionID uint  `json:"functionId"`
	MemberID   *uint `json:"memberId"` // nil => remove a escala daquela função no dia
}

func adminListScheduleHandler(db *gorm.DB) gin.HandlerFunc {
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

func adminGetScheduleDayHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := parseDateOnly(dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date (YYYY-MM-DD)"})
			return
		}

		var items []models.ScheduleAssignment
		if err := db.
			Preload("TeamFunction").
			Preload("Member").
			Where("date = ?", date).
			Order("team_function_id asc").
			Find(&items).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"date":  dateStr,
			"items": items,
		})
	}
}

func adminPutScheduleDayHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		dateStr := c.Param("date")
		date, err := parseDateOnly(dateStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid date (YYYY-MM-DD)"})
			return
		}

		var req daySaveReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		// valida itens
		seen := map[uint]struct{}{}
		for _, it := range req.Items {
			if it.FunctionID == 0 {
				c.JSON(http.StatusBadRequest, gin.H{"error": "functionId is required"})
				return
			}
			if _, ok := seen[it.FunctionID]; ok {
				c.JSON(http.StatusBadRequest, gin.H{"error": "duplicate functionId in items"})
				return
			}
			seen[it.FunctionID] = struct{}{}
		}

		// Transação: aplica todas as mudanças do dia
		err = db.Transaction(func(tx *gorm.DB) error {
			for _, it := range req.Items {
				// se memberId vier nil -> deletar a atribuição daquele function no dia
				if it.MemberID == nil || *it.MemberID == 0 {
					if err := tx.
						Where("date = ? AND team_function_id = ?", date, it.FunctionID).
						Delete(&models.ScheduleAssignment{}).Error; err != nil {
						return err
					}
					continue
				}

				// valida existência de function e member
				if err := tx.First(&models.TeamFunction{}, it.FunctionID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("invalid functionId")
					}
					return err
				}
				if err := tx.First(&models.Member{}, *it.MemberID).Error; err != nil {
					if errors.Is(err, gorm.ErrRecordNotFound) {
						return errors.New("invalid memberId")
					}
					return err
				}

				// upsert: (date + function) é único
				var existing models.ScheduleAssignment
				err := tx.Where("date = ? AND team_function_id = ?", date, it.FunctionID).First(&existing).Error
				if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
					return err
				}

				if errors.Is(err, gorm.ErrRecordNotFound) {
					newItem := models.ScheduleAssignment{
						Date:           date,
						TeamFunctionID: it.FunctionID,
						MemberID:       *it.MemberID,
					}
					if err := tx.Create(&newItem).Error; err != nil {
						return err
					}
				} else {
					existing.MemberID = *it.MemberID
					if err := tx.Save(&existing).Error; err != nil {
						return err
					}
				}
			}
			return nil
		})

		if err != nil {
			// mensagens amigáveis
			if err.Error() == "invalid functionId" || err.Error() == "invalid memberId" {
				c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save day"})
			return
		}

		// retorna o dia atualizado
		var items []models.ScheduleAssignment
		if err := db.Preload("TeamFunction").Preload("Member").
			Where("date = ?", date).
			Order("team_function_id asc").
			Find(&items).Error; err != nil {
			c.JSON(http.StatusOK, gin.H{"date": dateStr, "items": []models.ScheduleAssignment{}})
			return
		}

		c.JSON(http.StatusOK, gin.H{"date": dateStr, "items": items})
	}
}

func parseDateOnly(s string) (time.Time, error) {
	// date-only fixo, sem timezone (armazenado como DATE no Postgres)
	return time.Parse("2006-01-02", s)
}
