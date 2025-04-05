main: clean build

clean:
	rm -rf bin

build:
	go build -o bin/wallpicker main.go

install:
	cp bin/wallpicker ~/.local/bin/wallpicker
