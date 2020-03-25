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
	"bufio"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"
	"runtime/debug"

	"github.com/vmihailenco/msgpack"
)

// Start passes the handler all EdgeFunction requests for processing.
func Start(handler interface{}) {
	// fd3 is expected to be available
	file := os.NewFile(uintptr(3), "fd3")
	defer file.Close()

	handlerWrapper := NewHandler(handler)

	for {
		processRequest(handlerWrapper, file)
	}
}

func processRequest(handler Handler, file *os.File) {
	defer func() {
		if r := recover(); r != nil {
			stackTrace := debug.Stack()
			var err error
			switch x := r.(type) {
			case string:
				err = errors.New(x)
			case error:
				err = x
			default:
				// Fallback err (per specs, error strings should be lowercase w/o punctuation)
				err = errors.New("unknown panic")
			}

			writeErrorResponse(err, stackTrace, file)
		}
	}()

	// Attempt to read fd3, this will fail if it doesn't exist
	req := &invokeRequestWrapper{}
	decoder := msgpack.NewDecoder(bufio.NewReader(file))
	err := decoder.Decode(req)

	// unmarshal the msgpack message
	if err != nil {
		err = fmt.Errorf("failed to unmarshal msgpack request: %s", err.Error())
		writeErrorResponse(err, nil, file)
		return
	}

	// Execute the function
	respWrapper := new(invokeResponseWrapper)
	handlerCtx, cancel := req.Context.ToContext()
	defer cancel()

	respWrapper.Payload, err = handler.Invoke(handlerCtx, req.Payload)
	if err != nil {
		err = fmt.Errorf("failed to execute handler: %s", err.Error())
		writeErrorResponse(err, nil, file)
		return
	}

	// send the response
	writeResponse(respWrapper, file)
}

func writeErrorResponse(err error, stackTrace []byte, file *os.File) {
	resp := new(invokeResponseWrapper)
	resp.HandledError = true
	resp.Payload = executionError(err, stackTrace)
	writeResponse(resp, file)
}

func writeResponse(resp *invokeResponseWrapper, file *os.File) {
	// msgpack the response
	respBytes, err := msgpack.Marshal(resp)
	if err != nil {
		log.Println("Failed to marshal msgpack response:", err.Error())
		return
	}

	// send the response
	_, err = file.Write(respBytes)
	if err != nil {
		log.Println("Failed to write response:", err.Error())
		return
	}
}

func executionError(err error, stackTrace []byte) []byte {
	errPayload := &errorResult{
		ErrorType:    "errorString",
		ErrorMessage: err.Error(),
		StackTrace:   stackTrace,
	}

	// Should only get stack trace in a panic
	if len(stackTrace) > 0 {
		errPayload.ErrorType = "panic"
	}

	data, err := json.Marshal(errPayload)
	if err != nil {
		log.Println("Failed to marshal error payload:", err.Error())
		return []byte{}
	}

	return data
}
