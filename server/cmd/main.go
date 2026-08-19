package main

import (
	"embed"
	"net/http"
	"template/pkg/controllers"
	"template/pkg/models"

	"github.com/alpha-omega-corp/_core/app"
	"github.com/uptrace/bunrouter"
)

var (
	//go:embed config.yml fixtures/fixture.yml
	efs embed.FS
)

func main() {
	app.New(efs).Bootstrap(func(router *bunrouter.Router, deps app.Deps) {
		app.RegisterControllers(deps, controllers.Meat)

		router.GET("/meat", controllers.Meat.GetAll())
		router.GET("/meat/:id", controllers.Meat.GetOne())
		router.POST("/meat", controllers.Meat.Create())
		router.PUT("/meat/:id", controllers.Meat.Update())
		router.DELETE("/meat/:id", controllers.Meat.Delete())
		router.POST("/meat/:id/image", controllers.Meat.UploadImage())

		router.GET("/storage/*path", bunrouter.HTTPHandler(
			http.StripPrefix("/storage/", http.FileServer(http.Dir(controllers.StorageDir))),
		))

	}, []interface{}{
		(*models.Meat)(nil),
	}...)
}
