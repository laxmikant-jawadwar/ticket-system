package repository

import (
	"context"
	"database/sql"

	"ticket-system/internal/model"
)

type TicketRepository interface {
	CreateTicket(ctx context.Context, ticket *model.Ticket) error
	GetTicketsByUserID(ctx context.Context, userID int64) ([]model.Ticket, error)
	GetTicketByID(ctx context.Context, id int64) (*model.Ticket, error)
	UpdateTicketStatus(ctx context.Context, id int64, status string) error
}

type ticketRepository struct {
	db *sql.DB
}

func NewTicketRepository(db *sql.DB) TicketRepository {
	return &ticketRepository{
		db: db,
	}
}

func (r *ticketRepository) CreateTicket(
	ctx context.Context,
	ticket *model.Ticket,
) error {

	query := `
		INSERT INTO tickets (user_id, title, description)
		VALUES ($1, $2, $3)
		RETURNING id, status, created_at, updated_at
	`

	return r.db.QueryRowContext(
		ctx,
		query,
		ticket.UserID,
		ticket.Title,
		ticket.Description,
	).Scan(
		&ticket.ID,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)
}

func (r *ticketRepository) GetTicketsByUserID(
	ctx context.Context,
	userID int64,
) ([]model.Ticket, error) {

	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tickets
		WHERE user_id = $1
		ORDER BY created_at DESC
	`

	rows, err := r.db.QueryContext(ctx, query, userID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var tickets []model.Ticket

	for rows.Next() {
		var ticket model.Ticket

		err := rows.Scan(
			&ticket.ID,
			&ticket.UserID,
			&ticket.Title,
			&ticket.Description,
			&ticket.Status,
			&ticket.CreatedAt,
			&ticket.UpdatedAt,
		)

		if err != nil {
			return nil, err
		}

		tickets = append(tickets, ticket)
	}

	if err := rows.Err(); err != nil {
		return nil, err
	}

	return tickets, nil
}

func (r *ticketRepository) GetTicketByID(
	ctx context.Context,
	id int64,
) (*model.Ticket, error) {

	query := `
		SELECT id, user_id, title, description, status, created_at, updated_at
		FROM tickets
		WHERE id = $1
	`

	ticket := &model.Ticket{}

	err := r.db.QueryRowContext(ctx, query, id).Scan(
		&ticket.ID,
		&ticket.UserID,
		&ticket.Title,
		&ticket.Description,
		&ticket.Status,
		&ticket.CreatedAt,
		&ticket.UpdatedAt,
	)

	if err != nil {
		return nil, err
	}

	return ticket, nil
}

func (r *ticketRepository) UpdateTicketStatus(
	ctx context.Context,
	id int64,
	status string,
) error {

	query := `
		UPDATE tickets
		SET status = $1, updated_at = CURRENT_TIMESTAMP
		WHERE id = $2
	`

	_, err := r.db.ExecContext(ctx, query, status, id)

	return err
}
