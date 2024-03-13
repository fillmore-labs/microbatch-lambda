// Copyright 2023-2024 Oliver Eikemeier. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"fmt"
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
	fmt.Printf("Synthesized stacks %s\n", strings.Join(stackNames, ", "))
}
