package discord

import (
	"discord-event-bot/internal/storage"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

// handleSignup maneja las inscripciones
func handleSignup(s *discordgo.Session, i *discordgo.InteractionCreate, eventID, role, class string) {
	event, err := storage.Store.GetEvent(eventID)
	if err != nil {
		respondError(s, i, "Evento no encontrado")
		return
	}

	userID := i.Member.User.ID
	username := i.Member.User.Username

	// Verificar si ya está inscrito
	for r, signups := range event.Signups {
		for _, signup := range signups {
			if signup.UserID == userID {
				if r == role {
					respondError(s, i, "Ya estás inscrito en este rol")
					return
				}
				if !event.AllowMultiSignup {
					respondError(s, i, "Ya estás inscrito en otro rol. Cancela primero tu inscripción actual.")
					return
				}
			}
		}
	}

	// Verificar límite de rol
	confirmedCount := 0
	for _, signup := range event.Signups[role] {
		if signup.Status == "confirmed" {
			confirmedCount++
		}
	}

	var roleLimit int
	for _, r := range event.Roles {
		if r.Name == role {
			roleLimit = r.Limit
			break
		}
	}

	if roleLimit > 0 && confirmedCount >= roleLimit {
		respondError(s, i, fmt.Sprintf("El rol %s ya está lleno", role))
		return
	}

	// Agregar inscripción (con clase si aplica)
	var signupErr error
	if class != "" {
		signupErr = storage.Store.AddSignupWithClass(eventID, userID, username, role, class)
	} else {
		signupErr = storage.Store.AddSignup(eventID, userID, username, role)
	}
	if signupErr != nil {
		respondError(s, i, "Error procesando inscripción")
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
	event, err := storage.Store.GetEvent(eventID)
	if err != nil {
		respondError(s, i, "Evento no encontrado")
		return
	}

	userID := i.Member.User.ID
	removed := false

	for role := range event.Signups {
		if err := storage.Store.RemoveSignup(eventID, userID, role); err == nil {
			removed = true
		}
	}

	if !removed {
		respondError(s, i, "No estás inscrito en este evento")
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
