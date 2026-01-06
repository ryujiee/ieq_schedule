package services

import (
	"fmt"
	"log"
	"strings"
	"time"

	"ieq/backend/internal/models"
	"ieq/backend/internal/whatsapp"

	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type ReminderScheduler struct {
	db       *gorm.DB
	location *time.Location
	cron     *cron.Cron
	entryID  cron.EntryID

	wa *whatsapp.Manager
}

func NewReminderScheduler(db *gorm.DB, wa *whatsapp.Manager) *ReminderScheduler {
	loc, err := time.LoadLocation("America/Sao_Paulo")
	if err != nil {
		loc = time.FixedZone("BRT", -3*60*60)
	}
	c := cron.New(cron.WithLocation(loc))
	return &ReminderScheduler{db: db, location: loc, cron: c, entryID: 0, wa: wa}
}

func (s *ReminderScheduler) Start() {
	s.cron.Start()
	s.ReloadFromDB()
}

func (s *ReminderScheduler) Stop() {
	s.cron.Stop()
}

func (s *ReminderScheduler) ReloadFromDB() {
	// remove job anterior
	if s.entryID != 0 {
		s.cron.Remove(s.entryID)
		s.entryID = 0
	}

	var st models.Settings
	if err := s.db.Order("id asc").First(&st).Error; err != nil {
		st = *models.DefaultSettings()
		_ = s.db.Create(&st).Error
	}

	// compat: se você mudou pra hour/minute, adapte aqui
	if !st.ReminderEnabled {
		log.Println("[reminder] disabled")
		return
	}

	hhmm := strings.TrimSpace(st.ReminderTimeHHMM)
	if len(hhmm) != 5 || hhmm[2] != ':' {
		log.Println("[reminder] invalid reminderTime:", hhmm)
		return
	}
	hh := hhmm[0:2]
	mm := hhmm[3:5]

	// robfig/cron v3 padrão: 5 campos (min hora dia mes semana)
	spec := fmt.Sprintf("%s %s * * *", mm, hh)

	id, err := s.cron.AddFunc(spec, func() {
		s.runOnce()
	})
	if err != nil {
		log.Println("[reminder] cron add error:", err)
		return
	}

	s.entryID = id
	log.Printf("[reminder] scheduled daily at %s (America/Sao_Paulo) spec=%s\n", hhmm, spec)
}

func (s *ReminderScheduler) runOnce() {
	now := time.Now().In(s.location)
	day := time.Date(now.Year(), now.Month(), now.Day(), 0, 0, 0, 0, s.location)

	var st models.Settings
	if err := s.db.Order("id asc").First(&st).Error; err != nil || !st.ReminderEnabled {
		return
	}

	// se whatsapp não estiver conectado, não adianta
	if s.wa == nil {
		log.Println("[reminder] whatsapp manager not configured")
		return
	}
	waStatus := s.wa.Status()
	if !waStatus.LoggedIn || !waStatus.Connected {
		log.Println("[reminder] whatsapp not connected/logged in - skipping")
		return
	}

	var items []models.ScheduleAssignment
	if err := s.db.Preload("TeamFunction").Preload("Member").
		Where("date = ?", day).
		Find(&items).Error; err != nil {
		log.Println("[reminder] db error:", err)
		return
	}

	if len(items) == 0 {
		log.Println("[reminder] no assignments for today")
		return
	}

	for _, it := range items {
		// precisa ter membro + phone
		if it.Member.ID == 0 {
			continue
		}
		phone := strings.TrimSpace(it.Member.Phone)
		if phone == "" {
			log.Printf("[reminder] member=%s has no phone, skipping\n", it.Member.Name)
			continue
		}

		name := it.Member.Name
		fn := ""
		if it.TeamFunction.ID != 0 {
			fn = it.TeamFunction.Name
		}

		msg := st.ReminderTemplate
		if strings.TrimSpace(msg) == "" {
			msg = "Olá {Nome}! Hoje é o seu dia de servir na função {Funcao}. Deus abençoe! 🙏"
		}
		msg = strings.ReplaceAll(msg, "{Nome}", name)
		msg = strings.ReplaceAll(msg, "{Funcao}", fn)

		// backend recebe só dígitos tipo 5549...
		toDigits := onlyDigits(phone)

		if err := s.wa.SendText(toDigits, msg); err != nil {
			log.Printf("[reminder] send error to=%s name=%s: %v\n", toDigits, name, err)
			continue
		}

		log.Printf("[reminder] sent to=%s name=%s func=%s\n", toDigits, name, fn)
	}
}

func onlyDigits(s string) string {
	var b strings.Builder
	for _, r := range s {
		if r >= '0' && r <= '9' {
			b.WriteRune(r)
		}
	}
	return b.String()
}
