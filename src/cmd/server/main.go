package main

import (
	"log"
	"net/http"
	"os"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/chi/v5/middleware"

	"github.com/gandli/yctf/controllers"
	yctfMiddleware "github.com/gandli/yctf/middleware"
	"github.com/gandli/yctf/ws"
)

func main() {
	r := chi.NewRouter()

	r.Use(middleware.Logger)
	r.Use(middleware.Recoverer)
	r.Use(yctfMiddleware.CORSMiddleware(getCORSOrigins()...))

	r.Get("/health", healthHandler)

	r.Route("/api/v1", func(r chi.Router) {
		r.Group(func(r chi.Router) {
			r.Post("/auth/register", controllers.RegisterHandler)
			r.Post("/auth/login", controllers.LoginHandler)
			r.Get("/challenges", controllers.ListChallengesHandler)
			r.Get("/scoreboard", controllers.ScoreboardHandler)
			r.Get("/scoreboard/timeline", controllers.ScoreboardTimelineHandler)
		})

		r.Group(func(r chi.Router) {
			r.Get("/users/me", controllers.GetCurrentUserHandler)
			r.Post("/submit", controllers.SubmitFlagHandler)
			r.Get("/challenges/{id}", controllers.GetChallengeHandler)
			r.Get("/ws", ws.WebSocketHandler)
		})

		r.Group(func(r chi.Router) {
			r.Post("/challenges", controllers.CreateChallengeHandler)
			r.Put("/challenges/{id}", controllers.UpdateChallengeHandler)
			r.Delete("/challenges/{id}", controllers.DeleteChallengeHandler)
		})

		r.Group(func(r chi.Router) {
			r.Get("/admin/users", controllers.AdminListUsersHandler)
			r.Get("/admin/stats", controllers.AdminStatsHandler)
		})
	})

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("YCTF server starting on :%s", port)
	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatal(err)
	}
}

func getCORSOrigins() []string {
	origins := os.Getenv("CORS_ORIGINS")
	if origins == "" {
		return []string{"http://localhost:5173"}
	}
	return []string{origins}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}
