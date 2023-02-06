#!/bin/bash

# Remove the output directory
echo "Removing previous output directory..."
rm -rf output

# Create the output directory and binaries directory
echo "Creating output directory..."
mkdir -p output/bin

# Navigate to the service directory
echo "Navigating to service directory..."
cd service || exit

# Loop through all directories in the service directory
for d in */ ; do
  echo "Building $d..."
  # Navigate into the current directory
  cd "$d"

  # Remove the output directory
  echo "Removing previous output directory..."
  rm -rf output

  # Create the output directory
  echo "Creating output directory..."
  mkdir output

  # Check if the build.sh file exists, if not, skip the current directory
  if [ ! -f build.sh ]; then
    echo "build.sh file not found in $d, skipping..."
    continue
  fi

  # Execute the build.sh script
  echo "Executing build.sh..."
  ./build.sh

  # Copy the binary to the output/binaries directory
  echo "Copying binary to output/bin..."
  cp output/bin/"$(basename "$d")" ../../output/bin/

  # Check if the bootstrap.sh file exists
  if [ -f output/bootstrap.sh ]; then
    echo "bootstrap.sh file found in $d, copying to output directory..."
    # Rename the bootstrap.sh file and copy it to the output directory
    mv output/bootstrap.sh ../../output/bootstrap-"$(basename "$d")".sh
  fi

  # Navigate back to the service directory
  cd ..
  echo "Finished building $d"
done
