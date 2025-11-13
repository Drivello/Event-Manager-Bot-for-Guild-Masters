package storage

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"sync"

	"gopkg.in/yaml.v3"
)

const templatesDir = "data/templates"

// EventTemplate representa un template reutilizable para eventos
type EventTemplate struct {
	Name             string         `json:"name" yaml:"name"`
	Icon             string         `json:"icon" yaml:"icon"`
	MaxParticipants  int            `json:"max_participants" yaml:"max_participants"`
	Description      string         `json:"description" yaml:"description"`
	Roles            []TemplateRole `json:"roles" yaml:"roles"`
	AllowMultiSignup bool           `json:"allow_multi_signup" yaml:"allow_multi_signup"`
	CreatedAt        string         `json:"created_at" yaml:"created_at"`
	UpdatedAt        string         `json:"updated_at" yaml:"updated_at"`
}

// TemplateRole representa un rol dentro de un template
type TemplateRole struct {
	Name    string          `json:"name" yaml:"name"`
	Emoji   string          `json:"emoji" yaml:"emoji"`
	Limit   int             `json:"limit" yaml:"limit"`
	Classes []TemplateClass `json:"classes" yaml:"classes"`
}

// TemplateClass representa una clase/especialización dentro de un rol
type TemplateClass struct {
	Name        string `json:"name" yaml:"name"`
	Emoji       string `json:"emoji" yaml:"emoji"`
	Description string `json:"description,omitempty" yaml:"description,omitempty"`
}

// TemplateStore maneja el almacenamiento de templates
type TemplateStore struct {
	mu        sync.RWMutex
	templates map[string]*EventTemplate
}

var Templates *TemplateStore

// InitTemplateStore inicializa el almacenamiento de templates
func InitTemplateStore() error {
	Templates = &TemplateStore{
		templates: make(map[string]*EventTemplate),
	}

	// Crear directorio de templates si no existe
	if err := os.MkdirAll(templatesDir, 0755); err != nil {
		return fmt.Errorf("error creando directorio de templates: %w", err)
	}

	// Cargar templates existentes
	if err := Templates.LoadTemplates(); err != nil {
		log.Printf("Advertencia al cargar templates: %v", err)
	}

	// Crear templates por defecto si no existen
	if len(Templates.templates) == 0 {
		if err := Templates.CreateDefaultTemplates(); err != nil {
			log.Printf("Error creando templates por defecto: %v", err)
		}
	}

	log.Printf("✅ Sistema de templates inicializado con %d templates", len(Templates.templates))
	return nil
}

// SaveTemplate guarda un template en disco
func (ts *TemplateStore) SaveTemplate(template *EventTemplate) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Validar template
	if err := ts.validateTemplate(template); err != nil {
		return err
	}

	ts.templates[template.Name] = template

	// Guardar como JSON
	filename := filepath.Join(templatesDir, fmt.Sprintf("%s.json", sanitizeFilename(template.Name)))
	data, err := json.MarshalIndent(template, "", "  ")
	if err != nil {
		return fmt.Errorf("error serializando template: %w", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo: %w", err)
	}

	return nil
}

// SaveTemplateYAML guarda un template en formato YAML
func (ts *TemplateStore) SaveTemplateYAML(template *EventTemplate) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	// Validar template
	if err := ts.validateTemplate(template); err != nil {
		return err
	}

	ts.templates[template.Name] = template

	filename := filepath.Join(templatesDir, fmt.Sprintf("%s.yaml", sanitizeFilename(template.Name)))
	data, err := yaml.Marshal(template)
	if err != nil {
		return fmt.Errorf("error serializando template a YAML: %w", err)
	}

	if err := ioutil.WriteFile(filename, data, 0644); err != nil {
		return fmt.Errorf("error escribiendo archivo YAML: %w", err)
	}

	return nil
}

// GetTemplate obtiene un template por nombre
func (ts *TemplateStore) GetTemplate(name string) (*EventTemplate, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	template, exists := ts.templates[name]
	if !exists {
		return nil, fmt.Errorf("template no encontrado: %s", name)
	}

	return template, nil
}

// GetAllTemplates retorna todos los templates
func (ts *TemplateStore) GetAllTemplates() []*EventTemplate {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	templates := make([]*EventTemplate, 0, len(ts.templates))
	for _, template := range ts.templates {
		templates = append(templates, template)
	}

	return templates
}

// DeleteTemplate elimina un template
func (ts *TemplateStore) DeleteTemplate(name string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	delete(ts.templates, name)

	// Intentar eliminar tanto JSON como YAML
	jsonFile := filepath.Join(templatesDir, fmt.Sprintf("%s.json", sanitizeFilename(name)))
	yamlFile := filepath.Join(templatesDir, fmt.Sprintf("%s.yaml", sanitizeFilename(name)))

	os.Remove(jsonFile)
	os.Remove(yamlFile)

	return nil
}

// LoadTemplates carga todos los templates desde disco
func (ts *TemplateStore) LoadTemplates() error {
	files, err := ioutil.ReadDir(templatesDir)
	if err != nil {
		if os.IsNotExist(err) {
			return nil
		}
		return fmt.Errorf("error leyendo directorio de templates: %w", err)
	}

	for _, file := range files {
		ext := filepath.Ext(file.Name())
		if ext != ".json" && ext != ".yaml" && ext != ".yml" {
			continue
		}

		filename := filepath.Join(templatesDir, file.Name())
		data, err := ioutil.ReadFile(filename)
		if err != nil {
			log.Printf("Error leyendo archivo %s: %v", filename, err)
			continue
		}

		var template EventTemplate
		if ext == ".json" {
			if err := json.Unmarshal(data, &template); err != nil {
				log.Printf("Error parseando JSON %s: %v", filename, err)
				continue
			}
		} else {
			if err := yaml.Unmarshal(data, &template); err != nil {
				log.Printf("Error parseando YAML %s: %v", filename, err)
				continue
			}
		}

		ts.templates[template.Name] = &template
	}

	log.Printf("📦 Cargados %d templates desde disco", len(ts.templates))
	return nil
}

// validateTemplate valida que un template sea correcto
func (ts *TemplateStore) validateTemplate(template *EventTemplate) error {
	if template.Name == "" {
		return fmt.Errorf("el nombre del template es requerido")
	}

	if len(template.Roles) == 0 {
		return fmt.Errorf("el template debe tener al menos un rol")
	}

	totalLimit := 0
	for _, role := range template.Roles {
		if role.Name == "" {
			return fmt.Errorf("todos los roles deben tener nombre")
		}
		if role.Limit <= 0 {
			return fmt.Errorf("el límite del rol %s debe ser mayor a 0", role.Name)
		}
		totalLimit += role.Limit
	}

	if template.MaxParticipants > 0 && totalLimit > template.MaxParticipants {
		return fmt.Errorf("la suma de límites de roles (%d) excede el máximo de participantes (%d)", totalLimit, template.MaxParticipants)
	}

	return nil
}

// CreateDefaultTemplates crea templates por defecto
func (ts *TemplateStore) CreateDefaultTemplates() error {
	defaultTemplates := []*EventTemplate{
		{
			Name:            "Raid 20 jugadores",
			Icon:            "⚔️",
			MaxParticipants: 20,
			Description:     "Template estándar para raids de 20 jugadores",
			Roles: []TemplateRole{
				{
					Name:  "Tank",
					Emoji: "🛡️",
					Limit: 4,
					Classes: []TemplateClass{
						{Name: "Paladin", Emoji: "⚔️", Description: "Tank sagrado"},
						{Name: "Warrior", Emoji: "🪓", Description: "Guerrero defensor"},
						{Name: "Death Knight", Emoji: "💀", Description: "Caballero de la muerte"},
					},
				},
				{
					Name:  "DPS",
					Emoji: "🏹",
					Limit: 12,
					Classes: []TemplateClass{
						{Name: "Hunter", Emoji: "🎯", Description: "Cazador"},
						{Name: "Mage", Emoji: "❄️", Description: "Mago"},
						{Name: "Rogue", Emoji: "🗡️", Description: "Pícaro"},
						{Name: "Warlock", Emoji: "🔥", Description: "Brujo"},
					},
				},
				{
					Name:  "Support",
					Emoji: "💖",
					Limit: 4,
					Classes: []TemplateClass{
						{Name: "Priest", Emoji: "⛪", Description: "Sacerdote sanador"},
						{Name: "Druid", Emoji: "🌿", Description: "Druida restaurador"},
						{Name: "Shaman", Emoji: "⚡", Description: "Chamán"},
					},
				},
			},
		},
		{
			Name:            "Dungeon 5 jugadores",
			Icon:            "🏰",
			MaxParticipants: 5,
			Description:     "Template para mazmorras de 5 jugadores",
			Roles: []TemplateRole{
				{
					Name:  "Tank",
					Emoji: "🛡️",
					Limit: 1,
					Classes: []TemplateClass{
						{Name: "Paladin", Emoji: "⚔️"},
						{Name: "Warrior", Emoji: "🪓"},
					},
				},
				{
					Name:  "DPS",
					Emoji: "🏹",
					Limit: 3,
					Classes: []TemplateClass{
						{Name: "Hunter", Emoji: "🎯"},
						{Name: "Mage", Emoji: "❄️"},
						{Name: "Rogue", Emoji: "🗡️"},
					},
				},
				{
					Name:  "Healer",
					Emoji: "💚",
					Limit: 1,
					Classes: []TemplateClass{
						{Name: "Priest", Emoji: "⛪"},
						{Name: "Druid", Emoji: "🌿"},
					},
				},
			},
		},
		{
			Name:            "PvP Battleground",
			Icon:            "⚔️",
			MaxParticipants: 40,
			Description:     "Template para campos de batalla PvP",
			Roles: []TemplateRole{
				{
					Name:  "Melee DPS",
					Emoji: "🗡️",
					Limit: 15,
					Classes: []TemplateClass{
						{Name: "Warrior", Emoji: "🪓"},
						{Name: "Rogue", Emoji: "🗡️"},
						{Name: "Death Knight", Emoji: "💀"},
					},
				},
				{
					Name:  "Ranged DPS",
					Emoji: "🏹",
					Limit: 15,
					Classes: []TemplateClass{
						{Name: "Hunter", Emoji: "🎯"},
						{Name: "Mage", Emoji: "❄️"},
						{Name: "Warlock", Emoji: "🔥"},
					},
				},
				{
					Name:  "Healer",
					Emoji: "💚",
					Limit: 10,
					Classes: []TemplateClass{
						{Name: "Priest", Emoji: "⛪"},
						{Name: "Druid", Emoji: "🌿"},
						{Name: "Shaman", Emoji: "⚡"},
					},
				},
			},
		},
	}

	for _, template := range defaultTemplates {
		if err := ts.SaveTemplate(template); err != nil {
			return err
		}
	}

	log.Printf("✅ Creados %d templates por defecto", len(defaultTemplates))
	return nil
}

// sanitizeFilename limpia un nombre para usarlo como nombre de archivo
func sanitizeFilename(name string) string {
	// Reemplazar caracteres no válidos
	result := ""
	for _, char := range name {
		if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') || char == '-' || char == '_' {
			result += string(char)
		} else if char == ' ' {
			result += "_"
		}
	}
	return result
}

// CloneTemplate crea una copia de un template con un nuevo nombre
func (ts *TemplateStore) CloneTemplate(sourceName, newName string) error {
	ts.mu.Lock()
	defer ts.mu.Unlock()

	source, exists := ts.templates[sourceName]
	if !exists {
		return fmt.Errorf("template fuente no encontrado: %s", sourceName)
	}

	// Crear copia profunda
	data, err := json.Marshal(source)
	if err != nil {
		return err
	}

	var clone EventTemplate
	if err := json.Unmarshal(data, &clone); err != nil {
		return err
	}

	clone.Name = newName
	ts.templates[newName] = &clone

	// Guardar el clon
	filename := filepath.Join(templatesDir, fmt.Sprintf("%s.json", sanitizeFilename(newName)))
	cloneData, err := json.MarshalIndent(&clone, "", "  ")
	if err != nil {
		return err
	}

	return ioutil.WriteFile(filename, cloneData, 0644)
}

// ExportTemplate exporta un template a JSON
func (ts *TemplateStore) ExportTemplate(name string) ([]byte, error) {
	ts.mu.RLock()
	defer ts.mu.RUnlock()

	template, exists := ts.templates[name]
	if !exists {
		return nil, fmt.Errorf("template no encontrado: %s", name)
	}

	return json.MarshalIndent(template, "", "  ")
}

// ImportTemplate importa un template desde JSON
func (ts *TemplateStore) ImportTemplate(data []byte) error {
	var template EventTemplate
	if err := json.Unmarshal(data, &template); err != nil {
		return fmt.Errorf("error parseando JSON: %w", err)
	}

	return ts.SaveTemplate(&template)
}
