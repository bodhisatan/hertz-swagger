package main

import (
	hertzSwagger "github.com/bodhisatan/hertz-swagger"
	_ "github.com/bodhisatan/hertz-swagger/example/basic/docs"
	"github.com/bodhisatan/hertz-swagger/example/basic/handler"
	"github.com/cloudwego/hertz/pkg/app/server"
	swaggerFiles "github.com/swaggo/files"
)

// @title HertzTest
// @version 1.0
// @description This is a demo using Hertz.

// @contact.name Bodhisatan
// @contact.url https://github.com/bodhisatan
// @contact.email bodhisatanyao@gmail.com

// @license.name Apache 2.0
// @license.url http://www.apache.org/licenses/LICENSE-2.0.html

// @host localhost:8888
// @BasePath /
// @schemes http
func main() {
	h := server.Default()

	h.GET("/ping", handler.PingHandler)

	url := hertzSwagger.URL("http://localhost:8888/swagger/doc.json") // The url pointing to API definition
	h.GET("/swagger/*any", hertzSwagger.WrapHandler(swaggerFiles.Handler, url))

	h.Spin()
}
