package api

import (
	"context"
	"fmt"
	"time"
	"wechat_server/api/proto"
	"wechat_server/pkg/wechat"
)

func (s *WechatSvr) GetAuthAccessToken(ctx context.Context, req *proto.GetAuthAccessTokenReq) (ret *proto.GetAuthAccessTokenRes, err error) {
	ret = new(proto.GetAuthAccessTokenRes)
	cret := &proto.R{Code: 600}
	ret.Res = cret
	appid, secret := "", ""
	if req.Account == proto.Account_momo_za_huo_pu {
		appid = s.conf.Get().MomoZaHuoPuWechat.Appid
		secret = s.conf.Get().MomoZaHuoPuWechat.Secret
	}
	res, e := wechat.GetAuthAccessToken(appid, secret, req.Code)
	if e != nil {
		cret.Msg = e.Error()
		return
	}
	ret.AccessToken = res.AccessToken
	ret.Openid = res.Openid
	ret.RefreshToken = res.RefreshToken
	ret.Scope = res.Scope
	cret.Code = 200
	ret.Res = cret
	return
}

func (s *WechatSvr) GetUserInfo(ctx context.Context, req *proto.GetUserInfoReq) (ret *proto.GetUserInfoRes, err error) {
	ret = new(proto.GetUserInfoRes)
	cret := &proto.R{Code: 600}
	ret.Res = cret
	res, e := wechat.GetUserInfo(req.Openid, req.AuthAccessToken)
	if e != nil {
		cret.Msg = e.Error()
		return
	}
	ret.Openid = res.Openid
	ret.Nickname = res.Nickname
	ret.Headimgurl = res.Headimgurl
	cret.Code = 200
	ret.Res = cret
	return
}

func (s *WechatSvr) GetAccessToken(ctx context.Context, req *proto.GetAccessTokenReq) (ret *proto.GetAccessTokenRes, err error) {
	ret = new(proto.GetAccessTokenRes)
	cret := &proto.R{Code: 600}
	ret.Res = cret
	appid, secret := "", ""
	if req.Account == proto.Account_momo_za_huo_pu {
		appid = s.conf.Get().MomoZaHuoPuWechat.Appid
		secret = s.conf.Get().MomoZaHuoPuWechat.Secret
	}
	cacheToken, _ := s.redis.Get().Get(s.redis.GetAccessTokenKey(appid)).Result()
	if cacheToken != "" {
		cret.Code = 200
		ret.Res = cret
		ret.AccessToken = cacheToken
		return
	}
	res, e := wechat.GetAccessToken(appid, secret)
	if e != nil {
		cret.Msg = e.Error()
		return
	}
	cret.Code = 200
	ret.Res = cret
	ret.AccessToken = res.AccessToken
	_, _ = s.redis.Get().Set(s.redis.GetAccessTokenKey(appid), res.AccessToken, 7000*time.Second).Result()
	return
}

func (s *WechatSvr) GetTicket(ctx context.Context, req *proto.GetTicketReq) (ret *proto.GetTicketRes, err error) {
	ret = new(proto.GetTicketRes)
	cret := &proto.R{Code: 600}
	ret.Res = cret
	appid := ""
	if req.Account == proto.Account_momo_za_huo_pu {
		appid = s.conf.Get().MomoZaHuoPuWechat.Appid
	}
	cacheTicket, _ := s.redis.Get().Get(s.redis.GetAccessTicketKey(appid)).Result()
	if cacheTicket != "" {
		cret.Code = 200
		ret.Res = cret
		ret.Ticket = cacheTicket
		return
	}
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
	res, e := wechat.GetTicket(accessToken.AccessToken)
	if e != nil {
		cret.Msg = e.Error()
		return
	}
	cret.Code = 200
	ret.Res = cret
	ret.Ticket = res.Ticket
	_, _ = s.redis.Get().Set(s.redis.GetAccessTicketKey(appid), res.Ticket, 7000*time.Second).Result()
	return
}

func (s *WechatSvr) GetJssdk(ctx context.Context, req *proto.GetJssdkReq) (ret *proto.GetJssdkRes, err error) {
	ret = new(proto.GetJssdkRes)
	cret := &proto.R{Code: 600}
	ret.Res = cret
	appid := ""
	if req.Account == proto.Account_momo_za_huo_pu {
		appid = s.conf.Get().MomoZaHuoPuWechat.Appid
	}
	ticket, e := s.GetTicket(ctx, &proto.GetTicketReq{Account: req.Account})
	if e != nil {
		cret.Msg = e.Error()
		ret.Res = cret
		return
	}
	if ticket.Res.Code != 200 {
		cret.Msg = ticket.Res.Msg
		ret.Res = cret
		return
	}
	res, e := wechat.GetJssdk(req.Url, ticket.Ticket)
	if e != nil {
		cret.Msg = e.Error()
		return
	}
	cret.Code = 200
	ret.Res = cret
	ret.Appid = appid
	ret.Noncestr = res.Noncestr
	ret.Timestamp = fmt.Sprintf("%d", res.Timestamp)
	ret.Signature = res.Signature
	return
}
