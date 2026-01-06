package http

import (
	"net/http"
	"strings"

	"ieq/backend/internal/models"
	"ieq/backend/internal/services"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type updateSettingsReq struct {
	// NOVO (frontend atual)
	RemindersEnabled *bool  `json:"remindersEnabled"`
	ReminderHour     *int   `json:"reminderHour"`
	ReminderMinute   *int   `json:"reminderMinute"`
	ReminderMessage  string `json:"reminderMessage"`

	// LEGADO (se quiser aceitar também)
	ReminderEnabled  *bool  `json:"reminderEnabled"`
	ReminderTimeHHMM string `json:"reminderTime"`     // "HH:MM"
	ReminderTemplate string `json:"reminderTemplate"` // {Nome} {Funcao}

	ReminderChannelID *uint `json:"reminderChannelId"`
}

func adminGetSettingsHandler(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var s models.Settings
		if err := db.Order("id asc").First(&s).Error; err != nil {
			s = *models.DefaultSettings()
			_ = db.Create(&s).Error
		}

		// Responde no formato NOVO pro frontend
		hh, mm := parseHHMM(s.ReminderTimeHHMM)
		c.JSON(http.StatusOK, gin.H{
			"remindersEnabled": s.ReminderEnabled,
			"reminderHour":     hh,
			"reminderMinute":   mm,
			"reminderMessage":  s.ReminderTemplate,
		})
	}
}

func adminUpdateSettingsHandler(db *gorm.DB, rem *services.ReminderScheduler) gin.HandlerFunc {
	return func(c *gin.Context) {
		var req updateSettingsReq
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid body"})
			return
		}

		var s models.Settings
		if err := db.Order("id asc").First(&s).Error; err != nil {
			s = *models.DefaultSettings()
			_ = db.Create(&s).Error
		}

		// ====== NOVO (prioridade) ======
		if req.RemindersEnabled != nil {
			s.ReminderEnabled = *req.RemindersEnabled
		}

		if req.ReminderHour != nil || req.ReminderMinute != nil {
			hh, mm := parseHHMM(s.ReminderTimeHHMM)

			if req.ReminderHour != nil {
				if *req.ReminderHour < 0 || *req.ReminderHour > 23 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "reminderHour must be 0-23"})
					return
				}
				hh = *req.ReminderHour
			}
			if req.ReminderMinute != nil {
				if *req.ReminderMinute < 0 || *req.ReminderMinute > 59 {
					c.JSON(http.StatusBadRequest, gin.H{"error": "reminderMinute must be 0-59"})
					return
				}
				mm = *req.ReminderMinute
			}

			s.ReminderTimeHHMM = formatHHMM(hh, mm)
		}

		if req.ReminderMessage != "" {
			s.ReminderTemplate = req.ReminderMessage
		}

		// ====== LEGADO (se vier no body antigo) ======
		if req.ReminderEnabled != nil {
			s.ReminderEnabled = *req.ReminderEnabled
		}
		if strings.TrimSpace(req.ReminderTimeHHMM) != "" {
			t := strings.TrimSpace(req.ReminderTimeHHMM)
			if !validHHMM(t) {
				c.JSON(http.StatusBadRequest, gin.H{"error": "reminderTime must be HH:MM"})
				return
			}
			s.ReminderTimeHHMM = t
		}
		if req.ReminderTemplate != "" {
			s.ReminderTemplate = req.ReminderTemplate
		}

		// pode ser nil (desvincular)
		s.ReminderChannelID = req.ReminderChannelID

		if err := db.Save(&s).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "could not save settings"})
			return
		}

		if rem != nil {
			rem.ReloadFromDB()
		}

		hh, mm := parseHHMM(s.ReminderTimeHHMM)
		c.JSON(http.StatusOK, gin.H{
			"remindersEnabled": s.ReminderEnabled,
			"reminderHour":     hh,
			"reminderMinute":   mm,
			"reminderMessage":  s.ReminderTemplate,
		})
	}
}

func validHHMM(t string) bool {
	if len(t) != 5 || t[2] != ':' {
		return false
	}
	hh := t[0:2]
	mm := t[3:5]
	return isTwoDigits(hh) && isTwoDigits(mm) && hh <= "23" && mm <= "59"
}

func isTwoDigits(s string) bool {
	return len(s) == 2 && s[0] >= '0' && s[0] <= '9' && s[1] >= '0' && s[1] <= '9'
}

func parseHHMM(t string) (int, int) {
	// default seguro
	if !validHHMM(strings.TrimSpace(t)) {
		return 8, 0
	}
	h := int((t[0]-'0')*10 + (t[1] - '0'))
	m := int((t[3]-'0')*10 + (t[4] - '0'))
	return h, m
}

func formatHHMM(h, m int) string {
	return pad2(h) + ":" + pad2(m)
}

func pad2(v int) string {
	if v < 10 {
		return "0" + string(rune('0'+v))
	}
	return string(rune('0'+(v/10))) + string(rune('0'+(v%10)))
}
