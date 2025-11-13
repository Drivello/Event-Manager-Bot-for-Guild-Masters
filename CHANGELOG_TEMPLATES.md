# 📝 Changelog - Sistema de Templates

## Versión 2.0 - Sistema de Templates Personalizados

### 🎉 Nuevas Características

#### Backend (Go)

**Nuevos Archivos:**
- `internal/storage/templates.go` - Sistema completo de gestión de templates
  - Modelo de datos `EventTemplate`, `TemplateRole`, `TemplateClass`
  - Persistencia en JSON y YAML
  - Validaciones automáticas
  - Importar/Exportar funcionalidad
  - Sistema de clonado de templates
  
- `internal/web/templates_api.go` - API REST para templates
  - `GET /api/templates` - Listar todos los templates
  - `GET /api/templates/:name` - Obtener template específico
  - `POST /api/templates` - Crear nuevo template
  - `PUT /api/templates/:name` - Actualizar template
  - `DELETE /api/templates/:name` - Eliminar template
  - `POST /api/templates/:name/clone` - Clonar template
  - `GET /api/templates/:name/export` - Exportar a JSON
  - `POST /api/templates/import` - Importar desde JSON

**Archivos Modificados:**
- `internal/storage/events.go`
  - Extendido modelo `Event` con campo `TemplateName` y `MaxParticipants`
  - Extendido modelo `RoleSignup` con array de `Classes`
  - Nuevo modelo `ClassInfo` para clases/especializaciones
  - Extendido modelo `Signup` con campo `Class`
  - Nueva función `CreateEventFromTemplate()` - Crear eventos desde templates
  - Nueva función `AddSignupWithClass()` - Inscripciones con clase específica

- `internal/discord/handler.go`
  - Agregado parámetro `template` al comando `/create_event`
  - Actualizado `handleCreateEvent()` para soportar templates
  - Actualizado `buildSignupsText()` para mostrar desglose por clases

- `internal/web/server.go`
  - Registradas rutas de templates con `RegisterTemplateRoutes()`
  - Actualizado `handleCreateEventPage()` para pasar lista de templates
  - Actualizado `handleCreateEventPost()` para soportar creación desde templates

- `cmd/main.go`
  - Agregada inicialización de `storage.InitTemplateStore()`

#### Frontend (HTML/JavaScript)

**Nuevos Templates HTML:**
- `internal/web/templates/templates.html` - Página de gestión de templates
  - Grid responsive de templates
  - Acciones: Editar, Clonar, Exportar, Eliminar
  - Importar templates desde archivo
  - Vista de estadísticas por template

- `internal/web/templates/template_editor.html` - Editor visual de templates
  - Formulario interactivo para crear/editar templates
  - Gestión dinámica de roles y clases
  - Vista previa en tiempo real estilo Discord
  - Validación de formularios
  - Soporte para crear y editar templates

**Templates HTML Modificados:**
- `internal/web/templates/create_event.html`
  - Agregado selector de templates
  - Dropdown con templates disponibles
  - Opción de usar configuración por defecto

#### Documentación

**Nuevos Archivos:**
- `TEMPLATES_GUIDE.md` - Guía completa del sistema de templates
  - Conceptos básicos
  - Uso desde Discord y Panel Web
  - Estructura de templates (JSON/YAML)
  - Ejemplos prácticos
  - API REST completa
  - Solución de problemas
  - Guía de extensión

- `template_example.json` - Ejemplo completo de template
  - Raid de 15 jugadores
  - 4 roles diferentes
  - Múltiples clases por rol con descripciones

- `CHANGELOG_TEMPLATES.md` - Este archivo

**Archivos Modificados:**
- `README.md`
  - Actualizada sección de características
  - Agregada sección "Sistema de Templates"
  - Referencias a documentación de templates

### 🔧 Mejoras Técnicas

#### Persistencia
- Soporte dual JSON/YAML para templates
- Carga automática al iniciar el bot
- Creación de templates por defecto si no existen
- Sanitización de nombres de archivo

#### Validaciones
- Validación de límites de roles vs max_participants
- Validación de campos requeridos
- Prevención de templates duplicados
- Manejo robusto de errores

#### API REST
- Endpoints RESTful completos
- Autenticación mediante BasicAuth
- Respuestas JSON estructuradas
- Manejo de errores HTTP apropiado

### 📦 Templates Incluidos por Defecto

1. **Raid 20 jugadores**
   - 4 Tanks (Paladin, Warrior, Death Knight)
   - 12 DPS (Hunter, Mage, Rogue, Warlock)
   - 4 Support (Priest, Druid, Shaman)

2. **Dungeon 5 jugadores**
   - 1 Tank (Paladin, Warrior)
   - 3 DPS (Hunter, Mage, Rogue)
   - 1 Healer (Priest, Druid)

3. **PvP Battleground**
   - 15 Melee DPS (Warrior, Rogue, Death Knight)
   - 15 Ranged DPS (Hunter, Mage, Warlock)
   - 10 Healer (Priest, Druid, Shaman)

### 🎯 Casos de Uso

#### Para Organizadores de Eventos
- Crear templates para raids recurrentes
- Definir composiciones específicas de grupo
- Reutilizar configuraciones probadas
- Compartir templates entre guilds

#### Para Administradores
- Gestionar templates desde panel web
- Importar templates de otras comunidades
- Exportar templates para respaldo
- Clonar y modificar templates existentes

#### Para Jugadores
- Ver clases disponibles para cada rol
- Inscribirse con clase específica
- Mejor visibilidad de composición del grupo

### 🔄 Compatibilidad

#### Retrocompatibilidad
- ✅ Eventos existentes siguen funcionando sin cambios
- ✅ Configuración por defecto en `.env` se mantiene
- ✅ Comandos Discord existentes sin modificaciones obligatorias
- ✅ Panel web existente completamente funcional

#### Migración
- No se requiere migración de datos
- Templates son opcionales
- Sistema funciona con y sin templates

### 🚀 Próximas Mejoras Sugeridas

#### Funcionalidades Futuras
- [ ] Selector de clase en inscripción Discord (dropdown)
- [ ] Límites por clase individual
- [ ] Templates con requisitos (ilvl, logros, etc.)
- [ ] Estadísticas de uso de templates
- [ ] Compartir templates públicamente
- [ ] Versiones de templates
- [ ] Plantillas de mensajes personalizados

#### Optimizaciones
- [ ] Cache de templates en memoria
- [ ] Compresión de archivos de templates
- [ ] Búsqueda y filtrado de templates
- [ ] Tags/categorías para templates

### 📊 Estructura de Archivos

```
Guild-Master/
├── data/
│   ├── events/          # Eventos (sin cambios)
│   └── templates/       # ⭐ NUEVO: Templates
│       ├── Raid_20_jugadores.json
│       ├── Dungeon_5_jugadores.json
│       └── PvP_Battleground.json
├── internal/
│   ├── storage/
│   │   ├── events.go    # ✏️ MODIFICADO
│   │   └── templates.go # ⭐ NUEVO
│   ├── discord/
│   │   └── handler.go   # ✏️ MODIFICADO
│   └── web/
│       ├── server.go         # ✏️ MODIFICADO
│       ├── templates_api.go  # ⭐ NUEVO
│       └── templates/
│           ├── templates.html        # ⭐ NUEVO
│           ├── template_editor.html  # ⭐ NUEVO
│           └── create_event.html     # ✏️ MODIFICADO
├── cmd/
│   └── main.go          # ✏️ MODIFICADO
├── TEMPLATES_GUIDE.md   # ⭐ NUEVO
├── template_example.json # ⭐ NUEVO
└── README.md            # ✏️ MODIFICADO
```

### 🐛 Bugs Conocidos

#### Lint Warnings
- Warnings de JavaScript en `template_editor.html` línea 271
  - **Causa**: Sintaxis de Go templates dentro de JavaScript
  - **Impacto**: Solo warnings del IDE, el código funciona correctamente
  - **Estado**: Esperado y no requiere corrección

### 👥 Créditos

Sistema de templates diseñado e implementado para mejorar la gestión de eventos MMO en Discord, con enfoque en usabilidad y extensibilidad.

### 📞 Soporte

Para reportar bugs o sugerir mejoras al sistema de templates:
1. Revisa `TEMPLATES_GUIDE.md` para documentación completa
2. Verifica logs del bot: `journalctl -u discord-bot -f`
3. Abre un issue en el repositorio

---

**Fecha de Release**: Noviembre 2024  
**Versión**: 2.0.0  
**Compatibilidad**: Go 1.21+, Discord API v10
