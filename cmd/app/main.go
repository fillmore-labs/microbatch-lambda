package main

import (
	"log"
	"strings"

	"github.com/aws/jsii-runtime-go"
)

func main() {
	defer jsii.Close()

	assembly := NewAssembly()

	stacks := *assembly.Stacks()
	stackNames := make([]string, 0, len(stacks))
	for _, stack := range stacks {
		stackNames = append(stackNames, *stack.DisplayName())
	}
	log.Printf("Synthesized stacks %s", strings.Join(stackNames, ", "))
}
