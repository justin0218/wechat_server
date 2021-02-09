package api

import (
	"context"
	"wechat_server/api/proto"
	"wechat_server/pkg/wechat"
)

func (s *WechatSvr) MakeShortUrl(ctx context.Context, req *proto.MakeShortUrlReq) (ret *proto.MakeShortUrlRes, err error) {
	ret = new(proto.MakeShortUrlRes)
	shorUrl, _ := s.redis.Get().Get(s.redis.GetShortUrlKey(req.Url)).Result()
	if shorUrl != "" {
		ret.ShortUrl = shorUrl
		return
	}
	accessToken, e := s.GetAccessToken(ctx, &proto.GetAccessTokenReq{Account: req.Account})
	if e != nil {
		err = e
		return
	}
	res, e := wechat.GetShortUrl(req.Url, accessToken.AccessToken)
	if e != nil {
		err = e
		return
	}
	ret.ShortUrl = res.ShortUrl
	s.redis.Get().Set(s.redis.GetShortUrlKey(req.Url), res.ShortUrl, -1).Err()
	return
}
