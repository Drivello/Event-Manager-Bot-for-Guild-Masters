# 🤝 Guía de Contribución

¡Gracias por tu interés en contribuir al Discord Event Bot!

## 🔧 Configuración del Entorno de Desarrollo

### Requisitos

- Go 1.21 o superior
- Git
- Token de bot de Discord para testing
- Editor de código (VS Code recomendado)

### Instalación Local

```bash
# Clonar el repositorio
git clone <tu-repositorio>
cd discord-event-bot

# Instalar dependencias
GOPROXY=https://proxy.golang.org,direct go mod tidy

# Configurar variables de entorno
cp .env.example .env
# Editar .env con tus credenciales de testing

# Compilar
./build.sh

# Ejecutar
./discord-event-bot
```

## 📁 Estructura del Proyecto

```
discord-event-bot/
├── cmd/                    # Punto de entrada
├── config/                 # Configuración y .env
├── internal/
│   ├── discord/           # Lógica del bot
│   ├── storage/           # Persistencia de datos
│   └── web/               # Servidor web y templates
├── data/                  # Datos locales (gitignored)
└── scripts/               # Scripts de utilidad
```

## 🎨 Estándares de Código

### Go

- Seguir las convenciones de Go (gofmt, golint)
- Documentar funciones exportadas
- Usar nombres descriptivos
- Manejar errores explícitamente

```go
// ✅ Bueno
func CreateEvent(name string) (*Event, error) {
    if name == "" {
        return nil, fmt.Errorf("name is required")
    }
    // ...
}

// ❌ Malo
func ce(n string) *Event {
    // Sin manejo de errores
}
```

### Commits

Usar mensajes descriptivos siguiendo el formato:

```
tipo(alcance): descripción corta

Descripción larga opcional
```

Tipos:
- `feat`: Nueva funcionalidad
- `fix`: Corrección de bug
- `docs`: Documentación
- `style`: Formato de código
- `refactor`: Refactorización
- `test`: Tests
- `chore`: Tareas de mantenimiento

Ejemplos:
```
feat(discord): agregar comando /remind_all
fix(web): corregir error en confirmación de signups
docs(readme): actualizar instrucciones de instalación
```

## 🧪 Testing

### Ejecutar Tests

```bash
go test ./...
```

### Escribir Tests

Crear archivos `*_test.go` junto al código:

```go
func TestCreateEvent(t *testing.T) {
    event := &Event{
        Name: "Test Event",
        Type: "Raid",
    }
    
    err := storage.Store.SaveEvent(event)
    assert.NoError(t, err)
}
```

## 📋 Proceso de Contribución

1. **Fork** el repositorio
2. **Crear** una rama descriptiva (`feat/nueva-funcionalidad`)
3. **Hacer** commits atómicos y descriptivos
4. **Probar** los cambios localmente
5. **Push** a tu fork
6. **Crear** un Pull Request

### Pull Request

Tu PR debe incluir:
- ✅ Descripción clara del cambio
- ✅ Motivación (qué problema resuelve)
- ✅ Tests (si aplica)
- ✅ Documentación actualizada
- ✅ Screenshots (si hay cambios visuales)

## 🐛 Reportar Bugs

Abre un issue con:
- Descripción del problema
- Pasos para reproducir
- Comportamiento esperado vs actual
- Logs relevantes
- Información del sistema (OS, versión de Go)

## 💡 Sugerir Mejoras

Abre un issue con:
- Descripción de la funcionalidad
- Casos de uso
- Beneficios
- Posibles alternativas

## 🔍 Áreas que Necesitan Ayuda

- 📝 Mejorar documentación
- 🧪 Agregar tests unitarios
- 🌐 Internacionalización (i18n)
- 🎨 Mejorar UI del panel web
- ⚡ Optimización de rendimiento
- 🔒 Mejorar seguridad

## 📚 Recursos

- [Documentación de discordgo](https://github.com/bwmarrin/discordgo)
- [Documentación de Gin](https://gin-gonic.com/docs/)
- [Effective Go](https://golang.org/doc/effective_go)
- [Discord Developer Portal](https://discord.com/developers/docs)

## 📞 Contacto

¿Preguntas? Abre un issue o contacta a los mantenedores.

---

**¡Gracias por contribuir! 🎉**
