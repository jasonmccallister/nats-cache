package keygen

import (
	"testing"

	"github.com/jasonmccallister/nats-cache/internal/auth"
)

func TestFromToken(t *testing.T) {
	type args struct {
		t   auth.Token
		db  uint32
		key string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		want1   string
		wantErr bool
	}{
		{
			name: "should return the correct key",
			args: args{
				t: auth.Token{
					Subject: "test",
				},
				db:  1,
				key: "test",
			},
			want:    "test.1-test",
			want1:   "test",
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, got1, err := FromToken(tt.args.t, tt.args.db, tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("FromToken() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("FromToken() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.want1 {
				t.Errorf("FromToken() got1 = %v, want %v", got1, tt.want1)
			}
		})
	}
}
