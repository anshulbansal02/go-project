package middlewares

import (
	"anshulbansal02/scribbly/internal/user"

	"context"
	"net/http"
	"strings"
)

type RequestHandler = func(w http.ResponseWriter, r *http.Request)

type UserAuthContext struct {
	UserId   string
	Token    string
	FailCode string
}

type CtxKey string

const UserCtxKey CtxKey = "user"

func Authenticate(userService *user.UserService) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			userCtx := UserAuthContext{}

			authHeader := r.Header.Get("Authorization")
			if len(authHeader) == 0 {
				userCtx.FailCode = "EmptyAuthHeader"
			} else {
				auth := strings.SplitN(authHeader, " ", 2)

				if len(auth) != 2 || auth[0] != "Bearer" {
					userCtx.FailCode = "InvalidAuthHeader"
				} else {
					userCtx.Token = auth[1]
					claims, err := userService.VerifyUserToken(userCtx.Token)
					if err != nil {
						userCtx.FailCode = "InvalidAuthToken"
					} else {
						userCtx.UserId = claims.UserId
					}

				}
			}

			ctx := context.WithValue(r.Context(), UserCtxKey, userCtx)

			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}

func WithAuthorization(next RequestHandler) RequestHandler {

	return func(w http.ResponseWriter, r *http.Request) {

		userCtx := r.Context().Value(UserCtxKey).(UserAuthContext)

		var failReason string

		if userCtx.FailCode == "EmptyAuthHeader" {
			failReason = "Empty \"Authorization\" header"
		} else if userCtx.FailCode == "InvalidAuthHeader" {
			failReason = "Invalid \"Authorization\" header format. Expected \"Bearer <Token>\""
		} else if userCtx.FailCode == "InvalidAuthToken" {
			failReason = "Invalid Authentication Token"
		}

		if failReason != "" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte(failReason))
		} else {
			next(w, r)
		}

	}
}
