package web

import (
	"discord-event-bot/config"
	"discord-event-bot/internal/discord"
	"discord-event-bot/internal/storage"
	"encoding/json"
	"html/template"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
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
	authorized.GET("/config", handleConfigPage)

	// Rutas de templates
	RegisterTemplateRoutes(authorized)

	log.Printf("✅ Servidor web iniciado en http://localhost:%s", config.AppConfig.Port)
}

// StartWebServer inicia el servidor web
func StartWebServer() error {
	return router.Run(":" + config.AppConfig.Port)
}

// handleIndex muestra la página principal
func handleIndex(c *gin.Context) {
	events := storage.Store.GetActiveEvents()
	c.HTML(http.StatusOK, "index.html", gin.H{
		"title":  "Panel de Administración - Discord Event Bot",
		"events": events,
	})
}

// handleEventsList muestra la lista de eventos
func handleEventsList(c *gin.Context) {
	events := storage.Store.GetAllEvents()
	c.HTML(http.StatusOK, "events.html", gin.H{
		"title":  "Todos los Eventos",
		"events": events,
	})
}

// handleCreateEventPage muestra el formulario de creación
func handleCreateEventPage(c *gin.Context) {
	templates := storage.Templates.GetAllTemplates()
	c.HTML(http.StatusOK, "create_event.html", gin.H{
		"title":     "Crear Nuevo Evento",
		"roles":     config.AppConfig.DefaultRoles,
		"templates": templates,
	})
}

// handleCreateEventPost procesa la creación de un evento
func handleCreateEventPost(c *gin.Context) {
	nombre := c.PostForm("nombre")
	tipo := c.PostForm("tipo")
	fechaStr := c.PostForm("fecha")
	descripcion := c.PostForm("descripcion")
	channel := c.PostForm("channel")
	templateName := c.PostForm("template")
	repeatDaysStr := c.PostForm("repeat_days")
	createDiscordEvent := c.PostForm("discord_event") == "1"
	repeatEveryDays := 0
	if repeatDaysStr != "" {
		if v, err := strconv.Atoi(repeatDaysStr); err == nil && v > 0 {
			repeatEveryDays = v
		}
	}

	templates := storage.Templates.GetAllTemplates()

	// Parsear fecha
	loc, _ := time.LoadLocation(config.AppConfig.Timezone)
	fecha, err := time.ParseInLocation("2006-01-02T15:04", fechaStr, loc)
	if err != nil {
		c.HTML(http.StatusBadRequest, "create_event.html", gin.H{
			"title":     "Crear Nuevo Evento",
			"error":     "Formato de fecha inválido",
			"roles":     config.AppConfig.DefaultRoles,
			"templates": templates,
		})
		return
	}

	// Crear evento base
	event := &storage.Event{
		ID:                 uuid.New().String(),
		Name:               nombre,
		Type:               tipo,
		Description:        descripcion,
		DateTime:           fecha,
		Channel:            channel,
		Status:             "active",
		CreatedAt:          time.Now(),
		CreatedBy:          "admin_web",
		AllowMultiSignup:   false,
		Signups:            make(map[string][]storage.Signup),
		RepeatEveryDays:    repeatEveryDays,
		CreateDiscordEvent: createDiscordEvent,
	}

	// Si se especificó un template, usarlo
	if templateName != "" {
		event, err = storage.Store.CreateEventFromTemplate(templateName, event)
		if err != nil {
			c.HTML(http.StatusBadRequest, "create_event.html", gin.H{
				"title":     "Crear Nuevo Evento",
				"error":     "Error usando template: " + err.Error(),
				"roles":     config.AppConfig.DefaultRoles,
				"templates": templates,
			})
			return
		}
	} else {
		// Agregar roles por defecto
		for _, role := range config.AppConfig.DefaultRoles {
			event.Roles = append(event.Roles, storage.RoleSignup{
				Name:  role.Name,
				Emoji: role.Emoji,
				Limit: role.Limit,
			})
		}

		// Guardar evento
		if err := storage.Store.SaveEvent(event); err != nil {
			c.HTML(http.StatusInternalServerError, "create_event.html", gin.H{
				"title":     "Crear Nuevo Evento",
				"error":     "Error guardando evento: " + err.Error(),
				"roles":     config.AppConfig.DefaultRoles,
				"templates": templates,
			})
			return
		}
	}

	// Publicar en Discord
	if discord.Session != nil {
		if err := discord.PublishEventMessage(discord.Session, event); err != nil {
			log.Printf("Error publicando en Discord: %v", err)
		}

		if config.AppConfig.EnableDiscordEvents && event.CreateDiscordEvent {
			discord.CreateDiscordScheduledEvent(discord.Session, event)
		}
	}

	c.Redirect(http.StatusSeeOther, "/events/"+event.ID)
}

// handleEventDetail muestra los detalles de un evento
func handleEventDetail(c *gin.Context) {
	eventID := c.Param("id")
	event, err := storage.Store.GetEvent(eventID)

	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Evento no encontrado",
		})
		return
	}

	c.HTML(http.StatusOK, "event_detail.html", gin.H{
		"title": event.Name,
		"event": event,
	})
}

// handleCancelEvent cancela un evento
func handleCancelEvent(c *gin.Context) {
	eventID := c.Param("id")
	event, err := storage.Store.GetEvent(eventID)

	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Evento no encontrado"})
		return
	}

	event.Status = "cancelled"
	storage.Store.SaveEvent(event)

	// Eliminar mensaje de Discord si existe
	if discord.Session != nil && event.MessageID != "" {
		discord.Session.ChannelMessageDelete(event.Channel, event.MessageID)
	}

	c.Redirect(http.StatusSeeOther, "/")
}

// handleConfirmSignup confirma una inscripción
func handleConfirmSignup(c *gin.Context) {
	eventID := c.Param("id")
	userID := c.Param("userid")
	role := c.Param("role")

	if err := storage.Store.ConfirmSignup(eventID, userID, role, "admin_web"); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	// Actualizar mensaje en Discord
	event, _ := storage.Store.GetEvent(eventID)
	if discord.Session != nil && event != nil {
		discord.UpdateEventMessage(discord.Session, event)
	}

	c.Redirect(http.StatusSeeOther, "/events/"+eventID)
}

// handleConfigPage muestra la página de configuración
func handleConfigPage(c *gin.Context) {
	c.HTML(http.StatusOK, "config.html", gin.H{
		"title":  "Configuración",
		"config": config.AppConfig,
	})
}
