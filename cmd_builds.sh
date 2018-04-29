#!/bin/bash

mkdir -p build/arm
mkdir -p build/darwin

for i in cmd/*.go; do
  echo Building $i
  file=${i##*/}
  base=${file%.*}
  echo $base
  GOOS=linux GOARCH=arm go build -v $i
  mv $base build/arm
  go build -v $i
  mv $base build/darwin
done

ipfs add -r build/arm
