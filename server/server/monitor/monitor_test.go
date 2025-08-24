package monitor

import (
	"context"
	"logger"
	"reflect"
	"testing"
)

func Test_getAddressFills(t *testing.T) {
	type args struct {
		ctx     context.Context
		address string
	}
	tests := []struct {
		name    string
		args    args
		want    []UserFillResponse
		wantErr bool
	}{
		{
			name: "Test",
			args: args{
				ctx:     context.Background(),
				address: "0x1755e9a4e305f8528f0b0705fd4d0e0c860b5fB8",
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := getAddressFills(tt.args.ctx, tt.args.address)
			if (err != nil) != tt.wantErr {
				t.Errorf("getAddressFills() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			logger.Info("Getting fills", "fills", got)
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("getAddressFills() = %v, want %v", got, tt.want)
			}
		})
	}
}
