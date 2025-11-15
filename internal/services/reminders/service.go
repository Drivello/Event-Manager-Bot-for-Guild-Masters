package reminders

import (
	"discord-event-bot/config"
	"discord-event-bot/internal/storage"
	"log"
	"time"
)

// ProcessResult representa el resultado del procesamiento de recordatorios.
type ProcessResult struct {
	EventsToRemind         []*storage.Event
	EventsToUpdate         []*storage.Event
	EventsToDeleteMessages []*storage.Event
}

// ProcessReminders aplica la lógica de recordatorios sobre todos los eventos activos
// y devuelve qué eventos necesitan recordatorio y cuáles requieren actualización de mensaje.
func ProcessReminders(now time.Time) ProcessResult {
	result := ProcessResult{}

	events := storage.Store.GetActiveEvents()

	for _, event := range events {
		// Calcular offset de recordatorio (por evento o global)
		offsetMinutes := event.ReminderOffsetMinutes
		if offsetMinutes <= 0 {
			offsetMinutes = config.AppConfig.ReminderOffsetMinutes
		}
		if offsetMinutes <= 0 {
			offsetMinutes = 15
		}

		if event.RepeatEveryDays > 0 {
			changed, shouldRemind := processRecurringEvent(event, now, offsetMinutes)
			if changed {
				if err := storage.Store.SaveEvent(event); err != nil {
					log.Printf("Error guardando evento recurrente %s: %v", event.ID, err)
				} else {
					result.EventsToUpdate = append(result.EventsToUpdate, event)
				}
			}
			if shouldRemind {
				result.EventsToRemind = append(result.EventsToRemind, event)
			}
		} else {
			// Eventos no recurrentes
			reminderTime := event.DateTime.Add(-time.Duration(offsetMinutes) * time.Minute)

			if now.After(reminderTime) && !event.ReminderSent {
				// Marcar recordatorio como enviado y guardar
				event.ReminderSent = true
				if err := storage.Store.SaveEvent(event); err != nil {
					log.Printf("Error guardando evento %s al marcar recordatorio enviado: %v", event.ID, err)
				} else {
					result.EventsToRemind = append(result.EventsToRemind, event)
				}
			}

			// Marcar evento como completado 2h después
			if now.After(event.DateTime.Add(2*time.Hour)) && event.Status != "completed" {
				event.Status = "completed"
				if err := storage.Store.SaveEvent(event); err != nil {
					log.Printf("Error marcando evento %s como completado: %v", event.ID, err)
				}
			}
		}

		// Borrado automático de mensaje/hilo si está configurado
		if event.DeleteAfterHours > 0 && event.MessageID != "" {
			deleteTime := event.DateTime.Add(time.Duration(event.DeleteAfterHours) * time.Hour)
			if now.After(deleteTime) {
				result.EventsToDeleteMessages = append(result.EventsToDeleteMessages, event)
			}
		}
	}

	return result
}

// processRecurringEvent aplica la lógica específica para eventos recurrentes.
// Devuelve si hubo cambios persistentes en el evento y si corresponde enviar recordatorio ahora.
func processRecurringEvent(event *storage.Event, now time.Time, offsetMinutes int) (changed bool, shouldRemind bool) {
	// Avanzar la fecha del evento si ya pasó hace más de 2 horas
	for now.After(event.DateTime.Add(2 * time.Hour)) {
		event.DateTime = event.DateTime.Add(time.Duration(event.RepeatEveryDays) * 24 * time.Hour)
		event.ReminderSent = false
		// Recalcular AnnouncementTime para la nueva fecha si hay offset configurado
		if event.AnnouncementOffsetHours > 0 {
			event.AnnouncementTime = event.DateTime.Add(-time.Duration(event.AnnouncementOffsetHours) * time.Hour)
		}
		changed = true
	}

	if !event.ReminderSent {
		reminderTime := event.DateTime.Add(-time.Duration(offsetMinutes) * time.Minute)
		windowStart := reminderTime.Add(-1 * time.Minute)
		windowEnd := reminderTime.Add(1 * time.Minute)

		if now.After(windowStart) && now.Before(windowEnd) {
			// Toca enviar recordatorio
			event.ReminderSent = true
			changed = true
			shouldRemind = true
		}
	}

	return
}
