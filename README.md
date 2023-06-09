# Alpha Assistant

Alpha Assistant is a Large Language Model and Whisper.cpp powered personal assistant

## Setup
```bash
git clone https://github.com/Sam-Fatehmanesh/Alpha-Assistant.git
cd Alpha-Assistant
git clone https://github.com/ggerganov/whisper.cpp.git
cd whisper.cpp/bindings/go
make whisper
cd ../..
bash ./models/download-ggml-model.sh tiny.en
cd ..
CGO_CFLAGS=-I/abspath/to/whisper.cpp CGO_LDFLAGS=-L/abspath/to/whisper.cpp go build -o alpha
```
## Client
The alpha client is the user interface for the assistant, it uses whisper.cpp for speech to text and espeak for text to speech. The current interface method involves a hot key Voice Activity detection hybrid mechanism where the hotkey starts recording and either another press of the hotkey or no speech detected for 2 seconds ends the recording. The hotkey is currently super + X, customizabiltiy will be improved on a later date.

### Usage
Before runing client make sure server is running first as client needs to call the api started by server.
```java
$ ./alpha client --help
Start the alpha client, a client which records and transcribes audio which is then sent to the alpha server.

Usage:
alpha client [flags]

Flags:
-h, --help                help for client
-m, --model-path string   Directory for whisper STT model (default "./whisper.cpp/models/ggml-tiny.en.bin")
-s, --server-url string   URL for the alpha AI server (default "http://127.0.0.1:22589")
-v, --vad                 Use voice activity detection (default true)
    --vad-mode int        The strength of the voice activity detection, from 0, most sensitive, to 3, least sensitive (default 1)
```

## Server
### In Progress

### Usage
```java
$ ./alpha server --help
Start the alpha personal assistant server.

Usage:
  alpha server [flags]

Flags:
  -h, --help   help for server
```