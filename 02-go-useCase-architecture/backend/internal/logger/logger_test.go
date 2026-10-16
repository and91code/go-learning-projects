//go:build unit

package logger

import (
	"bytes"
	"log"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

// captureLog redireciona o logger padrão para um buffer durante o teste.
func captureLog(t *testing.T) *bytes.Buffer {
	t.Helper()

	var buffer bytes.Buffer
	previousOutput, previousFlags := log.Writer(), log.Flags()
	log.SetOutput(&buffer)
	log.SetFlags(0)
	t.Cleanup(func() {
		log.SetOutput(previousOutput)
		log.SetFlags(previousFlags)
	})
	return &buffer
}

// Cenário: cada nível de log adiciona seu prefixo e formata a mensagem.
func TestLogLevels_AddPrefix(t *testing.T) {
	buffer := captureLog(t)

	Infof("olá %s", "mundo")
	Warnf("atenção %d", 1)
	Errorf("falha %s", "x")

	assert.Equal(t, "INFO: olá mundo\nWARN: atenção 1\nERROR: falha x\n", buffer.String())
}

// Cenário: Init configura o logger padrão com timestamp UTC e arquivo/linha.
func TestInit_ConfiguresFlags(t *testing.T) {
	previousOutput, previousFlags := log.Writer(), log.Flags()
	t.Cleanup(func() {
		log.SetOutput(previousOutput)
		log.SetFlags(previousFlags)
	})

	Init()

	assert.Equal(t, log.Ldate|log.Ltime|log.LUTC|log.Lshortfile, log.Flags())
}

// Cenário: o middleware registra método, caminho e status da requisição sem alterar a resposta.
func TestGinMiddleware_LogsRequest(t *testing.T) {
	buffer := captureLog(t)
	gin.SetMode(gin.TestMode)
	router := gin.New()
	router.Use(GinMiddleware())
	router.GET("/ping", func(c *gin.Context) { c.Status(http.StatusTeapot) })

	response := httptest.NewRecorder()
	router.ServeHTTP(response, httptest.NewRequest(http.MethodGet, "/ping", nil))

	assert.Equal(t, http.StatusTeapot, response.Code)
	assert.Contains(t, buffer.String(), "INFO: GET /ping status=418")
}
