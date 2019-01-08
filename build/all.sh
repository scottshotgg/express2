#!/bin/sh

# This script builds a the Express compiler (ecc) for every supported architecture

mkdir -p build
mkdir -p build/backups

#archs = (amd64, 386, arm64, arm)

# Grab the date
date=$(date +%d.%m.%y-%H:%M)

echo "      Beginning build process"
echo "------------------------------------\n"
echo "Archiving old dist folder into backups/dist_$date ...\n"

# Archive the dist folder to preserve all of the current binaries
tar -zcvf "build/backups/dist_$date.tar.gz" "build/dist"
echo "\n"

# Build x86
arch=386
echo "Building AMD64 binaries ..."

echo "Building Linux ..."
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc"
echo "Building OSX ..."
GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc"
echo "Building Windows ..."
GOOS=windows GOARCH=$arch go build -o "build/dist/windows/$arch/ecc.exe"
echo "\n"

# Build AMD64
arch=amd64
echo "Building x86 binaries ..."

echo "Building Linux ..."
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc"
echo "Building OSX ..."
GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc"
echo "Building Windows ..."
GOOS=windows GOARCH=$arch go build -o "build/dist/windows/$arch/ecc.exe"
echo ""

# Build arm
arch=arm
echo "Building ARM binaries ..."

echo "Building Linux ..."
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc"
#GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc.exe"
#GOOS=windows GOARCH=$arch go build -o "build/dist/windows/$arch/ecc.exe"
echo ""

# Build arm64
arch=arm64
echo "Building ARM64 binaries ..."

echo "Building Linux ..."
GOOS=linux GOARCH=$arch go build -o "build/dist/linux/$arch/ecc"
#GOOS=darwin GOARCH=$arch go build -o "build/dist/osx/$arch/ecc.exe"
#GOOS=windows GOARCH=$arch go build -o "build/dist/windows/$arch/ecc.exe"
# echo ""