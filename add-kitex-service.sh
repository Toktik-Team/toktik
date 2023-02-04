#!/usr/bin/bash

mkdir -p kitex_gen
kitex -module "toktik" -I idl/ idl/"$1".proto

mkdir -p service/"$1"
cd service/"$1" && kitex -module "toktik" -service "$1" -use toktik/kitex_gen/ -I ../../idl/ ../../idl/"$1".proto

go mod tidy
