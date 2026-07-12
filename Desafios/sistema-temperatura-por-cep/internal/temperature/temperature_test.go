package temperature

import "testing"

func TestFromCelsius(t *testing.T) {
	tests := []struct {
		name    string
		celsius float64
		want    Temperatures
	}{
		{
			name:    "positive temperature",
			celsius: 28.5,
			want:    Temperatures{Celsius: 28.5, Fahrenheit: 83.3, Kelvin: 301.5},
		},
		{
			name:    "zero",
			celsius: 0,
			want:    Temperatures{Celsius: 0, Fahrenheit: 32, Kelvin: 273},
		},
		{
			name:    "negative temperature",
			celsius: -10,
			want:    Temperatures{Celsius: -10, Fahrenheit: 14, Kelvin: 263},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := FromCelsius(tt.celsius)
			if !almostEqual(got.Celsius, tt.want.Celsius) ||
				!almostEqual(got.Fahrenheit, tt.want.Fahrenheit) ||
				!almostEqual(got.Kelvin, tt.want.Kelvin) {
				t.Errorf("FromCelsius(%v) = %+v, want %+v", tt.celsius, got, tt.want)
			}
		})
	}
}

func almostEqual(a, b float64) bool {
	const epsilon = 1e-9
	diff := a - b
	if diff < 0 {
		diff = -diff
	}
	return diff < epsilon
}
