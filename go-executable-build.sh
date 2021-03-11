#!/usr/bin/env bash
# Modified from: https://www.digitalocean.com/community/tutorials/how-to-build-go-executables-for-multiple-platforms-on-ubuntu-16-04 

# This script requires a Golang file as a Command Line Argument (CLA)
# 	It then creates executables based on the platforms list
# 	or the platform provided as a second CLA

package=$1	# Get first command line argument
if [[ -z "$package" ]]; then	# Must provide Go file
	echo "usage: $0 <file-name>"	# (optional: <OS/Arch>)
	exit 1
fi
pack_split=(${package/./ })	# Splice file name and extension
package_name=${pack_split[0]}	# Get the name, without extension

platforms=$2	# Get the second command line argument
if [[ -z "$platforms" ]]; then	# Use default list:
	echo "No target platform provided, using default list."
	echo "usage: $0 <file-name> <OS/ARCH>"
	platforms=("windows/amd64" "windows/386" "darwin/amd64")
fi

# Loop through list
for platform in "${platforms[@]}"
do
	platform_split=(${platform//\// })	# Splice platform on slash
	OS=${platform_split[0]}		# First is the target OS
	ARCH=${platform_split[1]}	# Second is the computer architecture
	file_name=$package_name'-'$OS'-'$ARCH	# Create file name
	if [ $OS = "windows" ]; then	# Handle Windows
		file_name+='.exe'
	fi

	# Run build command
	env GOOS=$OS GOARCH=$ARCH go build -o $file_name $package
	if [ $? -ne 0 ]; then
		echo 'An error has occured! Aborting the script execution...'
		exit 1
	fi
	echo $file_name+' created successfully!'
done
