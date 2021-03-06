#!/bin/sh

npm install

compile_and_ts() {
  echo $(date +"%T") "circom circuit.circom --r1cs --wasm --sym"
  itime="$(date -u +%s)"
  ../node_modules/.bin/circom circuit.circom --r1cs --wasm --sym
  ftime="$(date -u +%s)"
  echo "	($(($(date -u +%s)-$itime))s)"

  echo $(date +"%T") "snarkjs info -r circuit.r1cs"
  ../node_modules/.bin/snarkjs info -r circuit.r1cs

  echo $(date +"%T") "snarkjs setup"
  itime="$(date -u +%s)"
  ../node_modules/.bin/snarkjs setup
  echo "	($(($(date -u +%s)-$itime))s)"
  echo $(date +"%T") "trusted setup generated"

  sed -i 's/null/["0","0","0"]/g' proving_key.json


  echo $(date +"%T") "snarkjs generateverifier"
  itime="$(date -u +%s)"
  ../node_modules/.bin/snarkjs generateverifier
  echo "	($(($(date -u +%s)-$itime))s)"
  echo $(date +"%T") "generateverifier generated"
}

echo "compile & trustesetup for circuit1"
cd circuit1
compile_and_ts

echo "compile & trustesetup for circuit2"
cd ../circuit2
compile_and_ts

echo "compile & trustesetup for circuit3"
cd ../circuit3
compile_and_ts
