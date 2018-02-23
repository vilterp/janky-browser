jankybrowser:
	go build -o jankybrowser ./main.go

deps:
	go get github.com/faiface/pixel
	go get github.com/faiface/glhf
	go get github.com/golang/freetype/truetype
	go get github.com/go-gl/glfw/v3.2/glfw

.PHONY deps jankybrowser
