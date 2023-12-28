package storage

import (
	"encoding/json"
	"reflect"
	"sync"
	"testing"
	"time"
)

func Test_inMemory_Delete(t *testing.T) {
	type fields struct {
		mu sync.RWMutex
		db map[string][]byte
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr bool
	}{
		{
			name: "should delete the key",
			fields: fields{
				mu: sync.RWMutex{},
				db: map[string][]byte{
					"test": marshalItem(t, Item{
						Value: []byte("test"),
						TTL:   0,
					}),
				},
			},
			args: args{
				key: "test",
			},
			wantErr: false,
		},
		{
			name: "should not error if the key does not exist",
			fields: fields{
				mu: sync.RWMutex{},
				db: map[string][]byte{},
			},
			args: args{
				key: "test",
			},
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &inMemory{
				mu: tt.fields.mu,
				db: tt.fields.db,
			}
			if err := s.Delete(tt.args.key); (err != nil) != tt.wantErr {
				t.Errorf("inMemory.Delete() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}

func marshalItem(t *testing.T, i Item) []byte {
	t.Helper()

	b, err := json.Marshal(i)
	if err != nil {
		t.Fatal(err)
	}
	return b
}

func Test_inMemory_Get(t *testing.T) {
	type fields struct {
		mu sync.RWMutex
		db map[string][]byte
	}
	type args struct {
		key string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    []byte
		wantErr bool
	}{
		{
			name: "should return the value",
			fields: fields{
				mu: sync.RWMutex{},
				db: map[string][]byte{
					"test": marshalItem(t, Item{
						Value: []byte("test"),
						TTL:   0,
					}),
				},
			},
			args: args{
				key: "test",
			},
			want:    []byte("test"),
			wantErr: false,
		},
		{
			name: "should return nil if the key does not exist",
			fields: fields{
				mu: sync.RWMutex{},
				db: map[string][]byte{},
			},
			args: args{
				key: "test",
			},
			want:    nil,
			wantErr: false,
		},
		{
			name: "should return nil if the key has expired",
			fields: fields{
				mu: sync.RWMutex{},
				db: map[string][]byte{
					"test": marshalItem(t, Item{
						Value: []byte("test"),
						TTL:   time.Now().Unix() - 20,
					}),
				},
			},
			args: args{
				key: "test",
			},
			want:    nil,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &inMemory{
				mu: tt.fields.mu,
				db: tt.fields.db,
			}
			got, err := s.Get(tt.args.key)
			if (err != nil) != tt.wantErr {
				t.Errorf("inMemory.Get() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("inMemory.Get() = %v, want %v", string(got), string(tt.want))
			}
		})
	}
}
