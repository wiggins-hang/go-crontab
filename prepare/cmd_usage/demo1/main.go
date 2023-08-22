package main

import (
	"fmt"
	"os/exec"
)

func main() {

	// cmd = exec.Command("/bin/bash", "-c", "echo 1;echo2;")

	cmd := exec.Command("bash", "-c", "echo 1")

	a, err := cmd.CombinedOutput()

	fmt.Println(string(a), err)
}
