package events

import (
	"discord-event-bot/config"
	"discord-event-bot/internal/storage"
	"fmt"
	"time"

	"github.com/google/uuid"
)

// CreateEventInput contiene los datos necesarios para crear un evento
// independientemente de si viene de la web, Discord, etc.
type CreateEventInput struct {
	Name                    string
	Type                    string
	Description             string
	DateTime                time.Time
	ChannelID               string
	RepeatEveryDays         int
	TemplateName            string
	CreateDiscordEvent      bool
	CreatedBy               string
	AnnounceHours           int
	ReminderOffsetMinutes   int
	DeleteAfterHours        int
	AnnouncementOffsetHours int
}

// CreateEvent aplica las reglas de negocio para crear un evento MMO
// (roles por defecto, templates, programación de anuncio, etc.)
func CreateEvent(input CreateEventInput) (*storage.Event, error) {
	if input.Name == "" {
		return nil, fmt.Errorf("el nombre del evento es obligatorio")
	}
	if input.Type == "" {
		return nil, fmt.Errorf("el tipo de evento es obligatorio")
	}
	if input.ChannelID == "" {
		return nil, fmt.Errorf("el canal es obligatorio")
	}

	announceHours := input.AnnounceHours
	if announceHours < 0 {
		announceHours = 0
	}

	event := &storage.Event{
		ID:                      uuid.New().String(),
		Name:                    input.Name,
		Type:                    input.Type,
		Description:             input.Description,
		DateTime:                input.DateTime,
		Channel:                 input.ChannelID,
		Status:                  "active",
		CreatedAt:               time.Now(),
		CreatedBy:               input.CreatedBy,
		AllowMultiSignup:        false,
		Signups:                 make(map[string][]storage.Signup),
		RepeatEveryDays:         input.RepeatEveryDays,
		CreateDiscordEvent:      input.CreateDiscordEvent,
		ReminderOffsetMinutes:   input.ReminderOffsetMinutes,
		DeleteAfterHours:        input.DeleteAfterHours,
		AnnouncementOffsetHours: 0,
	}

	if announceHours > 0 {
		event.AnnouncementTime = event.DateTime.Add(-time.Duration(announceHours) * time.Hour)
		event.AnnouncementOffsetHours = announceHours
	}

	// Si se especificó un template, delegar en el store
	if input.TemplateName != "" {
		var err error
		event, err = storage.Store.CreateEventFromTemplate(input.TemplateName, event)
		if err != nil {
			return nil, err
		}
	} else {
		// Roles por defecto
		for _, role := range config.AppConfig.DefaultRoles {
			event.Roles = append(event.Roles, storage.RoleSignup{
				Name:  role.Name,
				Emoji: role.Emoji,
				Limit: role.Limit,
			})
		}

		if err := storage.Store.SaveEvent(event); err != nil {
			return nil, err
		}
	}

	return event, nil
}
