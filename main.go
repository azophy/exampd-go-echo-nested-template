package main

import (
  "errors"
  "html/template"
  "net/http"
  "io"

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

func main() {
  // Echo instance
  e := echo.New()

  // Instantiate a template registry with an array of template set
  // Ref: https://gist.github.com/rand99/808e6e9702c00ce64803d94abff65678
  templates := make(map[string]*template.Template)
  templates["home.html"] = template.Must(template.ParseFiles("resources/view/home.html", "resources/view/base.html"))
  templates["about.html"] = template.Must(template.ParseFiles("resources/view/about.html", "resources/view/base.html"))
  e.Renderer = &TemplateRegistry{
    templates: templates,
  }

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
