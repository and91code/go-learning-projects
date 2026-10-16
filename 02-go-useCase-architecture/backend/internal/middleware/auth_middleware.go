package middleware

import (
	"errors"
	"fmt"
	"net/http"
	"os"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v5"
)

const UserIDKey = "userID"

func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 1. Extrai o cabeçalho Authorization
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Cabeçalho Authorization é obrigatório"})
			return
		}

		// 2. Valida o formato "Bearer <token>"
		parts := strings.Split(authHeader, " ")
		if len(parts) != 2 || parts[0] != "Bearer" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Formato de token inválido. Use 'Bearer <token>'"})
			return
		}

		tokenString := parts[1]
		secret := []byte(os.Getenv("JWT_SECRET"))

		if len(secret) == 0 {
			c.AbortWithStatusJSON(http.StatusInternalServerError, gin.H{"error": "Chave secreta JWT não configurada no servidor"})
			return
		}

		// 3. Valida e faz o parse do Token JWT
		token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
			// Garante que o método de assinatura é HMAC (HS256)
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("método de assinatura inesperado: %v", token.Header["alg"])
			}
			return secret, nil
		})

		if err != nil || !token.Valid {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Token JWT inválido ou expirado"})
			return
		}

		// 4. Extrai os claims e armazena o userID no contexto
		claims, ok := token.Claims.(jwt.MapClaims)

		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Payload do token inválido"})
			return
		}

		// No JSON/JWT, números são desserializados como float64
		userIDFloat, ok := claims["sub"].(float64)
		if !ok {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Identificador de usuário (sub) ausente ou inválido no token"})
			return
		}

		// Injeta o userID no contexto para os próximos Handlers utilizarem (c.Get / c.MustGet)
		c.Set(UserIDKey, int64(userIDFloat))

		// 5. Prossegue para o próximo Handler da rota
		c.Next()
	}
}

// GetUserIDFromContext é uma função utilitária para extrair com segurança o ID do usuário nos Handlers
func GetUserIDFromContext(c *gin.Context) (int64, error) {
	val, exists := c.Get(UserIDKey)
	if !exists {
		return 0, errors.New("id do usuário não encontrado no contexto da requisição")
	}

	userID, ok := val.(int64)
	if !ok {
		return 0, errors.New("tipo inválido para o id do usuário no contexto")
	}

	return userID, nil
}
