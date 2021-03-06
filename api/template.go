package api

import (
	"context"
	"wechat_server/api/proto"
	"wechat_server/pkg/wechat"
	"wechat_server/store"
)

func getTemplate(t proto.TemplateDefine) string {
	conf := new(store.Config)
	switch t {
	case proto.TemplateDefine_bill_notice:
		return conf.Get().MomoZaHuoPuWechat.BillNoticeTemplate
	}
	return ""
}

func (s *WechatSvr) SendTemplate(ctx context.Context, req *proto.SendTemplateReq) (ret *proto.SendTemplateRes, err error) {
	ret = new(proto.SendTemplateRes)
	accessToken, e := s.GetAccessToken(ctx, &proto.GetAccessTokenReq{Account: req.Account})
	if e != nil {
		err = e
		return
	}
	req.Template.TemplateId = getTemplate(req.TemplateDefine)
	e = wechat.SendTemplate(req.Template, accessToken.AccessToken)
	if e != nil {
		err = e
		return
	}
	return
}
