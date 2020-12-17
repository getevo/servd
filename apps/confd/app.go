package confd

// @doc type 		app
// @doc name		confd
// @doc description Confd Service
// @doc author		reza

import (
	"fmt"
	"github.com/getevo/evo"
	"github.com/getevo/evo/menu"
	"gorm.io/gorm"
	"sync"
)

func Register() {
	evo.Register(App{})
}

type App struct {
}

var db *gorm.DB
var mu sync.Mutex
var controller Controller

func (App) Register() {
	fmt.Println("Confd Registered")
}

// WhenReady called after setup all apps
func (App) WhenReady() {
	db = evo.GetDBO()

}

// Router setup routers
func (App) Router() {
	// @doc type 		meta
	// @doc prefix		/api/v1/
	var controller = Controller{}

	v1 := evo.Group("/v1")
	v1.Get("/token/get/:namespace.:env.:template",controller.getConfig)


	v1.Get("/admin/template/list",controller.listTemplate)
	v1.Post("/admin/template/create",controller.createTemplate)
	v1.Get("/admin/template/:id",controller.getTemplate)
	v1.Post("/admin/template/edit/:id",controller.editTemplate)
	v1.Post("/admin/template/delete/:id",controller.deleteTemplate)
	v1.Post("/admin/template/restore/:id",controller.restoreTemplate)
	v1.Get("/admin/template/revision/:id",controller.revisionTemplate)
	v1.Post("/admin/template/revert/:id",controller.revertTemplate)

	v1.Get("/admin/param/:namespace.:env.:name",controller.getParamValue)
	v1.Get("/admin/param/:namespace.:env",controller.getParams)
	v1.Post("/admin/param/:namespace.:env",controller.setParams)
	v1.Get("/admin/param/remove/:namespace.:env.:name",controller.removeParam)
	v1.Get("/admin/param/revision/:namespace.:env.",controller.revisionConfig)
	v1.Get("/admin/param/revert/:id",controller.revertParams)



}

// Permissions setup permissions of app
func (App) Permissions() []evo.Permission { return []evo.Permission{} }

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}

func (App) Pack() {}
