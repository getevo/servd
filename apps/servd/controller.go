package servd

import (
	"fmt"
	"getevo/servd/apps/models"
	"github.com/getevo/evo"
	"net"
	"net/http"
)

type Controller struct {}

func (Controller)JWTAuthorized(request *evo.Request) error {
	request.Next()
	return nil
}
func (Controller)TokenAuthorized(request *evo.Request) error {
	token := request.Query("token")

	IPAddress := net.ParseIP(request.Get("X-Real-Ip"))
	if request.Get("X-Real-Ip") == "" {
		IPAddress = net.ParseIP(request.Get("X-Forwarded-For"))
	}
	if request.Get("X-Forwarded-For") == "" {
		IPAddress = request.Context.Context().RemoteIP()
	}


	if token == ""{
		request.Status(http.StatusUnauthorized)
		return fmt.Errorf("token is not provided")
	}

	for _,item := range access{
		if item.NetMask.Contains(IPAddress) && (token == "*" || token == item.Token){
			request.Next()
			return nil
		}
	}


	return fmt.Errorf("access denied")
}


func (Controller)createNamespace(request *evo.Request)  {
	var ns models.Namespace
	request.BodyParser(&ns)
	err := db.Debug().Create(&ns).Error
	if err !=nil {
		request.WriteResponse(err)
		return
	}
	request.WriteResponse(&ns)
}

func (Controller)listNamespace(request *evo.Request)  {
	var trashed = request.Query("trash") != ""
	var ns = []models.Namespace{}

	if trashed{
		db.Debug().Unscoped().Where("deleted_at IS NOT NULL").Find(&ns)
	}else{
		db.Debug().Unscoped().Where("deleted_at IS NULL").Find(&ns)
	}
	request.WriteResponse(ns)
}

func (Controller)getNamespace(request *evo.Request)  {
	var ns models.Namespace
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_namespace = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	request.WriteResponse(ns)
}

func (Controller)editNamespace(request *evo.Request)  {
	var ns models.Namespace
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_namespace = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	request.BodyParser(&ns)
	ns.IDNamespace = id
	err := db.Debug().Where("id_namespace = ?",id).Updates(&ns).Error
	if err !=nil {
		request.WriteResponse(err)
		return
	}
	request.WriteResponse(&ns)
}

func (Controller)deleteNamespace(request *evo.Request)  {
	var ns models.Namespace
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_namespace = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	db.Debug().Delete(&ns)
	request.WriteResponse(true)
}

func (Controller)restoreNamespace(request *evo.Request)  {
	var ns models.Namespace
	id := request.ParamsI("id","number").Int()
	if db.Debug().Unscoped().Where("id_namespace = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	db.Model(&ns).Where("id_namespace = ?",id).Update("deleted_at",nil)
	request.WriteResponse(true)
}


func (Controller)createEnvironment(request *evo.Request)  {
	var ns models.Environment
	request.BodyParser(&ns)
	err := db.Debug().Create(&ns).Error
	if err !=nil {
		request.WriteResponse(err)
		return
	}
	request.WriteResponse(&ns)
}

func (Controller)listEnvironment(request *evo.Request)  {
	var trashed = request.Query("trash") != ""
	var ns = []models.Environment{}

	if trashed{
		db.Debug().Unscoped().Where("deleted_at IS NOT NULL").Find(&ns)
	}else{
		db.Debug().Unscoped().Where("deleted_at IS NULL").Find(&ns)
	}
	request.WriteResponse(ns)
}

func (Controller)getEnvironment(request *evo.Request)  {
	var ns models.Environment
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_environment = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	request.WriteResponse(ns)
}

func (Controller)editEnvironment(request *evo.Request)  {
	var ns models.Environment
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_environment = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	request.BodyParser(&ns)
	ns.IDEnvironment = id
	err := db.Debug().Where("id_environment = ?",id).Updates(&ns).Error
	if err !=nil {
		request.WriteResponse(err)
		return
	}
	request.WriteResponse(&ns)
}

func (Controller)deleteEnvironment(request *evo.Request)  {
	var ns models.Environment
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_environment = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	db.Debug().Delete(&ns)
	request.WriteResponse(true)
}

func (Controller)restoreEnvironment(request *evo.Request)  {
	var ns models.Environment
	id := request.ParamsI("id","number").Int()
	if db.Debug().Unscoped().Where("id_environment = ?",id).Find(&ns).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	db.Model(&ns).Where("id_environment = ?",id).Update("deleted_at",nil)
	request.WriteResponse(true)
}