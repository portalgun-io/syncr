// Code generated by https://github.com/abcum/tmpl
// Source file: gcs/gcs_test.go.tmpl
// DO NOT EDIT!

// Copyright © 2016 Abcum Ltd
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

package gcs

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"cloud.google.com/go/storage"
	"context"
)

var run = 5000
var sze int
var txt []byte
var out []byte

func TestMain(t *testing.T) {

	os.RemoveAll("syncr/")

	x := context.Background()
	c, e := storage.NewClient(x)
	if e != nil {
		panic(e)
	}
	b := c.Bucket("abcum-tests")
	i := b.Objects(x, &storage.Query{Prefix: "syncr/"})
	for {
		o, e := i.Next()
		if e != nil {
			break
		}
		b.Object(o.Name).Delete(x)
	}

	txt, _ = ioutil.ReadFile("../data.txt")

	defer func() {
		os.RemoveAll("syncr/")
	}()

	defer func() {
		x := context.Background()
		c, e := storage.NewClient(x)
		if e != nil {
			panic(e)
		}
		b := c.Bucket("abcum-tests")
		i := b.Objects(x, &storage.Query{Prefix: "syncr/"})
		for {
			o, e := i.Next()
			if e != nil {
				break
			}
			b.Object(o.Name).Delete(x)
		}
	}()

	var s *Storage

	Convey("Create a new gcs syncr instance", t, func() {
		h, e := New("abcum-tests/syncr/test.db", &Options{MinSize: 1})
		So(e, ShouldBeNil)
		s = h
	})

	Convey("Read data from the gcs syncr instance", t, func() {
		s.Seek(0, 0)
		b := make([]byte, 128)
		_, e := s.Read(b)
		So(e, ShouldEqual, io.EOF)
	})

	Convey("Write data to the gcs syncr instance", t, func() {
		s.Seek(0, 0)
		for i := 0; i < run; i++ {
			l, e := s.Write(txt)
			sze += l
			if e != nil {
				if e != io.EOF {
					panic(e)
				}
				break
			}
		}
		So(sze, ShouldEqual, run*len(txt))
	})

	Convey("Sync data to the gcs syncr instance", t, func() {
		So(s.Sync(), ShouldBeNil)
	})

	Convey("Read data from the gcs syncr instance", t, func() {
		s.Seek(0, 0)
		for i := 0; i <= run*50; i++ {
			b := make([]byte, 128)
			l, e := s.Read(b)
			out = append(out, b[:l]...)
			if e != nil {
				if e != io.EOF {
					panic(e)
				}
				break
			}
		}
		So(len(out), ShouldEqual, run*len(txt))
		So(out[:len(txt)], ShouldResemble, txt)
		So(out[len(out)-len(txt):], ShouldResemble, txt)
	})

	Convey("Attempt to seek to positions", t, func() {
		var e error
		_, e = s.Seek(0, 0)
		So(e, ShouldBeNil)
		_, e = s.Seek(1, 0)
		So(e, ShouldNotBeNil)
		_, e = s.Seek(2, 0)
		So(e, ShouldNotBeNil)
		_, e = s.Seek(0, 1)
		So(e, ShouldNotBeNil)
		_, e = s.Seek(0, 2)
		So(e, ShouldBeNil)
	})

	Convey("Close the gcs syncr instance", t, func() {
		So(s.Close(), ShouldBeNil)
	})

	Convey("Close again", t, func() {
		So(s.Close(), ShouldBeNil)
	})

}
