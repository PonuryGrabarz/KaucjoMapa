package handler

import (
	"reflect"
	"testing"

	"github.com/PonuryGrabarz/KaucjoMapa/backend/internal/entity"
)

func TestParseRefundMethods(t *testing.T) {
	tests := []struct {
		name    string
		value   string
		want    []entity.RefundMethod
		wantErr bool
	}{
		{name: "empty", value: "", want: nil},
		{
			name:  "multiple methods",
			value: "gotowka, karta,gotowka",
			want:  []entity.RefundMethod{entity.RefundMethodCash, entity.RefundMethodCard},
		},
		{name: "unknown method", value: "blik", wantErr: true},
		{name: "empty list item", value: "gotowka,,karta", wantErr: true},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := parseRefundMethods(tt.value)
			if (err != nil) != tt.wantErr {
				t.Fatalf("parseRefundMethods() error = %v, wantErr %v", err, tt.wantErr)
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("parseRefundMethods() = %v, want %v", got, tt.want)
			}
		})
	}
}
