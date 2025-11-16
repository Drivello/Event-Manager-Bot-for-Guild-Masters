#!/bin/bash

# Script de despliegue automático para Raspberry Pi

if [ -z "$1" ] || [ -z "$2" ]; then
	echo "❌ Uso: ./deploy-pi.sh <USUARIO_RPI> <IP_RPI>"
	echo "Ejemplo: ./deploy-pi.sh pi 192.168.0.82"
	exit 1
fi

PI_USER="$1"
PI_IP="$2"
PI_HOST="$PI_USER@$PI_IP"
PI_DIR="/home/$PI_USER/event-manager-bot"

echo "🚀 Desplegando Discord Event Bot en Raspberry Pi..."
echo "Usuario: $PI_USER"
echo "Host:   $PI_IP"
echo "Dir:    $PI_DIR"
echo ""

# 1️⃣ Compilar para ARM64
echo "1️⃣  Compilando para ARM64..."
./build-pi.sh
if [ $? -ne 0 ]; then
	echo "❌ Error en compilación"
	exit 1
fi

echo ""
echo "2️⃣  Creando estructura remota..."
ssh "$PI_HOST" "mkdir -p $PI_DIR/internal/web/templates $PI_DIR/data/events $PI_DIR/data/templates"

echo ""
echo "3️⃣  Transfiriendo archivos..."
scp discord-event-bot-arm64 "$PI_HOST:$PI_DIR/discord-event-bot"
scp -r internal/web/templates/* "$PI_HOST:$PI_DIR/internal/web/templates/"
scp discord-bot.service "$PI_HOST:$PI_DIR/"
scp raspi-install.sh "$PI_HOST:$PI_DIR/"

# Transferir .env si existe
if [ -f .env ]; then
	echo "⚠️  Encontrado archivo .env local. ¿Deseas transferirlo? (s/n)"
	read -r response
	if [[ "$response" =~ ^[Ss]$ ]]; then
		scp .env "$PI_HOST:$PI_DIR/"
		echo "✅ Archivo .env transferido"
	else
		echo "⚠️  Recuerda crear o ajustar el archivo .env en la Raspberry"
	fi
elif [ -f .env.example ]; then
	scp .env.example "$PI_HOST:$PI_DIR/.env.example"
	echo "⚠️  Transferido .env.example - recuerda configurarlo en la Raspberry"
else
	echo "⚠️  No se encontró .env ni .env.example en el proyecto"
fi

echo ""
echo "4️⃣  Dando permisos al instalador remoto..."
ssh "$PI_HOST" "chmod +x $PI_DIR/raspi-install.sh"

echo ""
echo "5️⃣  Ejecutando instalador en la Raspberry..."
ssh "$PI_HOST" "cd $PI_DIR && ./raspi-install.sh"

echo ""
echo "✅ Despliegue completado!"
echo ""
echo "📝 Comandos útiles en la Raspberry:" 
echo "   Ver estado:   sudo systemctl status discord-bot"
echo "   Ver logs:     sudo journalctl -u discord-bot -f"
echo "   Editar .env:  nano $PI_DIR/.env"
