package main

import (
	"encoding/json"
	"fmt"
	"log"
	"os/exec"
	"slices"
)

type Output struct {
	Name   string
	Active bool
}

func getOutputs() ([]Output, error) {
	cmd := exec.Command("swaymsg", "-t", "get_outputs")
	stdout, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	var outputs []Output

	err = json.Unmarshal(stdout, &outputs)
	if err != nil {
		return nil, err
	}

	return outputs, nil
}

func switchOn(out Output) error {
	cmd := exec.Command("swaymsg", "output", out.Name, "enable", "pos", "0", "0", "scale", "1")
	return cmd.Run()
}

func switchOff(out Output) error {
	cmd := exec.Command("swaymsg", "output", out.Name, "disable")
	return cmd.Run()
}

func main() {
	outputs, err := getOutputs()
	if err != nil {
		log.Fatalf("could not get outputs using swaymsg: %s", err)
	}

	activeIndex := slices.IndexFunc(outputs, func(out Output) bool { return out.Active })

	var newActiveIndex int
	if activeIndex == -1 {
		newActiveIndex = 0
	} else {
		newActiveIndex = (activeIndex + 1) % len(outputs)
	}

	for idx, out := range outputs {
		if idx == newActiveIndex {
			switchOn(out)
			fmt.Printf(" * %s\n", out.Name)
		} else {
			switchOff(out)
			fmt.Printf("   %s\n", out.Name)
		}
	}
}
