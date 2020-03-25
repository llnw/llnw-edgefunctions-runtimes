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

// EPInvokeRequest is used as the payload structure for EP Invokes
type EPInvokeRequest struct {
	Method            string              `json:"httpMethod"`
	RemoteAddr        string              `json:"remoteAddress"`
	Host              string              `json:"host"`
	Path              string              `json:"path"`
	Headers           map[string]string   `json:"headers"`
	MultiValueHeaders map[string][]string `json:"multiValueHeaders"`
	Queries           map[string]string   `json:"queries"`
	MultiValueQueries map[string][]string `json:"multiValueQueries"`
	Body              string              `json:"body"`
	IsBase64Encoded   bool                `json:"isBase64Encoded"`
}

// EPInvokeResponse is used as the payload structure for EP Invoke responses
type EPInvokeResponse struct {
	StatusCode         int                 `json:"statusCode"`
	Headers            map[string]string   `json:"headers"`
	MultiValueHeaders  map[string][]string `json:"multiValueHeaders"`
	Body               string              `json:"body"`
	ShouldBase64Decode bool                `json:"shouldBase64Decode"`
}
