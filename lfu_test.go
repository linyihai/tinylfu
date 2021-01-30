package lfu

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestLFU(t *testing.T) {
	lfu := NewLFU(2)
	lfu.Put(1, 1)
	lfu.Put(2, 2)

	Convey("should get 1", t, func() {
		So(lfu.Get(1), ShouldEqual, 1)
	})
	Convey("GetAll ok", t, func() {
		So(lfu.GetAll(), ShouldResemble, []interface{}{1, 2})
	})

	iter := lfu.GetIterator()
	var count int
	for {
		count++
		if v := iter(); v == nil {
			break
		} else {
			switch count {
				case 1:
					Convey("should get 1", t, func(){
						So(v.value, ShouldEqual, 1)
					})
				case 2:
					Convey("should get 2", t, func(){
						So(v.value, ShouldEqual, 2)
					})
			}
		}
	}

	lfu.Put(3, 3)
	Convey("should get -1", t, func(){
		So(lfu.Get(2), ShouldEqual, -1)
	})
	Convey("should get 3", t, func(){
		So(lfu.Get(3), ShouldEqual, 3)
	})

	Convey("GetAll failed", t, func(){
		So(lfu.GetAll(), ShouldResemble, []interface{}{3,1})
	})

	iter = lfu.GetIterator()
	count = 0
	for {
		count++
		if v := iter(); v == nil {
			break
		} else {
			switch count {
				case 1:
					Convey("should get 3", t, func(){
						So(v.value, ShouldEqual, 3)
					})
				case 2:
					Convey("should get 1", t, func(){
						So(v.value, ShouldEqual, 1)
					})
			}
		}
	}

	lfu.Put(4, 4)

	count = 0
	iter = lfu.GetIterator()
	for {
		count++
		if v := iter(); v == nil {
			break
		} else {
			switch count {
				case 1:
					Convey("should get 3", t, func(){
						So(v.value, ShouldEqual, 3)
					})
				case 2:
					Convey("should get 4", t, func(){
						So(v.value, ShouldEqual, 4)
					})
			}
		}
	}
}
