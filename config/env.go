package config

import (
	"encoding/json"
	"log"
	"os"
	"strconv"

	"github.com/joho/godotenv"
)

// Config contiene toda la configuración del bot
type Config struct {
	DiscordToken        string
	GuildID             string
	AdminUser           string
	AdminPass           string
	Port                string
	Timezone            string
	DefaultRoles        []Role
	EnableDiscordEvents bool
	MaxReactions        int
	DefaultLimits       map[string]int
}

// Role representa un rol/clase del MMO
type Role struct {
	Name  string `json:"name"`
	Emoji string `json:"emoji"`
	Limit int    `json:"limit"`
}

var AppConfig *Config

// LoadConfig carga la configuración desde el archivo .env
func LoadConfig() error {
	// Cargar archivo .env si existe
	if err := godotenv.Load(); err != nil {
		log.Println("No se encontró archivo .env, usando variables de entorno del sistema")
	}

	config := &Config{
		DiscordToken: getEnv("DISCORD_TOKEN", ""),
		GuildID:      getEnv("GUILD_ID", ""),
		AdminUser:    getEnv("ADMIN_USER", "admin"),
		AdminPass:    getEnv("ADMIN_PASS", "admin123"),
		Port:         getEnv("PORT", "8080"),
		Timezone:     getEnv("TIMEZONE", "America/Argentina/Buenos_Aires"),
		EnableDiscordEvents: getEnvAsBool("ENABLE_DISCORD_EVENTS", true),
		MaxReactions:        getEnvAsInt("MAX_REACTIONS", 50),
	}

	// Parsear roles por defecto
	rolesJSON := getEnv("DEFAULT_ROLES", `[
		{"name":"Tank","emoji":"🛡️","limit":2},
		{"name":"DPS","emoji":"⚔️","limit":6},
		{"name":"Healer","emoji":"💚","limit":2}
	]`)

	if err := json.Unmarshal([]byte(rolesJSON), &config.DefaultRoles); err != nil {
		log.Printf("Error parseando DEFAULT_ROLES, usando valores por defecto: %v", err)
		config.DefaultRoles = []Role{
			{Name: "Tank", Emoji: "🛡️", Limit: 2},
			{Name: "DPS", Emoji: "⚔️", Limit: 6},
			{Name: "Healer", Emoji: "💚", Limit: 2},
		}
	}

	// Parsear límites por defecto
	limitsJSON := getEnv("DEFAULT_LIMITS", `{"Tank":2,"DPS":6,"Healer":2}`)
	config.DefaultLimits = make(map[string]int)
	if err := json.Unmarshal([]byte(limitsJSON), &config.DefaultLimits); err != nil {
		log.Printf("Error parseando DEFAULT_LIMITS, usando valores por defecto: %v", err)
		config.DefaultLimits = map[string]int{
			"Tank":   2,
			"DPS":    6,
			"Healer": 2,
		}
	}

	// Validar configuración crítica
	if config.DiscordToken == "" {
		log.Fatal("DISCORD_TOKEN es requerido")
	}
	if config.GuildID == "" {
		log.Fatal("GUILD_ID es requerido")
	}

	AppConfig = config
	log.Println("✅ Configuración cargada exitosamente")
	return nil
}

// getEnv obtiene una variable de entorno o retorna un valor por defecto
func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

// getEnvAsInt convierte una variable de entorno a int
func getEnvAsInt(key string, defaultValue int) int {
	valueStr := getEnv(key, "")
	if value, err := strconv.Atoi(valueStr); err == nil {
		return value
	}
	return defaultValue
}

// getEnvAsBool convierte una variable de entorno a bool
func getEnvAsBool(key string, defaultValue bool) bool {
	valueStr := getEnv(key, "")
	if value, err := strconv.ParseBool(valueStr); err == nil {
		return value
	}
	return defaultValue
}
