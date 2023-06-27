package db

import (
	"context"
	"modules/internal/entities"
	"modules/pkg/dbclient"
	"modules/pkg/logger"
	"reflect"
	"testing"
)

func Test_db_GetAllGoods(t *testing.T) {
	type fields struct {
		client dbclient.Client
		logger *logger.Logger
	}
	type args struct {
		ctx context.Context
	}

	dbclient.GetConnection()
	log := logger.GetLogger()
	client, err := dbclient.NewClient(context.TODO(), 5, dbclient.PstgCon, log)
	if err != nil {
		t.Errorf("new client error (test_db_GetAllGoods): %s", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []entities.Good
		wantErr bool
	}{
		{
			name: "base test",
			fields: fields{
				client: client,
				logger: log,
			},
			args:    args{ctx: context.TODO()},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &db{
				client: tt.fields.client,
				logger: tt.fields.logger,
			}
			// TEST GetAllGoods func
			got, err := d.GetAllGoods(tt.args.ctx)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetAllGoods() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetAllGoods() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_db_GetGoodByCode(t *testing.T) {
	type fields struct {
		client dbclient.Client
		logger *logger.Logger
	}
	type args struct {
		ctx  context.Context
		code string
	}

	dbclient.GetConnection()
	log := logger.GetLogger()
	client, err := dbclient.NewClient(context.TODO(), 5, dbclient.PstgCon, log)
	if err != nil {
		t.Errorf("new client error (test_db_GetAllGoods): %s", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    entities.Good
		wantErr bool
	}{
		{
			name: "positive test",
			fields: fields{
				client: client,
				logger: log,
			},
			args: args{
				ctx:  context.TODO(),
				code: "1",
			},
			wantErr: false,
		},
		{
			name: "negative test",
			fields: fields{
				client: client,
				logger: log,
			},
			args: args{
				ctx:  context.TODO(),
				code: "9999+",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &db{
				client: tt.fields.client,
				logger: tt.fields.logger,
			}
			got, err := d.GetGoodByCode(tt.args.ctx, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGoodByCode() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("GetGoodByCode() got = %v, want %v", got, tt.want)
			}
		})
	}
}

func Test_db_GetGoodsCountByStockId(t *testing.T) {
	type fields struct {
		client dbclient.Client
		logger *logger.Logger
	}
	type args struct {
		ctx     context.Context
		stockId string
		code    string
	}

	dbclient.GetConnection()
	log := logger.GetLogger()
	client, err := dbclient.NewClient(context.TODO(), 5, dbclient.PstgCon, log)
	if err != nil {
		t.Errorf("new client error (test_db_GetAllGoods): %s", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		want    int64
		wantErr bool
	}{
		{
			name: "positive test",
			fields: fields{
				client: client,
				logger: log,
			},
			args: args{
				ctx:  context.TODO(),
				code: "1",
			},
			wantErr: false,
		},
		{
			name: "negative test",
			fields: fields{
				client: client,
				logger: log,
			},
			args: args{
				ctx:  context.TODO(),
				code: "9999+",
			},
			wantErr: true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &db{
				client: tt.fields.client,
				logger: tt.fields.logger,
			}
			got, err := d.GetGoodsCountByStockId(tt.args.ctx, tt.args.stockId, tt.args.code)
			if (err != nil) != tt.wantErr {
				t.Errorf("GetGoodsCountByStockId() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("GetGoodsCountByStockId() got = %v, want %v", got, tt.want)
			}
		})
	}
}
