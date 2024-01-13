package main

import (
	"net/http"

	"gof/pkg"
)

func main() {
	r := pkg.Default()
	r.Static("/assets", "./static")

	r.GET("/panic", func(c *pkg.Context) {
		names := []string{"gee"}
		c.String(http.StatusOK, names[100])
	})
	r.GET("/index", func(c *pkg.Context) {
		c.String(http.StatusOK, "this is a index page")
	})
	v1 := r.Group("/v1")
	{
		v1.GET("/hello/:name", func(c *pkg.Context) {
			c.JSON(http.StatusOK, pkg.H{
				"name":    c.Param("name"),
				"message": c.Query("msg"),
			})
		})
	}

	r.Run(":9999")
}
