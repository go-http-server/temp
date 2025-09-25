package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/go-http-server/core/api"
	database "github.com/go-http-server/core/internal/database/sqlc"
	"github.com/go-http-server/core/plugin/pkg/mailer"
	"github.com/go-http-server/core/utils"
	"github.com/go-http-server/core/worker"
	"github.com/hibiken/asynq"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"golang.org/x/sync/errgroup"
)

var interruptSignals = []os.Signal{
	os.Interrupt,
	syscall.SIGTERM,
	syscall.SIGINT,
}

// @title [Go Core API]: Core API written by Go
// @version 1.0
// @accept json
// @contact.name API Support
// @contact.email namphamhai7@gmail.com
// @license.name MIT
// @license.url https://opensource.org/licenses/MIT
// @description Core API này là mình - Phạm Hải Nam làm ra với mục đích ban đầu là làm lại 1 khóa học mà mình đã học trước đó về golang. Mình chỉ tập trung vào phần mà mình chắc chắn sẽ sử dụng. Cái tên Core cũng là từ đó, core ở đây chỉ đơn giản là phần chắc chắn API nào cũng có, chẳng có gì cao siêu cả
// @host localhost:8080
// @BasePath /api/v1
func main() {
	env, err := utils.LoadEnviromentVariables("./")
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot load environment variables")
	}
	if env.ENVIRONMENT != utils.ProductionEnvironment {
		log.Logger = log.Output(zerolog.ConsoleWriter{Out: os.Stderr})
	}

	signalContext, stop := signal.NotifyContext(context.Background(), interruptSignals...)
	defer stop()

	pool, err := pgxpool.New(signalContext, env.DBSource)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create pool to database")
	}
	store := database.NewStore(pool)

	redisOpts := asynq.RedisClientOpt{
		Addr:     env.RedisServerAddress,
		Password: env.RedisServerPassword,
	}

	emailSender := mailer.NewGmailSender(env.EmailUsernameSender, env.EmailAddressSender, env.EmailPasswordSender)
	taskDistributor := worker.NewRedisTaskDistributor(redisOpts)
	bot := utils.NewBotTelegramService(env.TelegramBotToken, env.TelegramChatID)

	waitGroup, signalContext := errgroup.WithContext(signalContext)

	// create tasks processor serve run queue.
	runTaskProcessor(signalContext, waitGroup, redisOpts, store, emailSender, bot)

	// Create and start new gin server.
	server, err := api.NewServer(signalContext, waitGroup, store, env, taskDistributor)
	if err != nil {
		log.Fatal().Err(err).Msg("Cannot create new server")
	}

	server.StartServer(signalContext, waitGroup, env.HTTPServerAddress)

	err = waitGroup.Wait()
	if err != nil {
		log.Fatal().Err(err).Msg("Error from wait group")
	}
}

func runTaskProcessor(ctx context.Context, waitGroup *errgroup.Group, redisOpts asynq.RedisClientOpt, store database.Store, sender mailer.EmailSender, bot utils.Bot) {
	taskProcessor := worker.NewRedisTaskProcessor(redisOpts, store, sender, bot)
	log.Info().Msg("Start task processor")

	err := taskProcessor.Start()
	if err != nil {
		log.Fatal().Err(err).Msg("Failed start task processor")
	}

	waitGroup.Go(func() error {
		<-ctx.Done()
		log.Info().Msg("Graceful shutdown task processor")

		taskProcessor.Shutdown()
		log.Info().Msg("Task processor is stopped")
		return nil
	})
}
