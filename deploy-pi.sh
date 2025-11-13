#!/bin/bash

# Script de despliegue automático para Raspberry Pi

if [ -z "$1" ]; then
    echo "❌ Uso: ./deploy-pi.sh <IP_RASPBERRY_PI>"
    echo "Ejemplo: ./deploy-pi.sh 192.168.1.100"
    exit 1
fi

PI_HOST="pi@$1"
PI_DIR="/home/pi/discord-event-bot"

echo "🚀 Desplegando Discord Event Bot en Raspberry Pi..."
echo "Host: $PI_HOST"
echo ""

# Compilar para ARM64
echo "1️⃣  Compilando para ARM64..."
./build-pi.sh
if [ $? -ne 0 ]; then
    echo "❌ Error en compilación"
    exit 1
fi

echo ""
echo "2️⃣  Creando directorio remoto..."
ssh $PI_HOST "mkdir -p $PI_DIR/internal/web/templates $PI_DIR/data/events"

echo ""
echo "3️⃣  Transfiriendo archivos..."
scp discord-event-bot-arm64 $PI_HOST:$PI_DIR/discord-event-bot
scp -r internal/web/templates/* $PI_HOST:$PI_DIR/internal/web/templates/
scp discord-bot.service $PI_HOST:$PI_DIR/

# Transferir .env si existe
if [ -f .env ]; then
    echo "⚠️  Encontrado archivo .env local. ¿Deseas transferirlo? (s/n)"
    read -r response
    if [[ "$response" =~ ^[Ss]$ ]]; then
        scp .env $PI_HOST:$PI_DIR/
        echo "✅ Archivo .env transferido"
    else
        echo "⚠️  Recuerda crear el archivo .env en el Raspberry Pi"
    fi
else
    scp .env.example $PI_HOST:$PI_DIR/
    echo "⚠️  Transferido .env.example - recuerda configurarlo"
fi

echo ""
echo "4️⃣  Configurando permisos..."
ssh $PI_HOST "chmod +x $PI_DIR/discord-event-bot"

echo ""
echo "5️⃣  Instalando servicio systemd..."
ssh $PI_HOST "sudo cp $PI_DIR/discord-bot.service /etc/systemd/system/ && sudo systemctl daemon-reload"

echo ""
echo "✅ Despliegue completado!"
echo ""
echo "📝 Próximos pasos:"
echo "   1. Configurar el archivo .env:"
echo "      ssh $PI_HOST"
echo "      cd $PI_DIR"
echo "      nano .env"
echo ""
echo "   2. Iniciar el servicio:"
echo "      sudo systemctl enable discord-bot"
echo "      sudo systemctl start discord-bot"
echo ""
echo "   3. Ver logs:"
echo "      sudo journalctl -u discord-bot -f"
