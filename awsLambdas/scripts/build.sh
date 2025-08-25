#!/bin/bash

mkdir -p bin
bin_folder=bin
bin_lambda_folder=$bin_folder/lambdas
lambdas_folder=./cmd/lambdas

echo "removing all files"
rm -Rv $bin_folder

echo "building files ..."

packages=$(go list -f '{{.ImportPath}}' ./$lambdas_folder/...)

# GOARCH=amd64 CGO_ENABLED=0 go build -C go -o ../$bin_folder/main main.go

# https://docs.aws.amazon.com/lambda/latest/dg/golang-package.html
# GOOS=linux GOARCH=amd64 go build -o bootstrap main.go

# generate a binary in a light version
# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -ldflags "-s -w" -o ../build/main-light go/main.go

# verbose
# go build -v -work -x -json -o bin/lambdas/
# GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -v -work -x -json -o $bin_lambda_folder/ $packages
GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o $bin_lambda_folder/ $packages

echo "zipping files ..."

for file in "$bin_lambda_folder"/*; do
  # Skip if it's not a regular file
  [ -f "$file" ] || continue

  # Get the base filename (e.g., hello-go)
  basename="${file##*/}"
  zipname="${basename}.zip"

  # rename to bootstrap
  cp "$file" bootstrap

  # Create the zip with bootstrap inside
  zip -q "$bin_lambda_folder/$zipname" bootstrap

  # Clean up the temporary bootstrap file
  rm bootstrap 
  rm "$file"

  echo "Created $zipname with contents renamed to bootstrap"
done