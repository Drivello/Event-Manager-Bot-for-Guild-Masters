package web

import (
	"discord-event-bot/internal/storage"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// RegisterTemplateRoutes registra las rutas de la API de templates
func RegisterTemplateRoutes(router *gin.RouterGroup) {
	// API REST para templates
	router.GET("/api/templates", handleGetAllTemplates)
	router.GET("/api/templates/:name", handleGetTemplate)
	router.POST("/api/templates", handleCreateTemplate)
	router.PUT("/api/templates/:name", handleUpdateTemplate)
	router.DELETE("/api/templates/:name", handleDeleteTemplate)
	router.POST("/api/templates/:name/clone", handleCloneTemplate)
	router.GET("/api/templates/:name/export", handleExportTemplate)
	router.POST("/api/templates/import", handleImportTemplate)

	// Páginas web para gestión de templates
	router.GET("/templates", handleTemplatesPage)
	router.GET("/templates/create", handleCreateTemplatePage)
	router.GET("/templates/:name/edit", handleEditTemplatePage)
}

// handleGetAllTemplates retorna todos los templates
func handleGetAllTemplates(c *gin.Context) {
	templates := storage.Templates.GetAllTemplates()
	c.JSON(http.StatusOK, gin.H{
		"templates": templates,
		"count":     len(templates),
	})
}

// handleGetTemplate retorna un template específico
func handleGetTemplate(c *gin.Context) {
	name := c.Param("name")
	template, err := storage.Templates.GetTemplate(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template no encontrado"})
		return
	}
	c.JSON(http.StatusOK, template)
}

// handleCreateTemplate crea un nuevo template
func handleCreateTemplate(c *gin.Context) {
	var template storage.EventTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Agregar timestamps
	now := time.Now().Format(time.RFC3339)
	template.CreatedAt = now
	template.UpdatedAt = now

	if err := storage.Templates.SaveTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":  "Template creado exitosamente",
		"template": template,
	})
}

// handleUpdateTemplate actualiza un template existente
func handleUpdateTemplate(c *gin.Context) {
	name := c.Param("name")

	// Verificar que el template existe
	existingTemplate, err := storage.Templates.GetTemplate(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Template no encontrado"})
		return
	}

	var template storage.EventTemplate
	if err := c.ShouldBindJSON(&template); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Datos inválidos: " + err.Error()})
		return
	}

	// Mantener el nombre original y createdAt
	template.Name = name
	template.CreatedAt = existingTemplate.CreatedAt
	template.UpdatedAt = time.Now().Format(time.RFC3339)

	if err := storage.Templates.SaveTemplate(&template); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message":  "Template actualizado exitosamente",
		"template": template,
	})
}

// handleDeleteTemplate elimina un template
func handleDeleteTemplate(c *gin.Context) {
	name := c.Param("name")

	if err := storage.Templates.DeleteTemplate(name); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Template eliminado exitosamente",
	})
}

// handleCloneTemplate clona un template existente
func handleCloneTemplate(c *gin.Context) {
	sourceName := c.Param("name")

	var req struct {
		NewName string `json:"new_name" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Nombre del nuevo template requerido"})
		return
	}

	if err := storage.Templates.CloneTemplate(sourceName, req.NewName); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Template clonado exitosamente",
		"name":    req.NewName,
	})
}

// handleExportTemplate exporta un template a JSON
func handleExportTemplate(c *gin.Context) {
	name := c.Param("name")

	data, err := storage.Templates.ExportTemplate(name)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.Header("Content-Disposition", "attachment; filename="+name+".json")
	c.Data(http.StatusOK, "application/json", data)
}

// handleImportTemplate importa un template desde JSON
func handleImportTemplate(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Archivo no proporcionado"})
		return
	}

	// Leer archivo
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo archivo"})
		return
	}
	defer f.Close()

	data, err := ioutil.ReadAll(f)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error leyendo contenido"})
		return
	}

	if err := storage.Templates.ImportTemplate(data); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Template importado exitosamente",
	})
}

// handleTemplatesPage muestra la página de gestión de templates
func handleTemplatesPage(c *gin.Context) {
	templates := storage.Templates.GetAllTemplates()
	c.HTML(http.StatusOK, "templates.html", gin.H{
		"title":     "Gestión de Templates",
		"templates": templates,
	})
}

// handleCreateTemplatePage muestra el formulario de creación de template
func handleCreateTemplatePage(c *gin.Context) {
	c.HTML(http.StatusOK, "template_editor.html", gin.H{
		"title":    "Crear Template",
		"mode":     "create",
		"template": nil,
	})
}

// handleEditTemplatePage muestra el formulario de edición de template
func handleEditTemplatePage(c *gin.Context) {
	name := c.Param("name")
	template, err := storage.Templates.GetTemplate(name)
	if err != nil {
		c.HTML(http.StatusNotFound, "error.html", gin.H{
			"title": "Error",
			"error": "Template no encontrado",
		})
		return
	}

	c.HTML(http.StatusOK, "template_editor.html", gin.H{
		"title":    "Editar Template: " + name,
		"mode":     "edit",
		"template": template,
	})
}
