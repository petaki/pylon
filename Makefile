VERSION := $(if $(RELEASE_VERSION),$(RELEASE_VERSION),"master")

all: pre_clean darwin darwin_arm64 linux linux_arm64 windows

pre_clean:
	rm -rf dist
	mkdir dist
	sed -i 's/Version:\s*"master"/Version: $(VERSION)/g' main.go

darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/pylon .
	cd dist && zip pylon_$(VERSION)_darwin_amd64.zip pylon
	rm -f dist/pylon

darwin_arm64:
	GOOS=darwin GOARCH=arm64 go build -o dist/pylon .
	cd dist && zip pylon_$(VERSION)_darwin_arm64.zip pylon
	rm -f dist/pylon

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/pylon .
	cd dist && zip pylon_$(VERSION)_linux_amd64.zip pylon
	rm -f dist/pylon

linux_arm64:
	GOOS=linux GOARCH=arm64 go build -o dist/pylon .
	cd dist && zip pylon_$(VERSION)_linux_arm64.zip pylon
	rm -f dist/pylon

windows:
	GOOS=windows GOARCH=amd64 go build -o dist/pylon.exe .
	cd dist && zip pylon_$(VERSION)_windows_amd64.zip pylon.exe
	rm -f dist/pylon.exe
