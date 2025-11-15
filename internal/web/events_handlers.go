package web

import (
	"discord-event-bot/config"
	"discord-event-bot/internal/discord"
	eventsvc "discord-event-bot/internal/services/events"
	"discord-event-bot/internal/storage"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/gin-gonic/gin"
)

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
	announceHoursStr := c.PostForm("announce_hours")
	reminderMinutesStr := c.PostForm("reminder_minutes")
	deleteAfterHoursStr := c.PostForm("delete_after_hours")
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

	announceHours := 0
	if announceHoursStr != "" {
		if v, err := strconv.Atoi(announceHoursStr); err == nil && v > 0 {
			announceHours = v
		}
	}

	reminderOffsetMinutes := 0
	if reminderMinutesStr != "" {
		if v, err := strconv.Atoi(reminderMinutesStr); err == nil && v > 0 {
			reminderOffsetMinutes = v
		}
	}

	deleteAfterHours := 0
	if deleteAfterHoursStr != "" {
		if v, err := strconv.Atoi(deleteAfterHoursStr); err == nil && v > 0 {
			deleteAfterHours = v
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

	event, err := eventsvc.CreateEvent(eventsvc.CreateEventInput{
		Name:                  nombre,
		Type:                  tipo,
		Description:           descripcion,
		DateTime:              fecha,
		ChannelID:             channel,
		RepeatEveryDays:       repeatEveryDays,
		TemplateName:          templateName,
		CreateDiscordEvent:    createDiscordEvent,
		CreatedBy:             "admin_web",
		AnnounceHours:         announceHours,
		ReminderOffsetMinutes: reminderOffsetMinutes,
		DeleteAfterHours:      deleteAfterHours,
	})
	if err != nil {
		c.HTML(http.StatusBadRequest, "create_event.html", gin.H{
			"title":     "Crear Nuevo Evento",
			"error":     "Error creando evento: " + err.Error(),
			"roles":     config.AppConfig.DefaultRoles,
			"templates": templates,
		})
		return
	}

	// Publicar en Discord
	if discord.Session != nil {
		if event.AnnouncementTime.IsZero() || !event.AnnouncementTime.After(time.Now()) {
			if err := discord.PublishEventMessage(discord.Session, event); err != nil {
				log.Printf("Error publicando en Discord: %v", err)
			}
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

	// Eliminar mensaje de Discord y cerrar hilo si existen
	if discord.Session != nil {
		if event.MessageID != "" {
			discord.Session.ChannelMessageDelete(event.Channel, event.MessageID)
		}
		if event.ThreadID != "" {
			archived := true
			locked := true
			if _, err := discord.Session.ChannelEdit(event.ThreadID, &discordgo.ChannelEdit{Archived: &archived, Locked: &locked}); err != nil {
				log.Printf("Error archivando hilo %s para evento %s: %v", event.ThreadID, event.ID, err)
			}
		}
	}

	c.Redirect(http.StatusSeeOther, "/")
}

func handleCleanupCancelledEvents(c *gin.Context) {
	deleted, err := storage.Store.DeleteCancelledEvents()
	if err != nil {
		log.Printf("Error eliminando eventos cancelados: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error eliminando eventos cancelados"})
		return
	}

	log.Printf("Eliminados %d eventos cancelados", deleted)
	c.Redirect(http.StatusSeeOther, "/events")
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
