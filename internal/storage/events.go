package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"
)

const eventsDir = "data/events"

// Event representa un evento del MMO
type Event struct {
	ID               string              `json:"id"`
	Name             string              `json:"name"`
	Type             string              `json:"type"`
	Description      string              `json:"description"`
	DateTime         time.Time           `json:"datetime"`
	Channel          string              `json:"channel"`
	MessageID        string              `json:"message_id"`
	DiscordEventID   string              `json:"discord_event_id,omitempty"`
	TemplateName     string              `json:"template_name,omitempty"`
	Roles            []RoleSignup        `json:"roles"`
	Signups          map[string][]Signup `json:"signups"`
	ReminderSent     bool                `json:"reminder_sent"`
	CreatedAt        time.Time           `json:"created_at"`
	CreatedBy        string              `json:"created_by"`
	AllowMultiSignup bool                `json:"allow_multi_signup"`
	Status           string              `json:"status"` // active, completed, cancelled
	MaxParticipants  int                 `json:"max_participants,omitempty"`
	RepeatEveryDays  int                 `json:"repeat_every_days,omitempty"`
}

// RoleSignup representa un rol disponible para el evento
type RoleSignup struct {
	Name    string      `json:"name"`
	Emoji   string      `json:"emoji"`
	Limit   int         `json:"limit"`
	Classes []ClassInfo `json:"classes,omitempty"`
}

// ClassInfo representa una clase/especialización dentro de un rol
type ClassInfo struct {
	Name        string `json:"name"`
	Emoji       string `json:"emoji"`
	Description string `json:"description,omitempty"`
}

// Signup representa una inscripción de usuario
type Signup struct {
	UserID      string    `json:"user_id"`
	Username    string    `json:"username"`
	Role        string    `json:"role"`
	Class       string    `json:"class,omitempty"`
	Status      string    `json:"status"` // pending, confirmed, declined
	SignedUpAt  time.Time `json:"signed_up_at"`
	ConfirmedBy string    `json:"confirmed_by,omitempty"`
}

// EventStore maneja el almacenamiento de eventos
type EventStore struct {
	mu     sync.RWMutex
	events map[string]*Event
}

var Store *EventStore

// InitStore inicializa el almacenamiento de eventos
func InitStore() error {
	Store = &EventStore{
		events: make(map[string]*Event),
	}

	// Crear directorio de datos si no existe
	if err := os.MkdirAll(eventsDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de datos: %w", err)
	}

	// Cargar eventos existentes
	if err := Store.LoadEvents(); err != nil {
		log.Printf("Advertencia al cargar eventos: %v", err)
	}

	log.Println("✅ Sistema de almacenamiento inicializado")
	return nil
}

// SaveEvent guarda un evento en disco
func (s *EventStore) SaveEvent(event *Event) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	s.events[event.ID] = event

	filename := filepath.Join(eventsDir, fmt.Sprintf("%s.json", event.ID))
	data, err := json.MarshalIndent(event, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando evento: %w", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return nil
}

// GetEvent obtiene un evento por ID
func (s *EventStore) GetEvent(id string) (*Event, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	event, exists := s.events[id]
	if !exists {
		return nil, fmt.Errorf("evento no encontrado: %s", id)
	}

	return event, nil
}

// GetAllEvents retorna todos los eventos
func (s *EventStore) GetAllEvents() []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*Event, 0, len(s.events))
	for _, event := range s.events {
		events = append(events, event)
	}

	return events
}

// GetActiveEvents retorna solo eventos activos
func (s *EventStore) GetActiveEvents() []*Event {
	s.mu.RLock()
	defer s.mu.RUnlock()

	events := make([]*Event, 0)
	for _, event := range s.events {
		if event.Status == "active" {
			events = append(events, event)
		}
	}

	return events
}

// DeleteEvent elimina un evento
func (s *EventStore) DeleteEvent(id string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.events, id)

	filename := filepath.Join(eventsDir, fmt.Sprintf("%s.json", id))
	if err := os.Remove(filename); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("error eliminando archivo: %w", err)
	}

	return nil
}

// LoadEvents carga todos los eventos desde disco
func (s *EventStore) LoadEvents() error {
	files, err := ioutil.ReadDir(eventsDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // Directorio no existe aún, no es error
		}
		return fmt.Errorf("error leyendo directorio: %w", err)
	}

	for _, file := range files {
		if filepath.Ext(file.Name()) != ".json" {
			continue
		}

		filename := filepath.Join(eventsDir, file.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Error leyendo archivo %s: %v", filename, err)
			continue
		}

		var event Event
		if err := json.Unmarshal(data, &event); err != nil {
			log.Printf("Error parseando archivo %s: %v", filename, err)
			continue
		}

		s.events[event.ID] = &event
	}

	log.Printf("📦 Cargados %d eventos desde disco", len(s.events))
	return nil
}

// AddSignup agrega una inscripción a un evento
func (s *EventStore) AddSignup(eventID, userID, username, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[eventID]
	if !exists {
		return fmt.Errorf("evento no encontrado")
	}

	if event.Signups == nil {
		event.Signups = make(map[string][]Signup)
	}

	signup := Signup{
		UserID:     userID,
		Username:   username,
		Role:       role,
		Status:     "pending",
		SignedUpAt: time.Now(),
	}

	event.Signups[role] = append(event.Signups[role], signup)

	return s.SaveEvent(event)
}

// RemoveSignup elimina una inscripción
func (s *EventStore) RemoveSignup(eventID, userID, role string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[eventID]
	if !exists {
		return fmt.Errorf("evento no encontrado")
	}

	signups := event.Signups[role]
	for i, signup := range signups {
		if signup.UserID == userID {
			event.Signups[role] = append(signups[:i], signups[i+1:]...)
			break
		}
	}

	return s.SaveEvent(event)
}

// ConfirmSignup confirma una inscripción
func (s *EventStore) ConfirmSignup(eventID, userID, role, confirmedBy string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[eventID]
	if !exists {
		return fmt.Errorf("evento no encontrado")
	}

	signups := event.Signups[role]
	for i, signup := range signups {
		if signup.UserID == userID {
			event.Signups[role][i].Status = "confirmed"
			event.Signups[role][i].ConfirmedBy = confirmedBy
			break
		}
	}

	return s.SaveEvent(event)
}

// CreateEventFromTemplate crea un evento basado en un template
func (s *EventStore) CreateEventFromTemplate(templateName string, eventData *Event) (*Event, error) {
	template, err := Templates.GetTemplate(templateName)
	if err != nil {
		return nil, fmt.Errorf("template no encontrado: %w", err)
	}

	// Copiar configuración del template al evento
	eventData.TemplateName = templateName
	eventData.MaxParticipants = template.MaxParticipants
	eventData.AllowMultiSignup = template.AllowMultiSignup

	// Convertir roles del template a roles del evento
	eventData.Roles = make([]RoleSignup, 0, len(template.Roles))
	for _, tRole := range template.Roles {
		role := RoleSignup{
			Name:  tRole.Name,
			Emoji: tRole.Emoji,
			Limit: tRole.Limit,
		}

		// Convertir clases
		if len(tRole.Classes) > 0 {
			role.Classes = make([]ClassInfo, 0, len(tRole.Classes))
			for _, tClass := range tRole.Classes {
				role.Classes = append(role.Classes, ClassInfo{
					Name:        tClass.Name,
					Emoji:       tClass.Emoji,
					Description: tClass.Description,
				})
			}
		}

		eventData.Roles = append(eventData.Roles, role)
	}

	// Inicializar signups
	if eventData.Signups == nil {
		eventData.Signups = make(map[string][]Signup)
	}

	// Guardar evento
	if err := s.SaveEvent(eventData); err != nil {
		return nil, err
	}

	return eventData, nil
}

// AddSignupWithClass agrega una inscripción con clase específica
func (s *EventStore) AddSignupWithClass(eventID, userID, username, role, class string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	event, exists := s.events[eventID]
	if !exists {
		return fmt.Errorf("evento no encontrado")
	}

	if event.Signups == nil {
		event.Signups = make(map[string][]Signup)
	}

	signup := Signup{
		UserID:     userID,
		Username:   username,
		Role:       role,
		Class:      class,
		Status:     "pending",
		SignedUpAt: time.Now(),
	}

	event.Signups[role] = append(event.Signups[role], signup)

	return s.SaveEvent(event)
}
