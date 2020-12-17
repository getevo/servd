package vcs

import (
	"encoding/json"
	"fmt"
	"reflect"
	"time"
)

// Change [...]
type Change struct {
	IDChange    int       `gorm:"primary_key;column:id_change;type:int(11);not null" json:"id_rev"`
	ResourceID  int       `gorm:"column:id_resource;type:int(11);not null" json:"-"`
	Type        string    `gorm:"column:type;type:varchar(16);not null" json:"type"`
	Instruction string    `gorm:"column:instruction;type:varchar(1024);not null" json:"-"`
	Original    string    `gorm:"column:original;type:text;not null" json:"-"`
	Changed     string    `gorm:"column:changed;type:text;not null" json:"-"`
	UpdatedAt   time.Time `gorm:"column:updated_at;type:timestamp;not null" json:"updated_at"`
	UpdatedBy   string    `gorm:"column:updated_by;type:varchar(256);not null" json:"updated_by"`
}

// TableName get sql table name.
func (this *Change) TableName() string {
	return "changes"
}

func (this *Change)SetType(t string) *Change  {
	this.Type = t
	return this
}

func (this *Change)SetId(t int) *Change  {
	this.ResourceID = t
	return this
}


func (this *Change)SetOriginal(t interface{}) error  {
	var kind = reflect.ValueOf(t).Kind()
	if kind == reflect.Struct  || kind == reflect.Slice{
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		this.Original = string(b)
	}else{
		this.Original = fmt.Sprint(t)
	}
	return nil
}


func (this *Change)SetChanged(t interface{}) error  {
	var kind = reflect.ValueOf(t).Kind()
	if kind == reflect.Struct  || kind == reflect.Slice{
		b, err := json.Marshal(t)
		if err != nil {
			return err
		}
		this.Changed = string(b)
	}else{
		this.Changed = fmt.Sprint(t)
	}
	return nil
}

func (this *Change)By(by string) *Change  {
	this.UpdatedBy = by
	return this
}

func (this *Change)SetInstruction(t interface{}) error  {
	b, err := json.Marshal(t)
	if err != nil {
		return err
	}
	this.Instruction = string(b)
	return nil
}

