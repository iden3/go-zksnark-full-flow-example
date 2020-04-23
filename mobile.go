package zkflowexample

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/iden3/go-circom-prover-verifier/parsers"
	"github.com/iden3/go-circom-prover-verifier/prover"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	witnesscalc "github.com/iden3/go-circom-witnesscalc"
)

type MobileZKFlow struct{}

func (m *MobileZKFlow) Run(storePath string) error {
	fmt.Println("Mobile full flow for circuit1")
	err := m.runCircuit(storePath, "http://161.35.72.58:9000/circuit1", IdStateInputs)
	if err != nil {
		return err
	}

	fmt.Println("\nMobile full flow for circuit2")
	err = m.runCircuit(storePath, "http://161.35.72.58:9000/circuit2", func() (string, error) {
		return `{"in":"1"}`, nil
	})
	if err != nil {
		return err
	}

	return nil
}

func (m *MobileZKFlow) runCircuit(storePath, filesServer string, funcInputs funcInputsGenerator) error {
	const wasmFilename = "circuit.wasm"
	const pkFilename = "proving_key.json"
	const vkFilename = "verification_key.json"
	// download artifacts: circuit in wasm, proving key and verification key.
	printT("=== Download artifacts:")
	printT("==> Downloading circuit...")
	if err := downloadFile(
		storePath+wasmFilename,
		filesServer+"/circuit.wasm",
	); err != nil {
		return err
	}
	printT("Done")
	printT("==> Downloading proving key...")
	if err := downloadFile(
		storePath+pkFilename,
		filesServer+"/proving_key.json",
	); err != nil {
		return err
	}
	printT("Done")
	printT("==> Downloading verification key...")
	if err := downloadFile(
		storePath+vkFilename,
		filesServer+"/verification_key.json",
	); err != nil {
		return err
	}
	printT("Done")
	fmt.Print("=============\n\n\n")

	inputsJson, err := funcInputs()
	if err != nil {
		return err
	}
	printT("Done")

	printT("=== Generate ZKP:")
	printT("==> Parsing inputs...")
	inputs, err := witnesscalc.ParseInputs([]byte(inputsJson))
	if err != nil {
		return err
	}
	printT("Done")
	printT("==> Calculating witness...")
	wasmBytes, err := ioutil.ReadFile(storePath + wasmFilename)
	if err != nil {
		return err
	}
	wit, err := ComputeWitness(wasmBytes, inputs)
	if err != nil {
		return err
	}
	witnessJson, err := json.Marshal(witnesscalc.WitnessJSON(wit))
	if err != nil {
		return err
	}
	w, err := parsers.ParseWitness(witnessJson)
	if err != nil {
		return err
	}
	printT("Done")

	printT("==> Loading proving key...")
	provingKeyJson, err := ioutil.ReadFile(storePath + pkFilename)
	if err != nil {
		return err
	}
	pk, err := parsers.ParsePk(provingKeyJson)
	if err != nil {
		return err
	}
	printT("Done")

	printT("==> Generating proof...")
	proof, pubSignals, err := prover.GenerateProof(pk, w)
	if err != nil {
		return err
	}
	printT("Done")

	fmt.Println("Proof generated successfuly!")
	proofStr, err := parsers.ProofToJson(proof)
	if err != nil {
		return err
	}
	fmt.Println("Proof: ", string(proofStr))
	publicStr, err := json.Marshal(parsers.ArrayBigIntToString(pubSignals))
	if err != nil {
		return err
	}
	fmt.Println("Public inputs: ", string(publicStr))
	fmt.Print("=============\n\n\n")

	printT("=== Verify ZKP:")
	printT("==> Loading verification key...")
	vkJson, err := ioutil.ReadFile(storePath + vkFilename)
	if err != nil {
		return err
	}
	vk, err := parsers.ParseVk(vkJson)
	if err != nil {
		return err
	}
	printT("Done")

	printT("==> Verifying ZKP...")
	v := verifier.Verify(vk, proof, pubSignals)
	printT("Done")

	fmt.Println("The result of the verification is: ", v)
	return nil
}
