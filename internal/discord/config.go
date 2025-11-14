package discord

import (
	"discord-event-bot/config"
	"fmt"

	"github.com/bwmarrin/discordgo"
)

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
