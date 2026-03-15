package equations

type System2Equation func(x float64, y float64) float64

type System2 [2]System2Equation

func NewSystem2(e1 System2Equation, e2 System2Equation) *System2 {
	return &System2{e1, e2}
}
