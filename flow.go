package zkflowexample

import (
	"math/big"

	witnesscalc "github.com/iden3/go-circom-witnesscalc"
	"github.com/iden3/go-wasm3"
)

type funcInputsGenerator func() (string, error)

func ComputeWitness(wasmBytes []byte, inputs map[string]interface{}) ([]*big.Int, error) {
	runtime := wasm3.NewRuntime(&wasm3.Config{
		Environment: wasm3.NewEnvironment(),
		StackSize:   64 * 1024,
	})
	defer runtime.Destroy()

	module, err := runtime.ParseModule(wasmBytes)
	if err != nil {
		return nil, err
	}
	module, err = runtime.LoadModule(module)
	if err != nil {
		return nil, err
	}

	witnessCalculator, err := witnesscalc.NewWitnessCalculator(runtime, module)
	if err != nil {
		return nil, err
	}

	wit, err := witnessCalculator.CalculateWitness(inputs, false)
	if err != nil {
		return nil, err
	}
	return wit, nil
}
