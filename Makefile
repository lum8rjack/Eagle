NAME := eagle
BUILD := go build -ldflags "-s -w" -trimpath

default: linux

clean:
	rm -f $(NAME)*.bin
	rm -f $(NAME)*.exe

linux:
	echo "Compiling for Linux x64"
	GOOS=linux GOARCH=amd64 $(BUILD) -o $(NAME)-Linux64.bin

windows:
	echo "Compiling for Windows x64"
	GOOS=windows GOARCH=amd64 $(BUILD) -o $(NAME)-Windows64.exe

mac:
	echo "Compiling for Mac x64"
	GOOS=darwin GOARCH=amd64 $(BUILD) -o $(NAME)-Darwin64.bin

m1:
	echo "Compiling for Mac M1"
	GOOS=darwin GOARCH=arm64 $(BUILD) -o $(NAME)-M1.bin

arm:
	echo "Compiling for Linux Arm64"
	GOOS=linux GOARCH=arm64 $(BUILD) -o $(NAME)-LinuxArm64.bin
