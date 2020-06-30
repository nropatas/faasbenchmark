package fission

import (
	"path/filepath"

	fissionstack "github.com/nropatas/faastest-stacks/fission"
	"github.com/nuweba/faasbenchmark/stack"
	fs "github.com/nuweba/faasbenchmark/stack/fission"
)

func (fission *Fission) NewStack(stackPath string) (stack.Stack, error) {
	stackPath = filepath.Join(stackPath, fission.Name())
	stack, err := fissionstack.New(stackPath)

	if err != nil {
		return nil, err
	}

	return &fs.Stack{FissionStack: stack}, nil
}
