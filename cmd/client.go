/*
Copyright © 2023 NAME HERE <EMAIL ADDRESS>
*/
package cmd

import (
	"fmt"
	"time"

	"bytes"
	"io"
	"io/ioutil"
	"net/http"

	"github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
	"github.com/gordonklaus/portaudio"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	"golang.design/x/hotkey"
)

// startClientCmd represents the startClient command
var startClientCmd = &cobra.Command{
	Use:   "client",
	Short: "start the alpha client",
	Long:  `Start the alpha client, a client which records and transcribes audio which is then sent to the alpha server.`,
	Run: func(cmd *cobra.Command, args []string) {
		go liseningForElaboration()

		use_vad, err := cmd.Flags().GetBool("vad")
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}

		vad_strength, err := cmd.Flags().GetInt("vad-mode")
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}

		model_path, err := cmd.Flags().GetString("model-path")
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}

		serverURL, err := cmd.Flags().GetString("server-url")

		err = hkClient(model_path, serverURL, use_vad, vad_strength)
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
		}

	},
}

func init() {
	rootCmd.AddCommand(startClientCmd)

	startClientCmd.Flags().StringP("model-path", "m", "./whisper.cpp/models/ggml-tiny.en.bin", "Directory for whisper STT model")
	startClientCmd.Flags().BoolP("vad", "v", true, "Use voice activity detection")
	startClientCmd.Flags().Int("vad-mode", 1, "The strength of the voice activity detection, from 0, most sensitive, to 3, least sensitive")
	startClientCmd.Flags().StringP("server-url", "s", "http://127.0.0.1:22589", "URL for the alpha AI server")
}

// hkClient is a function that initializes a HotKey, records audio input using PortAudio,
// transcribes the recorded audio using a specified model, sends the transcription to an AI server,
// and repeats the process until an error occurs or the function is interrupted.
func hkClient(modelpath string, serverURL string, vad bool, vad_strength int) error {
	// 0x40, 0x50 = super key
	// Initialize the HotKey. In this case, the hotkey is the letter "B" pressed with the "super" key.
	hk := hotkey.New([]hotkey.Modifier{hotkey.Modifier(0x50)}, hotkey.KeyB)
	// Register the HotKey to listen for key presses.
	if err := hk.Register(); err != nil {
		return err
	}
	defer hk.Unregister()

	// Initialize PortAudio for audio input
	if err := portaudio.Initialize(); err != nil {
		return err
	}
	defer portaudio.Terminate()

	// Create a buffer to hold the recorded audio data
	var buffer []float32
	// Open a default audio stream with 1 input channel, no output channels, 16kHz sample rate,
	// and a callback function to append the incoming audio data to the buffer
	stream, err := portaudio.OpenDefaultStream(1, 0, 16000, 0, func(in []float32) {
		buffer = append(buffer, in...)
	})
	if err != nil {
		return err
	}
	defer stream.Close()

	// Load the transcription model from the specified path
	model, err := whisper.New(modelpath)
	if err != nil {
		return err
	}
	defer model.Close()

	// Initialize VAD
	vadQuiet, err := VADSignal(buffer, vad_strength)
	if err != nil {
		return err
	}

	for {
		// Wait for the HotKey to be pressed to start recording
		<-hk.Keydown()

		fmt.Println("Start Recording")
		err = stream.Start()
		if err != nil {
			return err
		}

		// Wait for a short time
		time.Sleep(time.Millisecond * 100)

		// Record the starting index of the new audio data in the buffer
		start_index := len(buffer)

		// Check if Voice Activity Detection (VAD) is enabled
		if vad {
			// Wait for either the VAD to detect a quiet period or for the HotKey to be pressed again to end recording
			select {
			case <-vadQuiet:
			case <-hk.Keydown():
			}
		} else {
			// If VAD is not enabled, simply wait for the HotKey to be pressed again to end recording
			<-hk.Keydown()
		}

		fmt.Println("Stop Recording")
		err = stream.Stop()
		if err != nil {
			return err
		}

		// Transcribe the recorded audio data starting from the previously recorded starting index
		text, err := transcribe(buffer[start_index:], model)
		if err != nil {
			return err
		}

		// Reset the audio buffer to prepare for the next recording
		buffer = buffer[:0]

		// Print the transcription to the console for debugging purposes
		fmt.Println("transcription: " + text)

		err = textToSpeech("Understood")
		if err != nil {
			return err
		}

		// Send the transcription to the specified AI server
		err = sendToAI(text, serverURL)
		if err != nil {
			return err
		}
	}
}

func sendToAI(text string, server_url string) error {

	req, err := http.NewRequest("GET", server_url, nil)
	if err != nil {
		return err
	}

	req.Body = io.NopCloser(bytes.NewReader([]byte(text)))

	// Create a new HTTP client
	client := &http.Client{}

	// Send the request and get the response
	resp, err := client.Do(req)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	// Print the response body
	output := string(body)

	fmt.Println("output: " + output)

	err = textToSpeech(output)
	if err != nil {
		return err
	}

	return nil
}

func liseningForElaboration() error {
	// Define the HTTP endpoint that will receive requests
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		// Set the content type header to text/plain
		w.Header().Set("Content-Type", "text/plain")

		// Get the text from the request body
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, "Error reading request body", http.StatusBadRequest)
			return
		}

		receivedText := string(body)

		err = textToSpeech(receivedText)
		if err != nil {
			log.Fatal().Err(err).Msg(err.Error())
			return
		}
	})

	// Start the HTTP server and listen for incoming requests
	if err := http.ListenAndServe(":22588", nil); err != nil {
		log.Fatal().Err(err).Msg(err.Error())
		return err
	}

	return nil
}
