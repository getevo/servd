package servd

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

type App struct {}

var db *gorm.DB
var mu sync.Mutex
var controller Controller

func (App) Register() {
	fmt.Println("Servd Registered")
}

// WhenReady called after setup all apps
func (App) WhenReady() {
	db = evo.GetDBO()
	LoadAccess()

}

// Router setup routers
func (App) Router() {
	// @doc type 		meta
	// @doc prefix		/api/v1/
	var controller = Controller{}
	evo.Use("/*/token",controller.TokenAuthorized)
	evo.Use("/*/admin",controller.JWTAuthorized)

	v1 := evo.Group("/v1")

	v1.Get("/admin/namespace/list",controller.listNamespace)
	v1.Post("/admin/namespace/create",controller.createNamespace)
	v1.Get("/admin/namespace/:id",controller.getNamespace)
	v1.Post("/admin/namespace/edit/:id",controller.editNamespace)
	v1.Post("/admin/namespace/delete/:id",controller.deleteNamespace)
	v1.Post("/admin/namespace/restore/:id",controller.restoreNamespace)

	v1.Get("/admin/environment/list",controller.listEnvironment)
	v1.Post("/admin/environment/create",controller.createEnvironment)
	v1.Get("/admin/environment/:id",controller.getEnvironment)
	v1.Post("/admin/environment/edit/:id",controller.editEnvironment)
	v1.Post("/admin/environment/delete/:id",controller.deleteEnvironment)
	v1.Post("/admin/environment/restore/:id",controller.restoreEnvironment)
}

// Permissions setup permissions of app
func (App) Permissions() []evo.Permission { return []evo.Permission{} }

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}

func (App) Pack() {}
