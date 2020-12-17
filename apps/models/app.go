package models

// @doc type 		app
// @doc name		confd
// @doc description Confd Service
// @doc author		reza

import (
	"fmt"
	"github.com/getevo/evo"
	"github.com/getevo/evo/menu"
	"gorm.io/gorm"
)

func Register() {
	evo.Register(App{})
}

type App struct {
}

var db *gorm.DB


func (App) Register() {
	fmt.Println("Models Registered")
}

// WhenReady called after setup all apps
func (App) WhenReady() {
	db = evo.GetDBO()

}

// Router setup routers
func (App) Router() {


}

// Permissions setup permissions of app
func (App) Permissions() []evo.Permission { return []evo.Permission{} }

// Menus setup menus
func (App) Menus() []menu.Menu {
	return []menu.Menu{}
}

func (App) Pack() {}
