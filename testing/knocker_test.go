package testing

import (
	"fmt"
	"net/http"
	"testing"
	"yee"
	"yee/middleware"
)

type p struct {
	Test string `json:"test"`
}

func TestContext(t *testing.T) {
	r := yee.New()
	r.GET("/", func(c yee.Context) error {
		return c.String(http.StatusOK, "<h1>Hello Gee</h1>")
	})
	r.Static("/assets", "dist/assets")

	r.GET("/", func(c yee.Context) (err error) {
		return c.HTMLTml(http.StatusOK, "dist/index.html")
	})

	r.POST("/test", func(c yee.Context) (err error) {
		u := new(p)
		if err := c.Bind(u); err != nil {
			return c.JSON(http.StatusOK, err.Error())
		}
		return c.JSON(http.StatusOK, u.Test)
	})

	//r.GET("/:t/:x", func(c yee.Context) (err error) {
	//	j := c.Params("t")
	//	x := c.Params("x")
	//	fmt.Println(1)
	//	return c.JSON(http.StatusOK, fmt.Sprintf("%s-%s", j, x))
	//	//return c.JSON(http.StatusOK, map[string]interface{}{"name":"henry","age": 27})
	//})

	r.GET("/b/:name", func(c yee.Context) (err error) {
		j := c.Params("name")
		fmt.Println(2)
		return c.JSON(http.StatusOK, j)
		//return c.JSON(http.StatusOK, map[string]interface{}{"name":"henry","age": 27})
	})

	v1 := r.Group("/v1")
	v1.Use(middleware.Secure())
	v1.Use(middleware.JWTWithConfig(middleware.JwtConfig{SigningKey: "dbcjqheupqjsuwsm"}))
	//v1.Use(middleware.Cors())
	v1.GET("/fail", func(c yee.Context) (err error) {
		return c.String(http.StatusOK, "is_ok")
		//return c.JSON(http.StatusOK, map[string]interface{}{"name":"henry","age": 27})
	})
	_ = r.Start(":9999")
}