package api

import (
	"context"
	"wechat_server/api/proto"
	"wechat_server/pkg/wechat"
)

func (s *WechatSvr) DoPay(ctx context.Context, req *proto.DoPayReq) (ret *proto.DoPayRes, err error) {
	appid := s.conf.Get().MomoZaHuoPuWechat.Appid
	mchID := s.conf.Get().MomoZaHuoPuWechat.MchID
	mchApiKey := s.conf.Get().MomoZaHuoPuWechat.MchApiKey
	payRet, _, e := wechat.DoPay(appid, mchID, mchApiKey, req.Openid, req.OrderCode, req.Body, req.ClientIp, int(req.Price), req.NotifyUrl, req.TradeType)
	if e != nil {
		err = e
		return
	}
	jsapi := wechat.GetJsapiSign("prepay_id="+payRet.PrepayId, appid, mchApiKey)
	ret = new(proto.DoPayRes)
	ret.Timestamp = jsapi.TimeStamp
	ret.NonceStr = jsapi.NonceStr
	ret.Package = jsapi.Package
	ret.SignType = jsapi.SignType
	ret.PaySign = jsapi.PaySign
	return
}
