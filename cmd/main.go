package main

import (
	"discord-event-bot/config"
	"discord-event-bot/internal/discord"
	"discord-event-bot/internal/storage"
	"discord-event-bot/internal/web"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	log.Println("🚀 Iniciando Discord Event Bot...")

	// Cargar configuración
	if err := config.LoadConfig(); err != nil {
		log.Fatalf("Error cargando configuración: %v", err)
	}

	// Inicializar almacenamiento
	if err := storage.InitEventStore(); err != nil {
		log.Fatalf("Error inicializando almacenamiento: %v", err)
	}

	// Inicializar sistema de templates
	if err := storage.InitTemplateStore(); err != nil {
		log.Fatalf("Error inicializando templates: %v", err)
	}

	// Inicializar bot de Discord
	if err := discord.InitBot(); err != nil {
		log.Fatalf("Error inicializando bot de Discord: %v", err)
	}
	defer discord.Close()

	// Iniciar servicio de recordatorios
	discord.StartReminderService()

	// Inicializar servidor web
	web.InitWebServer()

	// Iniciar servidor web en goroutine
	go func() {
		log.Printf("🌐 Servidor web disponible en: http://localhost:%s", config.AppConfig.Port)
		log.Printf("👤 Usuario: %s", config.AppConfig.AdminUser)
		if err := web.StartWebServer(); err != nil {
			log.Fatalf("Error iniciando servidor web: %v", err)
		}
	}()

	log.Println("✅ Bot completamente operacional")
	log.Println("Presiona Ctrl+C para detener el bot")

	// Esperar señal de interrupción
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	log.Println("🛑 Deteniendo bot...")
	log.Println("👋 ¡Hasta luego!")
}
