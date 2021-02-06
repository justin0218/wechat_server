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
	cret := &proto.R{Code: 600}
	ret.Res = cret
	accessToken, e := s.GetAccessToken(ctx, &proto.GetAccessTokenReq{Account: req.Account})
	if e != nil {
		cret.Msg = e.Error()
		ret.Res = cret
		return
	}
	if accessToken.Res.Code != 200 {
		cret.Msg = accessToken.Res.Msg
		ret.Res = cret
		return
	}
	req.Template.TemplateId = getTemplate(req.TemplateDefine)
	e = wechat.SendTemplate(req.Template, accessToken.AccessToken)
	if e != nil {
		cret.Msg = e.Error()
		ret.Res = cret
		return
	}
	cret.Code = 200
	ret.Res = cret
	return
}
