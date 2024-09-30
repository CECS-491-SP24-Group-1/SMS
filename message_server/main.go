package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"time"

	"wraith.me/message_server/config"
	"wraith.me/message_server/consts"
	"wraith.me/message_server/db"
	"wraith.me/message_server/email"
	"wraith.me/message_server/mw"
	cr "wraith.me/message_server/redis"
	"wraith.me/message_server/router"
	"wraith.me/message_server/router/auth"
	"wraith.me/message_server/router/challenges"
	"wraith.me/message_server/task"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/redis/go-redis/v9"
)

const RQS_PER_MIN = 3

func main() {
	//Acquire a config instance, but cease further operation if an error occurred
	cfg, cfgErr := config.ConfigInit("")
	if cfgErr != nil {
		log.Panicf("Encountered unrecoverable error while loading config: %s\n", cfgErr.Error())
	}
	fmt.Printf("config:%+v\n", cfg)

	//Acquire an env instance, but cease further operation if an error occurred
	env, envErr := config.EnvInit("")
	if envErr != nil {
		log.Panicf("Encountered unrecoverable error while loading env: %s\n", envErr.Error())
	}
	fmt.Printf("env:%+v\n", env)

	//Setup scheduled tasks
	//TODO: Break this up into its own package
	//setupScheduledTasks()

	//Connect to MongoDB
	_, merr := db.GetInstance().Connect(&cfg.MongoDB)
	if merr != nil {
		panic(merr)
	}
	defer db.GetInstance().Disconnect()

	//Connect to Redis
	rclient, rerr := cr.GetInstance().Connect(&cfg.Redis)
	if rerr != nil {
		panic(rerr)
	}

	//Connect to the SMTP server
	_, eerr := email.GetInstance().Connect(&cfg.Email)
	if eerr != nil {
		panic(eerr)
	}

	//Setup scheduled tasks
	if err := setupScheduledTasks(rclient); err != nil {
		panic(err)
	}

	/*
		//Test listing
		names, _ := mclient.ListDatabaseNames(context.TODO(), bson.M{})
		for i, name := range names {
			fmt.Printf("DB #%d: %s\n", i, name)
		}
	*/

	//Setup Chi and start listening for connections
	r := setupServer(&cfg, &env)
	connStr := fmt.Sprintf("%s:%d", cfg.Server.BindAddr, cfg.Server.ListenPort)
	log.Printf("Listening on %s\n", connStr)
	http := http.Server{
		Addr:    connStr,
		Handler: r,
	}
	if err := http.ListenAndServe(); err != nil {
		panic(err)
	}
}

// TODO: Maybe add https://github.com/goware/firewall
func setupServer(cfg *config.Config, env *config.Env) chi.Router {
	//Setup Chi router
	r := chi.NewRouter()

	//Add essential 1p middlewares
	r.Use(middleware.RequestID)
	r.Use(middleware.RealIP)

	//Write the request ID to the response headers
	r.Use(mw.SendRequestID)

	//Add CORS
	//For more ideas, see: https://developer.github.com/v3/#cross-origin-resource-sharing
	r.Use(cors.Handler(cors.Options{
		AllowedOrigins: []string{
			"https://*",
			"http://*",
			//cfg.Client.BaseUrl,
		},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token", consts.TIMEZONE_OFFSET_HEADER}, // Ensure headers match client requests
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: true,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))

	//Perform access logging if its permitted
	/*
		if cfg.AccessLogs.Mode != alogs_t.OFF {
			r.Use(mw.NewZapMiddleware("router", &cfg.AccessLogs))
		}
	*/

	//HTTP rate-limiting
	/*
		limiter := mw.DefaultLimiter()
		limiter.Limit = 4
		//limiter.BurstLen = timex.Day * 6
		limiter.Debug = true
		r.Use(limiter.Impose)
	*/

	//Index route
	r.Get("/", router.Index)

	//Group subsequent routes under `/api`
	apir := chi.NewRouter()

	//Health route
	apir.Get("/heartbeat", router.Heartbeat)

	apir.Post("/send_message", router.SendMessage)

	//User auth routes
	apir.Mount("/auth", auth.AuthRoutes(cfg, env))

	//Challenge routes
	apir.Group(func(r chi.Router) {
		//authScopes := []token.TokenScope{token.TokenScopePOSTSIGNUP, token.TokenScopeUSER}
		//r.Use(mw.NewAuthMiddleware(authScopes))
		r.Mount("/challenges", challenges.ChallengeRoutes(env))
	})

	//Bind the API routes to the outgoing router
	r.Mount("/api", apir)

	//Return the built router for requests
	return r
}

func setupScheduledTasks(rc *redis.Client) error {
	//Create tasks
	purgeTask := task.PurgeOldUsersTask{
		TQ:  time.Minute * 5,
		RC:  rc,
		CTX: context.Background(),
	}

	//Setup the scheduler and run it
	sch := task.Scheduler{}
	if err := sch.Register(purgeTask); err != nil {
		return err
	}
	if err := sch.Start(); err != nil {
		return err
	}
	return nil
}

/*
func BookRoutes() chi.Router {
	r := chi.NewRouter()
	bookHandler := BookHandler{
		storage: BookStore{},
	}

	//Routes
	r.Get("/", bookHandler.ListBooks)
	r.Post("/", bookHandler.CreateBook)
	r.Get("/{id}", bookHandler.GetBooks)
	r.Put("/{id}", bookHandler.UpdateBook)
	r.Delete("/{id}", bookHandler.DeleteBook)
	return r
}
*/
