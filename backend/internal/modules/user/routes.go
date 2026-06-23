package user

import (
	"github.com/go-chi/chi/v5"
)

// Routes membungkus semua rute terkait autentikasi
func Routes() *chi.Mux {
	r := chi.NewRouter()

	r.Post("/register", RegisterHandler)
	r.Post("/login", LoginHandler)

	return r
}
