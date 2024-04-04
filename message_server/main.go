package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"wraith.me/message_server/config"
	"wraith.me/message_server/db"
	"wraith.me/message_server/email"
	"wraith.me/message_server/mw"
	"wraith.me/message_server/obj"
	credis "wraith.me/message_server/redis"
	"wraith.me/message_server/router"
	"wraith.me/message_server/router/challenges"
	"wraith.me/message_server/router/users"

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
	//TODO: pass config instance eventually
	//setupScheduledTasks()

	//Connect to MongoDB
	mclient, merr := db.GetInstance().Connect(&cfg.MongoDB)
	if merr != nil {
		panic(merr)
	}
	defer db.GetInstance().Disconnect()

	//Connect to Redis
	rclient, rerr := credis.GetInstance().Connect(&cfg.Redis)
	if rerr != nil {
		panic(rerr)
	}

	//Connect to the SMTP server
	_, eerr := email.GetInstance().Connect(&cfg.Email)
	if eerr != nil {
		panic(eerr)
	}

	//Test listing
	names, _ := mclient.ListDatabaseNames(context.TODO(), bson.M{})
	for i, name := range names {
		fmt.Printf("DB #%d: %s\n", i, name)
	}

	//Setup Chi and start listening for connections
	r := setupServer(&cfg, &env, mclient, rclient)
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
func setupServer(cfg *config.Config, env *config.Env, mclient *mongo.Client, rclient *redis.Client) chi.Router {
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
		AllowedOrigins:   []string{"https://*", "http://*"}, //Localhost and domain only
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"X-PINGOTHER", "Accept", "Authorization", "Content-Type", "X-CSRF-Token"}, // Ensure headers match client requests
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
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

	//Health route
	r.Get("/", router.Index)
	r.Get("/heartbeat", router.Heartbeat)

	r.Post("/send_message", router.SendMessage)

	//User routes
	r.Mount("/users", users.UserRoutes(cfg, env))

	//Challenge routes; protected by auth middleware (post_signup and users only)
	r.Group(func(r chi.Router) {
		authScopes := []obj.TokenScope{obj.TokenScopePOSTSIGNUP, obj.TokenScopeUSER}
		r.Use(mw.NewAuthMiddleware(authScopes, mclient, rclient))
		r.Mount("/challenges", challenges.ChallengeRoutes(mclient, rclient, env))
	})

	//AUTH TEST START
	scopes := []obj.TokenScope{obj.TokenScopeUSER}
	r.Group(func(r chi.Router) {
		r.Use(mw.NewAuthMiddleware(scopes, mclient, rclient))
		r.Get("/auth_test", router.AuthTest)
	})
	//AUTH TEST END

	//Return the built router for requests
	return r
}

/*
func setupScheduledTasks() {
	sch := gocron.NewScheduler(time.UTC)
	sch.Every(10).Seconds().Do(func() { fmt.Println("Sponsored by https://www.fittea.com") })
	sch.StartAsync()
}
*/

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
