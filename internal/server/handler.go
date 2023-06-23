package server

import (
	"context"
	"fmt"
	"modules/internal/db"
	"modules/internal/entities"
	"modules/pkg/logger"
	"net/http"
)

type Args struct {
	Goods []entities.Good
	Stock []entities.Stock
}

type Response struct {
	Message string
}

type Service struct {
	ctx context.Context
	db  db.Storage
	log *logger.Logger
}

func (s *Service) CreateNewGood(r *http.Request, args *Args, response *Response) error {
	newGood := make([]entities.Good, 0, len(args.Goods))

	for _, gd := range args.Goods {
		if err := s.db.CreateNewGood(s.ctx, gd); err != nil {
			return fmt.Errorf("error when creating a new good entity in db: %s", err)
		}
		newGood = append(newGood, gd) //add in res message
	}

	response.Message = "new good created"
	return nil
}

func (s *Service) CreateNewStock(r *http.Request, args *Args, response *Response) error {
	newStock := make([]entities.Stock, 0, len(args.Stock))

	for _, st := range args.Stock {
		if err := s.db.CreateNewStock(s.ctx, st); err != nil {
			return fmt.Errorf("error when creating a new stock entity in db: %s", err)
		}
		newStock = append(newStock, st) //add in res message
	}

	response.Message = "new stock created"
	return nil
}
