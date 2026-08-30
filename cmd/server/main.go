package main

import (
	"context"
	"fmt"
	"time"

	"github.com/jackc/pgx/v5/pgxpool"
	_ "github.com/lib/pq"

	"log"
	"os"

	"github.com/golang-migrate/migrate/v4"
	_ "github.com/golang-migrate/migrate/v4/database/postgres"
	_ "github.com/golang-migrate/migrate/v4/source/file"

	db "github.com/meads/notes-api/internal/db/sqlc"
	"github.com/meads/notes-api/internal/handler"
	"github.com/meads/notes-api/internal/repository"
	"github.com/meads/notes-api/internal/security"
	"github.com/meads/notes-api/internal/service"
)

func dbConnect(ctx context.Context, retries int, dbConnectionString string) *pgxpool.Pool {
	pool, err := pgxpool.New(ctx, dbConnectionString)
	if err != nil {
		log.Fatalf("Unable to connect to the database %v\n", err)
	}
	retryCount := 0
	for {
		err := pool.Ping(ctx)
		if err != nil {
			retryCount += 1
			time.Sleep(time.Second * 2)
			if retryCount == retries {
				log.Fatal(err)
			}
			fmt.Printf("Database connect retry attempt %d...\n", retryCount)
		} else {
			break
		}
	}

	return pool
}

func parseDurationFromEnv(name string, durationDefault time.Duration) time.Duration {
	valStr := os.Getenv(name)

	var dur time.Duration
	var err error

	if valStr == "" {
		dur = durationDefault
	} else {
		dur, err = time.ParseDuration(valStr)
		if err != nil {
			fmt.Printf("Invalid %s value: %v\n", name, err)
			dur = durationDefault
		}
	}
	return dur
}

func main() {
	dbConnectionString := os.Getenv("DATABASE_URL")
	secretKey := os.Getenv("SECRET_KEY")
	cookieDomain := os.Getenv("COOKIE_DOMAIN")
	accessTokenDuration := parseDurationFromEnv("ACCESS_TOKEN_DURATION", 15*time.Minute)
	refreshTokenDuration := parseDurationFromEnv("REFRESH_TOKEN_DURATION", 24*time.Hour)

	ctx := context.Background()
	retries := 10
	pool := dbConnect(ctx, retries, dbConnectionString)
	defer pool.Close()

	m, err := migrate.New("file://internal/db/migration", dbConnectionString)
	if err != nil {
		log.Fatalf("error calling New with sql-migration tool: %s", err)
		return
	}
	m.Up()
	fmt.Print("\nmigrations were a success. 🎉\n")

	tokener := security.NewTokenManager(secretKey)
	hasher := security.NewHasher()
	queries := db.New(pool)

	sessionRepo := repository.NewSessionRepository(queries)
	userRepo := repository.NewUserRepository(queries)

	authService := service.NewAuthService(
		userRepo, sessionRepo, tokener, hasher,
		accessTokenDuration, refreshTokenDuration)
	authHandler := handler.NewAuthHandler(authService, cookieDomain)

	userService := service.NewUserService(userRepo, hasher)
	userHandler := handler.NewUserHandler(userService)

	noteRepo := repository.NewNoteRepository(queries)
	noteService := service.NewNoteService(noteRepo, userRepo)
	noteHandler := handler.NewNoteHandler(noteService)

	router := handler.SetupRouter(authHandler, userHandler, noteHandler, tokener)

	err = router.Run(":" + os.Getenv("PORT"))
	if err != nil {
		log.Fatal("cannot start server: ", err)
	}
}
