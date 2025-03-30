package pac

import "github.com/dop251/goja"

const (
	resolveFuncName = "FindProxyForURL"
)

type resolveFunc func(url, host string) *string

func compile(source []byte) (resolveFunc, error) {
	vm := goja.New()
	declareBuiltins(vm)

	if _, err := vm.RunString(string(source)); err != nil {
		return nil, err
	}

	var resolve resolveFunc
	if err := vm.ExportTo(vm.Get(resolveFuncName), &resolve); err != nil {
		return nil, err
	}

	return resolve, nil
}
