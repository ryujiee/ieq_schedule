package whatsapp

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"sync"
	"time"

	"go.mau.fi/whatsmeow"
	"go.mau.fi/whatsmeow/store/sqlstore"
	"go.mau.fi/whatsmeow/types"
	"go.mau.fi/whatsmeow/types/events"
	waLog "go.mau.fi/whatsmeow/util/log"

	waE2E "go.mau.fi/whatsmeow/binary/proto"

	_ "github.com/jackc/pgx/v5/stdlib"
)

type Status struct {
	Connected bool   `json:"connected"`
	LoggedIn  bool   `json:"loggedIn"`
	Pairing   bool   `json:"pairing"`
	LastQR    string `json:"lastQr,omitempty"`
}

type Manager struct {
	mu sync.RWMutex

	db        *sql.DB
	container *sqlstore.Container
	client    *whatsmeow.Client

	pairing bool
	lastQR  string

	connectCancel context.CancelFunc
}

func NewManager() (*Manager, error) {
	dsn, err := buildPostgresDSNFromEnv()
	if err != nil {
		return nil, err
	}

	sqlDB, err := sql.Open("pgx", dsn)
	if err != nil {
		return nil, fmt.Errorf("open pg: %w", err)
	}

	logger := waLog.Stdout("wa", "INFO", true)

	ctx, cancel := context.WithTimeout(context.Background(), 20*time.Second)
	defer cancel()

	container, err := sqlstore.New(ctx, "pgx", dsn, logger)
	if err != nil {
		return nil, fmt.Errorf("sqlstore new: %w", err)
	}

	// cria/atualiza tabelas do store do whatsmeow
	if err := container.Upgrade(ctx); err != nil {
		return nil, fmt.Errorf("sqlstore upgrade: %w", err)
	}

	deviceStore, err := container.GetFirstDevice(ctx)
	if err != nil {
		return nil, fmt.Errorf("get device: %w", err)
	}

	client := whatsmeow.NewClient(deviceStore, logger)

	m := &Manager{
		db:        sqlDB,
		container: container,
		client:    client,
		pairing:   false,
		lastQR:    "",
	}

	// eventos essenciais (QR NÃO vem daqui; vem via channel)
	client.AddEventHandler(func(evt any) {
		switch evt.(type) {
		case *events.Connected:
			log.Println("[whatsapp] connected")
			m.setPairing(false, "")
		case *events.Disconnected:
			log.Println("[whatsapp] disconnected")
		case *events.LoggedOut:
			log.Println("[whatsapp] logged out")
			m.setPairing(false, "")
		}
	})

	return m, nil
}

func (m *Manager) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.connectCancel != nil {
		m.connectCancel()
		m.connectCancel = nil
	}
	if m.client != nil {
		m.client.Disconnect()
	}
	if m.db != nil {
		_ = m.db.Close()
	}
}

func (m *Manager) Status() Status {
	m.mu.RLock()
	defer m.mu.RUnlock()

	connected := false
	loggedIn := false

	if m.client != nil {
		connected = m.client.IsConnected()
		loggedIn = m.client.Store != nil && m.client.Store.ID != nil
	}

	return Status{
		Connected: connected,
		LoggedIn:  loggedIn,
		Pairing:   m.pairing,
		LastQR:    m.lastQR,
	}
}

func (m *Manager) LastQR() string {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.lastQR
}

func (m *Manager) ConnectOrRestore() {
	st := m.Status()
	if st.LoggedIn {
		log.Println("[whatsapp] session found -> connecting...")
		m.client.Connect()
		return
	}
	log.Println("[whatsapp] no session -> waiting admin connect (QR)")
}

func (m *Manager) StartPairing() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client == nil {
		return errors.New("whatsapp client not initialized")
	}

	// já tem sessão
	if m.client.Store != nil && m.client.Store.ID != nil {
		if !m.client.IsConnected() {
			m.client.Connect()
		}
		return nil
	}

	// já em pairing
	if m.pairing {
		return nil
	}

	// cancela pairing anterior
	if m.connectCancel != nil {
		m.connectCancel()
		m.connectCancel = nil
	}

	ctx, cancel := context.WithCancel(context.Background())
	m.connectCancel = cancel

	// QR channel precisa ser obtido antes do Connect()
	qrChan, err := m.client.GetQRChannel(ctx)
	if err != nil {
		return fmt.Errorf("get qr channel: %w", err)
	}

	m.pairing = true
	m.lastQR = ""

	go func() {
		defer func() {
			m.setPairing(false, m.LastQR())
		}()

		log.Println("[whatsapp] connecting for pairing...")
		m.client.Connect()

		for evt := range qrChan {
			// evt.Event: "code", "success", "timeout", ...
			if evt.Event == "code" && evt.Code != "" {
				m.setLastQR(evt.Code)
			}
			if evt.Event == "success" {
				log.Println("[whatsapp] pairing success")
				return
			}
			if evt.Event == "timeout" {
				log.Println("[whatsapp] pairing timeout")
				return
			}
		}
	}()

	return nil
}

func (m *Manager) Logout() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	if m.client == nil {
		return errors.New("whatsapp client not initialized")
	}

	if m.connectCancel != nil {
		m.connectCancel()
		m.connectCancel = nil
	}

	m.pairing = false
	m.lastQR = ""

	// Logout agora exige context
	return m.client.Logout(context.Background())
}

func (m *Manager) SendText(toDigits string, message string) error {
	toDigits = strings.TrimSpace(toDigits)
	toDigits = onlyDigits(toDigits)
	if toDigits == "" {
		return errors.New("empty phone")
	}

	if !m.client.IsConnected() {
		return errors.New("whatsapp not connected")
	}

	jid := types.NewJID(toDigits, "s.whatsapp.net")
	msg := &waE2E.Message{
		Conversation: &message,
	}

	_, err := m.client.SendMessage(context.Background(), jid, msg)
	return err
}

func (m *Manager) setLastQR(code string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.lastQR = code
}

func (m *Manager) setPairing(pairing bool, lastQR string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pairing = pairing
	if lastQR != "" {
		m.lastQR = lastQR
	}
}

func buildPostgresDSNFromEnv() (string, error) {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	pass := os.Getenv("DB_PASS")
	name := os.Getenv("DB_NAME")
	ssl := os.Getenv("DB_SSLMODE")

	if host == "" || port == "" || user == "" || name == "" {
		return "", errors.New("missing DB envs (DB_HOST, DB_PORT, DB_USER, DB_NAME)")
	}
	if ssl == "" {
		ssl = "disable"
	}

	// pgx dsn: postgres://user:pass@host:5432/db?sslmode=disable
	if pass != "" {
		return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
			urlEscape(user), urlEscape(pass), host, port, name, ssl,
		), nil
	}
	return fmt.Sprintf("postgres://%s@%s:%s/%s?sslmode=%s",
		urlEscape(user), host, port, name, ssl,
	), nil
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

func urlEscape(s string) string {
	replacer := strings.NewReplacer(
		":", "%3A",
		"/", "%2F",
		"@", "%40",
		"?", "%3F",
		"#", "%23",
		"[", "%5B",
		"]", "%5D",
	)
	return replacer.Replace(s)
}
