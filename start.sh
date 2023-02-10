#!/bin/bash

output_dir="output"
services_dir="service/"
services=$(ls "$services_dir")

while [[ $# -gt 0 ]]
do
  case "$1" in
    --service)
      service_name="$2"
      shift 2
      ;;
    *)
      echo "Unknown option: $1"
      exit 1
      ;;
  esac
done

if [ -z "$service_name" ]; then
  echo "Error: --service option is required."
  exit 1
fi

found=false
for s in $services; do
  if [ "$s" == "$service_name" ]; then
    found=true
    break
  fi
done

if [ "$found" = false ]; then
  echo "Error: Unrecognized service name: $service_name"
  printf 'Available service names:\n%s\n' "$services"
  exit 1
fi

command="$output_dir/bin/$service_name"

# Check if the bootstrap.sh file exists
if [ -f output/bootstrap-"${service_name}".sh ]; then
  command="$output_dir/bootstrap-${service_name}.sh"
fi

if [ ! -f "$command" ]; then
  echo "Error: Service binary not found: $command"
  exit 1
fi

"$command"