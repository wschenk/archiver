#!/bin/bash

mkdir -p build/arm

for i in cmd/*; do
  echo Building $i
  file=${i##*/}
  base=${file%.*}
  echo $base
  GOOS=linux GOARCH=arm go build -v $i
  mv $base build/arm
done

ipfs add -r build/arm
