package discord

import (
	"discord-event-bot/config"
	"discord-event-bot/internal/storage"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
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

// handleCreateEvent crea un nuevo evento
func handleCreateEvent(s *discordgo.Session, i *discordgo.InteractionCreate) {
	options := i.ApplicationCommandData().Options
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption, len(options))
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	nombre := optionMap["nombre"].StringValue()
	tipo := optionMap["tipo"].StringValue()
	fechaStr := optionMap["fecha"].StringValue()
	descripcion := optionMap["descripcion"].StringValue()

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
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "❌ Formato de fecha inválido. Usa: YYYY-MM-DD HH:MM (ej: 2024-12-25 20:00)",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Crear evento base
	event := &storage.Event{
		ID:               uuid.New().String(),
		Name:             nombre,
		Type:             tipo,
		Description:      descripcion,
		DateTime:         fecha,
		Channel:          channelID,
		Status:           "active",
		CreatedAt:        time.Now(),
		CreatedBy:        i.Member.User.ID,
		AllowMultiSignup: false,
		Signups:          make(map[string][]storage.Signup),
	}

	// Si se especificó un template, usarlo
	if templateName != "" {
		event, err = storage.Store.CreateEventFromTemplate(templateName, event)
		if err != nil {
			log.Printf("Error usando template: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ Template no encontrado: " + templateName,
					Flags:   discordgo.MessageFlagsEphemeral,
				},
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
			log.Printf("Error guardando evento: %v", err)
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "❌ Error guardando el evento",
					Flags:   discordgo.MessageFlagsEphemeral,
				},
			})
			return
		}
	}

	// Publicar mensaje del evento
	if err := PublishEventMessage(s, event); err != nil {
		log.Printf("Error publicando mensaje: %v", err)
	}

	// Crear evento oficial de Discord si está habilitado
	if config.AppConfig.EnableDiscordEvents {
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

// PublishEventMessage publica el mensaje del evento con botones de inscripción
func PublishEventMessage(s *discordgo.Session, event *storage.Event) error {
	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("📅 %s", event.Name),
		Description: event.Description,
		Color:       0x5865F2,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Tipo",
				Value:  event.Type,
				Inline: true,
			},
			{
				Name:   "Fecha y Hora",
				Value:  event.DateTime.Format("02/01/2006 15:04"),
				Inline: true,
			},
			{
				Name:   "ID del Evento",
				Value:  fmt.Sprintf("`%s`", event.ID),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Seleccioná tu rol para inscribirte",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Campo de inscripciones
	signupsText := buildSignupsText(event)
	if strings.TrimSpace(signupsText) == "" {
		signupsText = "Todavía no hay inscripciones."
	}
	embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
		Name:   "Inscripciones",
		Value:  signupsText,
		Inline: false,
	})

	// Crear botones para cada rol (máx 5 por fila)
	var components []discordgo.MessageComponent
	var currentRow discordgo.ActionsRow

	for i, role := range event.Roles {
		// Emoji va en el label, no en el campo Emoji del botón
		label := role.Name
		if role.Emoji != "" {
			label = fmt.Sprintf("%s %s", role.Emoji, role.Name)
		}

		button := discordgo.Button{
			Label:    label,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("signup_%s_%s", event.ID, role.Name),
		}

		currentRow.Components = append(currentRow.Components, button)

		// Si llegamos a 5 botones, cerramos la fila
		if len(currentRow.Components) == 5 {
			components = append(components, currentRow)
			currentRow = discordgo.ActionsRow{}
		}

		// Último rol: si quedó algo en la fila, la agregamos
		if i == len(event.Roles)-1 && len(currentRow.Components) > 0 {
			components = append(components, currentRow)
		}
	}

	// Botón para cancelar inscripción en una fila separada
	cancelRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Cancelar inscripción",
				Style:    discordgo.DangerButton,
				CustomID: fmt.Sprintf("cancel_%s", event.ID),
			},
		},
	}
	components = append(components, cancelRow)

	msg, err := s.ChannelMessageSendComplex(event.Channel, &discordgo.MessageSend{
		Embeds:     []*discordgo.MessageEmbed{embed},
		Components: components,
	})
	if err != nil {
		return fmt.Errorf("error enviando mensaje a Discord: %w", err)
	}

	event.MessageID = msg.ID
	if err := storage.Store.SaveEvent(event); err != nil {
		return fmt.Errorf("error guardando evento: %w", err)
	}

	return nil
}

// buildSignupsText construye el texto de inscripciones
func buildSignupsText(event *storage.Event) string {
	var builder strings.Builder

	for _, role := range event.Roles {
		signups := event.Signups[role.Name]
		confirmedCount := 0
		pendingCount := 0

		for _, signup := range signups {
			if signup.Status == "confirmed" {
				confirmedCount++
			} else {
				pendingCount++
			}
		}

		builder.WriteString(fmt.Sprintf("%s **%s**: %d/%d",
			role.Emoji, role.Name, confirmedCount, role.Limit))

		if pendingCount > 0 {
			builder.WriteString(fmt.Sprintf(" (%d pendiente)", pendingCount))
		}
		builder.WriteString("\n")

		// Si hay clases definidas, mostrar desglose por clase
		if len(role.Classes) > 0 {
			classCounts := make(map[string]int)
			for _, signup := range signups {
				if signup.Status == "confirmed" && signup.Class != "" {
					classCounts[signup.Class]++
				}
			}

			for _, class := range role.Classes {
				count := classCounts[class.Name]
				if count > 0 {
					builder.WriteString(fmt.Sprintf("  %s %s: %d\n", class.Emoji, class.Name, count))
				}
			}
		}
	}

	if builder.Len() == 0 {
		return "Sin inscripciones aún"
	}

	return builder.String()
}

// handleButtonClick maneja los clicks en botones
func handleButtonClick(s *discordgo.Session, i *discordgo.InteractionCreate) {
	customID := i.MessageComponentData().CustomID

	if strings.HasPrefix(customID, "signup_") {
		parts := strings.Split(customID, "_")
		if len(parts) >= 3 {
			eventID := parts[1]
			role := strings.Join(parts[2:], "_")
			handleSignup(s, i, eventID, role)
		}
	} else if strings.HasPrefix(customID, "cancel_") {
		eventID := strings.TrimPrefix(customID, "cancel_")
		handleCancelSignup(s, i, eventID)
	}
}

// handleSignup maneja las inscripciones
func handleSignup(s *discordgo.Session, i *discordgo.InteractionCreate, eventID, role string) {
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

	if confirmedCount >= roleLimit {
		respondError(s, i, fmt.Sprintf("El rol %s ya está lleno", role))
		return
	}

	// Agregar inscripción
	if err := storage.Store.AddSignup(eventID, userID, username, role); err != nil {
		respondError(s, i, "Error procesando inscripción")
		return
	}

	// Actualizar mensaje
	UpdateEventMessage(s, event)

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("✅ Te has inscrito como **%s**. Pendiente de confirmación.", role),
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

// UpdateEventMessage actualiza el mensaje del evento
func UpdateEventMessage(s *discordgo.Session, event *storage.Event) {
	// Recargar evento para obtener datos actualizados
	event, _ = storage.Store.GetEvent(event.ID)

	embed := &discordgo.MessageEmbed{
		Title:       fmt.Sprintf("📅 %s", event.Name),
		Description: event.Description,
		Color:       0x5865F2,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:   "Tipo",
				Value:  event.Type,
				Inline: true,
			},
			{
				Name:   "Fecha y Hora",
				Value:  event.DateTime.Format("02/01/2006 15:04"),
				Inline: true,
			},
			{
				Name:   "ID del Evento",
				Value:  fmt.Sprintf("`%s`", event.ID),
				Inline: false,
			},
			{
				Name:   "Inscripciones",
				Value:  buildSignupsText(event),
				Inline: false,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: "Selecciona tu rol para inscribirte",
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}

	// Crear botones para cada rol
	// Discord permite máximo 5 botones por fila
	var components []discordgo.MessageComponent
	var currentRow discordgo.ActionsRow

	for i, role := range event.Roles {
		button := discordgo.Button{
			Label:    role.Name,
			Style:    discordgo.PrimaryButton,
			CustomID: fmt.Sprintf("signup_%s_%s", event.ID, role.Name),
		}
		currentRow.Components = append(currentRow.Components, button)

		// Si llegamos a 5 botones o es el último rol, crear nueva fila
		if len(currentRow.Components) == 5 || i == len(event.Roles)-1 {
			components = append(components, currentRow)
			currentRow = discordgo.ActionsRow{}
		}
	}

	// Botón para cancelar inscripción en una fila separada
	cancelRow := discordgo.ActionsRow{
		Components: []discordgo.MessageComponent{
			discordgo.Button{
				Label:    "Cancelar Inscripción",
				Style:    discordgo.DangerButton,
				CustomID: fmt.Sprintf("cancel_%s", event.ID),
			},
		},
	}
	components = append(components, cancelRow)

	s.ChannelMessageEditComplex(&discordgo.MessageEdit{
		Channel:    event.Channel,
		ID:         event.MessageID,
		Embeds:     &[]*discordgo.MessageEmbed{embed},
		Components: &components,
	})
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

// handleConfig muestra la configuración actual
func handleConfig(s *discordgo.Session, i *discordgo.InteractionCreate) {
	rolesText := ""
	for _, role := range config.AppConfig.DefaultRoles {
		rolesText += fmt.Sprintf("%s %s (Límite: %d)\n", role.Emoji, role.Name, role.Limit)
	}

	embed := &discordgo.MessageEmbed{
		Title: "⚙️ Configuración del Bot",
		Color: 0x5865F2,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Guild ID", Value: config.AppConfig.GuildID, Inline: true},
			{Name: "Puerto Web", Value: config.AppConfig.Port, Inline: true},
			{Name: "Zona Horaria", Value: config.AppConfig.Timezone, Inline: true},
			{Name: "Eventos de Discord", Value: fmt.Sprintf("%v", config.AppConfig.EnableDiscordEvents), Inline: true},
			{Name: "Roles Disponibles", Value: rolesText, Inline: false},
		},
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
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
			Value:  fmt.Sprintf("ID: `%s`\nTipo: %s\nFecha: %s", event.ID, event.Type, event.DateTime.Format("02/01/2006 15:04")),
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

// handleMessageReactionAdd maneja cuando se agrega una reacción
func handleMessageReactionAdd(s *discordgo.Session, r *discordgo.MessageReactionAdd) {
	// Implementación opcional para sistema de reacciones legacy
}

// handleMessageReactionRemove maneja cuando se quita una reacción
func handleMessageReactionRemove(s *discordgo.Session, r *discordgo.MessageReactionRemove) {
	// Implementación opcional para sistema de reacciones legacy
}

// respondError responde con un mensaje de error
func respondError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: "❌ " + message,
			Flags:   discordgo.MessageFlagsEphemeral,
		},
	})
}
