#!/bin/sh

cd testdata
./clean-generated-files.sh

echo $(date +"%T") "./compile-and-trustedsetup.sh"
itime="$(date -u +%s)"
./compile-and-trustedsetup.sh
ftime="$(date -u +%s)"
echo "	Finish compile-and-trustedsetup.sh ($(($(date -u +%s)-$itime))s)"
cd ..

echo $(date +"%T") "go test -run TestFullFlow"
itime="$(date -u +%s)"
go test -run=TestFlowLocal
echo "	Finish go test ($(($(date -u +%s)-$itime))s)"
