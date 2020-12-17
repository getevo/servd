package servd

import (
	"getevo/servd/apps/models"
	"github.com/getevo/evo/lib/log"
	"net"
)

type Access struct {
	NetMask *net.IPNet
	Token   string
}
var access []Access
func LoadAccess()  {
	var items []models.Access
	db.Find(&items)
	var tmpAccess []Access
	for _,item := range items{
	    _,netmask,err := net.ParseCIDR(item.Netmask)
	    if err != nil{
	    	log.Error(err)
		}else {
			tmpAccess = append(tmpAccess, Access{
				NetMask:netmask,
				Token: item.Token,
			})
		}
	}
	access = tmpAccess
}
