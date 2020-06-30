package kubeless

import (
	"path/filepath"

	kubelessstack "github.com/nropatas/faastest-stacks/kubeless"
	"github.com/nuweba/faasbenchmark/stack"
	kl "github.com/nuweba/faasbenchmark/stack/kubeless"
)

func (kubeless *Kubeless) NewStack(stackPath string) (stack.Stack, error) {
	stackPath = filepath.Join(stackPath, kubeless.Name())
	stack, err := kubelessstack.New(stackPath)

	if err != nil {
		return nil, err
	}

	return &kl.Stack{KubelessStack: stack}, nil
}
