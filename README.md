# Alpha Assisant

Alpha Assisant is a Large Language Model and Whisper.cpp powered personal assistant

## Setup
```bash
    git clone https://github.com/Sam-Fatehmanesh/Alpha-Assistant.git
    cd Alpha-Assistant
    git clone https://github.com/ggerganov/whisper.cpp.git
    cd whisper.cpp/bindings/go
    make whisper
    cd ../../..
    bash ./models/download-ggml-model.sh tiny.en
    CGO_CFLAGS=-I/path/to/whisper.cpp CGO_LDFLAGS=-L/path/to/whisper.cpp go build -o alpha
```