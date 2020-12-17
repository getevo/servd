package confd

import (
	"bytes"
	"encoding/json"
	"fmt"
	"getevo/servd/apps/models"
	"getevo/servd/apps/vcs"
	"github.com/getevo/evo"
	"github.com/getevo/evo/lib/log"
	"text/template"
)

type Controller struct {}
type RecoverInstruction struct {
	IDEnvironment int `json:"id_environment"`
	IDNamespace   int `json:"id_namespace"`

}


var templates = map[string]*template.Template{}


func (Controller)getParams(request *evo.Request){
	var namespace models.Namespace
	var env models.Environment
	if err := db.Debug().Debug().Where("name = ?",request.Params("namespace")).Take(&namespace).Error; err != nil{
		request.Write("invalid namespace")
		return
	}
	if err := db.Debug().Debug().Where("name = ?",request.Params("env")).Take(&env).Error; err != nil{
		request.Write("invalid environment")
		return
	}
	var parameters = []models.Parameter{}
	db.Debug().Where("id_namespace = ? AND id_environment = ?",namespace.IDNamespace,env.IDEnvironment).Find(&parameters)
	request.WriteResponse(parameters)
}

func (Controller)getConfig(request *evo.Request){
	var namespace models.Namespace
	var env models.Environment
	var dbTemplate models.Template
	if err := db.Debug().Debug().Where("name = ?",request.Params("namespace")).Take(&namespace).Error; err != nil{
		request.Write("invalid namespace")
		return
	}
	if err := db.Debug().Debug().Where("name = ?",request.Params("env")).Take(&env).Error; err != nil{
		request.Write("invalid environment")
		return
	}
	if err := db.Debug().Debug().Where("name = ? AND id_namespace = ? AND id_environment = ?",request.Params("template"),namespace.IDNamespace,env.IDEnvironment).Take(&dbTemplate).Error; err != nil{
		request.Write("invalid template")
		return
	}
	key := namespace.Name+"."+env.Name+"."+request.Params("template")
	if _,ok := templates[key]; !ok{
		t, err := template.New(key).Parse(dbTemplate.Template)
		if err != nil{
			log.Error(err)
			request.Write("unable to get config")
			return
		}
		templates[key] = t
	}

	var parameters []models.Parameter
	db.Debug().Where("id_namespace = ? AND id_environment = ?",namespace.IDNamespace,env.IDEnvironment).Find(&parameters)
	var data  = map[string]string{}
	if parameters != nil{
		for _,param := range parameters{
			data[param.Name] = param.Value
		}
	}
	var tpl bytes.Buffer
	if err := templates[key].Execute(&tpl, data); err != nil {
		log.Error(err)
		request.Write("unable to render config")
		return
	}
	request.Set("Content-Disposition","attachment; filename=\""+dbTemplate.FileName+"\"")
	request.Write(tpl.Bytes())

}



func (Controller)editTemplate(request *evo.Request){
	var original models.Template
	var changed models.Template
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_template = ?",id).Find(&original).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	changed = original
	request.BodyParser(&changed)
	changed.IDTemplate = id
	err := db.Debug().Where("id_template = ?",id).Updates(&changed).Error
	if err !=nil {
		request.WriteResponse(err)
		return
	}
	request.WriteResponse(&changed)
}


func (Controller)setParams(request *evo.Request){
	var items []models.Parameter
	var namespace models.Namespace
	var environment models.Environment
	err := json.Unmarshal([]byte(request.Body()),&items)
	if err != nil{
		request.WriteResponse(err)
		return
	}

	if err := db.Debug().Debug().Where("name = ?",request.Params("namespace")).Take(&namespace).Error; err != nil{
		request.Write("invalid namespace")
		return
	}
	if err := db.Debug().Debug().Where("name = ?",request.Params("env")).Take(&environment).Error; err != nil{
		request.Write("invalid environment")
		return
	}

	var original = []models.Parameter{}
	var changed  = []models.Parameter{}
	db.Debug().Where("id_namespace = ? AND id_environment = ?",namespace.IDNamespace,environment.IDEnvironment).Find(&original)


	for _,item := range items{
		var param models.Parameter
		if db.Debug().Where("id_namespace = ? AND id_environment = ? AND name = ?",namespace.IDNamespace,environment.IDEnvironment,item.Name).Take(&param).Error == nil{
			if item.Value != param.Value{
				param.Value = item.Value
				db.Debug().Debug().Updates(&param)
			}
		}else{
			param = models.Parameter{
				IDNamespace: namespace.IDNamespace,
				IDEnvironment:environment.IDEnvironment,
				Name: item.Name,
				Value: item.Value,
				Description: item.Description,
			}
			db.Debug().Debug().Create(&param)
		}
	}

	db.Debug().Where("id_namespace = ? AND id_environment = ?",namespace.IDNamespace,environment.IDEnvironment).Find(&changed)

	var change = vcs.Change{}
	change.SetType("conf.param").By("system")
	change.SetId(environment.IDEnvironment)
	change.SetOriginal(original)
	change.SetChanged(changed)
	change.SetInstruction(RecoverInstruction{
		IDEnvironment:environment.IDEnvironment,
		IDNamespace:namespace.IDNamespace,
	})
	if change.Changed != change.Original {
		db.Debug().Create(&change)
	}

	request.WriteResponse(true)
}

func (Controller)revertParams(request *evo.Request){
	revertId := request.Params("id")
	var change vcs.Change
	var instruction RecoverInstruction
	if db.Debug().Where("id_change = ?",revertId).Find(&change).RowsAffected == 0{
		request.WriteResponse(false)
	}

	var parameters []models.Parameter
	err := json.Unmarshal([]byte(change.Changed),&parameters)
	if err != nil{
		request.WriteResponse(err)
		return
	}
	err = json.Unmarshal([]byte(change.Instruction),&instruction)
	if err != nil{
		request.WriteResponse(err)
		return
	}
	var original = []models.Parameter{}
	var changed  = []models.Parameter{}

	for _,item := range parameters{
		var param models.Parameter
		if db.Debug().Find("id_environment = ? AND id_namespace = ? AND name = ?",item.IDEnvironment,item.IDNamespace,item.Name).Take(&param).RowsAffected > 0{
			param.Value = item.Value
			db.Debug().Where("id_namespace = ? AND id_environment = ?",instruction.IDNamespace,instruction.IDEnvironment).Find(&original)
			db.Debug().Updates(&param)
		}else{
			param = models.Parameter{
				IDNamespace: item.IDNamespace,
				IDEnvironment:item.IDEnvironment,
				Name: item.Name,
				Value: item.Value,
				Description: item.Description,
			}
			db.Debug().Create(&param)
		}
	}
	db.Debug().Where("id_namespace = ? AND id_environment = ?",instruction.IDNamespace,instruction.IDEnvironment).Find(&changed)
	change = vcs.Change{
		Instruction:change.Instruction,
	}
	change.SetType("conf.param").By("system")
	change.SetId(instruction.IDEnvironment)
	change.SetOriginal(original)
	change.SetChanged(changed)
	if change.Changed != change.Original {
		db.Debug().Create(&change)
	}
	request.WriteResponse(true)
}

func (Controller)revertTemplate(request *evo.Request){
	revertId := request.Params("id")
	var change vcs.Change
	if db.Debug().Where("id_change = ?",revertId).Find(&change).RowsAffected == 0{
		request.WriteResponse(false)
	}

	var t models.Template
	err := json.Unmarshal([]byte(change.Changed),&t)
	if err != nil{
		request.WriteResponse(err)
		return
	}
	t.IDTemplate = change.ResourceID
	var item models.Template
	if db.Debug().Where("id_template = ?",t.IDTemplate).Find(&item).RowsAffected > 0{
		db.Debug().Where("id_template = ?",t.IDTemplate).Updates(&t)
	}else{
		db.Debug().Create(&t)
	}

	request.WriteResponse(true)
}

func (Controller)listTemplate(request *evo.Request){
	var trashed = request.Query("trash") != ""
	var templates = []models.Template{}
	if trashed{
		db.Debug().Unscoped().Where("deleted_at IS NOT NULL").Find(&templates)
	}else{
		db.Debug().Unscoped().Where("deleted_at IS NULL").Find(&templates)
	}
	for key,_ := range templates{
		templates[key].Template = ""
	}
	request.WriteResponse(templates)
}

func (Controller)createTemplate(request *evo.Request){
	var template models.Template
	request.BodyParser(&template)
	err := db.Debug().Create(&template).Error
	if err !=nil {
		request.WriteResponse(err)
		return
	}
	request.WriteResponse(&template)
}

func (Controller)getTemplate(request *evo.Request){
	var template models.Template
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_template = ?",id).Find(&template).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	request.WriteResponse(template)
}

func (Controller)deleteTemplate(request *evo.Request){
	var template models.Template
	id := request.ParamsI("id","number").Int()
	if db.Debug().Where("id_template = ?",id).Find(&template).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	db.Debug().Delete(&template)
	request.WriteResponse(true)
}

func (Controller)restoreTemplate(request *evo.Request){
	var template models.Template
	id := request.ParamsI("id","number").Int()
	if db.Unscoped().Where("id_template = ?",id).Find(&template).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	db.Unscoped().Model(&template).Where("id_template = ?",id).Update("deleted_at",nil)

	request.WriteResponse(true)
}

func (Controller)revisionTemplate(request *evo.Request){
	var changes = []vcs.Change{}
	db.Debug().Where("type = 'conf.template' AND id_resource = ?",request.ParamsI("id").Int()).Order("updated_at DESC").Limit(40).Find(&changes)
	request.WriteResponse(changes)
}

func (Controller)removeParam(request *evo.Request){
	var param models.Parameter
	id := request.ParamsI("id","number").Int()

	if db.Debug().Where("id = ?",id).Take(&param).RowsAffected == 0{
		request.WriteResponse(fmt.Errorf("invalid object"))
		return
	}
	var original = []models.Parameter{}
	var changed  = []models.Parameter{}
	db.Debug().Where("id_namespace = ? AND id_environment = ?",param.IDNamespace,param.IDEnvironment).Find(&original)
	db.Debug().Delete(&param)

	var change = vcs.Change{}
	change.SetType("conf.param").By("system")
	change.SetId(param.IDEnvironment)
	change.SetOriginal(original)
	change.SetChanged(changed)
	change.SetInstruction(RecoverInstruction{
		IDEnvironment:param.IDEnvironment,
		IDNamespace:param.IDNamespace,
	})

	request.WriteResponse(true)
}

func (Controller)revisionConfig(request *evo.Request){
	var changes = []vcs.Change{}
	db.Debug().Where("type = 'conf.param' AND id_resource = ?",request.ParamsI("environment").Int()).Order("updated_at DESC").Limit(40).Find(&changes)
	request.WriteResponse(changes)
}

func (Controller)getParamValue(request *evo.Request)  {
	var namespace models.Namespace
	var env models.Environment
	if err := db.Debug().Debug().Where("name = ?",request.Params("namespace")).Take(&namespace).Error; err != nil{
		request.Write("invalid namespace")
		return
	}
	if err := db.Debug().Debug().Where("name = ?",request.Params("env")).Take(&env).Error; err != nil{
		request.Write("invalid environment")
		return
	}
	var parameter models.Parameter
	db.Debug().Where("id_namespace = ? AND id_environment = ? AND name = ?",namespace.IDNamespace,env.IDEnvironment,request.Params("name")).Find(&parameter)
	request.WriteResponse(parameter)
}