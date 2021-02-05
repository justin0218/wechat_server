package api

import (
	"context"
	"wechat_server/api/proto"
	"wechat_server/pkg/wechat"
)

func (s *WechatSvr) MakeShortUrl(ctx context.Context, req *proto.MakeShortUrlReq) (ret *proto.MakeShortUrlRes, err error) {
	ret = new(proto.MakeShortUrlRes)
	cret := &proto.R{Code: 600}
	ret.Res = cret
	shorUrl, _ := s.redis.Get().Get(s.redis.GetShortUrlKey(req.Url)).Result()
	if shorUrl != "" {
		cret.Code = 200
		ret.Res = cret
		ret.ShortUrl = shorUrl
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
	res, e := wechat.GetShortUrl(req.Url, accessToken.AccessToken)
	if e != nil {
		cret.Msg = e.Error()
		ret.Res = cret
		return
	}
	ret.ShortUrl = res.ShortUrl
	cret.Code = 200
	ret.Res = cret
	s.redis.Get().Set(s.redis.GetShortUrlKey(req.Url), res.ShortUrl, -1).Err()
	return
}
