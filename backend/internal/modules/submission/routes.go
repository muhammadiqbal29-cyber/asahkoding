package submission

import (
	"github.com/go-chi/chi/v5"
	"github.com/muhammad-iqbal/leetcode-backend/internal/middleware"
)

// Routes mendaftarkan jalur-jalur untuk pengiriman kode
func Routes() *chi.Mux {
	r := chi.NewRouter()

	// Wajib melampirkan Token JWT untuk bisa mengirim jawaban
	r.Use(middleware.AuthMiddleware)
	r.Post("/", SubmitCode)

	return r
}
