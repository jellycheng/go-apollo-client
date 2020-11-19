package console

import (
	"fmt"
	"github.com/jellycheng/gosupport/ini"
	"github.com/sirupsen/logrus"
	"go-apollo-client/pkg/apollo"
	"strings"
)

func ApolloConsole(param map[string]string, iniCfg *ini.Config)  {
	if env, ok := param["env"];ok {
		env = strings.ToLower(env)
		apolloHost := iniCfg.MustGetValue(env, "apollo_host")
		appId := param["appid"]
		cluster := param["cluster"]
		namespace := param["namespace"]
		if resBody,err := apollo.GetConfig4Cache(apolloHost, appId,cluster,namespace,"");err== nil{
			fmt.Println(resBody)
			//保存env文件

			//调用shell同步env，重启服务

		}

	} else {
		logrus.Error("env值不存在")
	}

}
