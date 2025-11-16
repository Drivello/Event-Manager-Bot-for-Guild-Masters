#!/bin/bash

# Script de compilación cruzada para Raspberry Pi

echo "🔨 Compilando para Raspberry Pi (ARM64)..."

# Compilar para ARM64 (Raspberry Pi Zero 2 W)
GOOS=linux GOARCH=arm64 GOPROXY=https://proxy.golang.org,direct go build -o discord-event-bot-arm64 cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Compilación exitosa para ARM64"
    echo "📦 Binario: discord-event-bot-arm64"
    echo ""
    echo "Para desplegar en Raspberry Pi de forma automática puedes usar:"
    echo "  ./deploy-pi.sh <USUARIO_RPI> <IP_RPI>"
    echo "Ejemplo:"
    echo "  ./deploy-pi.sh pi 192.168.0.82"
else
    echo "❌ Error en la compilación"
    exit 1
fi

# También compilar para ARM (Raspberry Pi más antiguos)
echo ""
echo "🔨 Compilando para Raspberry Pi (ARM)..."
GOOS=linux GOARCH=arm GOARM=7 GOPROXY=https://proxy.golang.org,direct go build -o discord-event-bot-arm cmd/main.go

if [ $? -eq 0 ]; then
    echo "✅ Compilación exitosa para ARM"
    echo "📦 Binario: discord-event-bot-arm"
else
    echo "⚠️  Error en compilación ARM (opcional)"
fi
