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
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_FromContext(t *testing.T) {
	Convey("FromContext", t, func() {
		requestId := "1234"
		collectStdio := false
		functionName := "myFunc"
		functionVersion := "v1"
		qualifier := "alias"
		memoryLimitMB := int32(256)

		ctx := context.Background()

		Convey("with valid context", func() {
			ctx = context.WithValue(ctx, "RequestId", requestId)
			ctx = context.WithValue(ctx, "CollectStdio", collectStdio)
			ctx = context.WithValue(ctx, "FunctionName", functionName)
			ctx = context.WithValue(ctx, "FunctionVersion", functionVersion)
			ctx = context.WithValue(ctx, "Qualifier", qualifier)
			ctx = context.WithValue(ctx, "MemoryLimit", memoryLimitMB)
			efcontext, err := FromContext(ctx)

			Convey("should retrieve all values from context", func() {
				So(err, ShouldBeNil)
				So(efcontext.RequestID, ShouldEqual, requestId)
				So(efcontext.CollectStdio, ShouldEqual, collectStdio)
				So(efcontext.FunctionName, ShouldEqual, functionName)
				So(efcontext.FunctionVersion, ShouldEqual, functionVersion)
				So(efcontext.Qualifier, ShouldEqual, qualifier)
				So(efcontext.MemoryLimitMB, ShouldEqual, memoryLimitMB)
			})
		})

		Convey("with invalid context", func() {
			efcontext, err := FromContext(ctx)
			Convey("should return nil and false", func() {
				So(efcontext, ShouldBeNil)
				So(err, ShouldNotBeNil)
			})
		})
	})
}
