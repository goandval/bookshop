package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"github.com/go-chi/chi/v5"

	"github.com/yourorg/bookshop/internal/domain"
	"github.com/yourorg/bookshop/internal/service"
	"golang.org/x/exp/slog"
)

type Handler struct {
	Book     service.BookService
	Category service.CategoryService
	Cart     service.CartService
	Order    service.OrderService
	Logger   *slog.Logger
}

func NewHandler(book service.BookService, category service.CategoryService, cart service.CartService, order service.OrderService, logger *slog.Logger) *Handler {
	return &Handler{
		Book:     book,
		Category: category,
		Cart:     cart,
		Order:    order,
		Logger:   logger,
	}
}

// ListBooks godoc
// @Summary      Get list of books
// @Description  Returns a list of books, optionally filtered by category
// @Tags         books
// @Produce      json
// @Param        category_id  query     int  false  "Category ID"
// @Success      200  {array}  domain.Book
// @Failure      500  {object}  map[string]string
// @Router       /books [get]
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
		h.Logger.Error("failed to list books", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(books)
}

// GetBook godoc
// @Summary      Get book by ID
// @Description  Returns a book by its ID
// @Tags         books
// @Produce      json
// @Param        id   path      int  true  "Book ID"
// @Success      200  {object}  domain.Book
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Router       /books/{id} [get]
func (h *Handler) GetBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("invalid book id", "id", idStr, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book, err := h.Book.GetByID(r.Context(), id)
	if err != nil {
		h.Logger.Error("book not found", "id", id, "err", err)
		w.WriteHeader(http.StatusNotFound)
		return
	}
	json.NewEncoder(w).Encode(book)
}

// CreateBook godoc
// @Summary      Create a new book
// @Description  Creates a new book (admin only)
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        book  body      domain.Book  true  "Book to create"
// @Success      201  {object}  domain.Book
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /books [post]
func (h *Handler) CreateBook(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Title      string  `json:"title"`
		Author     string  `json:"author"`
		Year       int     `json:"year"`
		Price      float64 `json:"price"`
		CategoryID int     `json:"category_id"`
		Inventory  int     `json:"stock"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.Author == "" || req.CategoryID == 0 {
		h.Logger.Error("invalid book create request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book := &domain.Book{
		Title:      req.Title,
		Author:     req.Author,
		Year:       req.Year,
		Price:      req.Price,
		CategoryID: req.CategoryID,
		Inventory:  req.Inventory,
	}
	if err := h.Book.Create(r.Context(), book); err != nil {
		h.Logger.Error("failed to create book", "err", err)
		errStr := err.Error()
		if strings.Contains(errStr, "inventory must be >= 0") ||
			strings.Contains(errStr, "category required") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(book)
}

// UpdateBook godoc
// @Summary      Update a book
// @Description  Updates an existing book (admin only)
// @Tags         books
// @Accept       json
// @Produce      json
// @Param        id    path      int         true  "Book ID"
// @Param        book  body      domain.Book true  "Book to update"
// @Success      200  {object}  domain.Book
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /books/{id} [put]
func (h *Handler) UpdateBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("invalid book id", "id", idStr, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var req struct {
		Title      string  `json:"title"`
		Author     string  `json:"author"`
		Year       int     `json:"year"`
		Price      float64 `json:"price"`
		CategoryID int     `json:"category_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.Title == "" || req.Author == "" || req.CategoryID == 0 {
		h.Logger.Error("invalid book update request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	book := &domain.Book{
		ID:         id,
		Title:      req.Title,
		Author:     req.Author,
		Year:       req.Year,
		Price:      req.Price,
		CategoryID: req.CategoryID,
	}
	if err := h.Book.Update(r.Context(), book); err != nil {
		h.Logger.Error("failed to update book", "id", id, "err", err)
		if strings.Contains(err.Error(), "book not found") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(book)
}

// DeleteBook godoc
// @Summary      Delete a book
// @Description  Deletes a book by ID (admin only)
// @Tags         books
// @Param        id   path      int  true  "Book ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /books/{id} [delete]
func (h *Handler) DeleteBook(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("invalid book id for delete", "id", idStr, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Book.Delete(r.Context(), id); err != nil {
		h.Logger.Error("failed to delete book", "id", id, "err", err)
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ListCategories godoc
// @Summary      Get list of categories
// @Description  Returns a list of categories
// @Tags         categories
// @Produce      json
// @Success      200  {array}  domain.Category
// @Failure      500  {object}  map[string]string
// @Router       /categories [get]
func (h *Handler) ListCategories(w http.ResponseWriter, r *http.Request) {
	categories, err := h.Category.List(r.Context())
	if err != nil {
		h.Logger.Error("failed to list categories", "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(categories)
}

// CreateCategory godoc
// @Summary      Create a new category
// @Description  Creates a new category (admin only)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        category  body      domain.Category  true  "Category to create"
// @Success      201  {object}  domain.Category
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /categories [post]
func (h *Handler) CreateCategory(w http.ResponseWriter, r *http.Request) {
	var cat struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil || cat.Name == "" {
		h.Logger.Error("invalid category create request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := &domain.Category{Name: cat.Name}
	if err := h.Category.Create(r.Context(), c); err != nil {
		h.Logger.Error("failed to create category", "err", err)
		if strings.Contains(err.Error(), "name required") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// UpdateCategory godoc
// @Summary      Update a category
// @Description  Updates an existing category (admin only)
// @Tags         categories
// @Accept       json
// @Produce      json
// @Param        id        path      int             true  "Category ID"
// @Param        category  body      domain.Category true  "Category to update"
// @Success      200  {object}  domain.Category
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /categories/{id} [put]
func (h *Handler) UpdateCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("invalid category id", "id", idStr, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	var cat struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&cat); err != nil || cat.Name == "" {
		h.Logger.Error("invalid category update request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	c := &domain.Category{ID: id, Name: cat.Name}
	if err := h.Category.Update(r.Context(), c); err != nil {
		h.Logger.Error("failed to update category", "id", id, "err", err)
		if strings.Contains(err.Error(), "name required") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

// DeleteCategory godoc
// @Summary      Delete a category
// @Description  Deletes a category by ID (admin only)
// @Tags         categories
// @Param        id   path      int  true  "Category ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      404  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /categories/{id} [delete]
func (h *Handler) DeleteCategory(w http.ResponseWriter, r *http.Request) {
	idStr := chi.URLParam(r, "id")
	id, err := strconv.Atoi(idStr)
	if err != nil {
		h.Logger.Error("invalid category id for delete", "id", idStr, "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Category.Delete(r.Context(), id); err != nil {
		h.Logger.Error("failed to delete category", "id", id, "err", err)
		if strings.Contains(err.Error(), "not found") {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// GetCart godoc
// @Summary      Get user's cart
// @Description  Returns the cart for the authenticated user
// @Tags         cart
// @Produce      json
// @Success      200  {object}  domain.Cart
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /cart [get]
func (h *Handler) GetCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	cart, err := h.Cart.GetByUserID(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to get cart", "userID", userID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(cart)
}

// AddToCart godoc
// @Summary      Add book to cart
// @Description  Adds a book to the authenticated user's cart
// @Tags         cart
// @Accept       json
// @Produce      json
// @Param        item  body      map[string]int  true  "Book ID to add"
// @Success      201  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /cart [post]
func (h *Handler) AddToCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	var req struct {
		BookID int `json:"book_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil || req.BookID == 0 {
		h.Logger.Error("invalid add to cart request", "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Cart.AddItem(r.Context(), userID, req.BookID); err != nil {
		h.Logger.Error("failed to add item to cart", "userID", userID, "bookID", req.BookID, "err", err)
		errStr := err.Error()
		if strings.Contains(errStr, "book not found") ||
			strings.Contains(errStr, "out of stock") ||
			strings.Contains(errStr, "not enough books in stock") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
}

// RemoveFromCart godoc
// @Summary      Remove book from cart
// @Description  Removes a book from the authenticated user's cart
// @Tags         cart
// @Param        book_id   path      int  true  "Book ID"
// @Success      204  {object}  nil
// @Failure      400  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /cart/{book_id} [delete]
func (h *Handler) RemoveFromCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	bookID, err := strconv.Atoi(chi.URLParam(r, "book_id"))
	if err != nil {
		h.Logger.Error("invalid book id for remove from cart", "id", chi.URLParam(r, "book_id"), "err", err)
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	if err := h.Cart.RemoveItem(r.Context(), userID, bookID); err != nil {
		h.Logger.Error("failed to remove item from cart", "userID", userID, "bookID", bookID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// ClearCart godoc
// @Summary      Clear cart
// @Description  Clears the authenticated user's cart
// @Tags         cart
// @Success      204  {object}  nil
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /cart [delete]
func (h *Handler) ClearCart(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	if err := h.Cart.Clear(r.Context(), userID); err != nil {
		h.Logger.Error("failed to clear cart", "userID", userID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusNoContent)
}

// PlaceOrder godoc
// @Summary      Place an order
// @Description  Places an order for the authenticated user
// @Tags         orders
// @Produce      json
// @Success      201  {object}  domain.Order
// @Failure      400  {object}  map[string]string
// @Failure      409  {object}  map[string]string
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /orders [post]
func (h *Handler) PlaceOrder(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	order, err := h.Order.Create(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to place order", "userID", userID, "err", err)
		errStr := err.Error()
		if strings.Contains(errStr, "cart is empty") {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if strings.Contains(errStr, "book not found") || strings.Contains(errStr, "book out of stock") {
			w.WriteHeader(http.StatusConflict)
			return
		}
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(order)
}

// ListOrders godoc
// @Summary      List user's orders
// @Description  Returns a list of orders for the authenticated user
// @Tags         orders
// @Produce      json
// @Success      200  {array}  domain.Order
// @Failure      500  {object}  map[string]string
// @Security     ApiKeyAuth
// @Router       /orders [get]
func (h *Handler) ListOrders(w http.ResponseWriter, r *http.Request) {
	userID := r.Context().Value("userID").(string)
	orders, err := h.Order.ListByUser(r.Context(), userID)
	if err != nil {
		h.Logger.Error("failed to list orders", "userID", userID, "err", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	json.NewEncoder(w).Encode(orders)
}
