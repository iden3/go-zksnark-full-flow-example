module zksnark-full-flow-example

replace github.com/iden3/go-circom-prover-verifier => ../../iden3/go-circom-prover-verifier

replace github.com/iden3/go-circom-witnesscalc => ../../iden3/go-circom-witnesscalc

go 1.14

require (
	github.com/iden3/go-circom-prover-verifier v0.0.0-00010101000000-000000000000
	github.com/iden3/go-circom-witnesscalc v0.0.0-00010101000000-000000000000
	github.com/iden3/go-iden3-core v0.0.8-0.20200417110412-71c75d3b9482
	github.com/iden3/go-iden3-crypto v0.0.4
	github.com/iden3/go-wasm3 v0.0.0-20200407092348-656263e6984f
	github.com/stretchr/testify v1.5.1
)
