package discord

import (
	"discord-event-bot/config"
	"fmt"
	"log"

	"github.com/bwmarrin/discordgo"
)

// InitBot inicializa el bot de Discord
func InitBot() error {
	var err error
	Session, err = discordgo.New("Bot " + config.AppConfig.DiscordToken)
	if err != nil {
		return fmt.Errorf("error creando sesión de Discord: %w", err)
	}

	Session.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("✅ Bot conectado como: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	// Registrar handlers de interacciones
	Session.AddHandler(handleInteractionCreate)
	// Session.AddHandler(handleMessageReactionAdd)
	// Session.AddHandler(handleMessageReactionRemove)

	// Necesitamos permisos para intents
	Session.Identify.Intents = discordgo.IntentsGuildMessages |
		discordgo.IntentsGuildMessageReactions |
		discordgo.IntentsGuilds |
		discordgo.IntentsGuildScheduledEvents

	// Abrir conexión
	if err := Session.Open(); err != nil {
		return fmt.Errorf("error abriendo conexión: %w", err)
	}

	// Registrar comandos slash
	log.Println("📝 Registrando comandos slash...")
	for _, cmd := range commands {
		_, err := Session.ApplicationCommandCreate(Session.State.User.ID, config.AppConfig.GuildID, cmd)
		if err != nil {
			log.Printf("Error registrando comando %s: %v", cmd.Name, err)
		}
	}

	log.Println("✅ Bot de Discord inicializado correctamente")
	return nil
}

// Close cierra la sesión de Discord
func Close() {
	if Session != nil {
		Session.Close()
	}
}
