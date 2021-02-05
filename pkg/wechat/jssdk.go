package wechat

import (
	"crypto/sha1"
	"fmt"
	"time"
	"wechat_server/pkg/tool"
)

type Jssdk struct {
	Appid     string `json:"appid"`
	Noncestr  string `json:"noncestr"`
	Timestamp int64  `json:"timestamp"`
	Signature string `json:"signature"`
}

//jsapi_ticket=sM4AOVdWfPE4DxkXGEs8VMCPGGVi4C3VM0P37wVUCFvkVAy_90u5h9nbSlYy3-Sl-HhTdfl2fzFy1AOcHKP7qg&
//=Wm3WZYTPz0wzccnW&timestamp=1414587457&url=http://mp.weixin.qq.com?params=value

func GetJssdk(url, ticket string) (ret Jssdk, err error) {
	ret.Noncestr = tool.RandomStr(16)
	ret.Timestamp = time.Now().Unix()
	signStr := fmt.Sprintf("jsapi_ticket=%s&noncestr=%s&timestamp=%d&url=%s", ticket, ret.Noncestr, ret.Timestamp, url)
	ret.Signature = signature(signStr)
	return
}
func signature(signStr string) string {
	d := sha1.New()
	d.Write([]byte(signStr))
	l := fmt.Sprintf("%x", d.Sum(nil))
	return l
}
