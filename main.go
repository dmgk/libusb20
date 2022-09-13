package main

import (
	"fmt"
	"os"
)

func main() {
	backend, err := NewBackend()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	defer backend.Free()

	fmt.Printf("====> backend %#v\n", backend)

	for backend.Scan() {
		device := backend.Device()
		fmt.Printf("====> backend: %s, device: %s\n", device.BackendName(), device)
	}
}
