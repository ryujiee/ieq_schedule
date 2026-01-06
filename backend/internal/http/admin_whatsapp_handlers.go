package http

import (
	"net/http"

	"ieq/backend/internal/whatsapp"

	"github.com/gin-gonic/gin"
)

func adminWhatsAppStatusHandler(mgr *whatsapp.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.JSON(http.StatusOK, mgr.Status())
	}
}

func adminWhatsAppConnectHandler(mgr *whatsapp.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := mgr.StartPairing(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, mgr.Status())
	}
}

func adminWhatsAppQRHandler(mgr *whatsapp.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		st := mgr.Status()
		c.JSON(http.StatusOK, gin.H{
			"pairing":   st.Pairing,
			"qr":        st.LastQR,
			"loggedIn":  st.LoggedIn,
			"connected": st.Connected,
		})
	}
}

func adminWhatsAppLogoutHandler(mgr *whatsapp.Manager) gin.HandlerFunc {
	return func(c *gin.Context) {
		if err := mgr.Logout(); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"ok": true})
	}
}
