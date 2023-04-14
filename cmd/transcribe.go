package cmd

import "github.com/ggerganov/whisper.cpp/bindings/go/pkg/whisper"
import "github.com/rs/zerolog/log"

func transcribe(samples []float32, model whisper.Model) (string, error) {
	// Process samples
	model_context, err := model.NewContext()
	if err != nil {
		log.Fatal().Err(err)
	}
	model_context.SetLanguage("en")
	model_context.SetThreads(8)

	if err := model_context.Process(samples, nil); err != nil {
		log.Fatal().Err(err)
	}

	// Print out the results
	text := ""
	for {
		segment, err := model_context.NextSegment()
		if err != nil {
			break
		}
		text += segment.Text
	}

	return text, nil
}
