package web

import (
	"discord-event-bot/config"
	"encoding/json"
	"html/template"
	"log"

	"github.com/gin-gonic/gin"
)

var router *gin.Engine

// InitWebServer inicializa el servidor web
func InitWebServer() {
	gin.SetMode(gin.ReleaseMode)
	router = gin.Default()

	// Registrar funciones de template
	router.SetFuncMap(template.FuncMap{
		"json": func(v any) template.JS {
			b, err := json.Marshal(v)
			if err != nil {
				// En caso de error devolvemos un JSON válido
				return template.JS("[]")
			}
			// Lo marcamos como JS para que no escape comillas, etc.
			return template.JS(b)
		},
	})

	// Cargar templates HTML
	router.LoadHTMLGlob("internal/web/templates/*")

	// Middleware de autenticación básica
	authorized := router.Group("/", gin.BasicAuth(gin.Accounts{
		config.AppConfig.AdminUser: config.AppConfig.AdminPass,
	}))

	// Rutas de eventos
	authorized.GET("/", handleIndex)
	authorized.GET("/events", handleEventsList)
	authorized.GET("/events/create", handleCreateEventPage)
	authorized.POST("/events/create", handleCreateEventPost)
	authorized.GET("/events/:id", handleEventDetail)
	authorized.POST("/events/:id/cancel", handleCancelEvent)
	authorized.POST("/events/:id/confirm/:userid/:role", handleConfirmSignup)
	authorized.POST("/events/cleanup-cancelled", handleCleanupCancelledEvents)
	authorized.GET("/config", handleConfigPage)

	// Rutas de templates
	RegisterTemplateRoutes(authorized)

	log.Printf("✅ Servidor web iniciado en http://localhost:%s", config.AppConfig.Port)
}

// StartWebServer inicia el servidor web
func StartWebServer() error {
	return router.Run(":" + config.AppConfig.Port)
}
