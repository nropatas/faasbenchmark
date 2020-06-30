package openfaas

import (
	openfaasstack "github.com/nropatas/faastest-stacks/openfaas"
	"github.com/nuweba/faasbenchmark/stack"
)

type Stack struct {
	*openfaasstack.OpenFaaSStack
}

func New(path string, gatewayUrl string) (*Stack, error) {
	stack, err := openfaasstack.New(path, gatewayUrl)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}

func (s *Stack) ListFunctions() []stack.Function {
	var functions []stack.Function

	for _, f := range s.Functions {
		function := &Function{
			name:        f.Name,
			handler:     f.Handler,
			description: f.Description,
			runtime:     f.Runtime,
			memorySize:  f.MemorySize,
		}
		functions = append(functions, function)
	}

	return functions
}
