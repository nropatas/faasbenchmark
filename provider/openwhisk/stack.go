package openwhisk

import (
	"path/filepath"

	"github.com/nuweba/faasbenchmark/stack"
	"github.com/nuweba/faasbenchmark/stack/sls"
)

type Stack struct {
	*sls.Stack
}

func (openwhisk *OpenWhisk) NewStack(stackPath string) (stack.Stack, error) {
	stackPath = filepath.Join(stackPath, openwhisk.Name())
	stack, err := sls.New(openwhisk.Name(), stackPath)

	if err != nil {
		return nil, err
	}

	return &Stack{stack}, nil
}
