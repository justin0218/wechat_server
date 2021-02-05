package wechat

import (
	"encoding/json"
	"fmt"
	"github.com/parnurzeal/gorequest"
)

type AccessToken struct {
	Errmsg      string `json:"errmsg"`
	Errcode     int    `json:"errcode"`
	AccessToken string `json:"access_token"`
	ExpiresIn   int    `json:"expires_in"`
}

func GetAccessToken(appid, secret string) (ret AccessToken, err error) {
	rurl := fmt.Sprintf("https://api.weixin.qq.com/cgi-bin/token?grant_type=client_credential&appid=%s&secret=%s", appid, secret)
	_, bytesRes, errs := gorequest.New().Get(rurl).EndBytes()
	if len(errs) > 0 {
		err = fmt.Errorf("wechat get access_token err:%v", errs)
		return
	}
	err = json.Unmarshal(bytesRes, &ret)
	if err != nil {
		err = fmt.Errorf("wechat get access_token err:%v", err)
		return
	}
	if ret.Errcode != 0 {
		err = fmt.Errorf("wechat get access_token err msg:%s", ret.Errmsg)
		return
	}
	return
}
