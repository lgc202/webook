package middleware

import (
	"github.com/gin-contrib/sessions"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

// LoginMiddlewareBuilder 设计出 builder 模式增强扩展性
type LoginMiddlewareBuilder struct {
	paths []string
}

func NewLoginMiddlewareBuilder() *LoginMiddlewareBuilder {
	return &LoginMiddlewareBuilder{}
}

func (l *LoginMiddlewareBuilder) IgnorePaths(path string) *LoginMiddlewareBuilder {
	l.paths = append(l.paths, path)
	return l
}

func (l *LoginMiddlewareBuilder) Build() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// 不需要登录校验的
		for _, path := range l.paths {
			if ctx.Request.URL.Path == path {
				return
			}
		}
		sess := sessions.Default(ctx)
		id := sess.Get("userId")
		if id == nil {
			// 没有登录
			ctx.AbortWithStatus(http.StatusUnauthorized)
			return
		}

		// 支持 Session 刷新, 每次都需要访问 redis，所以才有后面的 jwt
		// 我怎么知道一分钟已经过去了
		updateTime := sess.Get("update_time")
		sess.Set("userId", id)
		sess.Options(sessions.Options{MaxAge: 60})
		now := time.Now().UnixMilli()
		// 刚登录，还没刷新过
		if updateTime == nil {
			sess.Set("update_time", now)
			sess.Options(sessions.Options{MaxAge: 60})
			sess.Save()
			return
		}

		updateTimeVal, ok := updateTime.(int64)
		if !ok {
			ctx.AbortWithStatus(http.StatusInternalServerError)
			return
		}

		// 大于 10 秒
		if now-updateTimeVal > 10*1000 {
			sess.Set("update_time", now)
			sess.Save()
			return
		}
	}
}
