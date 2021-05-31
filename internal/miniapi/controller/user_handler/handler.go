package user_handler

import (
	"github.com/xinliangnote/go-gin-api/internal/pkg/cache"
	"github.com/xinliangnote/go-gin-api/internal/pkg/core"
	"github.com/xinliangnote/go-gin-api/internal/pkg/db"
	"go.uber.org/zap"
)

var _ Handler = (*handler)(nil)

type Handler interface {

	i()

	//用户登录 微信登录
	// @Tags API.MiniAppUser
	// @Router /api/miniapp/loginbyweixin [post]
	LoginByWeiXin() core.HandlerFunc

}

type handler struct {
	logger       *zap.Logger
	db			 db.Repo
	cache        cache.Repo
}

func New(logger *zap.Logger, db db.Repo, cache cache.Repo) Handler {
	return &handler{
		logger:       logger,
		 db:    		db,
		cache:        cache,
	}
}

func (h *handler) i() {}