package config

//错误代号配置，常量首字母大写，驼峰式
const (
	Success             = 0
	Fail                = 1
	ErrorNotLogin       = 100
	ErrorRecordNotExist = 200
	ErrorParam          = 201
	ErrorIllegalRequest = 202
	ErrorDefault        = 10000 + iota
	ErrorDBNotFound
	ErrorRequestParam
	ErrorMysqlConnection
	ErrorRedisConnection
	ErrorOperationBusy
	ErrorSaveFailure
)

var codeMsg = map[int]string{
	Success:              "success",
	Fail:                 "fail",
	ErrorNotLogin:        "未登录",
	ErrorRecordNotExist:  "记录不存在",
	ErrorParam:           "参数不合法",
	ErrorIllegalRequest:  "非法请求",
	ErrorDefault:         "系统错误",
	ErrorDBNotFound:      "数据库不存在",
	ErrorRequestParam:    "请求参数错误",
	ErrorMysqlConnection: "mysql连接错误",
	ErrorRedisConnection: "redis连接错误",
	ErrorOperationBusy:   "操作过于频繁,请稍后再试",
	ErrorSaveFailure:     "数据保存失败,请稍后重试",
}

//通过代号获取错误信息
func GetMsg4Code(code int) string {
	if m, ok := codeMsg[code]; ok == true {
		return m
	}
	return codeMsg[ErrorDefault]
}

/**
调用示例
errMsg := config.GetMsg4Code(config.ErrorMysqlConnection)
fmt.Println(errMsg) //mysql连接错误
errMsg2 := config.GetMsg4Code(123)
fmt.Println(errMsg2) //系统错误
*/
