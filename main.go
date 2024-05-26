package main

import (
  "log"
  "errors"
  "html/template"
  "net/http"
  "io"
  "io/fs"
  "embed"
  "strings"

  "github.com/labstack/echo/v4"
)

//go:embed resources
var embededFiles embed.FS

// Define the template registry struct
type TemplateRegistry struct {
  templates map[string]*template.Template
  baseTemplatePaths map[string]string
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  tmpl, ok := t.templates[name]
  baseTemplatePath := t.baseTemplatePaths[name]
  if !ok {
    err := errors.New("Template not found -> " + name)
    log.Println(err)
    return err
  }
  err := tmpl.ExecuteTemplate(w, baseTemplatePath, data)
  if err != nil {
    log.Println(err)
    return err
  }

  return nil
}

func SetupTemplateRegistry(parentPath string) *TemplateRegistry {
  files, err := fs.Glob(embededFiles, parentPath)
  if err != nil {
      log.Println(err)
  }

  log.Printf("found %v files\n", len(files))
  templates := make(map[string]*template.Template)
  baseTemplatePaths := make(map[string]string)
  for _, filePath := range files {
    log.Printf("processing template file %v\n", filePath)
    tmpl, err := template.ParseFS(embededFiles, filePath)
    if err != nil {
        log.Println(err)
        continue
    }

    // by default render partial template
    baseTemplatePath := "body"
    // get basepath definition
    var renderedContent strings.Builder
    err = tmpl.ExecuteTemplate(&renderedContent, "base_template_path", nil)
    if err != nil {
        log.Println(err)
        log.Println("skipping...")
    } else {
        templatePath := renderedContent.String()
        res,_ := template.ParseFS(embededFiles, filePath, templatePath)
        if err != nil {
            log.Println(err)
            continue
        }

        tmpl = res
        baseTemplatePath = templatePath
    }

    templates[filePath] = tmpl
    baseTemplatePaths[filePath] = baseTemplatePath
  }

  return &TemplateRegistry{
    templates: templates,
    baseTemplatePaths: baseTemplatePaths,
  }
}

func main() {
  // Echo instance
  e := echo.New()

  e.Renderer = SetupTemplateRegistry("resources/view/*")
  // Route => handler
  e.GET("/", func (c echo.Context) error {
    return c.Render(http.StatusOK, "resources/view/home.html", nil)
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
