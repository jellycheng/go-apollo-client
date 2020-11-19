package cmd

import (
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/ini"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go-apollo-client/console"
	"os"
	"strings"
)

//job任务
var Artisan = cli.Command{
	Name:   "artisan",
	Usage:  "Start cron server",
	Action: runCron,
	Flags: []cli.Flag{
		stringFlag("config, c", "apollo_example.ini", "指定配置文件，默认apollo_example.ini"),
		stringFlag("appid", "", "应用AppId"),
		stringFlag("env", "", "环境代号，如dev、st、pre、prod等值"),
		stringFlag("cluster", "default", "集群名代号，默认default"),
		stringFlag("namespace", "application", "命名空间代号，默认application"),
		stringFlag("func,f", "ApolloConsole", "指定执行的方法"),
	},
}

func runCron(ctx *cli.Context) {
	var iniFile string
	iniFile = strings.TrimSpace(ctx.String("config"))
	if iniFile != "" {
		if gosupport.IsFile(iniFile) {
			globalCfg.Set("INI_FILE", iniFile) //项目加载的env文件
		} else {
			fmt.Println(gosupport.ToRed(fmt.Sprintf("配置文件%s不存在", iniFile)))
			os.Exit(0)
		}

	} else {
		ExitMsg("请输入指定配置文件")
	}
	commonActionInit(iniFile)
	logrus.Debug("启动artisan服务")
	var paramData = make(map[string]string)
	paramData["ini_file"] = iniFile
	paramData["appid"] = strings.TrimSpace(ctx.String("appid"))
	paramData["env"] = strings.TrimSpace(ctx.String("env"))
	paramData["cluster"] = strings.TrimSpace(ctx.String("cluster"))
	paramData["namespace"] = strings.TrimSpace(ctx.String("namespace"))
	if paramData["appid"] == "" {
		ExitMsg("-appid=缺少应用AppId值")
	}
	if paramData["env"] == "" {
		ExitMsg("-env=缺少环境代号值")
	}

	if f := ctx.String("func"); f != "" {
		paramData["func"] = f
		var jobFuncList = map[string]func(map[string]string, *ini.Config){
			"TestConsole":   console.TestConsole,
			"ApolloConsole": console.ApolloConsole,
		}
		if execFunc, ok := jobFuncList[f]; ok {
			execFunc(paramData, globalIniCfg)
		} else {
			fmt.Println("指定的方法不存在：", f)
		}

	} else {
		fmt.Println("指定的方法不存在，请正确指定要执行的方法")
	}

}
