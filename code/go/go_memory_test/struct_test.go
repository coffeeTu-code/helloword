package buildin

import (
	"testing"
	"unsafe"

	. "github.com/smartystreets/goconvey/convey"
)

func TestStructSize(t *testing.T) {
	Convey("基本数据类型大小检测", t, func() {
		So(unsafe.Sizeof(true), ShouldEqual, 1)
		So(unsafe.Sizeof(uint8(0)), ShouldEqual, 1)
		So(unsafe.Sizeof(int8(0)), ShouldEqual, 1)
		So(unsafe.Sizeof(int16(0)), ShouldEqual, 2)
		So(unsafe.Sizeof(int32(0)), ShouldEqual, 4)
		So(unsafe.Sizeof(int64(0)), ShouldEqual, 8)
		So(unsafe.Sizeof(int(0)), ShouldEqual, 8)
		So(unsafe.Sizeof(float32(0)), ShouldEqual, 4)
		So(unsafe.Sizeof(float64(0)), ShouldEqual, 8)
		So(unsafe.Sizeof(""), ShouldEqual, 16)
		So(unsafe.Sizeof("hello world"), ShouldEqual, 16)
		So(unsafe.Sizeof([]byte("hello world")), ShouldEqual, 24)
		So(unsafe.Sizeof([]int{}), ShouldEqual, 24)
		So(unsafe.Sizeof([]int{1, 2, 3}), ShouldEqual, 24)
		So(unsafe.Sizeof([3]int{1, 2, 3}), ShouldEqual, 24)
		So(unsafe.Sizeof(map[string]string{}), ShouldEqual, 8)
		So(unsafe.Sizeof(map[string]string{"1": "one", "2": "two"}), ShouldEqual, 8)
		So(unsafe.Sizeof(struct{}{}), ShouldEqual, 0)
	})

	Convey("自定义类型大小检测", t, func() {
		// |x---|
		So(unsafe.Sizeof(struct {
			i8 int8
		}{}), ShouldEqual, 1)

		// |x---|xxxx|xx--|
		So(unsafe.Sizeof(struct {
			i8  int8
			i32 int32
			i16 int16
		}{}), ShouldEqual, 12)

		// |x-xx|xxxx|
		So(unsafe.Sizeof(struct {
			i8  int8
			i16 int16
			i32 int32
		}{}), ShouldEqual, 8)

		// |x---|xxxx|xx--|----|xxxx|xxxx|
		So(unsafe.Sizeof(struct {
			i8  int8
			i32 int32
			i16 int16
			i64 int64
		}{}), ShouldEqual, 24)

		// |x-xx|xxxx|xxxx|xxxx|
		So(unsafe.Sizeof(struct {
			i8  int8
			i16 int16
			i32 int32
			i64 int64
		}{}), ShouldEqual, 16)

		type I8 int8
		type I16 int16
		type I32 int32

		So(unsafe.Sizeof(struct {
			i8  I8
			i16 I16
			i32 I32
		}{}), ShouldEqual, 8)
	})

	Convey("SCreativeAttr", t, func() {
		So(unsafe.Sizeof(struct {
			//系统自定义类型
			ResourceType    uint8
			SubResourceType uint8
			FormatType      uint8
			//素材通用类型
			DocId         string
			CreativeType  int //"201"
			CreativeKey   int64
			CreativeId    int64  //"creativeId" : NumberLong("2753977716"),
			AdvCreativeId string //eq: "1800060617"
			Mime          string //"mime" : "video/mp4",
			Attribute     string //"attribute" : "13，7，6"
			FMd5          string
			Source        int
			Ctime         int64
			Utime         int64
			//video素材独有
			VideoResolution string //"videoResolution" : "1280x720",
			VideoSize       string //"videoSize" : NumberLong(761894),
			VideoLength     int    //"videoLength" : NumberLong(18)
			Bitrate         int    //"bitrate" : NumberLong(334),
			Clarity         int    //"clarity" : NumberLong(1),
			Width           int    //"width" : NumberLong(1280),
			Height          int    //"height" : NumberLong(720),
			Orientation     int    //"orientation" : NumberLong(2),
			//pl zip素材
			MinOs    int64 //eq: 9000000
			Platform int32 //eq: 2
			//图片素材
			Resolution string //eq: "1000x560"
			Url        string //eq: "http://cdn-adn.rayjump.com/v3/v/19/11/04/13/Iqb90Oo0-e62cfe29-5aca-46e1-a5e5-023062967d9a.mp4"
			//pl js素材
			TagCode string //eq: "http://cdn-adn.rayjump.com/v3/v/19/11/04/13/Iqb90Oo0-e62cfe29-5aca-46e1-a5e5-023062967d9a.mp4"
			//common素材
			Value string
		}{}), ShouldEqual, 296)

		// |x-xx|xxxx|
		So(unsafe.Sizeof(struct {
			ResourceType    int8
			SubResourceType int8
			FormatType      int8
			Platform        int8 //eq: 2
			CreativeType    int8 //"201"
			Source          int8
			Clarity         int8  //"clarity" : NumberLong(1),
			Width           int8  //"width" : NumberLong(1280),
			Height          int8  //"height" : NumberLong(720),
			Orientation     int8  //"orientation" : NumberLong(2),
			VideoLength     int8  //"videoLength" : NumberLong(18)
			Bitrate         int8  //"bitrate" : NumberLong(334),
			MinOs           int32 //eq: 9000000
			CreativeKey     int64
			CreativeId      int64 //"creativeId" : NumberLong("2753977716"),
			Ctime           int64
			Utime           int64
			AdvCreativeId   int64 //eq: "1800060617"
			DocId           string
			Mime            string //"mime" : "video/mp4",
			Attribute       string //"attribute" : "13，7，6"
			FMd5            string
			VideoResolution string //"videoResolution" : "1280x720",
			VideoSize       string //"videoSize" : NumberLong(761894),
			Resolution      string //eq: "1000x560"
			Value           string
			Url             string //eq: "http://cdn-adn.rayjump.com/v3/v/19/11/04/13/Iqb90Oo0-e62cfe29-5aca-46e1-a5e5-023062967d9a.mp4"
			TagCode         string //eq: "http://cdn-adn.rayjump.com/v3/v/19/11/04/13/Iqb90Oo0-e62cfe29-5aca-46e1-a5e5-023062967d9a.mp4"
		}{}), ShouldEqual, 216)
	})
}

type SCreativeAttr struct {
	//系统自定义类型
	ResourceType    uint8
	SubResourceType uint8
	FormatType      uint8
	//素材通用类型
	DocId         string
	CreativeType  int //"201"
	CreativeKey   int64
	CreativeId    int64  //"creativeId" : NumberLong("2753977716"),
	AdvCreativeId string //eq: "1800060617"
	Mime          string //"mime" : "video/mp4",
	Attribute     string //"attribute" : "13，7，6"
	FMd5          string
	Source        int
	Ctime         int64
	Utime         int64
	//video素材独有
	VideoResolution string //"videoResolution" : "1280x720",
	VideoSize       string //"videoSize" : NumberLong(761894),
	VideoLength     int    //"videoLength" : NumberLong(18)
	Bitrate         int    //"bitrate" : NumberLong(334),
	Clarity         int    //"clarity" : NumberLong(1),
	Width           int    //"width" : NumberLong(1280),
	Height          int    //"height" : NumberLong(720),
	Orientation     int    //"orientation" : NumberLong(2),
	//pl zip素材
	MinOs    int64 //eq: 9000000
	Platform int32 //eq: 2
	//图片素材
	Resolution string //eq: "1000x560"
	Url        string //eq: "http://cdn-adn.rayjump.com/v3/v/19/11/04/13/Iqb90Oo0-e62cfe29-5aca-46e1-a5e5-023062967d9a.mp4"
	//pl js素材
	TagCode string //eq: "http://cdn-adn.rayjump.com/v3/v/19/11/04/13/Iqb90Oo0-e62cfe29-5aca-46e1-a5e5-023062967d9a.mp4"
	//common素材
	Value string
}
