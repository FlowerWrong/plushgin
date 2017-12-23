package plushgin

import (
	"io/ioutil"
	"net/http"
	"path"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/render"
	"github.com/gobuffalo/plush"
)

const htmlContentType = "text/html; charset=utf-8"

// RenderOptions is used to configure the renderer.
type RenderOptions struct {
	TemplateDir    string
	ContentType    string
	MaxCacheEnties int
}

// Plush2Render is a custom Gin template renderer using plush.
type Plush2Render struct {
	Options *RenderOptions
	Name    string
	Context plush.Context
	cache   *renderCache
}

// New creates a new Plush2Render instance with custom Options.
func New(options RenderOptions) *Plush2Render {
	return &Plush2Render{
		Options: &options,
		cache:   newRenderCache(options.MaxCacheEnties),
	}
}

// Default creates a Plush2Render instance with default options.
func Default() *Plush2Render {
	return New(RenderOptions{
		TemplateDir:    "templates",
		ContentType:    htmlContentType,
		MaxCacheEnties: 128,
	})
}

// Instance should return a new Plush2Render struct per request and prepare
// the template by either loading it from disk or using plush's cache.
func (p Plush2Render) Instance(name string, data interface{}) render.Render {
	return Plush2Render{
		Context: data.(plush.Context),
		Options: p.Options,
		cache:   p.cache,
		Name:    name,
	}
}

// Render should render the template to the response.
func (p Plush2Render) Render(w http.ResponseWriter) error {
	var err error
	var renderedStr string

	rendered := p.cache.Get(p.Name, p.Context)
	if rendered == nil || gin.Mode() == "debug" {
		filename := path.Join(p.Options.TemplateDir, p.Name)
		var buf []byte
		buf, err = ioutil.ReadFile(filename)
		if err != nil {
			panic(err)
		}
		renderedStr, err = plush.Render(string(buf), &p.Context)
		if err != nil {
			panic(err)
		}
		rendered = []byte(renderedStr)
		p.cache.Add(p.Name, p.Context, rendered)
	}
	p.WriteContentType(w)
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
