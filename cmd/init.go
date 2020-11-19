package cmd

import (
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/ini"
	"github.com/sirupsen/logrus"
	"go-apollo-client/constants"
	mylog "go-apollo-client/pkg/log"
	"io"
	"os"
	"time"
)

//cmd包下全局变量
var (
	START_TIME     time.Time   //启动时间
	globalCfg *gosupport.DataManage  //全局变量
	globalIniCfg *ini.Config //全局ini配置

)

func init()  {
	START_TIME = time.Now()

	globalCfg = gosupport.NewGlobalCfgSingleton()
	globalCfg.Set("SERVER_START_TIME", START_TIME) //服务启动时间
	globalCfg.Set("APP_CODE_VERSION", constants.APP_VERSION) //代码版本

}

//cmd公共的方法，在action方法中调用
func commonActionInit(iniFile string)  {
	//解析ini配置文件
	globalIniCfg = ini.NewIniConfig(iniFile)
	if err := globalIniCfg.ParseIniFile();err != nil {
		fmt.Println(err.Error())
	}
	globalCfg.Set("app_name", globalIniCfg.MustGetValue("default", "app_name"))

	//设置日志格式
	logrus.SetFormatter(new(mylog.LogFormatter))
	//设置日志级别
	logLevel, err := logrus.ParseLevel(globalIniCfg.MustGetValue("default", "log_level"))
	if err != nil {
		fmt.Println(err.Error())
	}
	logrus.SetLevel(logLevel)
	//设置日志输出多端
	targetDir := fmt.Sprintf("%s%s",globalIniCfg.MustGetValue("default", "log_dir"),globalIniCfg.MustGetValue("default", "app_name"))
	if !gosupport.IsDir(targetDir) {
		os.MkdirAll(targetDir, os.ModePerm)
	}
	logFileName := fmt.Sprintf("%s/%s.%s.log", targetDir,globalIniCfg.MustGetValue("default", "app_name"),time.Now().Format("2006-01-02"))
	writerF, err := os.OpenFile(logFileName, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0755)
	if err != nil {
		fmt.Printf("create file %s failed: %v", logFileName, err)
		os.Exit(0)
	}
	logrus.SetOutput(io.MultiWriter(os.Stdout, writerF))
	globalCfg.Set("SERVER_LOGRUS_LASTFILE_DATE", time.Now().Format("20060102"))


}
