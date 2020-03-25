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
	"fmt"
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func Test_NewHandler(t *testing.T) {
	Convey("NewHandler", t, func() {
		Convey("Invalid Handlers", func() {
			emptyBytes := []byte{}

			Convey("nil handler", func() {
				handlerWrapper := NewHandler(nil)
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "handler is nil")
			})

			Convey("handler is not a function", func() {
				handlerWrapper := NewHandler(struct{}{})
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "handler kind struct is not func")
			})

			Convey("handler declares too many arguments", func() {
				handler := func(n context.Context, x string, y string) error {
					return nil
				}
				handlerWrapper := NewHandler(handler)
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "handlers may not take more than two arguments, but handler takes 3")
			})

			Convey("two argument handler does not context as first argument", func() {
				handler := func(a string, x context.Context) error {
					return nil
				}
				handlerWrapper := NewHandler(handler)
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "handler takes two arguments, but the first is not Context. got string")
			})

			Convey("handler returns too many values", func() {
				handler := func() (error, error, error) {
					return nil, nil, nil
				}
				handlerWrapper := NewHandler(handler)
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "handler may not return more than two values")
			})

			Convey("handler returning two values does not declare error as the second return value", func() {
				handler := func() (error, string) {
					return nil, "hello"
				}
				handlerWrapper := NewHandler(handler)
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "handler returns two values, but the second does not implement error")
			})

			Convey("handler returning a single value does not implement error", func() {
				handler := func() string {
					return "hello"
				}
				handlerWrapper := NewHandler(handler)
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "handler returns a single value, but it does not implement error")
			})

			Convey("no return value should not result in error", func() {
				handlerWrapper := NewHandler(func() {})
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldBeNil)
			})

			Convey("handler returning a incompatible with json.Marshal return error", func() {
				handler := func() (interface{}, error) {
					return func() {}, nil
				}
				handlerWrapper := NewHandler(handler)
				_, err := handlerWrapper.Invoke(context.TODO(), emptyBytes)
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldContainSubstring, bad_json_response_message)
			})
		})

		Convey("Invocations", func() {
			Convey("func(string) (string, error)", func() {
				handler := func(input string) (string, error) {
					return fmt.Sprintf("faas %s", input), nil
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldBeNil)
				So(string(resp), ShouldEqual, `"faas test"`)
			})
			Convey("func(context, string) (string, error)", func() {
				handler := func(ctx context.Context, input string) (string, error) {
					return fmt.Sprintf("faas %s", input), nil
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldBeNil)
				So(string(resp), ShouldEqual, `"faas test"`)
			})
			Convey("func(*string) (*string, error)", func() {
				handler := func(input *string) (*string, error) {
					out := fmt.Sprintf("faas %s", *input)
					return &out, nil
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldBeNil)
				So(string(resp), ShouldEqual, `"faas test"`)
			})
			Convey("func(context, *string) (*string, error)", func() {
				handler := func(ctx context.Context, input *string) (*string, error) {
					out := fmt.Sprintf("faas %s", *input)
					return &out, nil
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldBeNil)
				So(string(resp), ShouldEqual, `"faas test"`)
			})
			Convey("func() error", func() {
				handler := func() error {
					return errors.New("no")
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no")
				So(string(resp), ShouldEqual, "")
			})
			Convey("func() (interface{}, error)", func() {
				handler := func() (interface{}, error) {
					return nil, errors.New("no")
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no")
				So(string(resp), ShouldEqual, "")
			})
			Convey("func(interface{}) (interface{}, error)", func() {
				handler := func(in interface{}) (interface{}, error) {
					return nil, errors.New("no")
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no")
				So(string(resp), ShouldEqual, "")
			})
			Convey("func(context, interface{}) (interface{}, error)", func() {
				handler := func(ctx context.Context, in interface{}) (interface{}, error) {
					return nil, errors.New("no")
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"test"`))
				So(err, ShouldNotBeNil)
				So(err.Error(), ShouldEqual, "no")
				So(string(resp), ShouldEqual, "")
			})
			Convey("func(struct{}) (string, error)", func() {
				handler := func(payload struct{ String string }) (string, error) {
					return payload.String, nil
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`{"String": "hello"}`))
				So(err, ShouldBeNil)
				So(string(resp), ShouldEqual, `"hello"`)
			})
			Convey("func(string) (struct{}, error)", func() {
				handler := func(payload string) (struct{ String string }, error) {
					return struct{ String string }{payload}, nil
				}

				wrapper := NewHandler(handler)
				resp, err := wrapper.Invoke(context.TODO(), []byte(`"hello"`))
				So(err, ShouldBeNil)
				So(string(resp), ShouldEqual, `{"String":"hello"}`)
			})
		})
	})
}
