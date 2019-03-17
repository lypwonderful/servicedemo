package esSearch

import (
	. "github.com/smartystreets/goconvey/convey"
	"testing"
)

func TestIsEsCluserOk(t *testing.T) {
	Convey("test is es cluster is ok", t, func() {
		Convey("case 1,succ", func() {
			isOK := IsEsCluserOk()
			So(isOK, ShouldBeTrue)
		})
	})
}

func TestNewEsClient(t *testing.T) {
	esCli := NewEsClient()
	if esCli == nil {
		t.Error("new es cli return error.")
	}
}

func TestUseEs(t *testing.T) {
	Convey("test  useEs func", t, func() {
		Convey("case 1,succ", func() {
			err := UseEs()
			So(err, ShouldBeNil)
		})
	})
}

func TestAddDocmToIndex(t *testing.T) {
	Convey("test add a doc in index", t, func() {
		Convey("case 1,succ", func() {
			err := AddDocmToIndex()
			So(err, ShouldBeNil)
		})
	})
}
func TestSearchEs(t *testing.T) {
	Convey("test search a name", t, func() {
		Convey("case 1,succ", func() {
			err := SearchEs()
			So(err, ShouldBeNil)
		})
	})
}

func TestReadLog(t *testing.T) {
	Convey("test readLog", t, func() {
		Convey("case 1,succ", func() {
			readLog()
		})
	})
}
