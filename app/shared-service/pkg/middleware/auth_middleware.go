package middleware

import (
	"net/http"
	"shared/pkg/models"
	"shared/pkg/utils"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

func AuthMiddleware(secretKey string) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		authHeader := ctx.GetHeader("Authorization")
		if authHeader == "" {
			utils.RespondWithError(ctx, http.StatusUnauthorized, "No token provided")
			ctx.Abort()
			return
		}

		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			utils.RespondWithError(ctx, http.StatusUnauthorized, "Invalid token format")
			ctx.Abort()
			return
		}

		tokenString := parts[1]

		token, err := jwt.ParseWithClaims(tokenString, &models.TokenClaims{}, func(token *jwt.Token) (any, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				utils.RespondWithError(ctx, http.StatusUnauthorized, "Invalid token format")
				ctx.Abort()
				return nil, nil
			}
			return []byte(secretKey), nil
		})

		if err != nil {
			utils.RespondWithError(ctx, http.StatusUnauthorized, "Invalid token: "+err.Error())
			ctx.Abort()
			return
		}

		if claims, ok := token.Claims.(*models.TokenClaims); ok && token.Valid {
			ctx.Set("userID", claims.UserID)
			ctx.Set("username", claims.Username)
			ctx.Next()
		} else {
			utils.RespondWithError(ctx, http.StatusUnauthorized, "Invalid token claims")
			ctx.Abort()
			return
		}
	}
}
