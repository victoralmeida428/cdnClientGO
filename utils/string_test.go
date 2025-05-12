package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestCamelCase(t *testing.T) {
	tests := []struct {
		in, want string
	}{
		{
			in:   "teste_fuzzy",
			want: "TesteFuzzy",
		},
	}
	
	for _, tt := range tests {
		t.Run(tt.in, func(t *testing.T) {
			assert.Equal(t, ToCamelCase(tt.in), tt.want)
		})
	}
}

func FuzzCamelCase(f *testing.F) {
	f.Add("teste_fuzzy")
	
	f.Fuzz(func(t *testing.T, in string) {
		ToCamelCase(in)
	})
}
