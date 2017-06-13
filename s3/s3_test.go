// Code generated by https://github.com/abcum/tmpl
// Source file: s3/s3_test.go.tmpl
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

package s3

import (
	"io"
	"io/ioutil"
	"os"
	"testing"

	. "github.com/smartystreets/goconvey/convey"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var run = 5000
var sze int
var txt []byte
var out []byte

func init() {

	os.RemoveAll("syncr/")

	n := session.Must(session.NewSession())
	s := s3.New(n, &aws.Config{Region: aws.String("eu-west-1")})
	s.ListObjectsPages(&s3.ListObjectsInput{
		Bucket: aws.String("abcum-tests"),
		Prefix: aws.String("syncr/"),
	}, func(p *s3.ListObjectsOutput, last bool) bool {
		for _, f := range p.Contents {
			s.DeleteObject(&s3.DeleteObjectInput{
				Bucket: aws.String("abcum-tests"),
				Key:    f.Key,
			})
		}
		return true
	})

	txt, _ = ioutil.ReadFile("../data.txt")

}

func TestMain(t *testing.T) {

	defer func() {
		os.RemoveAll("syncr/")
	}()

	defer func() {
		n := session.Must(session.NewSession())
		s := s3.New(n, &aws.Config{Region: aws.String("eu-west-1")})
		s.ListObjectsPages(&s3.ListObjectsInput{
			Bucket: aws.String("abcum-tests"),
			Prefix: aws.String("syncr/"),
		}, func(p *s3.ListObjectsOutput, last bool) bool {
			for _, f := range p.Contents {
				s.DeleteObject(&s3.DeleteObjectInput{
					Bucket: aws.String("abcum-tests"),
					Key:    f.Key,
				})
			}
			return true
		})
	}()

	var s *Storage

	Convey("Create a new syncr instance", t, func() {
		h, e := New("abcum-tests/syncr/test.db", &Options{MinSize: 1})
		So(e, ShouldBeNil)
		s = h
	})

	Convey("Read data from the syncr instance", t, func() {
		s.Seek(0, 0)
		b := make([]byte, 128)
		_, e := s.Read(b)
		So(e, ShouldEqual, io.EOF)
	})

	Convey("Write data to the syncr instance", t, func() {
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

	Convey("Sync data to the syncr instance", t, func() {
		So(s.Sync(), ShouldBeNil)
	})

	Convey("Read data from the syncr instance", t, func() {
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

	Convey("Close the syncr instance", t, func() {
		So(s.Close(), ShouldBeNil)
	})

	Convey("Close again", t, func() {
		So(s.Close(), ShouldBeNil)
	})

}
