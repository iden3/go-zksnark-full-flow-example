#!/bin/sh

echo $(date +"%T") "circom circuit.circom --r1cs --wasm --sym"
itime="$(date -u +%s)"
circom circuit.circom --r1cs --wasm --sym
ftime="$(date -u +%s)"
echo "	($(($(date -u +%s)-$itime))s)"

echo $(date +"%T") "snarkjs info -r circuit.r1cs"
snarkjs info -r circuit.r1cs

echo $(date +"%T") "snarkjs setup"
itime="$(date -u +%s)"
snarkjs setup
echo "	($(($(date -u +%s)-$itime))s)"
echo $(date +"%T") "trusted setup generated"

sed -i 's/null/["0","0","0"]/g' proving_key.json


echo $(date +"%T") "snarkjs generateverifier"
itime="$(date -u +%s)"
snarkjs generateverifier
echo "	($(($(date -u +%s)-$itime))s)"
echo $(date +"%T") "generateverifier generated"
