package main

import (
  "log"
  "errors"
  "html/template"
  "net/http"
  "io"
  "io/fs"
  "embed"

  "github.com/labstack/echo/v4"
)

//go:embed resources
var embededFiles embed.FS

// Define the template registry struct
type TemplateRegistry struct {
  templates map[string]*template.Template
  baseTemplatePath string
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  tmpl, ok := t.templates[name]
  if !ok {
    err := errors.New("Template not found -> " + name)
    log.Println(err)
    return err
  }
  return tmpl.ExecuteTemplate(w, t.baseTemplatePath, data)
}

func SetupTemplateRegistry(parentPath string, baseTemplatePath string) *TemplateRegistry {
  files, err := fs.Glob(embededFiles, parentPath)
  if err != nil {
      log.Println(err)
  }

  log.Printf("found %v files\n", len(files))
  templates := make(map[string]*template.Template)
  for _, filePath := range files {
    templates[filePath] = template.Must(template.ParseFS(embededFiles, filePath, baseTemplatePath))
  }

  return &TemplateRegistry{
    templates: templates,
    baseTemplatePath: baseTemplatePath,
  }
}

func main() {
  // Echo instance
  e := echo.New()

  e.Renderer = SetupTemplateRegistry("resources/view/*", "resources/view/base.html")
  // Route => handler
  e.GET("/", func (c echo.Context) error {
    return c.Render(http.StatusOK, "resources/view/home.html", map[string]interface{}{})
  })
  e.GET("/about", func (c echo.Context) error {
    return c.Render(http.StatusOK, "resources/view/about.html", map[string]interface{}{
      "name": "About",
      "msg": "All about azophy!",
    })
  })

  // Start the Echo server
  e.Logger.Fatal(e.Start(":3000"))
}
