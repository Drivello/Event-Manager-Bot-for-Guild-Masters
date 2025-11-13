# ⚡ Guía de Inicio Rápido

Esta guía te ayudará a tener el bot funcionando en **5 minutos**.

## 📋 Pre-requisitos

- ✅ Go 1.21+ instalado ([descargar](https://golang.org/dl/))
- ✅ Bot de Discord creado ([tutorial](#crear-bot-de-discord))
- ✅ Git instalado

## 🚀 Pasos

### 1. Descargar el Proyecto

```bash
git clone <tu-repositorio>
cd discord-event-bot
```

### 2. Configurar Variables de Entorno

```bash
cp .env.example .env
nano .env  # o usa tu editor favorito
```

**Configuración mínima requerida:**
```env
DISCORD_TOKEN=tu_token_aqui
GUILD_ID=tu_guild_id_aqui
ADMIN_USER=admin
ADMIN_PASS=tu_password_seguro
```

### 3. Compilar y Ejecutar

```bash
# Instalar dependencias y compilar
./build.sh

# Ejecutar el bot
./discord-event-bot
```

¡Eso es todo! El bot debería estar corriendo ahora.

## 🌐 Acceder al Panel Web

Abre tu navegador en: **http://localhost:8080**

- Usuario: `admin` (o el que configuraste)
- Contraseña: la que configuraste en `.env`

## 🎮 Probar en Discord

1. Ve a tu servidor de Discord
2. Escribe `/` para ver los comandos disponibles
3. Usa `/create_event` para crear tu primer evento

## 📱 Crear Bot de Discord

Si aún no tienes un bot:

1. Ve a https://discord.com/developers/applications
2. Click en "New Application"
3. Dale un nombre y crea
4. Ve a la sección "Bot" → "Add Bot"
5. **Copia el token** (guárdalo de forma segura)
6. Habilita estos intents:
   - ✅ Presence Intent
   - ✅ Server Members Intent
   - ✅ Message Content Intent
7. Ve a "OAuth2" → "URL Generator"
8. Selecciona scopes: `bot` y `applications.commands`
9. Selecciona permisos: `Administrator` (o permisos específicos)
10. Copia la URL generada y ábrela para invitar el bot

## 🆔 Obtener Guild ID

1. Abre Discord y ve a Configuración de Usuario
2. Avanzado → Habilita "Modo Desarrollador"
3. Click derecho en tu servidor → "Copiar ID"
4. Pega ese ID en `GUILD_ID` en tu `.env`

## ✅ Verificación

Si todo está bien, deberías ver:

```
🚀 Iniciando Discord Event Bot...
✅ Configuración cargada exitosamente
✅ Sistema de almacenamiento inicializado
✅ Bot conectado como: TuBot#1234
📝 Registrando comandos slash...
✅ Bot de Discord inicializado correctamente
✅ Servicio de recordatorios iniciado
🌐 Servidor web disponible en: http://localhost:8080
✅ Bot completamente operacional
```

## 🐛 Problemas Comunes

### "DISCORD_TOKEN es requerido"
→ No configuraste el token en `.env`

### "Invalid authentication"
→ El token es incorrecto, verifica que lo copiaste completo

### "Missing Access"
→ El bot no tiene permisos en el servidor

### Los comandos no aparecen
→ Espera 1-5 minutos para que Discord sincronice los comandos

### Puerto 8080 en uso
→ Cambia el `PORT` en `.env` a otro número (ej: 8081)

## 📚 Siguientes Pasos

- 📖 Lee el [README completo](README.md) para características avanzadas
- ⚙️ Personaliza los roles en `.env`
- 🎨 Personaliza los templates HTML en `internal/web/templates/`
- 🖥️ Sigue la [guía de despliegue en Raspberry Pi](README.md#-instalación-en-raspberry-pi)

## 💡 Comandos Útiles

```bash
# Ver logs en tiempo real
tail -f logs.txt

# Detener el bot
Ctrl+C

# Recompilar después de cambios
./build.sh

# Compilar para Raspberry Pi
./build-pi.sh
```

## 🆘 ¿Necesitas Ayuda?

- 📖 Revisa el [README completo](README.md)
- 🐛 Reporta bugs en GitHub Issues
- 💬 Únete a nuestro servidor de Discord [enlace]

---

**¡Disfruta usando el bot! 🎉**
