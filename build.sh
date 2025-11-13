#!/bin/bash

# Script de compilación para Discord Event Bot

echo "🔨 Compilando Discord Event Bot..."

# Compilar para la arquitectura actual
GOPROXY=https://proxy.golang.org,direct go build -o discord-event-bot cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Compilación exitosa"
    echo "📦 Binario: discord-event-bot"
    echo ""
    echo "Para ejecutar:"
    echo "  ./discord-event-bot"
else
    echo "❌ Error en la compilación"
    exit 1
fi
