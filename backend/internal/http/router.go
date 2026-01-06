package http

import (
	"time"

	"ieq/backend/internal/config"
	"ieq/backend/internal/middlewares"
	"ieq/backend/internal/services"
	"ieq/backend/internal/whatsapp"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func NewRouter(db *gorm.DB, cfg config.Config, waMgr *whatsapp.Manager, rem *services.ReminderScheduler) *gin.Engine {
	r := gin.Default()

	// CORS (frontend local)
	r.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:9001",
			"http://127.0.0.1:9001",
		},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))

	// Health
	r.GET("/health", health)

	// Public schedule (sem login)
	pub := r.Group("/public")
	pub.GET("/schedule", publicListScheduleHandler(db))
	pub.GET("/schedule/next", publicNextScheduleHandler(db))

	// Services
	jwtSvc, err := services.NewJWTService(cfg.JWTSecret)
	if err != nil {
		panic(err)
	}
	authSvc := services.NewAuthService(db, jwtSvc)

	// Auth
	r.POST("/auth/login", loginHandler(authSvc))

	// Protected
	protected := r.Group("/")
	protected.Use(middlewares.AuthRequired(jwtSvc))
	protected.GET("/me", meHandler())

	// Admin-only
	admin := protected.Group("/admin")
	admin.Use(middlewares.AdminOnly())

	admin.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"ok": true})
	})

	// WhatsApp (admin-only)
	admin.GET("/whatsapp/status", adminWhatsAppStatusHandler(waMgr))
	admin.POST("/whatsapp/connect", adminWhatsAppConnectHandler(waMgr))
	admin.GET("/whatsapp/qr", adminWhatsAppQRHandler(waMgr))
	admin.POST("/whatsapp/logout", adminWhatsAppLogoutHandler(waMgr))

	// Team Functions CRUD
	admin.GET("/functions", listFunctionsHandler(db))
	admin.POST("/functions", createFunctionHandler(db))
	admin.PUT("/functions/:id", updateFunctionHandler(db))
	admin.DELETE("/functions/:id", deleteFunctionHandler(db))

	// Members CRUD
	admin.GET("/members", listMembersHandler(db))
	admin.GET("/members/:id", getMemberHandler(db))
	admin.POST("/members", createMemberHandler(db))
	admin.PUT("/members/:id", updateMemberHandler(db))
	admin.DELETE("/members/:id", deleteMemberHandler(db))

	// Schedule (admin-only)
	admin.GET("/schedule", adminListScheduleHandler(db))
	admin.GET("/schedule/day/:date", adminGetScheduleDayHandler(db))
	admin.PUT("/schedule/day/:date", adminPutScheduleDayHandler(db))

	// Settings
	admin.GET("/settings", adminGetSettingsHandler(db))
	admin.PUT("/settings", adminUpdateSettingsHandler(db, rem))

	return r
}
