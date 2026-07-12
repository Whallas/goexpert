package temperature

// Temperatures holds a temperature value expressed in the three scales
// required by the challenge contract.
type Temperatures struct {
	Celsius    float64 `json:"temp_C"`
	Fahrenheit float64 `json:"temp_F"`
	Kelvin     float64 `json:"temp_K"`
}

// FromCelsius converts a Celsius reading into all supported scales.
//
//	F = C × 1.8 + 32
//	K = C + 273
func FromCelsius(celsius float64) Temperatures {
	return Temperatures{
		Celsius:    celsius,
		Fahrenheit: celsius*1.8 + 32,
		Kelvin:     celsius + 273,
	}
}
