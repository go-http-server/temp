package database

import (
	"context"
	"os"
	"testing"
	"time"

	"github.com/go-http-server/core/utils"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog/log"
)

var testStore Store

func TestMain(m *testing.M) {
	loc, err := time.LoadLocation("Asia/Ho_Chi_Minh")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load location timezone")
	}
	time.Local = loc
	env, err := utils.LoadEnviromentVariables("../../../")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load environment variables file")
	}

	connectionPool, err := pgxpool.New(context.Background(), env.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create connection pool into database")
	}
	defer connectionPool.Close()
	testStore = NewStore(connectionPool)
	os.Exit(m.Run())
}
