package solver

type Point struct {
	X float64 `json:"x"`
	Y float64 `json:"y"`
}

type Result struct {
	Solution   float64  `json:"solution"`
	Value      float64  `json:"value"`
	Iterations int      `json:"iterations"`
	Steps      []Point  `json:"steps"`
	ArgError   float64  `json:"arg_error"`
	Messages   []string `json:"messages"`
}

type System2Result struct {
	Solution   Point    `json:"solution"`
	Iterations int      `json:"iterations"`
	Messages   []string `json:"messages"`
	Steps      []Point  `json:"steps"`
}
