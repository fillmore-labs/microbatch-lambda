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

import pb "github.com/fillmore-labs/microbatch-lambda/api/proto/v1alpha1"

type remoteError struct {
	msg string
}

func (r *remoteError) Error() string {
	return r.msg
}

var errMissingResult = &remoteError{"missing result"}

func unwrap(result *pb.JobResult, err error) (string, error) {
	if err != nil {
		return "", err
	}

	switch r := result.GetResult().(type) {
	case *pb.JobResult_Body:
		return r.Body, nil

	case *pb.JobResult_Error:
		return "", &remoteError{r.Error}

	default:
		return "", errMissingResult
	}
}
