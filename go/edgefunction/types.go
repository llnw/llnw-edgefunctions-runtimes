// Copyright 2020 Limelight Networks, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
package edgefunction

import (
	"context"
	"time"
)

type invokeRequestWrapper struct {
	Payload []byte
	Context invokeContext
}

type invokeContext struct {
	RequestId       string
	CollectStdio    bool
	FunctionName    string
	FunctionVersion string
	Qualifier       string // request qualifier
	MemoryLimit     int32  // function memory limit
	Deadline        int64  // epoch time at which function will timeout
}

func (ic *invokeContext) ToContext() (context.Context, context.CancelFunc) {
	ctx := context.Background()
	ctx = context.WithValue(ctx, "RequestId", ic.RequestId)
	ctx = context.WithValue(ctx, "CollectStdio", ic.CollectStdio)
	ctx = context.WithValue(ctx, "FunctionName", ic.FunctionName)
	ctx = context.WithValue(ctx, "FunctionVersion", ic.FunctionVersion)
	ctx = context.WithValue(ctx, "Qualifier", ic.Qualifier)
	ctx = context.WithValue(ctx, "MemoryLimit", ic.MemoryLimit)

	deadlinens := ic.Deadline * int64(time.Millisecond)
	return context.WithDeadline(ctx, time.Unix(0, deadlinens))
}

type invokeResponseWrapper struct {
	HandledError bool
	Payload      []byte
}

// errorResults represents data sent back on an error
type errorResult struct {
	ErrorType    string `json:"error"`
	ErrorMessage string `json:"message"`
	StackTrace   []byte `json:"stackTrace,omitempty"`
}
