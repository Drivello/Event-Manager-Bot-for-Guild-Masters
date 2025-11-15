package discord

import (
	signupsvc "discord-event-bot/internal/services/signups"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// handleSignup maneja las inscripciones
func handleSignup(s *discordgo.Session, i *discordgo.InteractionCreate, eventID, role, class string) {
	userID := i.Member.User.ID
	username := i.Member.User.Username

	event, err := signupsvc.SignupToEvent(signupsvc.SignupInput{
		EventID:  eventID,
		UserID:   userID,
		Username: username,
		Role:     role,
		Class:    class,
	})
	if err != nil {
		respondError(s, i, err.Error())
		return
	}

	// Actualizar mensaje
	UpdateEventMessage(s, event)

	label := role
	if class != "" {
		label = fmt.Sprintf("%s - %s", role, class)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Te has inscrito como **%s**. Tu inscripción está confirmada.", label),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// handleCancelSignup maneja la cancelación de inscripción
func handleCancelSignup(s *discordgo.Session, i *discordgo.InteractionCreate, eventID string) {
	event, err := signupsvc.CancelSignup(signupsvc.CancelInput{
		EventID: eventID,
		UserID:  i.Member.User.ID,
	})
	if err != nil {
		respondError(s, i, err.Error())
		return
	}

	UpdateEventMessage(s, event)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "✅ Tu inscripción ha sido cancelada",
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
