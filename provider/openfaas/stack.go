package openfaas

import (
	"fmt"
	"path/filepath"

	openfaasstack "github.com/nropatas/faastest-stacks/openfaas"
	"github.com/nuweba/faasbenchmark/stack"
	faas "github.com/nuweba/faasbenchmark/stack/openfaas"
)

func (openfaas *OpenFaaS) NewStack(stackPath string) (stack.Stack, error) {
	stackPath = filepath.Join(stackPath, openfaas.Name())
	stack, err := openfaasstack.New(stackPath, fmt.Sprintf("http://%s", openfaas.IngressUrl()))

	if err != nil {
		return nil, err
	}

	return &faas.Stack{OpenFaaSStack: stack}, nil
}
