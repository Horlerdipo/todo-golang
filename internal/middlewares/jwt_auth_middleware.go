package middlewares

import (
	"github.com/golang-jwt/jwt/v5"
	"github.com/horlerdipo/todo-golang/env"
	"github.com/horlerdipo/todo-golang/internal/database"
	"github.com/horlerdipo/todo-golang/utils"
	"golang.org/x/net/context"
	"net/http"
	"strings"
)

type contextKey string

type AuthDetails struct {
	UserId            uint
	JwtToken          string
	JwtExpirationTime *jwt.NumericDate
}

const UserKey contextKey = "user"

func JwtAuthMiddleware(tokenBlacklistRepository database.TokenBlacklistRepository) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			//log.Println("JwtAuthMiddleware hit:", r.URL.Path, r.Header.Get("Authorization"))
			//check if auth header exists
			header := r.Header.Get("Authorization")
			isQueryBasedAuth := false
			if header == "" {
				header = r.URL.Query().Get("_token")
				isQueryBasedAuth = true
			}

			if header == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
				return
			}

			//take code out of bearer string
			tokenString := ""
			if !isQueryBasedAuth {
				tokens := strings.Split(header, " ")
				if len(tokens) > 1 {
					tokenString = tokens[1]
				}

			} else {
				tokenString = header
			}

			if tokenString == "" {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
				return
			}

			//check code validity
			claim, err := utils.ValidateJwtToken(tokenString, env.FetchString("JWT_SECRET"))
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
				return
			}

			if claim == nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
				return
			}

			//check if token is not blacklisted
			isTokenBlackListed := tokenBlacklistRepository.CheckTokenExistence(r.Context(), tokenString)
			if isTokenBlackListed {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
				return
			}

			//add user details to Context
			expTime, err := claim.GetExpirationTime()
			if err != nil {
				utils.RespondWithError(w, http.StatusUnauthorized, "Unauthenticated", struct{}{})
				return
			}

			ctx := context.WithValue(r.Context(), UserKey, AuthDetails{
				UserId:            uint(claim["data"].(float64)),
				JwtToken:          tokenString,
				JwtExpirationTime: expTime,
			})
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}

}
