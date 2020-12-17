package main

import (
	"getevo/servd/apps/confd"
	"getevo/servd/apps/models"
	"getevo/servd/apps/servd"
	"github.com/getevo/evo"
)

func main()  {
	evo.Setup()
	servd.Register()
	models.Register()
	confd.Register()

	evo.Run()
}

