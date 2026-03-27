package dto

type CalcIntegralRequestDto struct {
	ID       int     `json:"id"`
	A        float64 `json:"a"`
	B        float64 `json:"b"`
	MethodID int     `json:"method_id"`
	Eps      float64 `json:"eps"`
	N        *int    `json:"n,omitempty"`
}

type CalcIntegralResponseDto struct {
	Value    float64  `json:"value"`
	N        int      `json:"n"`
	Messages []string `json:"messages"`
}
