package server

import (
	"fmt"
	"modules/internal/entities"
	"net/http"
	"strconv"
)

type DefaultArgs struct {
	Message string
}
type CreateArgs struct {
	Goods []entities.Good
	Stock []entities.Stock
}

type GoodArgs struct {
	Code    string
	Name    string
	Size    float64
	Value   int64
	StockId string
	Dynamic bool
}

type StockArgs struct {
	Id        string
	Nmae      string
	Available bool
}

type Response struct {
	Message string
}

func (s *Service) CreateNewGood(r *http.Request, args *CreateArgs, response *Response) error {
	newGood := make([]entities.Good, 0, len(args.Goods))

	for _, gd := range args.Goods {
		if err := s.db.CreateNewGood(s.ctx, gd); err != nil {
			return fmt.Errorf("error when creating a new good entity in db: %s", err)
		}
		newGood = append(newGood, gd)
	}

	var resp string = "the following entities have been created: "
	for _, v := range newGood {
		resp = resp + v.Name + " "
	}

	response.Message = resp
	return nil
}

func (s *Service) CreateNewStock(r *http.Request, args *CreateArgs, response *Response) error {
	newStock := make([]entities.Stock, 0, len(args.Stock))

	for _, st := range args.Stock {
		if err := s.db.CreateNewStock(s.ctx, st); err != nil {
			return fmt.Errorf("error when creating a new stock entity in db: %s", err)
		}
		newStock = append(newStock, st)
	}

	var resp string = "the following entities have been created: "
	for _, v := range newStock {
		resp = resp + v.Name + " "
	}

	response.Message = resp
	return nil
}

func (s *Service) AddGood(r *http.Request, args *GoodArgs, response *Response) error {
	if err := s.db.AddGood(s.ctx, args.Code, args.StockId, args.Value, args.Dynamic); err != nil {
		return fmt.Errorf("error when adding: %s", err)
	}

	response.Message = "done"
	return nil
}

func (s *Service) ReservationGood(r *http.Request, args *GoodArgs, response *Response) error {
	if err := s.db.ReservationGood(s.ctx, args.Code, args.StockId, args.Value); err != nil {
		return fmt.Errorf("error when reservation: %s", err)
	}

	response.Message = "done"
	return nil
}

func (s *Service) GetAllGoods(r *http.Request, args *DefaultArgs, response *Response) error {
	goods, err := s.db.GetAllGoods(s.ctx)
	if err != nil {
		return fmt.Errorf("error when getting all goods: %s", err)
	}

	var resp string = "all goods: "
	for _, v := range goods {
		resp = resp + v.Name
	}

	response.Message = resp

	return nil
}

func (s *Service) GetGoodByCode(r *http.Request, args *GoodArgs, response *Response) error {
	g, err := s.db.GetGoodByCode(s.ctx, args.Code)
	if err != nil {
		return fmt.Errorf("error when getting good by id: %s", err)
	}

	response.Message = ("code: " + string(g.Code) + " name: " + g.Name + " size: " + strconv.FormatFloat(g.Size, 'f', 2, 64) + " value: " + string(g.Value))

	return nil
}

func (s *Service) GetGoodsCountByStockId(r *http.Request, args *GoodArgs, response *Response) error {
	count, err := s.db.GetGoodsCountByStockId(s.ctx, args.StockId, args.Code)
	if err != nil {
		return fmt.Errorf("error when getting goods count by stock id: %s", err)
	}

	response.Message = string(count)

	return nil
}
