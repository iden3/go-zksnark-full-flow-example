# go-zksnark-full-flow-example

zkSNARK full flow example in Go

- Compile the circuit using [circom](https://github.com/iden3/circom)
- Generate the trusted setup using [snarkjs](https://github.com/iden3/snarkjs)
- Calculate Witness from Go using [go-circom-witnesscalc](https://github.com/iden3/go-circom-witnesscalc)
- Generate zkProof from Go using [go-circom-prover-verifier](https://github.com/iden3/go-circom-prover-verifier)
- Verify zkProof from Go using [go-circom-prover-verifier](https://github.com/iden3/go-circom-prover-verifier)

![](go-zksnark-flow.png)

## Usage
- compile the circuit and generate trusted setup
```
cd testdata
./compile-and-trustedsetup.sh
```

- the rest of the flow
```
go test
```
