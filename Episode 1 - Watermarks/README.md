# Episode 1 - Watermarks

Video Link: [https://www.youtube.com/watch?v=TKq6ypmcnb8](https://www.youtube.com/watch?v=TKq6ypmcnb8)

### About

This application demonstrates the steganography technique of watermarking. It can embed a watermark into a base image with the requested bit-depth, which can later be extracted. If a small enough bit depth (around 8 bits or less seems to work well), then the watermark will be undetectable in the base image by the human eye.

### Using the program

This application is written in Go. Use `go build -o <program>` to build. Use `./<program> -h` to get a list and description of parameters.
