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
	"errors"
)

// EdgeFunctionContext is an easy to use extraction of a context
// passed into the EdgeFunction.
type EdgeFunctionContext struct {
	RequestID       string
	CollectStdio    bool
	FunctionName    string
	FunctionVersion string
	Qualifier       string
	MemoryLimitMB   int32
}

// FromContext parses out key fields from the passed in ctx
// and creates an easy to use EdgeFunctionContext.
// Returns an error if the conversion failed due to a missing value.
func FromContext(ctx context.Context) (*EdgeFunctionContext, error) {
	efContext := new(EdgeFunctionContext)

	ok := true
	efContext.RequestID, ok = ctx.Value("RequestId").(string)
	if !ok {
		return nil, errors.New("Context did not contain the key RequestId")
	}

	efContext.CollectStdio, ok = ctx.Value("CollectStdio").(bool)
	if !ok {
		return nil, errors.New("Context did not contain the key CollectStdio")
	}

	efContext.FunctionName, ok = ctx.Value("FunctionName").(string)
	if !ok {
		return nil, errors.New("Context did not contain the key FunctionName")
	}

	efContext.FunctionVersion, ok = ctx.Value("FunctionVersion").(string)
	if !ok {
		return nil, errors.New("Context did not contain the key FunctionVersion")
	}

	efContext.Qualifier, ok = ctx.Value("Qualifier").(string)
	if !ok {
		return nil, errors.New("Context did not contain the key Qualifier")
	}

	efContext.MemoryLimitMB, ok = ctx.Value("MemoryLimit").(int32)
	if !ok {
		return nil, errors.New("Context did not contain the key MemoryLimit")
	}

	return efContext, nil
}
