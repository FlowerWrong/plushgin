# Plush template for gin

Package plushgin is a plush template renderer that can be used with the [Gin web framework](https://github.com/gin-gonic/gin).
It uses the [plush template library](https://github.com/gobuffalo/plush).

## Basic Usage

```golang
router.HTMLRender = plushgin.Default()
// or
router.HTMLRender = plushgin.New(plushgin.RenderOptions{
	TemplateDir: 		"templates",
	ContentType: 		"text/html; charset=utf-8",
	MaxCacheEnties: 128,
})

router.GET("/", func(c *gin.Context) {
	c.HTML(200, "index.html", *plush.NewContextWith(map[string]interface{}{
		"name": "value",
		"a": 1,
		"bool": true,
	}))
})
```

## TODO

* [ ] layout
