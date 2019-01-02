#!/bin/sh

# This script builds a the Express compiler (ecc) 
# for every supported architecture

#archs = (amd64, 386, arm64, arm)

date=$(date +%d.%m.%y-%H:%M)

# Build AMD64
arch=amd64
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc_"$arch"_$date.exe"
GOOS=windows GOARCH=$arch go build -o "build/dist/windows/$arch/ecc_"$arch"_$date.exe"
GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc_"$arch"_$date.exe"

# Build x86
arch=386
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc_"$arch"_$date.exe"
GOOS=windows GOARCH=$arch go build -o "build/dist/windows/$arch/ecc_"$arch"_$date.exe"
GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc_"$arch"_$date.exe"

# Build arm
arch=arm
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc_"$arch"_$date.exe"
#GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc_"$arch"_$date.exe"

# Build arm64
arch=arm64
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc_"$arch"_$date.exe"
#GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc_"$arch"_$date.exe"