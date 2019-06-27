package main

import (
	"log"
	"os"
	"os/exec"
	"fmt"
	"strings"
)


func main() {
	log.Println("Counter")

	output, err := exec.Command("sudo", "./dnpds40", "-n").CombinedOutput()
	if err != nil {
	os.Stderr.WriteString(err.Error())
	}

	m := strings.Split(string(output), "\n")
	g := strings.Split(m[12], ":")

	fmt.Println(strings.TrimLeft(g[2], " "))
}

