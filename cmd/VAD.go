package cmd

import (
	"fmt"
	"github.com/maxhawkins/go-webrtcvad"
	"github.com/rs/zerolog/log"
	"time"
)

func VADSignal(buffer []float32, vad_strength int) (<-chan struct{}, error) {

	// Initialize VAD
	vad, err := webrtcvad.New()
	if err != nil {
		return nil, err
	}

	if err := vad.SetMode(vad_strength); err != nil {
		return nil, err
	}

	rate := 16000        // Hz
	frame_len_time := 30 // in milliseconds
	frame := make([]byte, (rate*frame_len_time)/1000)
	frame_len := len(frame)

	if ok := vad.ValidRateAndFrameLength(rate, (rate*frame_len_time)/1000); !ok {
		return nil, fmt.Errorf("invalid rate or frame length")
	}

	// Create a channel for signaling when speech is detected
	signalCh := make(chan struct{})

	active_buffer := make([]bool, 0)

	go func() {
		for {
			for len(buffer) < (2*frame_len)+1 {
				time.Sleep(time.Millisecond * time.Duration(frame_len_time*2))
			}

			time.Sleep(time.Millisecond * time.Duration(frame_len_time*2))

			// Convert buffer to bytes for VAD processing
			frame, err = float32ToInt8Bytes(buffer[len(buffer)-(2*frame_len):])
			if err != nil {
				log.Printf("error converting buffer to bytes: %v", err)
				continue
			}

			// Run VAD on frame
			active, err := vad.Process(rate, frame)
			if err != nil {
				log.Printf("error running VAD: %v", err)
				continue
			}
			active_buffer = append(active_buffer, active)

			// Truncate buffer if it exceeds 30 seconds
			if len(buffer) > 30*rate {
				log.Print("Truncated Buffer")
				buffer = buffer[30*rate:]
			}

			// Check if speech has stopped for more than 2 seconds
			if lastNFalse(active_buffer, int(2000/frame_len_time)) {
				signalCh <- struct{}{}
			}
		}
	}()

	return signalCh, nil
}
