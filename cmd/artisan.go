package cmd

import (
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/robfig/cron/v3"
	"github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"go-apollo-client/console"
	"os"
)

//crontab计划任务
var Artisan = cli.Command{
	Name:   "artisan",
	Usage:  "Start cron server",
	Action: runCron,
	Flags: []cli.Flag{
		stringFlag("config, c", "apollo_example.ini", "指定配置文件，默认apollo_example.ini"),
		stringFlag("appid", "", "应用AppId"),
		stringFlag("env", "", "环境代号，如dev、st、pre、prod等值"),
		stringFlag("cluster", "application", "集群名代号，默认application"),
		stringFlag("func,f", "", "指定执行的方法"),
	},
}

func runCron(ctx *cli.Context) {
	var iniFile string
	iniFile = ctx.String("config")
	if iniFile != "" {
		if gosupport.IsFile(iniFile) {
			globalCfg.Set("INI_FILE", iniFile) //项目加载的env文件
		} else {
			fmt.Println(gosupport.ToRed(fmt.Sprintf("配置文件%s不存在", iniFile)))
			os.Exit(0)
		}

	} else {
		fmt.Println(gosupport.ToYellow("请输入指定配置文件"))
		os.Exit(0)
	}
	commonActionInit(iniFile)
	logrus.Debug("启动artisan服务")
	if f := ctx.String("func"); f != "" {
		var jobFuncList = map[string]func(){
			"TestConsole":           console.TestConsole,
		}
		if execFunc, ok := jobFuncList[f]; ok {
			execFunc()
		} else {
			fmt.Println("指定的方法不存在：", f)
		}

		return
	}

	c := cron.New(cron.WithSeconds())

	//每1秒执行一次，测试代码
	//c.AddFunc("* * * * * *", console.TestConsole)

	//每隔2分钟执行一次： "0 0/2 * * * *"

	c.Start()

	select {}

}
