package console

import (
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/ini"
	"github.com/sirupsen/logrus"
	"go-apollo-client/pkg/apollo"
	"os"
	"path/filepath"
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

			envContent := apollo.JsonToEnvContent(resBody)
			envDir := iniCfg.MustGetValue("default", "env_dir")
			envFile := fmt.Sprintf("%s/%s/.env", envDir, appId)
			targetDir := filepath.Dir(envFile)
			if !gosupport.IsDir(targetDir) {
				os.MkdirAll(targetDir, os.ModePerm)
			}
			//保存env文件
			if _, err := gosupport.FilePutContents(envFile, envContent);err == nil {
				//调用shell同步env，重启服务
				fmt.Println(envFile)

			} else {
				fmt.Println("env文件保存失败：", envFile)
			}


		}

	} else {
		logrus.Error("env值不存在")
	}

}
