package middlewares

import (
	"context"
	"net/http"
	"strings"

	"github.com/DEV-BC/backend_chatapp/internal/utils"
)

const (
	CtxUserID          string = "user_id"
	CtxUserDisplayName string = "name"
	CtxPlatform        string = "X-Platform"
	CtxAuthorization   string = "Authorization"
	PlatformWeb        string = "web"
	PlatformMobile     string = "mobile"
)

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		authHeader := strings.TrimSpace(r.Header.Get(string(CtxAuthorization)))
		if authHeader == "" || !strings.HasPrefix(strings.ToLower(authHeader), "bearer ") {
			utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		platform := strings.ToLower(strings.TrimSpace(r.Header.Get(string(CtxPlatform))))
		if platform != PlatformWeb && platform != PlatformMobile {
			utils.JSON(w, http.StatusBadRequest, false, "invalid platform", nil)
			return
		}

		accessToken := strings.TrimSpace(authHeader[7:])

		userId, name, tokenPlatform, err := utils.VerifyJWT(accessToken)
		if err != nil {
			utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		if tokenPlatform != platform {
			utils.JSON(w, http.StatusUnauthorized, false, "Unauthorized", nil)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxUserID, userId)
		ctx = context.WithValue(ctx, CtxUserDisplayName, name)
		ctx = context.WithValue(ctx, CtxPlatform, tokenPlatform)

		next.ServeHTTP(w, r.WithContext(ctx))
	})

}
