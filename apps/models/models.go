package models

import (
	"getevo/servd/apps/vcs"
	"gorm.io/gorm"
	"time"
)

// Access [...]
type Access struct {
	IDAccess      int         `gorm:"primary_key;column:id_access;type:int(11);not null" json:"-"`
	Netmask       string      `gorm:"column:netmask;type:varchar(18);not null" json:"netmask"`
	Token         string      `gorm:"column:token;type:varchar(32);not null" json:"token"`
	Label         string      `gorm:"column:label;type:varchar(255);not null" json:"label"`
	Description   string      `gorm:"column:description;type:varchar(1024);not null" json:"description"`
	IDEnvironment *int        `gorm:"index:fk_access_id_environment;column:id_environment;type:int(11)" json:"id_environment"`
	Environment   Environment `gorm:"association_foreignkey:id_environment;foreignkey:id_environment" json:"-"`
	IDNamespace   *int        `gorm:"index:fk_access_id_namespace;column:id_namespace;type:int(11)" json:"id_namespace"`
	Namespace     Namespace   `gorm:"association_foreignkey:id_namespace;foreignkey:id_namespace" json:"-"`
}

// TableName get sql table name.
func (m *Access) TableName() string {
	return "access"
}

// Environment [...]
type Environment struct {
	IDEnvironment int        `gorm:"primary_key;column:id_environment;type:int(11);not null" json:"id_environment"`
	IDNamespace   int        `gorm:"index:fk_environment_id_namespace;column:id_namespace;type:int(11);not null" form:"id_namespace" json:"id_namespace"`
	Namespace     Namespace  `gorm:"association_foreignkey:id_namespace;foreignkey:id_namespace" json:"-"`
	Name          string     `gorm:"unique;column:name;type:varchar(32);not null" json:"name" form:"name"`
	Label         string     `gorm:"column:label;type:varchar(255);not null" json:"label" form:"label"`
	Description   string     `gorm:"column:description;type:varchar(512);not null" json:"description" form:"description"`
	CreatedAt     time.Time  `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt     time.Time  `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
}

// TableName get sql table name.
func (m *Environment) TableName() string {
	return "environment"
}

// Namespace [...]
type Namespace struct {
	IDNamespace int        `gorm:"primary_key;column:id_namespace;type:int(11);not null" json:"id_namespace"`
	Name        string     `gorm:"unique;column:name;type:varchar(32);not null"  form:"name" json:"name"`
	Label       string     `gorm:"column:label;type:varchar(255);not null"   form:"label" json:"label"`
	Description string     `gorm:"column:description;type:varchar(1024);not null"   form:"description" json:"description"`
	CreatedAt   time.Time  `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt   time.Time  `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	DeletedAt   gorm.DeletedAt `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
}

// TableName get sql table name.
func (m *Namespace) TableName() string {
	return "namespace"
}

// Parameter [...]
type Parameter struct {
	gorm.Model
	IDNamespace   int         `gorm:"index:fk_parameter_id_namespace;column:id_namespace;type:int(11);not null" form:"id_namespace" json:"id_namespace"`
	Namespace     Namespace   `gorm:"association_foreignkey:id_namespace;foreignkey:id_namespace" json:"-"`
	IDEnvironment int         `gorm:"index:fk_parameter_id_environment;column:id_environment;type:int(11);not null" form:"id_environment" json:"id_environment"`
	Environment   Environment `gorm:"association_foreignkey:id_environment;foreignkey:id_environment" json:"-"`
	Name          string      `gorm:"column:name;type:varchar(255);not null" json:"name" form:"name"`
	Value         string      `gorm:"column:value;type:varchar(2048);not null" json:"value" form:"value"`
	Description   string      `gorm:"column:description;type:varchar(1024);not null" json:"description" form:"description"`
}

// TableName get sql table name.
func (m *Parameter) TableName() string {
	return "parameter"
}

// Template [...]
type Template struct {
	IDTemplate    int         `gorm:"primary_key;column:id_template;type:int(11);not null" json:"id_template"`
	IDNamespace   int         `gorm:"index:fk_template_id_namespace;column:id_namespace;type:int(11);not null" form:"id_namespace" json:"id_namespace"`
	Namespace     Namespace   `gorm:"association_foreignkey:id_namespace;foreignkey:id_namespace" json:"-"`
	IDEnvironment int         `gorm:"index:fk_template_id_environment;column:id_environment;type:int(11);not null" form:"id_environment" json:"id_environment"`
	Environment   Environment `gorm:"association_foreignkey:id_environment;foreignkey:id_environment" json:"-"`
	FileName      string      `gorm:"column:filename;type:varchar(255);not null" json:"filename" form:"filename"`
	Name          string      `gorm:"unique;column:name;type:varchar(32);not null" json:"name" form:"name"`
	Label         string      `gorm:"column:label;type:varchar(255);not null" json:"label" form:"label"`
	Description   string      `gorm:"column:description;type:varchar(1024);not null" json:"description" form:"description"`
	Type          string      `gorm:"column:type;type:varchar(32);not null" json:"type" form:"type"`
	Template      string      `gorm:"column:template;type:text;not null" json:"template" form:"template"`
	CreatedAt     time.Time   `gorm:"column:created_at;type:timestamp;not null" json:"created_at"`
	UpdatedAt     time.Time   `gorm:"column:updated_at;type:datetime;not null" json:"updated_at"`
	DeletedAt     gorm.DeletedAt  `gorm:"column:deleted_at;type:timestamp" json:"deleted_at"`
}

// TableName get sql table name.获取数据库表名
func (m *Template) TableName() string {
	return "template"
}


func (t *Template) BeforeSave(tx *gorm.DB) (err error) {
	t.Version(tx)
	return
}

func (t *Template) AfterCreate(tx *gorm.DB) (err error) {
	if t.IDTemplate > 0 {
		t.Version(tx)
	}
	return
}

func (t *Template) BeforeDelete(tx *gorm.DB) (err error) {
	t.Version(tx)
	return
}


func (t *Template) Version(tx *gorm.DB){
	var original Template
	tx.Debug().Where("id_template = ?",t.IDTemplate).Find(&original)
	var change = vcs.Change{}
	change.SetType("conf.template").By("system")
	change.SetId(original.IDTemplate)
	change.SetOriginal(original)
	change.SetChanged(*t)
	if change.Changed != change.Original {
		db.Debug().Create(&change)
	}

}

