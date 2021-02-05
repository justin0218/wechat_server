package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

type Ticket struct {
	Errmsg  string `json:"errmsg"`
	Errcode int    `json:"errcode"`
	Ticket  string `json:"ticket"`
}

func GetTicket(accessToken string) (ret Ticket, err error) {
	rurl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/ticket/getticket?access_token=%s&type=jsapi", accessToken)
	_, bytesRes, errs := gorequest.New().Get(rurl).EndBytes()
	if len(errs) > 0 {
		err = fmt.Errorf("wechat get ticket err:%v", errs)
		return
	}
	err = json.Unmarshal(bytesRes, &ret)
	if err != nil {
		return
	}
	if ret.Errcode != 0 {
		err = fmt.Errorf(ret.Errmsg)
		return
	}
	return
}
