APP=fb_crm_audience

build:
	rm -rf electron/dist
	rm -rf electron/binary
	cd frontend && yarn build
	CGO_ENABLED=1 CC=x86_64-w64-mingw32-gcc CXX=x86_64-w64-mingw32-g++ GOOS=windows GOARCH=amd64 go build -x -v -ldflags "-s -w" -o electron/binary/fb_crm_audience.exe
	cd electron && yarn package
	rm -rf electron/binary/fb_crm_audience.exe
	CGO_ENABLED=1 GOOS=darwin GOARCH=amd64 go build -x -v -ldflags "-s -w" -o electron/binary/fb_crm_audience
	cd electron && yarn package-mac
