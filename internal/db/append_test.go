package db

import (
	"context"
	"modules/pkg/dbclient"
	"modules/pkg/logger"
	"testing"
)

func Test_db_AddGood(t *testing.T) {
	type fields struct {
		client dbclient.Client
		logger *logger.Logger
	}
	type args struct {
		ctx     context.Context
		code    string
		stockId string
		value   int64
		dynamic bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "positive test",
			args: args{
				ctx:     context.TODO(),
				code:    "1",
				stockId: "08452345-b231-49a1-ae71-81cf8cb857d9",
				value:   10,
				dynamic: true,
			},
			wantErr: false,
		},
		{
			name: "error test",
			args: args{
				ctx:     context.TODO(),
				code:    "1",
				stockId: "08452345-b231-49a1-ae71-81cf8cb857d9",
				value:   10000,
				dynamic: false,
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
			if err := d.AddGood(tt.args.ctx, tt.args.code, tt.args.stockId, tt.args.value, tt.args.dynamic); (err != nil) != tt.wantErr {
				t.Errorf("AddGood() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
