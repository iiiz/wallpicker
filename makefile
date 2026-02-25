main: clean build

clean:
	rm -rf bin

build:
	go build -o bin/wallpicker main.go

install:
	sudo cp ./bin/wallpicker /usr/bin/wallpicker
