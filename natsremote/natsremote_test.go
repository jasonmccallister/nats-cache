package natsremote

import (
	"github.com/nats-io/nats-server/v2/server"
	"os"
	"reflect"
	"testing"
)

func TestRemoteLeafFromEnv(t *testing.T) {
	tests := []struct {
		name    string
		envs    map[string]string
		want    *server.RemoteLeafOpts
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			for k, v := range tt.envs {
				if err := os.Setenv(k, v); err != nil {
					t.Fatalf("failed to set env var %s: %v", k, err)
				}
			}

			got, err := RemoteLeafFromEnv()
			if (err != nil) != tt.wantErr {
				t.Errorf("RemoteLeafFromEnv() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("RemoteLeafFromEnv() got = %v, want %v", got, tt.want)
			}
		})
	}
}
