package auth

import (
	"reflect"
	"testing"
)

func Test_pasetoPublicKey_Authorize(t *testing.T) {
	type fields struct {
		PublicKey string
	}
	type args struct {
		token string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    *Token
		wantErr bool
	}{
		{
			name: "should return the correct token",
			fields: fields{
				PublicKey: "58d50a037ae847a27006146bb95cf96cf520cac1409ee56bb8db297b6b4ceb6f",
			},
			args: args{
				token: "Bearer v4.public.eyJzdWIiOiJ0ZXN0LXN1YmplY3QifRFVZ6vjT01viOnDR0n670LzD6gRmsbvTHYyK05fc9vzwT8cc-9mk26R1C2gT1xrfZl2qBAe6aQCBvfdgH2X7wU",
			},
			want: &Token{
				Subject: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			p := &pasetoPublicKey{
				PublicKey: tt.fields.PublicKey,
			}
			got, err := p.Authorize(tt.args.token)
			if (err != nil) != tt.wantErr {
				t.Errorf("pasetoPublicKey.Authorize() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("pasetoPublicKey.Authorize() = %v, want %v", got, tt.want)
			}
		})
	}
}
