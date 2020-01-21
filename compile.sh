#!/bin/bash
OsArray=("linux" "darwin" "windows")

for os in ${OsArray[*]}; do
  name="go-file-copies"
  if [ "$os" = "windows" ]; then
    name+=".exe"
  fi
  GOOS=$os go build -o build/"$os"/$name
  cp config-sample.json build/"$os"/
done
