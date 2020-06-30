package knative

import (
	"github.com/nropatas/faastest-stacks/knative"
	"github.com/nuweba/faasbenchmark/stack"
)

type Stack struct {
	*knativestack.KnativeStack
}

func New(path string) (*Stack, error) {
	stack, err := knativestack.New(path)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}

func (s *Stack) ListFunctions() []stack.Function {
	var functions []stack.Function

	for _, f := range s.Functions {
		function := &Function{
			name: 			f.Name,
			handler: 		f.Handler,
			description:	f.Description,
			runtime: 		f.Runtime,
			memorySize: 	f.MemorySize,
		}
		functions = append(functions, function)
	}

	return functions
}
