package api

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	database "github.com/go-http-server/core/internal/database/sqlc"
	"github.com/go-http-server/core/utils"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store database.Store) *Server {
	env := utils.EnvironmentVariables{
		DBSource:          "",
		HTTPServerAddress: "",
		TimeExpiredToken:  30 * time.Minute,
	}

	testServer, err := NewServer(context.Background(), nil, store, env, nil)
	require.NoError(t, err)
	require.NotEmpty(t, testServer)

	return testServer
}

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	os.Exit(m.Run())
}
