package middlewares

import (
	"fmt"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/utils"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type contextKey string

const UserKey contextKey = "user"

func JwtAuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("attempting to authenticate")
		//check if auth header exists
		header := r.Header.Get("Authorization")
		if header == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
		}

		//take code out of bearer string
		tokenString := strings.Split(header, " ")[1]
		if tokenString == "" {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
		}

		//check code validity
		claim, err := utils.ValidateJwtToken(tokenString, env.FetchString("JWT_SECRET"))
		if err != nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
		}

		if claim == nil {
			utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
		}

		//add user id to Context
		ctx := context.WithValue(r.Context(), UserKey, uint(claim["data"].(float64)))
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
