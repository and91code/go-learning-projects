// Package logger centraliza o logging da aplicação e das requisições HTTP.
package logger

import (
	"log"
	"os"
	"time"

	"github.com/gin-gonic/gin"
)

// Init configura o logger padrão com timestamp e localização no código.
func Init() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.LUTC | log.Lshortfile)
}

// Infof registra uma mensagem informativa formatada.
func Infof(format string, args ...any) {
	log.Printf("INFO: "+format, args...)
}

// Warnf registra um aviso formatado.
func Warnf(format string, args ...any) {
	log.Printf("WARN: "+format, args...)
}

// Errorf registra um erro formatado.
func Errorf(format string, args ...any) {
	log.Printf("ERROR: "+format, args...)
}

// GinMiddleware registra método, caminho, status e duração de cada requisição.
func GinMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		startedAt := time.Now()
		c.Next()

		Infof("%s %s status=%d latency=%s client_ip=%s",
			c.Request.Method,
			c.Request.URL.Path,
			c.Writer.Status(),
			time.Since(startedAt),
			c.ClientIP(),
		)
	}
}
