package web

import (
	"discord-event-bot/config"

	"github.com/gin-gonic/gin"
)

// handleConfigPage muestra la página de configuración
func handleConfigPage(c *gin.Context) {
	c.HTML(200, "config.html", gin.H{
		"title":  "Configuración",
		"config": config.AppConfig,
	})
}
