# 🎨 Guía de Templates para Eventos MMO

## 📋 Índice
- [Introducción](#introducción)
- [Conceptos Básicos](#conceptos-básicos)
- [Uso desde Discord](#uso-desde-discord)
- [Uso desde Panel Web](#uso-desde-panel-web)
- [Estructura de Templates](#estructura-de-templates)
- [Creación de Templates](#creación-de-templates)
- [Gestión de Templates](#gestión-de-templates)
- [API REST](#api-rest)
- [Ejemplos](#ejemplos)
- [Solución de Problemas](#solución-de-problemas)

---

## 🎯 Introducción

El sistema de templates permite crear modelos reutilizables para eventos MMO, definiendo de antemano:
- **Roles** disponibles (Tank, DPS, Healer, etc.)
- **Cupos** por rol
- **Clases/Especializaciones** dentro de cada rol
- **Emojis** personalizados para cada rol y clase
- **Límites** de participantes totales

Esto facilita la organización de eventos recurrentes sin tener que configurar manualmente los roles cada vez.

---

## 📚 Conceptos Básicos

### Template
Un **template** es una plantilla que define la estructura de un evento:
```json
{
  "name": "Raid 20 jugadores",
  "icon": "⚔️",
  "max_participants": 20,
  "description": "Template estándar para raids de 20 jugadores",
  "roles": [...]
}
```

### Rol
Un **rol** representa una función dentro del evento (Tank, DPS, Support):
```json
{
  "name": "Tank",
  "emoji": "🛡️",
  "limit": 4,
  "classes": [...]
}
```

### Clase
Una **clase** es una especialización dentro de un rol:
```json
{
  "name": "Paladin",
  "emoji": "⚔️",
  "description": "Tank sagrado"
}
```

---

## 🎮 Uso desde Discord

### Crear Evento con Template

Usa el comando `/create_event` con el parámetro `template`:

```
/create_event 
  nombre: Raid Semanal
  tipo: Raid
  fecha: 2024-12-20 20:00
  descripcion: Raid mítica del viernes
  template: Raid 20 jugadores
```

### Listar Templates Disponibles

Los templates disponibles se pueden consultar desde el panel web en `/templates`.

### Inscripción a Eventos

Cuando un evento usa un template con clases:
1. Haz clic en el botón del rol deseado (ej: 🛡️ Tank)
2. El sistema registrará tu inscripción
3. Los organizadores pueden ver qué clase elegiste

---

## 🌐 Uso desde Panel Web

### Acceder a Templates

1. Inicia sesión en el panel web: `http://localhost:8080`
2. Navega a **Templates** en el menú
3. Verás todos los templates disponibles

### Crear Template desde Web

1. Click en **"➕ Crear Nuevo Template"**
2. Completa los datos básicos:
   - Nombre del template
   - Icono (emoji)
   - Máximo de participantes
   - Descripción
3. Agrega roles con **"➕ Agregar Rol"**
4. Para cada rol, define:
   - Nombre y emoji
   - Límite de jugadores
   - Clases disponibles (opcional)
5. Visualiza en tiempo real en el panel de **Vista Previa**
6. Click en **"💾 Guardar Template"**

### Crear Evento con Template

1. Ve a **"Crear Nuevo Evento"**
2. Selecciona un template del dropdown **"Template (Opcional)"**
3. Completa los datos del evento
4. El evento heredará automáticamente los roles y configuración del template

---

## 🏗️ Estructura de Templates

### Formato JSON

```json
{
  "name": "Nombre del Template",
  "icon": "🎯",
  "max_participants": 20,
  "description": "Descripción del template",
  "allow_multi_signup": false,
  "roles": [
    {
      "name": "Tank",
      "emoji": "🛡️",
      "limit": 4,
      "classes": [
        {
          "name": "Paladin",
          "emoji": "⚔️",
          "description": "Tank sagrado"
        },
        {
          "name": "Warrior",
          "emoji": "🪓",
          "description": "Guerrero defensor"
        }
      ]
    },
    {
      "name": "DPS",
      "emoji": "🏹",
      "limit": 12,
      "classes": [
        {
          "name": "Hunter",
          "emoji": "🎯"
        },
        {
          "name": "Mage",
          "emoji": "❄️"
        }
      ]
    }
  ],
  "created_at": "2024-01-15T10:00:00Z",
  "updated_at": "2024-01-15T10:00:00Z"
}
```

### Formato YAML

```yaml
name: Raid 20 jugadores
icon: ⚔️
max_participants: 20
description: Template estándar para raids
allow_multi_signup: false
roles:
  - name: Tank
    emoji: 🛡️
    limit: 4
    classes:
      - name: Paladin
        emoji: ⚔️
        description: Tank sagrado
      - name: Warrior
        emoji: 🪓
        description: Guerrero defensor
  - name: DPS
    emoji: 🏹
    limit: 12
    classes:
      - name: Hunter
        emoji: 🎯
      - name: Mage
        emoji: ❄️
```

---

## ✨ Creación de Templates

### Templates por Defecto

El sistema incluye 3 templates predefinidos:

1. **Raid 20 jugadores** - Para raids estándar
2. **Dungeon 5 jugadores** - Para mazmorras
3. **PvP Battleground** - Para campos de batalla de 40 jugadores

### Crear Template Personalizado

#### Opción 1: Desde el Editor Web

Usa el editor visual en `/templates/create` que incluye:
- Formulario interactivo
- Vista previa en tiempo real
- Validación automática

#### Opción 2: Importar JSON/YAML

1. Crea un archivo JSON o YAML con la estructura del template
2. Ve a `/templates`
3. Click en **"📥 Importar Template"**
4. Selecciona tu archivo

#### Opción 3: Clonar Template Existente

1. Ve a `/templates`
2. En el template que quieres clonar, click en **"📋 Clonar"**
3. Ingresa el nombre del nuevo template
4. Edita el clon según necesites

---

## 🔧 Gestión de Templates

### Editar Template

1. Ve a `/templates`
2. Click en **"✏️ Editar"** en el template deseado
3. Modifica los campos necesarios
4. Guarda los cambios

### Exportar Template

Para compartir o respaldar un template:

1. Ve a `/templates`
2. Click en **"💾 Exportar"** en el template
3. Se descargará un archivo JSON

### Eliminar Template

⚠️ **Precaución**: Eliminar un template no afecta eventos ya creados.

1. Ve a `/templates`
2. Click en **"🗑️"** en el template
3. Confirma la eliminación

### Ubicación de Archivos

Los templates se almacenan en:
```
data/templates/
  ├── Raid_20_jugadores.json
  ├── Dungeon_5_jugadores.json
  └── PvP_Battleground.json
```

---

## 🔌 API REST

### Endpoints Disponibles

#### Listar Templates
```http
GET /api/templates
```

**Respuesta:**
```json
{
  "templates": [...],
  "count": 3
}
```

#### Obtener Template
```http
GET /api/templates/:name
```

#### Crear Template
```http
POST /api/templates
Content-Type: application/json

{
  "name": "Mi Template",
  "icon": "🎯",
  "max_participants": 10,
  "roles": [...]
}
```

#### Actualizar Template
```http
PUT /api/templates/:name
Content-Type: application/json

{
  "icon": "🎮",
  "max_participants": 15,
  ...
}
```

#### Eliminar Template
```http
DELETE /api/templates/:name
```

#### Clonar Template
```http
POST /api/templates/:name/clone
Content-Type: application/json

{
  "new_name": "Copia de Template"
}
```

#### Exportar Template
```http
GET /api/templates/:name/export
```

#### Importar Template
```http
POST /api/templates/import
Content-Type: multipart/form-data

file: template.json
```

---

## 💡 Ejemplos

### Template para Raid Mítica 10 Jugadores

```json
{
  "name": "Raid Mítica 10",
  "icon": "⚔️",
  "max_participants": 10,
  "description": "Raid mítica de 10 jugadores",
  "roles": [
    {
      "name": "Tank",
      "emoji": "🛡️",
      "limit": 2,
      "classes": [
        {"name": "Protection Warrior", "emoji": "🪓"},
        {"name": "Guardian Druid", "emoji": "🐻"}
      ]
    },
    {
      "name": "Healer",
      "emoji": "💚",
      "limit": 2,
      "classes": [
        {"name": "Holy Priest", "emoji": "⛪"},
        {"name": "Restoration Druid", "emoji": "🌿"}
      ]
    },
    {
      "name": "DPS",
      "emoji": "🏹",
      "limit": 6,
      "classes": [
        {"name": "Hunter", "emoji": "🎯"},
        {"name": "Mage", "emoji": "❄️"},
        {"name": "Rogue", "emoji": "🗡️"}
      ]
    }
  ]
}
```

### Template para Arena 3v3

```json
{
  "name": "Arena 3v3",
  "icon": "⚔️",
  "max_participants": 3,
  "description": "Equipo de arena 3v3",
  "roles": [
    {
      "name": "DPS",
      "emoji": "🗡️",
      "limit": 2,
      "classes": []
    },
    {
      "name": "Healer",
      "emoji": "💚",
      "limit": 1,
      "classes": []
    }
  ]
}
```

---

## 🔍 Solución de Problemas

### Template no aparece en Discord

**Problema**: El template no se muestra en el comando `/create_event`

**Solución**: 
- Los templates se seleccionan por nombre exacto
- Verifica que el template existe en `/templates`
- Asegúrate de escribir el nombre correctamente

### Error al guardar template

**Problema**: "Error: la suma de límites excede el máximo"

**Solución**:
- Verifica que la suma de límites de todos los roles no exceda `max_participants`
- Ejemplo: Si `max_participants: 10`, los límites de roles deben sumar ≤ 10

### Template no se carga al iniciar

**Problema**: Los templates no aparecen después de reiniciar el bot

**Solución**:
- Verifica que los archivos existan en `data/templates/`
- Revisa los logs del bot para errores de parseo
- Valida el formato JSON/YAML del archivo

### Clases no se muestran en Discord

**Problema**: Las clases definidas no aparecen en el mensaje del evento

**Solución**:
- Las clases se muestran solo cuando hay inscripciones confirmadas
- Verifica que el template tenga clases definidas en los roles
- Actualiza el mensaje del evento después de confirmar inscripciones

---

## 🚀 Extensión del Sistema

### Agregar Nuevos Campos a Templates

Para extender la funcionalidad de templates:

1. Actualiza la estructura en `internal/storage/templates.go`:
```go
type EventTemplate struct {
    // ... campos existentes
    MinLevel int `json:"min_level,omitempty"`
}
```

2. Actualiza el editor web en `template_editor.html`

3. Actualiza la lógica de creación de eventos en `events.go`

### Crear Validaciones Personalizadas

Edita `validateTemplate()` en `templates.go`:

```go
func (ts *TemplateStore) validateTemplate(template *EventTemplate) error {
    // Validaciones existentes...
    
    // Nueva validación
    if template.MinLevel < 1 || template.MinLevel > 80 {
        return fmt.Errorf("nivel mínimo debe estar entre 1 y 80")
    }
    
    return nil
}
```

---

## 📞 Soporte

Para más ayuda:
- Revisa los logs del bot: `journalctl -u discord-bot.service -f`
- Consulta el código fuente en `internal/storage/templates.go`
- Abre un issue en el repositorio del proyecto

---

**Última actualización**: Noviembre 2024  
**Versión del sistema**: 2.0
