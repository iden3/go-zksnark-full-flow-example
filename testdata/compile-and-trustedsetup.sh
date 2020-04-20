#!/bin/sh

echo "circom circuit.circom --r1cs --wasm --sym"
circom circuit.circom --r1cs --wasm --sym

echo "snarkjs info -r circuit.r1cs"
snarkjs info -r circuit.r1cs

date +"%T"
echo "snarkjs setup"
snarkjs setup
date +"%T"

sed -i 's/null/["0","0","0"]/g' proving_key.json
