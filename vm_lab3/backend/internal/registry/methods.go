package registry

import (
	"fmt"
	"vm_lab3/internal/service"
)

var (
	leftRectSolver      *service.LeftRectSolver      = service.NewLeftRectSolver()
	rightRectSolver     *service.RightRectSolver     = service.NewRightRectSolver()
	centralRectSolver   *service.CentralRectSolver   = service.NewCentralRectSolver()
	trapezoidRectSolver *service.TrapezoidRectSolver = service.NewTrapezoidRectSolver()
	simpsonSolver       *service.SimpsonSolver       = service.NewSimpsonSolver()
)

var (
	id2method = map[int]service.IntegralSolver{
		1: leftRectSolver,
		2: rightRectSolver,
		3: centralRectSolver,
		4: trapezoidRectSolver,
		5: simpsonSolver,
	}
)

func GetMethod(id int) (*service.IntegralSolver, error) {
	method := id2method[id]
	if method == nil {
		return nil, fmt.Errorf("no method with id=%d", id)
	}
	return &method, nil
}
