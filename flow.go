package zkflowexample

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"strconv"
	"time"

	"github.com/iden3/go-circom-prover-verifier/parsers"
	"github.com/iden3/go-circom-prover-verifier/prover"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	witnesscalc "github.com/iden3/go-circom-witnesscalc"
	"github.com/iden3/go-wasm3"
)

type funcInputsGenerator func() (string, error)

func ExecuteFlow(path, inputsJson string) (string, error) {
	fmt.Println("path:", path)

	printT("- Generate inputs")
	err := ioutil.WriteFile(path+"/inputs.json", []byte(inputsJson), 0644)
	if err != nil {
		return "", err
	}

	printT("- Parse witness file")
	wasmBytes, err := ioutil.ReadFile(path + "/circuit.wasm")
	if err != nil {
		return "", err
	}
	inputsBytes, err := ioutil.ReadFile(path + "/inputs.json")
	if err != nil {
		return "", err
	}
	inputs, err := witnesscalc.ParseInputs(inputsBytes)
	if err != nil {
		return "", err
	}
	printT("- Calculate witness")
	wit, err := ComputeWitness(wasmBytes, inputs)
	if err != nil {
		return "", err
	}

	printT("- Write witness output files")
	wJSON, err := json.Marshal(witnesscalc.WitnessJSON(wit))
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(path+"/witness.json", []byte(wJSON), 0644)
	if err != nil {
		return "", err
	}

	witnessJson, err := json.Marshal(witnesscalc.WitnessJSON(wit))
	if err != nil {
		return "", err
	}
	w, err := parsers.ParseWitness(witnessJson)
	if err != nil {
		return "", err
	}

	printT("- Load proving_key.json")
	provingKeyJson, err := ioutil.ReadFile(path + "/proving_key.json")
	if err != nil {
		return "", err
	}
	pk, err := parsers.ParsePk(provingKeyJson)
	if err != nil {
		return "", err
	}

	printT("- Generate zkProof")
	beforeT := time.Now()
	proof, pubSignals, err := prover.GenerateProof(pk, w)
	if err != nil {
		return "", err
	}
	proofGenTime := time.Since(beforeT)
	fmt.Println("proof generation time elapsed:", proofGenTime)

	proofStr, err := parsers.ProofToJson(proof)
	if err != nil {
		return "", err
	}
	publicStr, err := json.Marshal(parsers.ArrayBigIntToString(pubSignals))
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(path+"/proof.json", proofStr, 0644)
	if err != nil {
		return "", err
	}
	err = ioutil.WriteFile(path+"/public.json", publicStr, 0644)
	if err != nil {
		return "", err
	}

	printT("- Loading verification_key.json")
	vkJson, err := ioutil.ReadFile(path + "/verification_key.json")
	if err != nil {
		return "", err
	}
	vk, err := parsers.ParseVk(vkJson)
	if err != nil {
		return "", err
	}

	printT("- Verify zkProof")
	v := verifier.Verify(vk, proof, pubSignals)
	fmt.Println("verifier.Verify", v)
	if !v {
		fmt.Errorf("zkProof verification failed")
	}

	elapsedSec := float64(proofGenTime) / float64(time.Second)
	return strconv.FormatFloat(elapsedSec, 'f', 6, 64) + "s", nil
}

type MobileWrapper struct{}

func (m *MobileWrapper) ExecuteFlowDownloading(path, filesServer, generatedInputs string) (string, error) {
	downloadFiles(path, filesServer)

	r, err := ExecuteFlow(path, generatedInputs)
	if err != nil {
		return "", err
	}

	return r, nil
}

func ComputeWitness(wasmBytes []byte, inputs []witnesscalc.Input) ([]*big.Int, error) {
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
