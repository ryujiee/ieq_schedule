package main

import (
	"log"

	"ieq/backend/internal/config"
	"ieq/backend/internal/db"
	"ieq/backend/internal/http"
	"ieq/backend/internal/services"
	"ieq/backend/internal/whatsapp"
)

func main() {
	cfg := config.Load()

	database, err := db.Connect(cfg)
	if err != nil {
		log.Fatalf("db connect error: %v", err)
	}

	if err := db.AutoMigrateAndSeed(database, db.SeedConfig{
		AdminEmail:    "admin@ieq.local",
		AdminPassword: "ieq25032010",
		AdminName:     "Admin IEQ",
	}); err != nil {
		log.Fatalf("migration/seed error: %v", err)
	}

	// WhatsApp manager (single instance)
	waMgr, err := whatsapp.NewManager()
	if err != nil {
		log.Fatalf("whatsapp init error: %v", err)
	}
	waMgr.ConnectOrRestore()
	defer waMgr.Close()

	// Scheduler
	rem := services.NewReminderScheduler(database, waMgr)
	rem.Start()
	defer rem.Stop()

	// Router (passa waMgr + rem)
	r := http.NewRouter(database, cfg, waMgr, rem)

	log.Printf("server running on :%s", cfg.AppPort)
	if err := r.Run(":" + cfg.AppPort); err != nil {
		log.Fatalf("server error: %v", err)
	}
}
