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
	"testing"
	"time"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_invokeContext_ToContext(t *testing.T) {
	Convey("ToContext", t, func() {
		ic := invokeContext{
			RequestId:       "1234",
			CollectStdio:    false,
			FunctionName:    "myFunc",
			FunctionVersion: "v1",
			Qualifier:       "alias",
			MemoryLimit:     int32(256),
			Deadline:        time.Now().UnixNano() / int64(time.Millisecond),
		}

		ctx, _ := ic.ToContext()
		Convey("should copy all basic values", func() {
			So(ctx.Value("RequestId"), ShouldEqual, ic.RequestId)
			So(ctx.Value("CollectStdio"), ShouldEqual, ic.CollectStdio)
			So(ctx.Value("FunctionName"), ShouldEqual, ic.FunctionName)
			So(ctx.Value("FunctionVersion"), ShouldEqual, ic.FunctionVersion)
			So(ctx.Value("Qualifier"), ShouldEqual, ic.Qualifier)
			So(ctx.Value("MemoryLimit"), ShouldEqual, ic.MemoryLimit)
		})

		Convey("should recreate deadline", func() {
			ctxDeadline, ok := ctx.Deadline()
			So(ok, ShouldBeTrue)
			So(ctxDeadline.UnixNano(), ShouldEqual, ic.Deadline*int64(time.Millisecond))
		})
	})
}
