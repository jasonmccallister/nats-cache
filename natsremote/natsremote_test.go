package natsremote

import (
	"github.com/nats-io/nats-server/v2/server"
	"net/url"
	"os"
	"reflect"
	"testing"
)

func TestRemoteLeafFromEnv(t *testing.T) {
	tests := []struct {
		name      string
		envVars   map[string]string
		want      *server.RemoteLeafOpts
		wantCreds string
		wantErr   bool
	}{
		{
			name:      "no env vars require nats seed and key",
			envVars:   make(map[string]string),
			want:      nil,
			wantCreds: "",
			wantErr:   true,
		},
		{
			name: "nats seed and key set",
			envVars: map[string]string{
				"NATS_SEED": "seed",
				"NATS_JWT":  "jwt",
			},
			want: &server.RemoteLeafOpts{
				URLs: []*url.URL{&url.URL{Scheme: "tls", Host: "connect.ngs.global"}},
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envVars {
				if err := os.Setenv(k, v); err != nil {
					t.Errorf("failed to set env var %s: %v", k, err)
				}
			}

			got, got1, err := RemoteLeafFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteLeafFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteLeafFromEnv() got = %v, want %v", got, tt.want)
			}
			if got1 != tt.wantCreds {
				t.Errorf("RemoteLeafFromEnv() got1 = %v, want %v", got1, tt.wantCreds)
			}
		})
	}
}
