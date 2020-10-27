package middleware

import (
	"context"
	"fmt"
	"github.com/google/logger"
	"github.com/technoZoomers/MasterHubBackend/models"
	"github.com/technoZoomers/MasterHubBackend/useCases"
	"github.com/technoZoomers/MasterHubBackend/utils"
	"net/http"
)

type AuthMidleware struct {
	 middlewares *Middlewares
	 UserUC         useCases.UsersUCInterface
}

func (authMiddleware *AuthMidleware) Auth(httpHandler http.HandlerFunc, passNext bool) http.Handler {
	return http.HandlerFunc(func(writer http.ResponseWriter, req *http.Request) {
		ctx := req.Context()
		cookie, err := req.Cookie(authMiddleware.middlewares.cookieString)
		fmt.Println(cookie)
		if err != nil {
			ctx = context.WithValue(ctx, authMiddleware.middlewares.contextAuthorisedKey, false)
			authMiddleware.passNext(passNext, httpHandler, writer, req, ctx)
			return
		}
		var user models.User
		err = authMiddleware.UserUC.GetUserByCookie(cookie.Value, &user)
		if err != nil {
			ctx = context.WithValue(ctx, authMiddleware.middlewares.contextAuthorisedKey, false)
			authMiddleware.passNext(passNext, httpHandler, writer, req, ctx)
			return
		}
		ctx = context.WithValue(ctx, authMiddleware.middlewares.contextAuthorisedKey, true)
		ctx = context.WithValue(ctx, authMiddleware.middlewares.contextUserKey, user)
		ctx = context.WithValue(ctx, authMiddleware.middlewares.contextCookieKey, cookie.Value)

		httpHandler.ServeHTTP(writer, req.WithContext(ctx))
	})
}

func (authMiddleware *AuthMidleware) passNext(passNext bool, httpHandler http.HandlerFunc, writer http.ResponseWriter, req *http.Request, ctx context.Context) {
	if passNext {
		httpHandler.ServeHTTP(writer, req.WithContext(ctx))
	} else {
		sessionError := fmt.Errorf("failed to authorise user")
		logger.Errorf(sessionError.Error())
		utils.CreateErrorAnswerJson(writer, http.StatusUnauthorized, models.CreateMessage(sessionError.Error()))
	}
}