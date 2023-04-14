package cmd

import "os/exec"

func textToSpeech(text string) error {
	cmd := exec.Command("espeak", "-p", "10", "-s", "100", text)
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
