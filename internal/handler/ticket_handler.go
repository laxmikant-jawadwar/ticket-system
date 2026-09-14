package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"strconv"
	"strings"

	"ticket-system/internal/middleware"
	"ticket-system/internal/model"
	"ticket-system/internal/service"
)

type TicketHandler struct {
	ticketService service.TicketService
}

func NewTicketHandler(ticketService service.TicketService) *TicketHandler {
	return &TicketHandler{
		ticketService: ticketService,
	}
}

// POST /tickets
func (h *TicketHandler) CreateTicket(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	var req model.CreateTicketRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	ticket, err := h.ticketService.CreateTicket(
		r.Context(),
		userID,
		req,
	)

	if err != nil {
		if errors.Is(err, service.ErrInvalidInput) {
			http.Error(w, "Invalid input", http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	json.NewEncoder(w).Encode(ticket)
}

// GET /tickets
func (h *TicketHandler) GetTickets(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	tickets, err := h.ticketService.GetTickets(
		r.Context(),
		userID,
	)

	if err != nil {
		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(tickets)
}

// GET /tickets/{id}
func (h *TicketHandler) GetTicketByID(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodGet {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ticketID, err := getTicketID(r.URL.Path)

	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	ticket, err := h.ticketService.GetTicketByID(
		r.Context(),
		userID,
		ticketID,
	)

	if err != nil {

		if errors.Is(err, service.ErrTicketNotFound) {
			http.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, service.ErrUnauthorized) {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	json.NewEncoder(w).Encode(ticket)
}

// PATCH /tickets/{id}/status
func (h *TicketHandler) UpdateTicketStatus(w http.ResponseWriter, r *http.Request) {

	if r.Method != http.MethodPatch {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	userID, ok := getUserIDFromContext(r)
	if !ok {
		http.Error(w, "Unauthorized", http.StatusUnauthorized)
		return
	}

	ticketID, err := getTicketID(r.URL.Path)

	if err != nil {
		http.Error(w, "Invalid ticket ID", http.StatusBadRequest)
		return
	}

	var req model.UpdateTicketStatusRequest

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	err = h.ticketService.UpdateTicketStatus(
		r.Context(),
		userID,
		ticketID,
		req.Status,
	)

	if err != nil {

		if errors.Is(err, service.ErrTicketNotFound) {
			http.Error(w, "Ticket not found", http.StatusNotFound)
			return
		}

		if errors.Is(err, service.ErrUnauthorized) {
			http.Error(w, "Unauthorized", http.StatusForbidden)
			return
		}

		if errors.Is(err, service.ErrInvalidStatus) {
			http.Error(w, "Invalid status", http.StatusBadRequest)
			return
		}

		if errors.Is(err, service.ErrInvalidStatusFlow) {
			http.Error(w, "Invalid status transition", http.StatusBadRequest)
			return
		}

		http.Error(w, "Internal server error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")

	json.NewEncoder(w).Encode(map[string]string{
		"message": "Ticket status updated successfully",
	})
}

// Extract authenticated user ID from request context.
func getUserIDFromContext(r *http.Request) (int64, bool) {

	value := r.Context().Value(middleware.UserIDKey)

	userID, ok := value.(int64)

	return userID, ok
}

// Extract ticket ID from URL.
func getTicketID(path string) (int64, error) {

	parts := strings.Split(strings.Trim(path, "/"), "/")

	if len(parts) < 2 || parts[0] != "tickets" {
		return 0, errors.New("invalid ticket path")
	}

	return strconv.ParseInt(parts[1], 10, 64)
}

func (h *TicketHandler) Tickets(w http.ResponseWriter, r *http.Request) {

	if r.URL.Path == "/tickets" {

		switch r.Method {
		case http.MethodPost:
			h.CreateTicket(w, r)
		case http.MethodGet:
			h.GetTickets(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	if strings.HasPrefix(r.URL.Path, "/tickets/") {

		switch r.Method {
		case http.MethodGet:
			h.GetTicketByID(w, r)

		case http.MethodPatch:
			h.UpdateTicketStatus(w, r)

		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

		return
	}

	http.NotFound(w, r)
}
