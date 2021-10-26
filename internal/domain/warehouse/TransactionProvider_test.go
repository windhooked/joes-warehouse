package warehouse

import (
	"testing"

	"github.com/averageflow/joes-warehouse/internal/infrastructure"
)

func TestCreateTransaction(t *testing.T) {
	t.Parallel()

	type args struct {
		db infrastructure.ApplicationDatabase
	}

	tests := []struct {
		name    string
		args    args
		want    int64
		wantErr bool
	}{
		{
			name:    "test creation returns no errors",
			args:    args{db: &infrastructure.MockApplicationDatabase{}},
			want:    0,
			wantErr: false,
		},
	}

	for _, tt := range tests {
		mockDB := &infrastructure.MockApplicationDatabase{}
		tt.args.db = mockDB

		t.Run(tt.name, func(t *testing.T) {
			got, err := CreateTransaction(tt.args.db)
			if (err != nil) != tt.wantErr {
				t.Errorf("CreateTransaction() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("CreateTransaction() = %v, want %v", got, tt.want)
			}
		})
	}
}

func TestCreateTransactionProductRelation(t *testing.T) {
	t.Parallel()

	type args struct {
		db            infrastructure.ApplicationDatabase
		transactionID int64
		productData   map[int64]int64
	}

	tests := []struct {
		name    string
		args    args
		wantErr bool
	}{
		{
			name: "test creation returns no errors",
			args: args{
				db:            &infrastructure.MockApplicationDatabase{},
				transactionID: 1,
				productData: map[int64]int64{
					1: 1, 2: 2, 3: 4,
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if err := CreateTransactionProductRelation(tt.args.db, tt.args.transactionID, tt.args.productData); (err != nil) != tt.wantErr {
				t.Errorf("CreateTransactionProductRelation() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
