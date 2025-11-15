package signups

import (
	"discord-event-bot/internal/storage"
	"fmt"
)

// SignupInput representa los datos necesarios para inscribir a un usuario en un evento.
type SignupInput struct {
	EventID  string
	UserID   string
	Username string
	Role     string
	Class    string
}

// CancelInput representa los datos necesarios para cancelar la inscripción de un usuario.
type CancelInput struct {
	EventID string
	UserID  string
}

// SignupToEvent aplica las reglas de negocio para inscribir a un usuario en un evento.
func SignupToEvent(input SignupInput) (*storage.Event, error) {
	event, err := storage.Store.GetEvent(input.EventID)
	if err != nil {
		return nil, fmt.Errorf("Evento no encontrado")
	}

	// Verificar si ya está inscrito
	for r, signups := range event.Signups {
		for _, signup := range signups {
			if signup.UserID == input.UserID {
				if r == input.Role {
					return nil, fmt.Errorf("Ya estás inscrito en este rol")
				}
				if !event.AllowMultiSignup {
					return nil, fmt.Errorf("Ya estás inscrito en otro rol. Cancela primero tu inscripción actual.")
				}
			}
		}
	}

	// Verificar límite de rol
	confirmedCount := 0
	for _, signup := range event.Signups[input.Role] {
		if signup.Status == "confirmed" {
			confirmedCount++
		}
	}

	var roleLimit int
	for _, r := range event.Roles {
		if r.Name == input.Role {
			roleLimit = r.Limit
			break
		}
	}

	if roleLimit > 0 && confirmedCount >= roleLimit {
		return nil, fmt.Errorf("El rol %s ya está lleno", input.Role)
	}

	// Agregar inscripción (con clase si aplica)
	if input.Class != "" {
		if err := storage.Store.AddSignupWithClass(input.EventID, input.UserID, input.Username, input.Role, input.Class); err != nil {
			return nil, fmt.Errorf("Error procesando inscripción")
		}
	} else {
		if err := storage.Store.AddSignup(input.EventID, input.UserID, input.Username, input.Role); err != nil {
			return nil, fmt.Errorf("Error procesando inscripción")
		}
	}

	// El puntero event apunta a la misma instancia que se actualiza en el store
	return event, nil
}

// CancelSignup aplica las reglas de negocio para cancelar la inscripción de un usuario.
func CancelSignup(input CancelInput) (*storage.Event, error) {
	event, err := storage.Store.GetEvent(input.EventID)
	if err != nil {
		return nil, fmt.Errorf("Evento no encontrado")
	}

	userID := input.UserID
	removed := false

	for role := range event.Signups {
		if err := storage.Store.RemoveSignup(input.EventID, userID, role); err == nil {
			removed = true
		}
	}

	if !removed {
		return nil, fmt.Errorf("No estás inscrito en este evento")
	}

	return event, nil
}
