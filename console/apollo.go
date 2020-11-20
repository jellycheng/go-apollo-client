package console

import (
	"encoding/json"
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/ini"
	"github.com/sirupsen/logrus"
	"go-apollo-client/pkg/apollo"
	"os"
	"path/filepath"

)
//拉去env内容，保存文件、回调shell
func ApolloConsole(param map[string]string, iniCfg *ini.Config)  {

	if env, ok := param["env"];ok {
		apolloHost := iniCfg.MustGetValue(env, "apollo_host")
		if apolloHost == "" {
			logrus.Error(fmt.Sprintf("请配置%s环境的apollo host", env))
			return
		}
		appId := param["appid"]
		cluster := param["cluster"]
		namespace := param["namespace"]
		//格式：阿波罗项目名appid_集群名cluster_命名空间名application=env文件地址
		envFile := iniCfg.MustGetValue(env, fmt.Sprintf("%s_%s_%s", appId, cluster, namespace))

		if resBody,err := apollo.GetConfig4Cache(apolloHost, appId,cluster,namespace,"");err== nil{
			envContent := apollo.JsonToEnvContent(resBody)
			if envFile == "" {
				envDir := iniCfg.MustGetValue("default", "env_dir")
				envFile = fmt.Sprintf("%s/%s/.env", envDir, appId)
			}
			targetDir := filepath.Dir(envFile)
			if !gosupport.IsDir(targetDir) {
				os.MkdirAll(targetDir, os.ModePerm)
			}
			//保存env文件
			if err := apollo.CleanContentWrite(envFile, envContent);err == nil {
				//调用shell同步env，重启服务
				logrus.Debug("保存env文件：", envFile)
				cmdHook := iniCfg.MustGetValue(env, "cmd_hook")
				if cmdHook != "" {
					//钩子接收的参数：环境代号、阿波罗项目名appid、集群名、命名空间名
					if okmsg, failMsg, err := gosupport.ExecCmd(cmdHook, env, appId, cluster, namespace);err == nil {
						if okmsg != "" {
							fmt.Println(fmt.Sprintf("ok:%s", okmsg))
						}
						if failMsg != "" {
							fmt.Println(fmt.Sprintf("fail: %s ", failMsg))
						}

					} else {
						logrus.Error("命令执行失败：" , err.Error())
					}
				}
			} else {
				logrus.Error("env文件保存失败：", envFile)
			}

		} else {
			logrus.Error(err.Error())
		}

	} else {
		logrus.Error("env值不存在")
	}

}


//监听项目
func MonitorApolloConsole(param map[string]string, iniCfg *ini.Config)  {
	if env, ok := param["env"];ok {
		apolloHost := iniCfg.MustGetValue(env, "apollo_host")
		if apolloHost == "" {
			logrus.Error(fmt.Sprintf("请配置%s环境的apollo host", env))
			return
		}
		appId := param["appid"]
		cluster := param["cluster"]
		namespace := param["namespace"]
		var paramData = make(map[string]string)
		paramData["appid"] = appId
		paramData["env"] = env
		paramData["cluster"] = cluster
		paramData["namespace"] = namespace
		var noId int64
		for {

			notificationsByte, _ := json.Marshal([]map[string]interface{}{
												{
													"namespaceName":  namespace,
													"notificationId": noId,
												},
											})
			notifications := string(notificationsByte)
			isUpdate, notificationId := apollo.Notifications(apolloHost, appId, cluster, notifications)
			if isUpdate {
				ApolloConsole(param, iniCfg)
				noId = notificationId
				logrus.Info(fmt.Sprintf("同步成功，%s,%s,%s,%s,%v", apolloHost, appId, cluster, namespace, notificationId))
			}
		}

	}  else {
		logrus.Error("env值不存在")
	}

}
