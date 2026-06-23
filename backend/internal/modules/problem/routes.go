package problem

import (
	"github.com/go-chi/chi/v5"
	"github.com/muhammad-iqbal/leetcode-backend/internal/middleware"
)

// Routes membungkus semua rute terkait soal (Problems)
func Routes() *chi.Mux {
	r := chi.NewRouter()

	// Rute Publik (Bisa diakses tanpa login)
	r.Get("/", GetProblems)
	r.Get("/{id}", GetProblemByID)

	// Rute Privat (Hanya Admin yang bisa membuat soal)
	r.Group(func(r chi.Router) {
		r.Use(middleware.AuthMiddleware)
		r.Use(middleware.AdminOnlyMiddleware)
		
		r.Post("/", CreateProblem)
	})

	return r
}
