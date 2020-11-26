package apollo

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jellycheng/gosupport"
	"github.com/jellycheng/gosupport/curl"
	"github.com/sirupsen/logrus"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
	"time"
)


/*
   通过带缓存的Http接口从Apollo读取配置
*/
func GetConfig4Cache(apolloHost, appId, clusterName, namespace, ip string) (string, error) {
	body := ""
	url := fmt.Sprintf("%s/configfiles/json/%s/%s/%s",
						apolloHost, appId, clusterName, namespace)
	if ip != "" {
		url += "?ip=" + ip
	}
	req := curl.NewHttpRequest()
	//设置超时
	req.SetTimeout(int64(15 * time.Second))
	resObj, err := req.SetUrl(url).Get()
	if err==nil {
		body = resObj.GetBody()
		statusCode := resObj.GetRaw().StatusCode
		if statusCode == 404 {
			var envCfgMap map[string]interface{}
			errMsg := "获取失败"
			if err := json.Unmarshal([]byte(body), &envCfgMap);err == nil{
				errMsg = envCfgMap["message"].(string)
			}
			return body, errors.New(errMsg)
		}
		return body, nil
	} else {
		return body, err
	}

}

/**
 * 应用感知配置更新,服务端会hold住请求60秒
 * [{"namespaceName":"application","notificationId":-1}]
 * http://10.30.60.129:8080/notifications/v2?appId=SampleApp&cluster=default&notifications=%5B%7B%22namespaceName%22%3A%22application%22%2C%22notificationId%22%3A-1%7D%5D
 *
curl -X GET \
  'http://10.30.60.129:8080/notifications/v2?appId=SampleApp&cluster=default&notifications=%5B%7B%22namespaceName%22%3A%22application%22%2C%22notificationId%22%3A-1%7D%5D'

返回示例：
[{"namespaceName":"application","notificationId":4,"messages":{"details":{"SampleApp+default+application":4}}}]
*/
func Notifications(apolloHost, appId, clusterName, notifications string) (bool, int64, error) {
	query := map[string]string{
								"appId":         appId,
								"cluster":       clusterName,
								"notifications": notifications,
							}
	url := fmt.Sprintf("%s/notifications/v2?%s", apolloHost, FormQuery(query))
	response, err := http.Get(url)
	if err != nil {
		logrus.Error("Notifications Get#" + err.Error())
		return false, 0, err
	} else {
		logrus.Debug("apollo url: " + url)
	}
	var body []struct {
		Namespace      string `json:"namespace"`
		NotificationId int64  `json:"notificationId"`
		Messages       struct {
			Details map[string]int64 `json:"details"`
		} `json:"messages"`
	}
	//200有变化、304没有变化、401未授权
	if response.StatusCode == 200 {
		err = json.NewDecoder(response.Body).Decode(&body)
		if err != nil {
			logrus.Error("Notifications Decode#" + err.Error())
			return false, 0,err
		}
		return true, body[0].NotificationId,nil
	}
	logrus.Debug("apollo响应状态码：" + gosupport.ToStr(response.StatusCode))
	return false, 0,nil
}


func FormQuery(queryArr map[string]string) string {
	var qu []string
	for key, value := range queryArr {
		qu = append(qu, key+"="+value)
	}
	query, _ := url.ParseQuery(strings.Join(qu, "&"))
	return query.Encode()
}



func JsonToEnvContent(jsonStr string) string  {
	var retStr string
	if jsonStr == "" {
		return retStr
	}
	var envCfgMap map[string]interface{}
	var sbObj strings.Builder

	if err := json.Unmarshal([]byte(jsonStr), &envCfgMap);err == nil{
		for k, v:= range envCfgMap {
			line := fmt.Sprintf("%s=%v%s", k, v, LineBreak)
			sbObj.WriteString(line)
		}
		retStr = sbObj.String()
	}

	return retStr
}

func CleanContentWrite(file, wireteString string) error {
	var c = []byte(wireteString)
	err := ioutil.WriteFile(file, c, 0666);
	return err
}
