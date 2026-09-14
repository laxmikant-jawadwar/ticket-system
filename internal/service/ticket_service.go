package service

import (
	"context"
	"database/sql"
	"errors"
	"strings"

	"ticket-system/internal/model"
	"ticket-system/internal/repository"
)

var (
	ErrTicketNotFound    = errors.New("ticket not found")
	ErrUnauthorized      = errors.New("unauthorized")
	ErrInvalidStatus     = errors.New("invalid status")
	ErrInvalidStatusFlow = errors.New("invalid status transition")
)

type TicketService interface {
	CreateTicket(ctx context.Context, userID int64, req model.CreateTicketRequest) (*model.Ticket, error)
	GetTickets(ctx context.Context, userID int64) ([]model.Ticket, error)
	GetTicketByID(ctx context.Context, userID int64, ticketID int64) (*model.Ticket, error)
	UpdateTicketStatus(ctx context.Context, userID int64, ticketID int64, status string) error
}

type ticketService struct {
	ticketRepository repository.TicketRepository
}

func NewTicketService(ticketRepository repository.TicketRepository) TicketService {
	return &ticketService{
		ticketRepository: ticketRepository,
	}
}

func (s *ticketService) CreateTicket(
	ctx context.Context,
	userID int64,
	req model.CreateTicketRequest,
) (*model.Ticket, error) {

	title := strings.TrimSpace(req.Title)
	description := strings.TrimSpace(req.Description)

	if title == "" || description == "" {
		return nil, ErrInvalidInput
	}

	ticket := &model.Ticket{
		UserID:      userID,
		Title:       title,
		Description: description,
		Status:      "open",
	}

	if err := s.ticketRepository.CreateTicket(ctx, ticket); err != nil {
		return nil, err
	}

	return ticket, nil
}

func (s *ticketService) GetTickets(
	ctx context.Context,
	userID int64,
) ([]model.Ticket, error) {

	return s.ticketRepository.GetTicketsByUserID(ctx, userID)
}

func (s *ticketService) GetTicketByID(
	ctx context.Context,
	userID int64,
	ticketID int64,
) (*model.Ticket, error) {

	ticket, err := s.ticketRepository.GetTicketByID(ctx, ticketID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, ErrTicketNotFound
		}

		return nil, err
	}

	// Ownership check
	if ticket.UserID != userID {
		return nil, ErrUnauthorized
	}

	return ticket, nil
}

func (s *ticketService) UpdateTicketStatus(
	ctx context.Context,
	userID int64,
	ticketID int64,
	status string,
) error {

	// Only these three statuses are allowed.
	if status != "open" &&
		status != "in_progress" &&
		status != "closed" {
		return ErrInvalidStatus
	}

	ticket, err := s.ticketRepository.GetTicketByID(ctx, ticketID)

	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return ErrTicketNotFound
		}

		return err
	}

	// Ownership check
	if ticket.UserID != userID {
		return ErrUnauthorized
	}

	// Closed tickets cannot be reopened or changed.
	if ticket.Status == "closed" {
		return ErrInvalidStatusFlow
	}

	// Status flow:
	// open -> in_progress
	// in_progress -> closed

	if ticket.Status == "open" && status != "in_progress" {
		return ErrInvalidStatusFlow
	}

	if ticket.Status == "in_progress" && status != "closed" {
		return ErrInvalidStatusFlow
	}

	return s.ticketRepository.UpdateTicketStatus(
		ctx,
		ticketID,
		status,
	)
}
