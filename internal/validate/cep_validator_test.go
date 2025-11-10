package validate

import "testing"

func TestCepValidator(t *testing.T) {
	tests := []struct {
		name string
		cep  string
		want bool
	}{
		{name: "valid cep", cep: "12345000", want: true},
		{name: "too short", cep: "1234567", want: false},
		{name: "too long", cep: "123456789", want: false},
		{name: "non numeric zone characters", cep: "12A45B78", want: false},
		{name: "non numeric suffix characters", cep: "12345AB!", want: false},
		{name: "zone too small", cep: "00999000", want: false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := CepValidator(tt.cep); got != tt.want {
				t.Fatalf("CepValidator(%q) = %v, want %v", tt.cep, got, tt.want)
			}
		})
	}
}
