package discord

import (
	"discord-event-bot/internal/storage"
	"fmt"
	"log"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

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
				Value:  fmt.Sprintf("<t:%d:F>", event.DateTime.Unix()),
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

	if event.RepeatEveryDays > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Recurrencia",
			Value:  fmt.Sprintf("Cada %d días", event.RepeatEveryDays),
			Inline: true,
		})
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

	for _, role := range event.Roles {
		if len(role.Classes) > 0 {
			// Botones por clase dentro del rol
			for _, class := range role.Classes {
				emojiComponent, isCustomEmoji := parseComponentEmoji(class.Emoji)

				label := class.Name
				if !isCustomEmoji && class.Emoji != "" {
					label = fmt.Sprintf("%s %s", class.Emoji, class.Name)
				}

				button := discordgo.Button{
					Label:    label,
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("signup_%s_%s__%s", event.ID, role.Name, class.Name),
				}
				if emojiComponent != nil {
					button.Emoji = emojiComponent
				}

				currentRow.Components = append(currentRow.Components, button)
				if len(currentRow.Components) == 5 {
					components = append(components, currentRow)
					currentRow = discordgo.ActionsRow{}
				}
			}
		} else {
			// Botón único por rol
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
			if len(currentRow.Components) == 5 {
				components = append(components, currentRow)
				currentRow = discordgo.ActionsRow{}
			}
		}
	}

	if len(currentRow.Components) > 0 {
		components = append(components, currentRow)
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

	threadName := fmt.Sprintf("Chat - %s", event.Name)
	if thread, err := s.MessageThreadStart(event.Channel, msg.ID, threadName, 1440); err != nil {
		log.Printf("Error creando hilo para evento %s: %v", event.ID, err)
	} else if thread != nil {
		event.ThreadID = thread.ID
	}

	if err := storage.Store.SaveEvent(event); err != nil {
		return fmt.Errorf("error guardando evento: %w", err)
	}

	return nil
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
				Value:  fmt.Sprintf("<t:%d:F>", event.DateTime.Unix()),
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

	if event.RepeatEveryDays > 0 {
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   "Recurrencia",
			Value:  fmt.Sprintf("Cada %d días", event.RepeatEveryDays),
			Inline: true,
		})
	}

	// Crear botones para cada rol
	// Discord permite máximo 5 botones por fila
	var components []discordgo.MessageComponent
	var currentRow discordgo.ActionsRow

	for _, role := range event.Roles {
		if len(role.Classes) > 0 {
			// Botones por clase dentro del rol
			for _, class := range role.Classes {
				emojiComponent, isCustomEmoji := parseComponentEmoji(class.Emoji)

				label := class.Name
				if !isCustomEmoji && class.Emoji != "" {
					label = fmt.Sprintf("%s %s", class.Emoji, class.Name)
				}

				button := discordgo.Button{
					Label:    label,
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("signup_%s_%s__%s", event.ID, role.Name, class.Name),
				}
				if emojiComponent != nil {
					button.Emoji = emojiComponent
				}

				currentRow.Components = append(currentRow.Components, button)
				if len(currentRow.Components) == 5 {
					components = append(components, currentRow)
					currentRow = discordgo.ActionsRow{}
				}
			}
		} else {
			// Botón único por rol
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
			if len(currentRow.Components) == 5 {
				components = append(components, currentRow)
				currentRow = discordgo.ActionsRow{}
			}
		}
	}

	if len(currentRow.Components) > 0 {
		components = append(components, currentRow)
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

// buildSignupsText construye el texto de inscripciones
func buildSignupsText(event *storage.Event) string {
	var builder strings.Builder

	for _, role := range event.Roles {
		signups := event.Signups[role.Name]

		// Cabecera del rol con contador simple de inscriptos
		limitText := "∞"
		if role.Limit > 0 {
			limitText = fmt.Sprintf("%d", role.Limit)
		}
		builder.WriteString(fmt.Sprintf("%s **%s**: %d/%s\n",
			role.Emoji, role.Name, len(signups), limitText))

		// Listado de nombres debajo del rol
		for _, signup := range signups {
			builder.WriteString(fmt.Sprintf("- %s\n", signup.Username))
		}

		// Si hay clases definidas, mostrar desglose por clase
		if len(role.Classes) > 0 {
			classCounts := make(map[string]int)
			for _, signup := range signups {
				if signup.Class != "" {
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

func parseComponentEmoji(raw string) (*discordgo.ComponentEmoji, bool) {
	raw = strings.TrimSpace(raw)
	if raw == "" {
		return nil, false
	}

	if strings.HasPrefix(raw, "<") && strings.HasSuffix(raw, ">") {
		inner := strings.Trim(raw, "<>")
		parts := strings.Split(inner, ":")
		if len(parts) == 3 {
			animated := parts[0] == "a"
			name := parts[1]
			id := parts[2]
			if name != "" && id != "" {
				return &discordgo.ComponentEmoji{
					Name:     name,
					ID:       id,
					Animated: animated,
				}, true
			}
		} else if len(parts) == 2 {
			name := parts[0]
			id := parts[1]
			if name != "" && id != "" {
				return &discordgo.ComponentEmoji{
					Name: name,
					ID:   id,
				}, true
			}
		}
	}

	return nil, false
}
