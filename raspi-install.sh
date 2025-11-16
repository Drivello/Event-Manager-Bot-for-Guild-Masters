#!/bin/bash
set -e

SERVICE_NAME="discord-bot"
INSTALL_DIR="$(cd "$(dirname "$0")" && pwd)"

echo "🔧 Configurando Discord Event Bot en Raspberry Pi..."
echo "Directorio de instalación: $INSTALL_DIR"
echo ""

# Asegurar permisos de ejecución del binario
if [ -f "$INSTALL_DIR/discord-event-bot" ]; then
  chmod +x "$INSTALL_DIR/discord-event-bot"
else
  echo "❌ No se encontró $INSTALL_DIR/discord-event-bot"
  exit 1
fi

# Instalar servicio systemd
SERVICE_FILE="$INSTALL_DIR/discord-bot.service"
if [ ! -f "$SERVICE_FILE" ]; then
  echo "❌ No se encontró $SERVICE_FILE"
  exit 1
fi

echo "📦 Instalando servicio systemd en /etc/systemd/system/$SERVICE_NAME.service..."
sudo cp "$SERVICE_FILE" "/etc/systemd/system/$SERVICE_NAME.service"
sudo systemctl daemon-reload

# Habilitar y reiniciar servicio
if ! systemctl is-enabled --quiet "$SERVICE_NAME" 2>/dev/null; then
  echo "🔁 Habilitando servicio..."
  sudo systemctl enable "$SERVICE_NAME"
fi

echo "🚀 Reiniciando servicio..."
sudo systemctl restart "$SERVICE_NAME"

echo ""
echo "✅ Instalación/actualización completa."
echo "   Ver logs con: sudo journalctl -u $SERVICE_NAME -f"
