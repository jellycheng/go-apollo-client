.PHONY: all build devrun clean help go_build

# 定义变量
OUTPUTDIR="./output/"
BINARY="go-apollo-client"

all: clean go_build

build:
	go build -o ${OUTPUTDIR}${BINARY} main.go
	cp apollo_example.ini ${OUTPUTDIR}/apollo.ini

go_build: clean
	export GO111MODULE=on
	export GOPROXY=https://goproxy.cn
	env GOOS=linux GOARCH=amd64 go build -o ${OUTPUTDIR}${BINARY} main.go

devrun:
	@go run ./main.go artisan -c=${APOLLO_CONFIG} -env=${ENVCODE} -appid=${APPID} -cluster=${CLUSTER} -namespace=${NAMESPACE}

clean:
	@if [ -f ${OUTPUTDIR}${BINARY} ] ; then rm ${OUTPUTDIR}${BINARY} ; fi

pm2_start_job:
	pm2 start "${OUTPUTDIR}${BINARY} artisan -c=$(APOLLO_CONFIG) -f=MonitorApolloConsole -env=${ENVCODE} -appid=${APPID} -cluster=${CLUSTER} -namespace=${NAMESPACE} " --name=${PM2NAME}_apollo

pm2_restart_job:
	pm2 restart ${PM2NAME}_apollo

pm2_stop_job:
	pm2 stop ${PM2NAME}_apollo

pm2_delete_job:
	pm2 delete ${PM2NAME}_apollo

#显示命令帮助，如 make help
help:
	@echo "make - 删除原二进制文件后再编译生成新的linux二进制文件"
	@echo "make build - 本地开发环境编译 Go 代码, 生成二进制文件-适用于dev环境"
	@echo "make go_build - 编译Go代码-适用于linux环境的代码打包"
	@echo "APOLLO_CONFIG=apollo_example.ini ENVCODE=DEV APPID=SampleApp CLUSTER=default NAMESPACE=application make devrun - 直接运行 Go 代码-适用于dev环境"
	@echo "make clean - 移除二进制文件"
	@echo "PM2NAME=SampleApp make pm2_stop_job - 停止job服务"
	@echo "PM2NAME=SampleApp make pm2_restart_job - 重启job服务"
	@echo "APOLLO_CONFIG=./output/apollo.ini ENVCODE=DEV APPID=SampleApp CLUSTER=default NAMESPACE=application PM2NAME=SampleApp make pm2_start_job - 启动job服务"
