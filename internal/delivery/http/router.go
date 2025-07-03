package http

import (
	"net/http"

	"github.com/go-chi/chi/v5"
)

func (h *Handler) Router(auth *AuthMiddleware) http.Handler {
	r := chi.NewRouter()

	// --- Публичные ---
	r.Get("/books", h.ListBooks)
	r.Get("/books/{id}", h.GetBook)
	r.Get("/categories", h.ListCategories)

	// --- Только для админов ---
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuth)
		r.Use(auth.RequireRole("admin"))
		r.Post("/categories", h.CreateCategory)
		r.Put("/categories/{id}", h.UpdateCategory)
		r.Delete("/categories/{id}", h.DeleteCategory)
		r.Post("/books", h.CreateBook)
		r.Put("/books/{id}", h.UpdateBook)
		r.Delete("/books/{id}", h.DeleteBook)
	})

	// --- Для аутентифицированных пользователей ---
	r.Group(func(r chi.Router) {
		r.Use(auth.JWTAuth)
		r.Get("/cart", h.GetCart)
		r.Post("/cart", h.AddToCart)
		r.Delete("/cart/{book_id}", h.RemoveFromCart)
		r.Delete("/cart", h.ClearCart)
		r.Post("/orders", h.PlaceOrder)
		r.Get("/orders", h.ListOrders)
	})

	return r
}
