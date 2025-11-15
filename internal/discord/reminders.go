package discord

import (
	remindersvc "discord-event-bot/internal/services/reminders"
	"discord-event-bot/internal/storage"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

// StartReminderService inicia el servicio de recordatorios automáticos
func StartReminderService() {
	ticker := time.NewTicker(1 * time.Minute)
	go func() {
		for range ticker.C {
			checkAndSendReminders()
			checkAndPublishScheduledEvents()
		}
	}()
	log.Println("✅ Servicio de recordatorios iniciado")
}

// checkAndSendReminders verifica y envía recordatorios
func checkAndSendReminders() {
	result := remindersvc.ProcessReminders(time.Now())

	deliverReminders(result)
	updateEventMessages(result)
	cleanupEventMessages(result)
}

func deliverReminders(result remindersvc.ProcessResult) {
	// Enviar recordatorios a través de Discord
	for _, event := range result.EventsToRemind {
		sendReminder(Session, event)
	}
}

func updateEventMessages(result remindersvc.ProcessResult) {
	// Actualizar mensajes en Discord para eventos recurrentes que cambiaron
	for _, event := range result.EventsToUpdate {
		if Session != nil && event.MessageID != "" {
			UpdateEventMessage(Session, event)
		}
	}
}

func cleanupEventMessages(result remindersvc.ProcessResult) {
	// Borrar mensajes/hilos cuando corresponde
	for _, event := range result.EventsToDeleteMessages {
		if Session == nil {
			continue
		}
		// Eliminar mensaje principal
		if event.MessageID != "" {
			if err := Session.ChannelMessageDelete(event.Channel, event.MessageID); err != nil {
				log.Printf("Error borrando mensaje del evento %s: %v", event.ID, err)
			}
		}
		// Cerrar hilo asociado si existe
		if event.ThreadID != "" {
			archived := true
			locked := true
			if _, err := Session.ChannelEdit(event.ThreadID, &discordgo.ChannelEdit{Archived: &archived, Locked: &locked}); err != nil {
				log.Printf("Error archivando hilo %s para evento %s: %v", event.ThreadID, event.ID, err)
			}
		}

		// Limpiar referencias para permitir republicación futura (especialmente en recurrentes)
		event.MessageID = ""
		event.ThreadID = ""
		if err := storage.Store.SaveEvent(event); err != nil {
			log.Printf("Error guardando evento %s tras borrar mensaje/hilo: %v", event.ID, err)
		}
	}
}

func checkAndPublishScheduledEvents() {
	if Session == nil {
		return
	}

	events := storage.Store.GetActiveEvents()
	now := time.Now()

	for _, event := range events {
		if !shouldPublishScheduledEvent(event, now) {
			continue
		}

		if err := PublishEventMessage(Session, event); err != nil {
			log.Printf("Error publicando mensaje programado para evento %s: %v", event.ID, err)
		}
	}
}

func shouldPublishScheduledEvent(event *storage.Event, now time.Time) bool {
	if event.MessageID != "" {
		return false
	}
	if event.AnnouncementTime.IsZero() {
		return false
	}
	if now.Before(event.AnnouncementTime) {
		return false
	}
	// Si el tiempo de borrado automático ya pasó, no volver a publicar el mensaje
	if event.DeleteAfterHours > 0 {
		deleteTime := event.DateTime.Add(time.Duration(event.DeleteAfterHours) * time.Hour)
		if now.After(deleteTime) {
			return false
		}
	}
	return true
}

// handleRemindEvent envía recordatorio de un evento
func handleRemindEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	eventID := options[0].StringValue()

	event, err := storage.Store.GetEvent(eventID)
	if err != nil {
		respondError(s, i, "Evento no encontrado")
		return
	}

	sendReminder(s, event)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ Recordatorio enviado",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// sendReminder envía un recordatorio del evento
func sendReminder(s *discordgo.Session, event *storage.Event) {
	var mentions []string
	for _, signups := range event.Signups {
		for _, signup := range signups {
			if signup.Status == "confirmed" {
				mentions = append(mentions, fmt.Sprintf("<@%s>", signup.UserID))
			}
		}
	}

	totalConfirmed := 0
	for _, signups := range event.Signups {
		for _, signup := range signups {
			if signup.Status == "confirmed" {
				totalConfirmed++
			}
		}
	}

	maxParticipants := event.MaxParticipants
	if maxParticipants == 0 {
		for _, role := range event.Roles {
			maxParticipants += role.Limit
		}
	}

	prefix := ""
	if maxParticipants == 0 || totalConfirmed < maxParticipants {
		prefix = "@here "
	}

	content := fmt.Sprintf("%s🔔 **Recordatorio**: El evento **%s** comienza <t:%d:R>\n\n%s",
		prefix,
		event.Name,
		event.DateTime.Unix(),
		strings.Join(mentions, " "))

	// Enviar al hilo del evento si existe, con fallback al canal principal
	targetChannelID := event.Channel
	if event.ThreadID != "" {
		if _, err := s.ChannelMessageSend(event.ThreadID, content); err != nil {
			log.Printf("Error enviando recordatorio al hilo %s para evento %s: %v", event.ThreadID, event.ID, err)
		} else {
			return
		}
	}

	s.ChannelMessageSend(targetChannelID, content)
}
