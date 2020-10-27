package middleware

import (
	"github.com/technoZoomers/MasterHubBackend/useCases"
)

type Middlewares struct {
	AuthMiddleware *AuthMidleware
	cookieString string
	contextUserKey string
	contextCookieKey string
}

func (middlewares *Middlewares) Init(usersUC useCases.UsersUCInterface) error {
	middlewares.AuthMiddleware = &AuthMidleware{middlewares: middlewares, UserUC: usersUC}
	middlewares.cookieString = "user_session"
	middlewares.contextCookieKey = "cookie_key"
	middlewares.contextUserKey = "user_key"
	return nil
}