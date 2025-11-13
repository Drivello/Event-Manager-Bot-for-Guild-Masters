package discord

import "github.com/bwmarrin/discordgo"

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
