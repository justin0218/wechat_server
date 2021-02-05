package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
	"time"
	"wechat_server/api/proto"
)

type TemplateRes struct {
	Errcode int    `json:"errcode"`
	Errmsg  string `json:"errmsg"`
}

func SendTemplate(data *proto.Template, accessToken string) (err error) {
	sendData, e := json.Marshal(data)
	if e != nil {
		err = e
		return
	}
	rurl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/message/template/send?access_token=%s", accessToken)
	res := TemplateRes{}
	_, _, errs := gorequest.New().Post(rurl).Timeout(time.Second*30).Set("Content-Type", "application/json").Send(string(sendData)).EndStruct(&res)
	if len(errs) != 0 {
		return errs[0]
	}
	if res.Errcode != 0 {
		err = fmt.Errorf(res.Errmsg)
		return
	}
	return nil
}
