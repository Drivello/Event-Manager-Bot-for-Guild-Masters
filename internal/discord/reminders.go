package discord

import (
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
		}
	}()
	log.Println("✅ Servicio de recordatorios iniciado")
}

// checkAndSendReminders verifica y envía recordatorios
func checkAndSendReminders() {
	events := storage.Store.GetActiveEvents()
	now := time.Now()

	for _, event := range events {
		if event.RepeatEveryDays > 0 {
			handleRecurringEventReminder(event, now)
			continue
		}

		reminderTime := event.DateTime.Add(-15 * time.Minute)

		if now.After(reminderTime) && !event.ReminderSent {
			sendReminder(Session, event)
			event.ReminderSent = true
			storage.Store.SaveEvent(event)
		}

		if now.After(event.DateTime.Add(2 * time.Hour)) {
			event.Status = "completed"
			storage.Store.SaveEvent(event)
		}
	}
}

func handleRecurringEventReminder(event *storage.Event, now time.Time) {
	changed := false

	for now.After(event.DateTime.Add(2 * time.Hour)) {
		event.DateTime = event.DateTime.Add(time.Duration(event.RepeatEveryDays) * 24 * time.Hour)
		event.ReminderSent = false
		changed = true
	}

	if !event.ReminderSent {
		windowStart := event.DateTime.Add(-1 * time.Minute)
		windowEnd := event.DateTime.Add(1 * time.Minute)

		if now.After(windowStart) && now.Before(windowEnd) {
			sendReminder(Session, event)
			event.ReminderSent = true
			changed = true
		}
	}

	if changed {
		storage.Store.SaveEvent(event)
		if Session != nil && event.MessageID != "" {
			UpdateEventMessage(Session, event)
		}
	}
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
