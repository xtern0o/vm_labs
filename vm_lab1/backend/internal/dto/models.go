package dto

type SolveRequest struct {
	Matrix       [][]float64 `json:"matrix" validate:"required"`
	Vector       []float64   `json:"vector" validate:"required"`
	InitialGuess []float64   `json:"initial_guess,omitempty"`
	Epsilon      *float64    `json:"epsilon,omitempty"`
	MaxIter      *int        `json:"max_iter,omitempty"`
}

type SolveResponse struct {
	Solution   []float64 `json:"solution"`
	Errors     []float64 `json:"errors"`
	Iterations int       `json:"iterations"`
}

type ErrorResponse struct {
	Error   string `json:"error"`
	Details string `json:"details,omitempty"`
}
