package main

import (
	"github.com/urfave/cli"
	"go-apollo-client/cmd"
	"go-apollo-client/constants"
	"os"
	"runtime"
)

func init()  {
	runtime.GOMAXPROCS(runtime.NumCPU())

}

func main()  {
	var AppVer string = constants.APP_VERSION
	app := cli.NewApp()
	app.Name = "go-apollo-client"
	app.Usage = "阿波罗客户端"
	app.Version = AppVer
	app.Author = "cjs"
	app.Email = "42282367@qq.com"
	app.Commands = []cli.Command{
		cmd.Artisan,
	}
	app.Flags = append(app.Flags, []cli.Flag{}...)
	app.Run(os.Args)

}
