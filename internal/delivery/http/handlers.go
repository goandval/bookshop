package http

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/go-chi/chi/v5"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/service"
)

type Handler struct {
	Book     service.BookService
	Category service.CategoryService
	Cart     service.CartService
	Order    service.OrderService
}

func NewHandler(book service.BookService, category service.CategoryService, cart service.CartService, order service.OrderService) *Handler {
	return &Handler{
		Book:     book,
		Category: category,
		Cart:     cart,
		Order:    order,
	}
}

// --- Book ---
func (h *Handler) ListBooks(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	var categoryIDs []int
	for _, v := range q["category_id"] {
		id, err := strconv.Atoi(v)
		if err == nil {
			categoryIDs = append(categoryIDs, id)
		}
	}
	books, err := h.Book.List(r.Context(), categoryIDs, 100, 0)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book, err := h.Book.GetByID(r.Context(), id)
	if err != nil {
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	// TODO: реализовать
}

// --- Category ---
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Category.List(r.Context())
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil || cat.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := &domain.Category{Name: cat.Name}
	if err := h.Category.Create(r.Context(), c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var cat struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil || cat.Name == "" {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := &domain.Category{ID: id, Name: cat.Name}
	if err := h.Category.Update(r.Context(), c); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Category.Delete(r.Context(), id); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Cart ---
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	cart, err := h.Cart.GetByUserID(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	var req struct {
		BookID int `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.BookID == 0 {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Cart.AddItem(r.Context(), userID, req.BookID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	bookID, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Cart.RemoveItem(r.Context(), userID, bookID); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	if err := h.Cart.Clear(r.Context(), userID); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// --- Order ---
func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	order, err := h.Order.Create(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	orders, err := h.Order.ListByUser(r.Context(), userID)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
