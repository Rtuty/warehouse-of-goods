package db

import (
	"context"
	"modules/internal/entities"
	"modules/pkg/dbclient"
	"modules/pkg/logger"
	"testing"
)

func Test_db_CreateNewStock(t *testing.T) {
	type fields struct {
		client dbclient.Client
		logger *logger.Logger
	}
	type args struct {
		ctx context.Context
		s   entities.Stock
	}

	dbclient.GetConnection()
	log := logger.GetLogger()
	client, err := dbclient.NewClient(context.TODO(), 5, dbclient.PstgCon, log)
	if err != nil {
		t.Errorf("New client error (test_db_GetAllGoods): %s", err)
	}

	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "base test",
			fields: fields{
				client: client,
				logger: log,
			},
			args: args{
				ctx: context.TODO(),
				s: entities.Stock{
					Name:      "test stock",
					Available: false,
				},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			d := &db{
				client: tt.fields.client,
				logger: tt.fields.logger,
			}

			if err := d.CreateNewStock(tt.args.ctx, tt.args.s); (err != nil) != tt.wantErr {
				t.Errorf("CreateNewStock() error = %v, wantErr %v", err, tt.wantErr)
			}

			if _, err := tt.fields.client.Exec(tt.args.ctx, "delete from stocks where name = 'test stock' and available = false"); err != nil {
				t.Errorf("deleting test stock error: %s", err)
			}
		})
	}
}
