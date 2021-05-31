package user_handler

import "github.com/xinliangnote/go-gin-api/internal/pkg/core"

type loginRequest struct {
	OpenID string `form:"openid"` // openid
	Avatar string `form:"avatar"` // 头像
	NickName string `form:"nickname"` // 昵称
	PhoneNumber string `form:"phonenumber"` //手机号码
}

type loginReponse struct {
	Code int32 `json:"code"` //错误码
	Data string `json:"data"` // 内容
	Msg  string `json:"msg"`
}

func (h *handler) LoginByWeiXin() core.HandlerFunc {

	return func(c core.Context) {

		//req := new(loginRequest)
		res := new(loginReponse)
		res.Data = "login success mock"
		res.Code = 0
		res.Msg = "请求成功"

		c.Payload(res)
	}
}
