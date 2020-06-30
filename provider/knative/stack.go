package knative

import (
	"path/filepath"

	knativestack "github.com/nropatas/faastest-stacks/knative"
	"github.com/nuweba/faasbenchmark/stack"
	kn "github.com/nuweba/faasbenchmark/stack/knative"
)

func (knative *Knative) NewStack(stackPath string) (stack.Stack, error) {
	stackPath = filepath.Join(stackPath, knative.Name())
	stack, err := knativestack.New(stackPath)

	if err != nil {
		return nil, err
	}

	return &kn.Stack{KnativeStack: stack}, nil
}
