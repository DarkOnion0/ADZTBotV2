#!/usr/bin/env bash

#######################################
## Global variables and string style ##
#######################################

bold=$(tput bold)

red="\033[0;31m"
stop_color='\033[0m'

##########################################
## Set the execution mode of the script ##
##########################################

if [ -z $1 ]; then
	VERSION="latest"
	echo "Running in manual mode, version=$VERSION"
else
	VERSION=$(echo $1 | sed -e 's/\./-/g')
	echo "Running in release mode, version=$VERSION"
fi

###############################################
## Cross-compile ADZTBotV2 && ZIP the output ##
###############################################

mkdir "bin"
cd ./bin/

for os in linux windows darwin
do
	# Linux is seperated from the other os due to the fact that it support more architechture
	if [ $os == linux ]; then
		echo -e "\n$red${bold}Building linux binary...${bold}$stop_color"
		echo -e "$red${bold}===========================${bold}$stop_color"

		for arch in amd64 386 arm64 arm
		do
			echo "${bold}$os/$arch...${bold}"
			env GOOS=$os GOARCH=$arch go build -ldflags="-X 'github.com/DarkOnion0/ADZTBotV2/config.RawVersion=$VERSION'"$VERSION'" -o adztbotv2_$os-$arch-$VERSION ./../main.go
			sha256sum adztbotv2_$os-$arch-$VERSION > adztbotv2_$os-$arch-$VERSION-sha256sum.txt
			zip adztbotv2_$os-$arch-$VERSION adztbotv2_$os-$arch-$VERSION adztbotv2_$os-$arch-$VERSION-sha256sum.txt
		done
	elif [ $os == darwin ]; then
		echo -e "\n$red${bold}Building $os binary...${bold}$stop_color"
		echo -e "$red${bold}===========================${bold}$stop_color"

		arch=amd64
		echo "${bold}$os/$arch...${bold}"
		env GOOS=$os GOARCH=$arch go build -ldflags="-X 'github.com/DarkOnion0/ADZTBotV2/config.RawVersion=$VERSION'"$VERSION'" -o adztbotv2_$os-$arch-$VERSION ./../main.go
		sha256sum adztbotv2_$os-$arch-$VERSION > adztbotv2_$os-$arch-$VERSION-sha256sum.txt
		zip adztbotv2_$os-$arch-$VERSION adztbotv2_$os-$arch-$VERSION adztbotv2_$os-$arch-$VERSION-sha256sum.txt
	else
		echo -e "\n$red${bold}Building $os binary...${bold}$stop_color"
		echo -e "$red${bold}===========================${bold}$stop_color"
		
		for arch in amd64 386
		do 
			echo "${bold}$os/$arch...${bold}"
			env GOOS=$os GOARCH=$arch go build -ldflags="-X 'github.com/DarkOnion0/ADZTBotV2/config.RawVersion=$VERSION'"$VERSION'" -o adztbotv2_$os-$arch-$VERSION ./../main.go
			sha256sum adztbotv2_$os-$arch-$VERSION > adztbotv2_$os-$arch-$VERSION-sha256sum.txt
			zip adztbotv2_$os-$arch-$VERSION adztbotv2_$os-$arch-$VERSION adztbotv2_$os-$arch-$VERSION-sha256sum.txt
		done
	fi
done
