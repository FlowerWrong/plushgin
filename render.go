package plushgin

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin/render"
	"github.com/gobuffalo/plush"
)

const htmlContentType = "text/html; charset=utf-8"

// RenderOptions is used to configure the renderer.
type RenderOptions struct {
	TemplateDir string
	ContentType string
}

// Plush2Render is a custom Gin template renderer using plush.
type Plush2Render struct {
	Options  *RenderOptions
	Template *plush.Template
	Context  plush.Context
}

// New creates a new Plush2Render instance with custom Options.
func New(options RenderOptions) *Plush2Render {
	return &Plush2Render{
		Options: &options,
	}
}

// Default creates a Plush2Render instance with default options.
func Default() *Plush2Render {
	return New(RenderOptions{
		TemplateDir: "templates",
		ContentType: htmlContentType,
	})
}

// Instance should return a new Plush2Render struct per request and prepare
// the template by either loading it from disk or using plush's cache.
func (p Plush2Render) Instance(name string, data interface{}) render.Render {
	var template *plush.Template
	filename := path.Join(p.Options.TemplateDir, name)

	buf, err := ioutil.ReadFile(filename)
	if err != nil {
		panic(err)
	}
	template, err = plush.NewTemplate(string(buf))
	if err != nil {
		panic(err)
	}

	return Plush2Render{
		Template: template,
		Context:  data.(plush.Context),
		Options:  p.Options,
	}
}

// Render should render the template to the response.
func (p Plush2Render) Render(w http.ResponseWriter) error {
	p.WriteContentType(w)
	rendered, err := plush.Render(p.Template.Input, &p.Context)
	w.Write([]byte(rendered))
	return err
}

// WriteContentType should add the Content-Type header to the response when not set yet.
func (p Plush2Render) WriteContentType(w http.ResponseWriter) {
	header := w.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = []string{p.Options.ContentType}
	}
}
