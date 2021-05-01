VERSION := $(if $(RELEASE_VERSION),$(RELEASE_VERSION),"master")

all: pre_clean darwin linux windows

pre_clean:
	rm -rf dist

darwin:
	GOOS=darwin GOARCH=amd64 go build -o dist/pylon .
	cd dist && zip pylon_$(VERSION)_darwin_amd64.zip pylon
	rm -f dist/pylon

linux:
	GOOS=linux GOARCH=amd64 go build -o dist/pylon .
	cd dist && zip pylon_$(VERSION)_linux_amd64.zip pylon
	rm -f dist/pylon

windows:
	GOOS=windows GOARCH=amd64 go build -o dist/pylon.exe .
	cd dist && zip pylon_$(VERSION)_windows_amd64.zip pylon.exe
	rm -f dist/pylon.exe
