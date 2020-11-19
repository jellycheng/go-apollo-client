package apollo

import (
	"encoding/json"
	"fmt"
	"github.com/sirupsen/logrus"
	"github.com/jellycheng/gosupport/curl"
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
		return body, nil
	} else {
		return body, err
	}

}

/**
应用感知配置更新
*/
func Notifications(apolloHost, appId, clusterName, notifications string) (bool, int64) {
	query := map[string]string{
		"appId":         appId,
		"cluster":       clusterName,
		"notifications": notifications,
	}
	url := fmt.Sprintf("%s/notifications/v2?%s",
						apolloHost, FormQuery(query))
	response, err := http.Get(url)
	if err != nil {
		logrus.Error("Notifications Get#" + err.Error())
	}
	var body []struct {
		Namespace      string `json:"namespace"`
		NotificationId int64  `json:"notificationId"`
		Messages       struct {
			Details map[string]int64 `json:"details"`
		} `json:"messages"`
	}

	if response.StatusCode == 200 {
		err = json.NewDecoder(response.Body).Decode(&body)
		if err != nil {
			logrus.Error("Notifications Decode#" + err.Error())
		}
		return true, body[0].NotificationId
	}
	return false, 0
}


func FormQuery(queryArr map[string]string) string {
	var qu []string
	for key, value := range queryArr {
		qu = append(qu, key+"="+value)
	}
	query, _ := url.ParseQuery(strings.Join(qu, "&"))
	return query.Encode()
}
