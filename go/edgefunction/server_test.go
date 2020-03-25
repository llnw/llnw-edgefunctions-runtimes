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
	"context"
	"encoding/json"
	"errors"
	"io/ioutil"
	"os"
	"runtime/debug"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
	"github.com/vmihailenco/msgpack"
)

func TestExecutionError(t *testing.T) {
	Convey("formatting a response with error, no stacktrace", t, func() {
		passedErr := errors.New("something bad")
		data := executionError(passedErr, nil)

		var errorPayload errorResult
		err := json.Unmarshal(data, &errorPayload)
		So(err, ShouldBeNil)

		Convey(`ErrorType is "errorString"`, func() {
			So(errorPayload.ErrorType, ShouldEqual, "errorString")
		})
		Convey(`ErrorMessage is error message`, func() {
			So(errorPayload.ErrorMessage, ShouldEqual, passedErr.Error())
		})
		Convey(`StackTrace is nil`, func() {
			So(errorPayload.StackTrace, ShouldBeNil)
		})
	})
	Convey("formatting a response with error and stacktrace", t, func() {
		passedErr := errors.New("something bad")
		stackTrace := debug.Stack()
		data := executionError(passedErr, stackTrace)

		var errorPayload errorResult
		err := json.Unmarshal(data, &errorPayload)
		So(err, ShouldBeNil)

		Convey(`ErrorType is "panic"`, func() {
			So(errorPayload.ErrorType, ShouldEqual, "panic")
		})
		Convey(`ErrorMessage is error message`, func() {
			So(errorPayload.ErrorMessage, ShouldEqual, passedErr.Error())
		})
		Convey(`StackTrace is nil`, func() {
			So(errorPayload.StackTrace, ShouldResemble, stackTrace)
		})
	})
}

func TestListen(t *testing.T) {
	Convey("process a request", t, func() {
		req := new(invokeRequestWrapper)
		req.Context = invokeContext{}
		req.Payload = []byte{}

		var handler Handler

		// Open a temp file
		testFileHandle, err := ioutil.TempFile("", "faas-go-sdk-test")
		So(err, ShouldBeNil)
		defer os.Remove(testFileHandle.Name())

		Convey("run handler with no errors", func() {
			type resultStuff struct {
				Message string
			}

			// Format a request
			writeTestRequest(req, testFileHandle)
			message := "yay"
			handler = NewHandler(func(ctx context.Context) (*resultStuff, error) {
				return &resultStuff{Message: message}, nil
			})

			// Create a new file handle for process request
			otherHandle, err := os.OpenFile(testFileHandle.Name(), os.O_RDWR|os.O_APPEND, os.ModePerm)
			So(err, ShouldBeNil)

			processRequest(handler, otherHandle)

			// Check response
			resp := new(invokeResponseWrapper)
			readTestResponse(resp, testFileHandle)
			So(resp.HandledError, ShouldBeFalse)

			var resultOut resultStuff
			err = json.Unmarshal(resp.Payload, &resultOut)
			So(err, ShouldBeNil)
			So(resultOut.Message, ShouldEqual, message)

		})

		Convey("bad request", func() {
			writeTestRequest("nope", testFileHandle)

			// Create a new file handle for process request
			otherHandle, err := os.OpenFile(testFileHandle.Name(), os.O_RDWR|os.O_APPEND, os.ModePerm)
			So(err, ShouldBeNil)

			// It's ok we don't have a handler shouldn't get here
			processRequest(handler, otherHandle)

			// Check response
			Convey("should return error response", func() {
				resp := new(invokeResponseWrapper)
				readTestResponse(resp, testFileHandle)
				validateErrResponse(resp, "failed to unmarshal msgpack request", false)
			})
		})

		Convey("handler panics", func() {
			// Format a request
			writeTestRequest(req, testFileHandle)

			panicString := "ahhh"
			handler = NewHandler(func(ctx context.Context) error {
				panic(panicString)
			})

			// Create a new file handle for process request
			otherHandle, err := os.OpenFile(testFileHandle.Name(), os.O_RDWR|os.O_APPEND, os.ModePerm)
			So(err, ShouldBeNil)

			processRequest(handler, otherHandle)

			// Check response
			Convey("should return error response with panic message", func() {
				resp := new(invokeResponseWrapper)
				readTestResponse(resp, testFileHandle)
				validateErrResponse(resp, panicString, true)
			})
		})

		Convey("handler error", func() {
			// Format a request
			writeTestRequest(req, testFileHandle)

			errorString := "ahhh"
			handler = NewHandler(func(ctx context.Context) error {
				return errors.New(errorString)
			})
			// Create a new file handle for process request
			otherHandle, err := os.OpenFile(testFileHandle.Name(), os.O_RDWR|os.O_APPEND, os.ModePerm)
			So(err, ShouldBeNil)

			processRequest(handler, otherHandle)

			// Check response
			Convey("should return error response with error message", func() {
				resp := new(invokeResponseWrapper)
				readTestResponse(resp, testFileHandle)
				validateErrResponse(resp, errorString, false)
			})
		})
	})
}

func writeTestRequest(req interface{}, file *os.File) {
	reqBytes, err := msgpack.Marshal(req)
	So(err, ShouldBeNil)

	_, err = file.Write(reqBytes)
	So(err, ShouldBeNil)
}

func readTestResponse(resp *invokeResponseWrapper, file *os.File) {
	decoder := msgpack.NewDecoder(bufio.NewReader(file))
	err := decoder.Decode(resp)
	So(err, ShouldBeNil)
}

func validateErrResponse(resp *invokeResponseWrapper, message string, hasStackTrace bool) {
	So(resp.HandledError, ShouldBeTrue)

	var errPayload errorResult
	err := json.Unmarshal(resp.Payload, &errPayload)
	So(err, ShouldBeNil)

	So(errPayload.ErrorMessage, ShouldContainSubstring, message)

	if hasStackTrace {
		So(errPayload.ErrorType, ShouldEqual, "panic")
		So(errPayload.StackTrace, ShouldNotBeEmpty)
	} else {
		So(errPayload.ErrorType, ShouldEqual, "errorString")
		So(errPayload.StackTrace, ShouldBeNil)
	}
}
