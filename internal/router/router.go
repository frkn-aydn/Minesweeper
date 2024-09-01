package router

import (
	"mine-game/internal/database"
	handlers "mine-game/internal/handler"
	repository "mine-game/internal/repositories"
	"mine-game/internal/services"
	"net/http"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"
	"github.com/go-chi/cors"
)

func NewRouter(mongodb *database.MongoDB) *chi.Mux {

	// Repositories
	gameRepository := repository.NewGameRepository(mongodb)

	// Services
	gameService := services.NewGameService(gameRepository)

	// Handlers
	gameHandler := handlers.NewGameHandler(gameService)

	router := chi.NewRouter()

	// Middleware stack
	router.Use(middleware.Logger)
	router.Use(middleware.Recoverer)
	router.Use(middleware.RequestID)
	router.Use(middleware.RealIP)
	router.Use(middleware.NoCache)
	router.Use(middleware.Timeout(60 * time.Second))

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PATCH", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "User-Agent"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	router.Get("/health", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte("OK"))
	})

	router.Route("/v1/mine", func(r chi.Router) {
		r.Post("/", gameHandler.CreateGame)
		r.Get("/", gameHandler.MakeMove)
	})
	return router
}
