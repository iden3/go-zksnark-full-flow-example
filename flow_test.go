package zkflowexample

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"testing"

	"github.com/iden3/go-circom-prover-verifier/parsers"
	"github.com/iden3/go-circom-prover-verifier/prover"
	"github.com/iden3/go-circom-prover-verifier/verifier"
	witnesscalc "github.com/iden3/go-circom-witnesscalc"
	"github.com/stretchr/testify/require"
)

type funcInputsGenerator func() (string, error)

func TestFullFlow(t *testing.T) {
	// compile circuits & compute trusted setup using
	// compile-and-trustedsetup.sh

	testFullFlowCircuit(t, "circuit1", IdStateInputs)
	testFullFlowCircuit(t, "circuit2", func() (string, error) {
		return `{"in":"1"}`, nil
	})
}

func testFullFlowCircuit(t *testing.T, circuitdir string, funcInputs funcInputsGenerator) {
	fmt.Println("test full flow for circuit:", circuitdir)

	// build the testing environment: claims, identities, merkletrees, etc
	printT("- Generate inputs")
	inputsJson, err := funcInputs()
	require.Nil(t, err)
	err = ioutil.WriteFile("testdata/"+circuitdir+"/inputs.json", []byte(inputsJson), 0644)
	require.Nil(t, err)

	// calculate witness
	printT("- Parse witness files")
	wasmFilename := "testdata/" + circuitdir + "/circuit.wasm"
	inputsFilename := "testdata/" + circuitdir + "/inputs.json"

	wasmBytes, err := ioutil.ReadFile(wasmFilename)
	require.Nil(t, err)
	inputsBytes, err := ioutil.ReadFile(inputsFilename)
	require.Nil(t, err)
	inputs, err := witnesscalc.ParseInputs(inputsBytes)
	require.Nil(t, err)

	printT("- Calculate witness")
	wit, err := ComputeWitness(wasmBytes, inputs)
	require.Nil(t, err)

	printT("- Write witness output files")
	wJSON, err := json.Marshal(witnesscalc.WitnessJSON(wit))
	require.Nil(t, err)
	err = ioutil.WriteFile("testdata/"+circuitdir+"/witness.json", []byte(wJSON), 0644)
	require.Nil(t, err)

	// generate zk proof
	// read ProvingKey & Witness files
	printT("- Read proving_key.json & witness.json files")
	provingKeyJson, err := ioutil.ReadFile("testdata/" + circuitdir + "/proving_key.json")
	require.Nil(t, err)
	witnessJson, err := ioutil.ReadFile("testdata/" + circuitdir + "/witness.json")
	require.Nil(t, err)

	// parse Proving Key
	pk, err := parsers.ParsePk(provingKeyJson)
	require.Nil(t, err)
	// parse Witness
	w, err := parsers.ParseWitness(witnessJson)
	require.Nil(t, err)

	// generate the proof
	printT("- Generate zkProof")
	proof, pubSignals, err := prover.GenerateProof(pk, w)
	require.Nil(t, err)

	// print proof & publicSignals
	proofStr, err := parsers.ProofToJson(proof)
	require.Nil(t, err)
	publicStr, err := json.Marshal(parsers.ArrayBigIntToString(pubSignals))
	require.Nil(t, err)

	err = ioutil.WriteFile("testdata/"+circuitdir+"/proof.json", proofStr, 0644)
	require.Nil(t, err)
	err = ioutil.WriteFile("testdata/"+circuitdir+"/public.json", publicStr, 0644)
	require.Nil(t, err)

	// verify zk proof
	// read proof & verificationKey & publicSignals
	proofJson, err := ioutil.ReadFile("testdata/" + circuitdir + "/proof.json")
	require.Nil(t, err)
	printT("- Read verification_key.json & public.json files")
	vkJson, err := ioutil.ReadFile("testdata/" + circuitdir + "/verification_key.json")
	require.Nil(t, err)
	publicJson, err := ioutil.ReadFile("testdata/" + circuitdir + "/public.json")
	require.Nil(t, err)

	// parse proof & verificationKey & publicSignals
	public, err := parsers.ParsePublicSignals(publicJson)
	require.Nil(t, err)
	proofParsed, err := parsers.ParseProof(proofJson)
	require.Nil(t, err)
	vk, err := parsers.ParseVk(vkJson)
	require.Nil(t, err)

	// verify the proof with the given verificationKey & publicSignals
	printT("- Verify zkProof")
	v := verifier.Verify(vk, proof, pubSignals)
	fmt.Println("verifier.Verify", v)

	// verify but using the stored files
	vkJson, err = ioutil.ReadFile("testdata/" + circuitdir + "/verification_key.json")
	require.Nil(t, err)
	vk, err = parsers.ParseVk(vkJson)
	require.Nil(t, err)
	v2 := verifier.Verify(vk, proofParsed, public)
	fmt.Println("verifier.Verify", v2)
}
