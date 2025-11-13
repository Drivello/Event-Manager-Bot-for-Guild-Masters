# 🎮 Discord Event Bot para MMO Guilds

Bot profesional de Discord especializado en la gestión de eventos para guilds de juegos MMO (WoW, FFXIV, etc.), con panel web de administración. Optimizado para ejecutarse en dispositivos de bajo consumo como **Raspberry Pi Zero 2 W**.

## ✨ Características

### Bot de Discord
- ✅ Comandos slash para gestión completa de eventos
- 🎯 Sistema de inscripciones con botones interactivos
- 👥 Roles personalizables (Tank, DPS, Healer, etc.)
- 🎨 **Sistema de templates reutilizables** con clases/especializaciones
- 🔔 Recordatorios automáticos programables
- ✅ Confirmación manual de inscritos por administradores
- 📅 Integración opcional con eventos oficiales de Discord
- 💾 Almacenamiento local en archivos JSON/YAML (sin base de datos externa)

### Panel Web
- 🌐 Interfaz web responsive accesible en LAN
- 🔐 Autenticación básica con usuario/contraseña
- 📝 Creación y gestión de eventos desde el navegador
- 🎨 **Editor visual de templates** con vista previa en tiempo real
- 👥 Visualización de inscripciones en tiempo real
- 📥 Importar/Exportar templates en JSON
- ⚙️ Página de configuración del sistema
- 📱 Diseño optimizado para móviles

## 📋 Requisitos

- Go 1.21 o superior
- Token de bot de Discord
- Servidor Discord con permisos de administrador
- (Opcional) Raspberry Pi Zero 2 W o similar para despliegue

## 🚀 Instalación

### 1. Clonar o descargar el proyecto

```bash
git clone <tu-repositorio>
cd discord-event-bot
```

### 2. Configurar variables de entorno

Copia el archivo de ejemplo y edítalo con tus credenciales:

```bash
cp .env.example .env
nano .env
```

Variables requeridas:
- `DISCORD_TOKEN`: Token de tu bot de Discord
- `GUILD_ID`: ID de tu servidor de Discord
- `ADMIN_USER` y `ADMIN_PASS`: Credenciales del panel web

### 3. Obtener el Token de Discord

1. Ve a https://discord.com/developers/applications
2. Crea una nueva aplicación
3. En la sección "Bot", crea un bot y copia el token
4. Habilita los siguientes **Privileged Gateway Intents**:
   - Server Members Intent
   - Message Content Intent
5. En "OAuth2 > URL Generator", selecciona:
   - Scopes: `bot`, `applications.commands`
   - Bot Permissions: `Administrator` (o permisos específicos)
6. Usa la URL generada para invitar el bot a tu servidor

### 4. Obtener el Guild ID

1. Habilita el modo desarrollador en Discord (Ajustes > Avanzado > Modo desarrollador)
2. Click derecho en tu servidor > Copiar ID

### 5. Compilar e instalar dependencias

```bash
go mod tidy
go build -o discord-event-bot cmd/main.go
```

### 6. Ejecutar el bot

```bash
./discord-event-bot
```

El bot estará disponible en:
- Discord: Automáticamente conectado
- Panel Web: http://localhost:8080

## 📦 Estructura del Proyecto

```
discord-event-bot/
├── cmd/
│   └── main.go                 # Punto de entrada principal
├── config/
│   └── env.go                  # Gestión de configuración
├── internal/
│   ├── discord/
│   │   └── handler.go          # Lógica del bot de Discord
│   ├── storage/
│   │   └── events.go           # Sistema de almacenamiento JSON
│   └── web/
│       ├── server.go           # Servidor web
│       └── templates/          # Templates HTML
│           ├── index.html
│           ├── create_event.html
│           ├── event_detail.html
│           ├── events.html
│           ├── config.html
│           └── error.html
├── data/
│   └── events/                 # Archivos JSON de eventos
├── go.mod                      # Dependencias de Go
├── .env.example                # Plantilla de configuración
├── discord-bot.service         # Archivo de servicio systemd
└── README.md                   # Este archivo
```

## 🎯 Comandos de Discord

### Comandos Slash Disponibles

- `/create_event` - Crear un nuevo evento
  - `nombre`: Nombre del evento
  - `tipo`: Tipo (Raid, Dungeon, PvP, Social, etc.)
  - `fecha`: Fecha y hora (formato: YYYY-MM-DD HH:MM)
  - `descripcion`: Descripción del evento
  - `canal`: Canal donde publicar (opcional)

- `/delete_event` - Eliminar un evento existente
  - `id`: ID del evento

- `/remind_event` - Enviar recordatorio inmediato
  - `id`: ID del evento

- `/list_events` - Listar todos los eventos activos

- `/config` - Mostrar configuración actual del bot

## 🌐 Panel Web

### Acceso

Navega a `http://localhost:8080` (o la IP de tu dispositivo si accedes desde otro equipo en la LAN)

Credenciales por defecto (cámbialas en `.env`):
- Usuario: `admin`
- Contraseña: `admin123`

### Funcionalidades

- **Dashboard**: Vista de eventos activos
- **Crear Evento**: Formulario para crear eventos desde el navegador
- **Ver Eventos**: Lista completa de todos los eventos
- **Detalles de Evento**: Ver inscripciones y confirmar participantes
- **Configuración**: Ver ajustes actuales del bot

## 🔧 Configuración Avanzada

### Personalizar Roles

Edita la variable `DEFAULT_ROLES` en `.env`:

```env
DEFAULT_ROLES=[{"name":"Tank","emoji":"🛡️","limit":2},{"name":"Healer","emoji":"💚","limit":3},{"name":"DPS","emoji":"⚔️","limit":8},{"name":"Support","emoji":"🔮","limit":2}]
```

### Zona Horaria

Cambia la zona horaria según tu ubicación:

```env
TIMEZONE=America/Argentina/Buenos_Aires
```

Opciones comunes:
- `America/New_York`
- `Europe/Madrid`
- `America/Mexico_City`
- `America/Santiago`

## 🖥️ Instalación en Raspberry Pi

### 1. Compilar para ARM

En tu PC (compilación cruzada):

```bash
GOOS=linux GOARCH=arm64 go build -o discord-event-bot cmd/main.go
```

### 2. Transferir archivos

```bash
scp discord-event-bot pi@tu-raspberry-pi:/home/pi/
scp .env pi@tu-raspberry-pi:/home/pi/
```

### 3. Configurar como servicio systemd

```bash
sudo cp discord-bot.service /etc/systemd/system/
sudo systemctl daemon-reload
sudo systemctl enable discord-bot
sudo systemctl start discord-bot
```

### 4. Verificar estado

```bash
sudo systemctl status discord-bot
sudo journalctl -u discord-bot -f
```

## 📊 Logs y Monitoreo

Ver logs en tiempo real:

```bash
sudo journalctl -u discord-bot -f
```

Ver logs históricos:

```bash
sudo journalctl -u discord-bot --since "1 hour ago"
```

## 🔒 Seguridad

- ✅ El panel web usa autenticación básica HTTP
- ✅ Solo accesible desde LAN por defecto
- ✅ Tokens y contraseñas en archivo `.env` (no versionado)
- ✅ Servicio systemd con restricciones de seguridad
- ⚠️ Para acceso remoto, usa un túnel SSH o VPN

### Túnel SSH para acceso remoto

```bash
ssh -L 8080:localhost:8080 pi@tu-raspberry-pi
```

Luego accede desde tu navegador a `http://localhost:8080`

## 🎨 Sistema de Templates

El bot incluye un sistema completo de templates para eventos reutilizables. Ver **[TEMPLATES_GUIDE.md](TEMPLATES_GUIDE.md)** para documentación detallada.

### Características de Templates
- 📝 Crear templates personalizados con roles y clases
- 🎯 Definir cupos específicos por rol y clase
- 🎨 Emojis personalizados para cada elemento
- 💾 Almacenamiento en JSON o YAML
- 📥 Importar/Exportar templates
- 🔄 Clonar y modificar templates existentes
- 👁️ Vista previa en tiempo real en el editor web

### Templates Incluidos
- **Raid 20 jugadores** - Template estándar para raids
- **Dungeon 5 jugadores** - Para mazmorras pequeñas
- **PvP Battleground** - Para campos de batalla de 40 jugadores

### Uso Rápido

**Desde Discord:**
```
/create_event nombre:"Raid Semanal" tipo:Raid fecha:"2024-12-20 20:00" 
  descripcion:"Raid del viernes" template:"Raid 20 jugadores"
```

**Desde Panel Web:**
1. Ve a `/templates` para gestionar templates
2. Crea eventos en `/events/create` seleccionando un template

## 🛠️ Solución de Problemas

### El bot no se conecta a Discord

1. Verifica que el token sea correcto en `.env`
2. Asegúrate que el bot esté invitado al servidor
3. Revisa los logs: `journalctl -u discord-bot`

### No aparecen los comandos slash

1. Espera unos minutos (Discord puede tardar en sincronizar)
2. Reinicia el bot
3. Verifica que el `GUILD_ID` sea correcto
4. Confirma que el bot tenga permisos de `applications.commands`

### El panel web no carga

1. Verifica que el puerto no esté en uso: `netstat -tuln | grep 8080`
2. Comprueba que los templates HTML estén en `internal/web/templates/`
3. Revisa los logs para errores

### Error de permisos en Raspberry Pi

```bash
chmod +x discord-event-bot
chown pi:pi discord-event-bot
```

## 🔄 Actualización

```bash
git pull
go build -o discord-event-bot cmd/main.go
sudo systemctl restart discord-bot
```

## 📝 Formato de Fechas

Al crear eventos, usa el formato: `YYYY-MM-DD HH:MM`

Ejemplos:
- `2024-12-25 20:00` - 25 de diciembre a las 8 PM
- `2024-01-15 14:30` - 15 de enero a las 2:30 PM

## 🤝 Contribución

Este proyecto es de código abierto. Si encuentras bugs o quieres agregar features:

1. Crea un fork del repositorio
2. Haz tus cambios en una rama nueva
3. Envía un pull request

## 📄 Licencia

Este proyecto está bajo licencia MIT. Ver archivo `LICENSE` para más detalles.

## 🙏 Créditos

Desarrollado con:
- [discordgo](https://github.com/bwmarrin/discordgo) - Biblioteca de Discord para Go
- [gin](https://github.com/gin-gonic/gin) - Framework web
- [godotenv](https://github.com/joho/godotenv) - Gestión de variables de entorno

## 📞 Soporte

Para reportar problemas o sugerencias, abre un issue en el repositorio.

---

**¡Disfruta organizando eventos para tu guild! 🎮**
