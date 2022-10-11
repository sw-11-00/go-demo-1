package router

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"

	"go-demo-1/handler"
)

func Init() {
	r := gin.Default()
	// middleware
	r.Use(PrintMiddleware)
	r.Use(ValidateToken())

	v1 := r.Group("/v1")
	{
		// http://localhost:8000/v1/hello
		v1.GET("/hello", handler.HelloPage)

		// http://localhost:8000/v1/hello/michael
		v1.GET("/hello/:name", func(context *gin.Context) {
			name := context.Param("name")
			context.String(http.StatusOK, "Hello %s", name)
		})

		// http://localhost:8000/v1/value/michael/male
		v1.GET("/value/:a/:b", func(context *gin.Context) {
			valueA := context.Param("a")
			valueB := context.Param("b")
			context.String(http.StatusOK, "value a:%s b:%s", valueA, valueB)
		})

		// http://localhost:8000/v1/value/
		v1.GET("/value/", func(context *gin.Context) {
			context.String(http.StatusOK, "value 没有参数")
		})

		// http://localhost:8000/v1/welcome?lastname=muzico&firstname=yococo
		v1.GET("/welcome", func(context *gin.Context) {
			firstname := context.DefaultQuery("firstname", "Guest")
			lastname := context.Query("lastname")
			context.String(http.StatusOK, "Hello %s %s", firstname, lastname)
		})

		r.LoadHTMLGlob(getCurrentPath() + "/router/templates/*")
		v2 := r.Group("/v2")
		{
			// http://localhost:8000/v2/index
			v2.GET("/index", func(context *gin.Context) {
				context.HTML(http.StatusOK, "index.html", gin.H{
					"title": "hello world!",
				})
			})
		}

		r.NoRoute(func(context *gin.Context) {
			context.JSON(http.StatusNotFound, gin.H{
				"status": 404,
				"error":  "404, page not exists",
			})
		})

	}

	r.Run(":8000")
}

func getCurrentPath() string {
	str, _ := os.Getwd()
	return str
}

func PrintMiddleware(context *gin.Context) {
	fmt.Print("before http")
	context.Next()
}

func ValidateToken() gin.HandlerFunc {
	return func(context *gin.Context) {
		token := context.Request.FormValue("token")
		if token == "" {
			context.JSON(401, gin.H{
				"message": "Token required",
			})
			context.Abort()
			return
		}
		if token != "accesstoken" {
			context.JSON(http.StatusOK, gin.H{
				"message": "Invalid Token",
			})
			context.Abort()
			return
		}
		context.Next()
	}
}
