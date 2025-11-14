package discord

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

var (
	Session  *discordgo.Session
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "create_event",
			Description: "Crear un nuevo evento para el guild",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "nombre",
					Description: "Nombre del evento",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "tipo",
					Description: "Tipo de evento (Raid, Dungeon, PvP, Social, etc.)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "fecha",
					Description: "Fecha y hora (formato: YYYY-MM-DD HH:MM)",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "descripcion",
					Description: "Descripción del evento",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "template",
					Description: "Template a usar (opcional)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionChannel,
					Name:        "canal",
					Description: "Canal donde se publicará el evento",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionBoolean,
					Name:        "discord_event",
					Description: "Crear también el evento oficial de Discord (Guild Scheduled Event)",
					Required:    false,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "repeat_days",
					Description: "Cada cuántos días se repite el evento (0 o vacío = no se repite)",
					Required:    false,
				},
			},
		},
		{
			Name:        "delete_event",
			Description: "Eliminar un evento existente",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID del evento a eliminar",
					Required:    true,
				},
			},
		},
		{
			Name:        "remind_event",
			Description: "Enviar recordatorio inmediato de un evento",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "id",
					Description: "ID del evento",
					Required:    true,
				},
			},
		},
		{
			Name:        "config",
			Description: "Mostrar la configuración actual del bot",
		},
		{
			Name:        "list_events",
			Description: "Listar todos los eventos activos",
		},
	}
)

// handleInteractionCreate maneja las interacciones de comandos slash
func handleInteractionCreate(s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		handleSlashCommand(s, i)
	case discordgo.InteractionMessageComponent:
		handleButtonClick(s, i)
	}
}

// handleSlashCommand procesa los comandos slash
func handleSlashCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	commandName := i.ApplicationCommandData().Name

	switch commandName {
	case "create_event":
		handleCreateEvent(s, i)
	case "delete_event":
		handleDeleteEvent(s, i)
	case "remind_event":
		handleRemindEvent(s, i)
	case "config":
		handleConfig(s, i)
	case "list_events":
		handleListEvents(s, i)
	}
}

// handleButtonClick maneja los clicks en botones
func handleButtonClick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID

	if strings.HasPrefix(customID, "signup_") {
		payload := strings.TrimPrefix(customID, "signup_")
		underscoreIdx := strings.Index(payload, "_")
		if underscoreIdx == -1 {
			return
		}

		eventID := payload[:underscoreIdx]
		rest := payload[underscoreIdx+1:]

		role := rest
		class := ""
		if sep := strings.Index(rest, "__"); sep != -1 {
			role = rest[:sep]
			class = rest[sep+2:]
		}

		handleSignup(s, i, eventID, role, class)
	} else if strings.HasPrefix(customID, "cancel_") {
		eventID := strings.TrimPrefix(customID, "cancel_")
		handleCancelSignup(s, i, eventID)
	}
}
