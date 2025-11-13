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
		// Enviar recordatorio 15 minutos antes
		reminderTime := event.DateTime.Add(-15 * time.Minute)

		if now.After(reminderTime) && !event.ReminderSent {
			sendReminder(Session, event)
			event.ReminderSent = true
			storage.Store.SaveEvent(event)
		}

		// Marcar como completado si ya pasó
		if now.After(event.DateTime.Add(2 * time.Hour)) {
			event.Status = "completed"
			storage.Store.SaveEvent(event)
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

	content := fmt.Sprintf("🔔 **Recordatorio**: El evento **%s** comienza <t:%d:R>\n\n%s",
		event.Name,
		event.DateTime.Unix(),
		strings.Join(mentions, " "))

	s.ChannelMessageSend(event.Channel, content)
}
