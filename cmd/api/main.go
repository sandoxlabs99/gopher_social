package main

import (
	"expvar"
	"gopher_social/internal/auth"
	"gopher_social/internal/db"
	"gopher_social/internal/env"
	"gopher_social/internal/mailer"
	"gopher_social/internal/ratelimiter"
	"gopher_social/internal/store"
	"gopher_social/internal/store/cache"
	"runtime"
	"time"

	"github.com/joho/godotenv"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

const version = "0.0.1" // semvar (semantic versioning)

func main() {
	// Logger
	logger := zap.Must(zap.NewProduction()).Sugar()
	defer logger.Sync()

	err := godotenv.Load()
	if err != nil {
		logger.Fatal("Error loading .env file")
	}

	cfg := config{
		serverAddr:  env.GetString("SERVER_ADDR", ":8080"),
		apiURL:      env.GetString("EXTERNAL_URL", "localhost:8080"),
		frontendURL: env.GetString("FRONTEND_URL", "http://localhost:3000"),
		namespace:   env.GetString("NAMESPACE_ENV", "development"),
		db: dbConfig{
			addr:         env.GetString("DATABASE_URL", "postgres://admin:adm1npass@localhost/social?sslmode=disable"),
			maxOpenConns: env.GetInt("DB_MAX_OPEN_CONNS", 30),
			maxIdleConns: env.GetInt("DB_MAX_IDLE_CONNS", 10),
			maxIdleTime:  env.GetDuration("DB_MAX_IDLE_TIME", "15m"),
			maxLifeTime:  env.GetDuration("DB_MAX_LIFE_TIME", "1h"),
		},
		mail: mailConfig{
			exp:       time.Minute * 15, // 15 minutes,
			fromEmail: env.GetString("FROM_EMAIL", "Acme <onboarding@resend.dev>"),
			sendGrid: sendGridConfig{
				apiKey: env.GetString("SENDGRID_API_KEY", ""),
			},
			mailTrap: mailTrapConfig{
				apiKey: env.GetString("MAILTRAP_API_KEY", ""),
			},
			resend: resendConfig{
				apiKey: env.GetString("RESEND_API_KEY", ""),
			},
		},
		auth: authConfig{
			basic: basicConfig{
				user: env.GetString("AUTH_BASIC_USER", "admin"),
				pass: env.GetString("AUTH_BASIC_PASS", ""),
			},
			token: tokenConfig{
				secret: env.GetString("AUTH_TOKEN_SECRET", "example"),
				exp:    time.Hour * 24, // 1 day
				iss:    "gophersocial",
			},
		},
		redis: redisConfig{
			isEnabled: env.GetBool("IS_REDIS_ENABLED", true),
			addr:      env.GetString("REDIS_ADDR", "localhost:6379"),
			pwd:       env.GetString("REDIS_PWD", ""),
			db:        env.GetInt("REDIS_DB", 0),
		},
		rateLimiter: ratelimiter.Config{
			RequestsPerTimeFrame: env.GetInt("RATELIMITER_REQS_COUNT", 20),
			TimeFrame:            5 * time.Second,
			IsEnabled:            env.GetBool("IS_RATELIMITER_ENABLED", true),
		},
	}

	// Database
	db, err := db.NewConn(
		cfg.db.addr,
		cfg.db.maxOpenConns,
		cfg.db.maxIdleConns,
		cfg.db.maxIdleTime,
		cfg.db.maxLifeTime,
	)
	if err != nil {
		logger.Fatal(err)
	}

	defer db.Close()
	logger.Info("Database connection pool established")

	store := store.NewStorage(db)

	// Redis Cache
	var redisDB *redis.Client

	if cfg.redis.isEnabled {
		redisDB = cache.NewRedisClient(cfg.redis.addr, cfg.redis.pwd, cfg.redis.db)
		logger.Info("Redis connection established")
		defer redisDB.Close()
	}

	redisStore := cache.NewRedisStorage(redisDB)

	// mailer := mailer.NewSendGrid(cfg.mail.sendGrid.apiKey, cfg.mail.fromEmail, logger)
	// mailTrap, err := mailer.NewMailTrapClient(cfg.mail.mailTrap.apiKey, cfg.mail.fromEmail)
	// if err != nil {
	// 	logger.Fatal(err)
	// }
	resend, err := mailer.NewResendClient(cfg.mail.resend.apiKey, cfg.mail.fromEmail, logger)
	if err != nil {
		logger.Fatal(err)
	}

	// rate limiter
	rateLimiter := ratelimiter.NewFixedWindowLimiter(
		cfg.rateLimiter.RequestsPerTimeFrame,
		cfg.rateLimiter.TimeFrame,
	)

	jwtAuthenticator := auth.NewJWTAuthenticator(cfg.auth.token.secret, cfg.auth.token.iss, cfg.auth.token.iss)

	app := &application{
		config:        cfg,
		store:         store,
		logger:        logger,
		mailer:        resend,
		authenticator: jwtAuthenticator,
		cacheStorage:  redisStore,
		rateLimiter:   rateLimiter,
	}

	// expvar metrics collected
	expvar.NewString("version").Set(version)
	expvar.Publish("database", expvar.Func(func() any {
		return db.Stats()
	}))

	expvar.Publish("goroutines", expvar.Func(func() any {
		return runtime.NumGoroutine()
	}))

	mux := app.mount()
	logger.Fatal(app.run(mux))
}
