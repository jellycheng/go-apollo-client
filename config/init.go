package config

import (
	"github.com/jellycheng/gosupport"
)

var globalenv *gosupport.DataManage


func init()  {
	globalenv = gosupport.NewGlobalEnvSingleton()

}

