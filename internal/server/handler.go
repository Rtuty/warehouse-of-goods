package server

import (
	"context"
	"modules/internal/db"
	"modules/internal/entities"
	"modules/pkg/logger"
)

type Handler struct {
	ctx context.Context
	db  db.Storage
	log *logger.Logger
}

func (h *Handler) CreateNewGood(args []entities.Good, resp *[]entities.Good) error {
	newGood := make([]entities.Good, 0, len(args))

loop:
	for _, gd := range args {
		if err := h.db.CreateNewGood(h.ctx, gd); err != nil {
			h.log.Errorf("error when creating a new entity product in db: %s", err)
			continue loop
		}
		newGood = append(newGood, gd)
	}
	*resp = newGood

	return nil
}

func (h *Handler) CreateNewStock(args []entities.Stock, resp *[]entities.Stock) error {
	newStock := make([]entities.Stock, 0, len(args))

loop:
	for _, st := range args {
		if err := h.db.CreateNewStock(h.ctx, st); err != nil {
			h.log.Errorf("error when creating a new entity stock in db: %s", err)
			continue loop
		}
		newStock = append(newStock, st)
	}

	*resp = newStock

	return nil
}
