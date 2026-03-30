package dto

type CalcIntegralRequestDto struct {
	FuncID int     `json:"func_id"`
	A      float64 `json:"a"`
	B      float64 `json:"b"`
	Eps    float64 `json:"eps"`
	N      *int    `json:"n,omitempty"`
}

type CalcIntegralResponseDto struct {
	Value    float64  `json:"value"`
	N        int      `json:"n"`
	Messages []string `json:"messages"`
	RungeR   float64  `json:"runge_r"`
}

type ErrorResponse struct {
	Error string `json:"error"`
}
