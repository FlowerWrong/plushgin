package plushgin

import (
	"github.com/gin-gonic/gin"
	"github.com/gobuffalo/plush"
)

// NewContext create a plush.Context
func NewContext(c gin.H) plush.Context {
	return *plush.NewContextWith(c)
}
