package main

import (
  "errors"
  "html/template"
  "net/http"
  "io"
  "path/filepath"

  "github.com/labstack/echo/v4"
)

// Define the template registry struct
type TemplateRegistry struct {
  templates map[string]*template.Template
}

// Implement e.Renderer interface
func (t *TemplateRegistry) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
  tmpl, ok := t.templates[name]
  if !ok {
    err := errors.New("Template not found -> " + name)
    return err
  }
  return tmpl.ExecuteTemplate(w, "base.html", data)
}

func SetupTemplateRegistry(parentPath string, baseTemplatePath string) *TemplateRegistry {
  files, err := filepath.Glob("testFolder/*.go")
  if err != nil {
      log.Fatal(err)
  }

  for _, file := range files {
      fmt.Println(file)
  }

  templates := make(map[string]*template.Template)
  templates["home.html"] = template.Must(template.ParseFiles("resources/view/home.html", baseTemplatePath))
  templates["about.html"] = template.Must(template.ParseFiles("resources/view/about.html", baseTemplatePath))
  return &TemplateRegistry{
    templates: templates,
  }
}

func main() {
  // Echo instance
  e := echo.New()

  // Instantiate a template registry with an array of template set
  // Ref: https://gist.github.com/rand99/808e6e9702c00ce64803d94abff65678

  e.Renderer = SetupTemplateRegistry("resources/*", "resources/view/base.html")

  // Route => handler
  e.GET("/", func (c echo.Context) error {
    // Please note the the second parameter "about.html" is the template name and should
    // be equal to one of the keys in the TemplateRegistry array defined in main.go
    return c.Render(http.StatusOK, "home.html", map[string]interface{}{
      "name": "About",
      "msg": "All about Boatswain!",
    })
  })
  e.GET("/about", func (c echo.Context) error {
    // Please note the the second parameter "about.html" is the template name and should
    // be equal to one of the keys in the TemplateRegistry array defined in main.go
    return c.Render(http.StatusOK, "about.html", map[string]interface{}{
      "name": "About",
      "msg": "All about Boatswain!",
    })
  })

  // Start the Echo server
  e.Logger.Fatal(e.Start(":3000"))
}
