# 🎮 Discord Event Bot para MMO Guilds

Bot profesional de Discord especializado en la gestión de eventos para guilds de juegos MMO (WoW, FFXIV, etc.), con panel web de administración. 
Se creó con el objetivo de ayudarme a gestionar eventos en mi guild ejecutándose en una simple Raspberry Pi Zero 2.

## ✨ Características

### Bot de Discord
- ✅ Comandos slash para gestión completa de eventos
- 🎯 Sistema de inscripciones con botones interactivos por rol
- 🧬 Botones por clase dentro de cada rol, con emojis personalizados
- 👥 Roles personalizables (Tank, DPS, Healer, etc.)
- 🎨 Sistema de templates reutilizables con clases/especializaciones
- 📊 Límites opcionales por rol y globales (0 = sin límite, se muestra como ∞)
- 🧵 Creación automática de hilos de discusión por evento 
- 🔁 Soporte para eventos recurrentes
- 🔔 Recordatorios automáticos programables
- 📅 Integración opcional con eventos oficiales de Discord
- 💾 Almacenamiento local en archivos JSON/YAML

### Panel Web
- 🌐 Interfaz web responsive accesible en LAN
- 🔐 Autenticación básica con usuario/contraseña
- 📝 Creación y gestión de eventos desde el navegador
- 🎨 Editor visual de templates con vista previa en tiempo real
- 👥 Visualización de inscripciones en tiempo real
- 📥 Importar/Exportar templates en JSON
- ⚙️ Página de configuración del sistema
- 🧹 Botón para limpiar eventos cancelados del historial
- 📱 Diseño optimizado para móviles

## 📋 Requisitos

- Go 1.21 o superior
- Token de bot de Discord
- Servidor Discord con permisos de administrador
- (Opcional) Raspberry Pi Zero 2 W o similar para despliegue

## 🚀 Instalación

### 1. Clonar o descargar el proyecto

```bash
git clone https://github.com/Drivello/Event-Manager-Bot-for-Guild-Masters
cd Event-Manager-Bot-for-Guild-Masters
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
./build.sh
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
Event-Manager-Bot-for-Guild-Masters/
├── cmd/
│   └── main.go                 # Punto de entrada principal
├── config/
│   └── env.go                  # Gestión de configuración y variables de entorno
├── internal/
│   ├── discord/
│   │   ├── botInit.go          # Inicialización del bot de Discord y registro de handlers
│   │   ├── config.go           # Configuración específica del bot de Discord
│   │   ├── interactions.go     # Comandos slash y ruteo de interacciones
│   │   ├── events.go           # Lógica de creación/listado/eliminación de eventos
│   │   ├── messages.go         # Publicación y actualización de mensajes y botones
│   │   ├── signup.go           # Manejo de inscripciones y cancelaciones
│   │   ├── errors.go           # Helpers para respuestas de error
│   │   └── reminders.go        # Servicio de recordatorios
│   ├── storage/
│   │   ├── events.go           # Sistema de almacenamiento JSON de eventos
│   │   └── templates.go        # Sistema de almacenamiento de templates
│   └── web/
│       ├── server.go           # Servidor web (panel de administración)
│       └── templates/          # Templates HTML del panel
│           ├── index.html
│           ├── create_event.html
│           ├── event_detail.html
│           ├── events.html
│           ├── templates.html
│           ├── template_editor.html
│           ├── config.html
│           └── error.html
├── data/
│   ├── events/                 # Archivos JSON de eventos
│   └── templates/              # Archivos de templates (JSON/YAML)
├── go.mod                      # Dependencias de Go
├── .env.example                # Plantilla de configuración
├── discord-bot.service         # Archivo de servicio systemd
└── README.md                   # Este archivo
```

## 🎯 Comandos de Discord

### Comandos Slash Disponibles

- `/create_event` - Crear un nuevo evento (y su hilo de discusión)
  - `nombre`: Nombre del evento
  - `tipo`: Tipo de evento (Raid, Dungeon, PvP, Social, etc.)
  - `fecha`: Fecha y hora (formato: YYYY-MM-DD HH:MM)
  - `descripcion`: Descripción del evento
  - `template`: Nombre del template a usar (opcional, debe coincidir con un template existente)
  - `canal`: Canal donde se publicará el evento (opcional)
  - `discord_event`: `true` para crear también el evento oficial de Discord (Guild Scheduled Event) si está habilitado globalmente
  - `repeat_days`: Cada cuántos días se repite el evento (0 o vacío = no se repite)

- `/delete_event` - Eliminar un evento existente (borra el mensaje y archiva/cierra el hilo asociado)
  - `id`: ID del evento

- `/remind_event` - Enviar recordatorio inmediato en el hilo del evento (o en el canal si no hay hilo)
  - `id`: ID del evento

- `/list_events` - Listar todos los eventos activos

- `/config` - Mostrar configuración actual del bot (roles por defecto, zona horaria, etc.)

## 🌐 Panel Web

### Acceso

Navega a `http://localhost:8080` (o la IP de tu dispositivo si accedes desde otro equipo en la LAN)

Credenciales por defecto (cámbialas en `.env`):
- Usuario: `admin`
- Contraseña: `admin123`

### Funcionalidades

- **Dashboard**: Vista de eventos activos
- **Crear Evento**: Formulario para crear eventos desde el navegador
- **Ver Eventos**: Lista completa de todos los eventos (incluidos cancelados y completados)
- **Detalles de Evento**: Ver inscripciones, confirmar participantes y ver el hilo asociado
- **Templates**: Crear, editar, clonar, importar y exportar templates
- **Limpieza de cancelados**: Botón para eliminar del sistema todos los eventos con estado *cancelled*
- **Configuración**: Ver ajustes actuales del bot

## 🔧 Configuración Avanzada

### Personalizar Roles

Edita la variable `DEFAULT_ROLES` en `.env` (roles usados cuando creas un evento sin template):

```env
DEFAULT_ROLES=[{"name":"Tank","emoji":"🛡️","limit":2},{"name":"Healer","emoji":"💚","limit":3},{"name":"DPS","emoji":"⚔️","limit":8},{"name":"Support","emoji":"🔮","limit":2}]
```

Notas:
- `limit` define el máximo de jugadores por rol.
- Si `limit` es `0` o se omite, ese rol no tiene límite de jugadores (se muestra como `∞` / "Sin límite").
- Los límites globales y por rol también pueden configurarse en los templates desde el panel web.

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

### Eventos oficiales de Discord

Controla si el bot puede crear **Guild Scheduled Events** cuando usas `/create_event` con `discord_event: true`:

```env
ENABLE_DISCORD_EVENTS=true
```

- `true`: permite crear eventos oficiales de Discord.
- `false`: ignora la opción `discord_event` en los comandos y desde el panel web.

## 🖥️ Instalación en Raspberry Pi

La guía detallada de despliegue en Raspberry Pi (incluyendo `systemd`, estructura de carpetas y troubleshooting) se encuentra en:

`docs/rapsberry-docs.md`

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
- 🎯 Definir cupos específicos por rol (con desglose de inscripciones por clase)
- ♾️ Soportar límites opcionales: `max_participants` y `limit` de rol en `0` = sin límite
- 🎨 Emojis personalizados para cada elemento (incluyendo emojis personalizados de Discord en los botones)
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
./build.sh
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
