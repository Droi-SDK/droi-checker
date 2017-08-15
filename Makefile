OUTPUT=./_build
SRC=$(shell find . -iname "*.go")
LDFLAGS='-X main.pkgType=binary -s -w'
RESOURCES=$(wildcard ./console/resources/*.html)

all: binaries

binaries: $(SRC)
	GOOS=darwin GOARCH=amd64 POSTFIX= make $(OUTPUT)/droi-checker-darwin-amd64
	GOOS=windows GOARCH=386 POSTFIX=.exe make $(OUTPUT)/droi-checker-windows-386.exe
	GOOS=windows GOARCH=amd64 POSTFIX=.exe make $(OUTPUT)/droi-checker-windows-amd64.exe
	GOOS=linux GOARCH=amd64 POSTFIX= make $(OUTPUT)/droi-checker-linux-amd64
	GOOS=linux GOARCH=386 POSTFIX= make $(OUTPUT)/droi-checker-linux-386

msi:
	wixl -a x86 packaging/msi/droi-checker-x86.wxs -o $(OUTPUT)/droi-checker-setup-x86.msi
	wixl -a x64 packaging/msi/droi-checker-x64.wxs -o $(OUTPUT)/droi-checker-setup-x64.msi

deb:
	mkdir -p $(OUTPUT)/x86-deb/DEBIAN/
	mkdir -p $(OUTPUT)/x86-deb/usr/bin/
	cp $(OUTPUT)/droi-checker-linux-386 $(OUTPUT)/x86-deb/usr/bin/droi-checker
	cp packaging/deb/control-x86 $(OUTPUT)/x86-deb/DEBIAN/control
	dpkg-deb --build $(OUTPUT)/x86-deb $(OUTPUT)/droi-checker-x86.deb
	mkdir -p $(OUTPUT)/x64-deb/DEBIAN/
	mkdir -p $(OUTPUT)/x64-deb/usr/bin/
	cp $(OUTPUT)/droi-checker-linux-amd64 $(OUTPUT)/x64-deb/usr/bin/droi-checker
	cp packaging/deb/control-x64 $(OUTPUT)/x64-deb/DEBIAN/control
	dpkg-deb --build $(OUTPUT)/x64-deb $(OUTPUT)/droi-checker-x64.deb

$(OUTPUT)/droi-checker-$(GOOS)-$(GOARCH)$(POSTFIX): $(SRC)
	GOOS=$(GOOS) GOARCH=$(GOARCH) go build -o $@ -ldflags=$(LDFLAGS) github.com/Droi-SDK/droi-checker

install: resources
	GOOS=$(GOOS) go install github.com/Droi-SDK/droi-checker

resources:
	(cd console; $(MAKE))

cdroi-checker:
	rm -rf $(OUTPUT)

.PHONY: test msi deb install clean resources