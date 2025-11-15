package discord

import (
	"discord-event-bot/config"
	eventsvc "discord-event-bot/internal/services/events"
	"discord-event-bot/internal/storage"
	"fmt"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

// CreateDiscordScheduledEvent crea un evento oficial de Discord
func CreateDiscordScheduledEvent(s *discordgo.Session, event *storage.Event) {
	// Calcular hora de fin (2 horas después del inicio)
	endTime := event.DateTime.Add(2 * time.Hour)

	params := &discordgo.GuildScheduledEventParams{
		Name:               event.Name,
		Description:        event.Description,
		ScheduledStartTime: &event.DateTime,
		ScheduledEndTime:   &endTime,
		PrivacyLevel:       discordgo.GuildScheduledEventPrivacyLevelGuildOnly,
		EntityType:         discordgo.GuildScheduledEventEntityTypeExternal,
		EntityMetadata: &discordgo.GuildScheduledEventEntityMetadata{
			Location: "In-Game",
		},
	}
	discordEvent, err := s.GuildScheduledEventCreate(config.AppConfig.GuildID, params)
	if err != nil {
		log.Printf("Error creando evento de Discord: %v", err)
		return
	}

	event.DiscordEventID = discordEvent.ID
	storage.Store.SaveEvent(event)
}

// handleCreateEvent crea un nuevo evento
func handleCreateEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	input, err := buildCreateEventInputFromInteraction(i)
	if err != nil {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Formato de fecha inválido. Usa: YYYY-MM-DD HH:MM (ej: 2024-12-25 20:00)",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	event, err := eventsvc.CreateEvent(input)
	if err != nil {
		log.Printf("Error creando evento: %v", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Error creando el evento: " + err.Error(),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Publicar mensaje del evento (inmediato o programado)
	if event.AnnouncementTime.IsZero() || !event.AnnouncementTime.After(time.Now()) {
		if err := PublishEventMessage(s, event); err != nil {
			log.Printf("Error publicando mensaje: %v", err)
		}
	}

	// Crear evento oficial de Discord solo si está habilitado globalmente y el evento lo requiere
	if config.AppConfig.EnableDiscordEvents && event.CreateDiscordEvent {
		CreateDiscordScheduledEvent(s, event)
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Evento creado exitosamente! ID: `%s`", event.ID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

func buildCreateEventInputFromInteraction(i *discordgo.InteractionCreate) (eventsvc.CreateEventInput, error) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	nombre := optionMap["nombre"].StringValue()
	tipo := optionMap["tipo"].StringValue()
	fechaStr := optionMap["fecha"].StringValue()
	descripcion := optionMap["descripcion"].StringValue()

	announceHours := 0
	if ahOpt, ok := optionMap["announce_hours"]; ok {
		announceHours = int(ahOpt.IntValue())
		if announceHours < 0 {
			announceHours = 0
		}
	}

	repeatEveryDays := 0
	if repeatOpt, ok := optionMap["repeat_days"]; ok {
		repeatEveryDays = int(repeatOpt.IntValue())
		if repeatEveryDays < 0 {
			repeatEveryDays = 0
		}
	}

	createDiscordEvent := false
	if deOpt, ok := optionMap["discord_event"]; ok {
		createDiscordEvent = deOpt.BoolValue()
	}

	reminderOffsetMinutes := 0
	if rhOpt, ok := optionMap["reminder_minutes"]; ok {
		v := int(rhOpt.IntValue())
		if v > 0 {
			reminderOffsetMinutes = v
		}
	}

	deleteAfterHours := 0
	if dahOpt, ok := optionMap["delete_after_hours"]; ok {
		v := int(dahOpt.IntValue())
		if v > 0 {
			deleteAfterHours = v
		}
	}

	// Template opcional
	templateName := ""
	if tmpl, ok := optionMap["template"]; ok {
		templateName = tmpl.StringValue()
	}

	// Canal por defecto es el canal actual
	channelID := i.ChannelID
	if canal, ok := optionMap["canal"]; ok {
		channelID = canal.StringValue()
	}

	// Parsear fecha
	loc, _ := time.LoadLocation(config.AppConfig.Timezone)
	fecha, err := time.ParseInLocation("2006-01-02 15:04", fechaStr, loc)
	if err != nil {
		return eventsvc.CreateEventInput{}, err
	}

	return eventsvc.CreateEventInput{
		Name:                  nombre,
		Type:                  tipo,
		Description:           descripcion,
		DateTime:              fecha,
		ChannelID:             channelID,
		RepeatEveryDays:       repeatEveryDays,
		TemplateName:          templateName,
		CreateDiscordEvent:    createDiscordEvent,
		CreatedBy:             i.Member.User.ID,
		AnnounceHours:         announceHours,
		ReminderOffsetMinutes: reminderOffsetMinutes,
		DeleteAfterHours:      deleteAfterHours,
	}, nil
}

// handleDeleteEvent elimina un evento
func handleDeleteEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	eventID := options[0].StringValue()

	event, err := storage.Store.GetEvent(eventID)
	if err != nil {
		respondError(s, i, "Evento no encontrado")
		return
	}

	// Eliminar mensaje
	if event.MessageID != "" {
		s.ChannelMessageDelete(event.Channel, event.MessageID)
	}

	// Cerrar hilo asociado si existe
	if event.ThreadID != "" {
		archived := true
		locked := true
		if _, err := s.ChannelEdit(event.ThreadID, &discordgo.ChannelEdit{Archived: &archived, Locked: &locked}); err != nil {
			log.Printf("Error archivando hilo %s para evento %s: %v", event.ThreadID, event.ID, err)
		}
	}

	// Eliminar evento
	storage.Store.DeleteEvent(eventID)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Evento `%s` eliminado", eventID),
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}

// handleListEvents lista todos los eventos activos
func handleListEvents(s *discordgo.Session, i *discordgo.InteractionCreate) {
	events := storage.Store.GetActiveEvents()

	if len(events) == 0 {
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No hay eventos activos",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	embed := &discordgo.MessageEmbed{
		Title: "📋 Eventos Activos",
		Color: 0x5865F2,
	}

	for _, event := range events {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   event.Name,
			Value:  fmt.Sprintf("ID: `%s`\nTipo: %s\nFecha: <t:%d:F>", event.ID, event.Type, event.DateTime.Unix()),
			Inline: false,
		})
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
}
