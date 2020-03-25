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
package events

import (
	"encoding/json"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_EPInvoke(t *testing.T) {
	Convey("EPInvokeRequest", t, func() {
		headers := map[string]string{
			"foo": "bar",
		}
		multiHeaders := map[string][]string{
			"hello": []string{"world", "world2"},
		}

		multiQueries := map[string][]string{
			"hello": []string{"world", "world2"},
		}
		queries := map[string]string{
			"foo": "bar",
		}

		request := EPInvokeRequest{
			Method:            "GET",
			RemoteAddr:        "123.456.789",
			Host:              "host1",
			Path:              "/shortname/redirect",
			Headers:           headers,
			MultiValueHeaders: multiHeaders,
			Queries:           queries,
			MultiValueQueries: multiQueries,
			Body:              "{\"property\": \"one\"}",
			IsBase64Encoded:   false,
		}

		requestJson := `
{
    "httpMethod":"GET",
    "remoteAddress":"123.456.789",
    "host":"host1",
    "path":"/shortname/redirect",
    "headers": {
        "foo": "bar"
    },
    "multiValueHeaders":{
        "hello": ["world", "world2"]
    },
    "queries": {
        "foo": "bar"
    },
    "multiValueQueries":{
        "hello": ["world", "world2"]
    },
    "body":"{\"property\": \"one\"}",
    "isBase64Encoded": false
}`

		Convey("Marshal & Unmarshal", func() {
			// Test to make sure we can unmarshal a json string payload to what we expect
			unmarshalled := &EPInvokeRequest{}
			err := json.Unmarshal([]byte(requestJson), unmarshalled)
			So(err, ShouldBeNil)

			So(unmarshalled.Method, ShouldEqual, request.Method)
			So(unmarshalled.RemoteAddr, ShouldEqual, request.RemoteAddr)
			So(unmarshalled.Host, ShouldEqual, request.Host)
			So(unmarshalled.Path, ShouldEqual, request.Path)
			So(unmarshalled.Body, ShouldEqual, `{"property": "one"}`)
			So(unmarshalled.IsBase64Encoded, ShouldBeFalse)

			So(unmarshalled.MultiValueHeaders["hello"], ShouldResemble, request.MultiValueHeaders["hello"])
			So(unmarshalled.Headers["foo"], ShouldEqual, request.Headers["foo"])

			So(unmarshalled.MultiValueQueries["hello"], ShouldResemble, request.MultiValueQueries["hello"])
			So(unmarshalled.Queries["foo"], ShouldEqual, request.Queries["foo"])

			// Test to make sure we can marshal to the same json string
			backToJson, err := json.Marshal(unmarshalled)
			So(err, ShouldBeNil)

			payload, err := json.Marshal(request)
			So(err, ShouldBeNil)
			So(string(payload), ShouldEqual, string(backToJson))
		})
	})

	Convey("EPInvokeResponse", t, func() {
		headers := map[string]string{
			"foo": "bar",
		}
		multiHeaders := map[string][]string{
			"hello": []string{"world", "world2"},
		}

		response := EPInvokeResponse{
			StatusCode:         207,
			Headers:            headers,
			MultiValueHeaders:  multiHeaders,
			Body:               `{"property": "one"}`,
			ShouldBase64Decode: false,
		}

		responseJson := `
{
    "statusCode": 207,
    "headers": {
        "foo": "bar"
    },
    "multiValueHeaders":{
        "hello": ["world", "world2"]
    },
    "body":"{\"property\": \"one\"}",
    "shouldBase64Decode": false
}`

		Convey("Marshal & Unmarshal", func() {
			// Test to make sure we can unmarshal a raw json string to what we expect
			unmarshalled := &EPInvokeResponse{}
			err := json.Unmarshal([]byte(responseJson), unmarshalled)
			So(err, ShouldBeNil)

			So(unmarshalled.StatusCode, ShouldEqual, response.StatusCode)
			So(unmarshalled.Body, ShouldEqual, response.Body)
			So(unmarshalled.ShouldBase64Decode, ShouldBeFalse)

			So(unmarshalled.MultiValueHeaders["hello"], ShouldResemble, response.MultiValueHeaders["hello"])
			So(unmarshalled.Headers["foo"], ShouldEqual, response.Headers["foo"])

			// Test to make sure we can marshal to the same json string
			backToJson, err := json.Marshal(unmarshalled)
			So(err, ShouldBeNil)

			payload, err := json.Marshal(response)
			So(err, ShouldBeNil)
			So(string(payload), ShouldEqual, string(backToJson))
		})
	})
}
