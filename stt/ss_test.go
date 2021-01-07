package stt

import (
	"reflect"
	"testing"
)

type SSTTestCase struct {
	a string
	b interface{}
}

var decodeCaseList = []SSTTestCase{
	{a: "aa@=11/bb@=123/", b: map[string]interface{}{"aa": "11", "bb": "123"}},
	{a: "a@Sa@=1@A1/", b: map[string]interface{}{"a/a": "1@1"}},
	{a: "aa@=11", b: map[string]interface{}{"aa": "11"}},
	{a: "aa/bb/cc/", b: []interface{}{"aa", "bb", "cc"}},
	{a: "aa/bb/cc", b: []interface{}{"aa", "bb", "cc"}},
	{
		a: "type@=rnewbc/rid@=0/gid@=0/bt@=1/uid@=172172172/unk@=啊啊啊/uic@=avanew@Sface@S201712@S19@S07@S4a1d3c24cf37e54a9ba2c4f3ceea4118/drid@=123123123/donk@=SS水水水水/nl@=3/",
		b: map[string]interface{}{
			"type": "rnewbc",
			"rid":  "0",
			"gid":  "0",
			"bt":   "1",
			"uid":  "172172172",
			"unk":  "啊啊啊",
			"uic":  "avanew/face/201712/19/07/4a1d3c24cf37e54a9ba2c4f3ceea4118",
			"drid": "123123123",
			"donk": "SS水水水水",
			"nl":   "3",
		},
	},
	{
		a: "bt@=1/donk@=SS水水水水/drid@=123123123/gid@=0/nl@=3/rid@=0/type@=rnewbc/uic@=avanew@Sface@S201712@S19@S07@S4a1d3c24cf37e54a9ba2c4f3ceea4118/uid@=172172172/unk@=啊啊啊/",
		b: map[string]interface{}{
			"type": "rnewbc",
			"rid":  "0",
			"gid":  "0",
			"bt":   "1",
			"uid":  "172172172",
			"unk":  "啊啊啊",
			"uic":  "avanew/face/201712/19/07/4a1d3c24cf37e54a9ba2c4f3ceea4118",
			"drid": "123123123",
			"donk": "SS水水水水",
			"nl":   "3",
		},
	},
}

var encodeCaseList = []SSTTestCase{
	{a: "aa@=11/bb@=123/", b: map[string]interface{}{"aa": "11", "bb": "123"}},
	{a: "a@Sa@=1@A1/", b: map[string]interface{}{"a/a": "1@1"}},
	{a: "aa/bb/cc/", b: []interface{}{"aa", "bb", "cc"}},
	{
		a: "bt@=1/donk@=SS水水水水/drid@=123123123/gid@=0/nl@=3/rid@=0/type@=rnewbc/uic@=avanew@Sface@S201712@S19@S07@S4a1d3c24cf37e54a9ba2c4f3ceea4118/uid@=172172172/unk@=啊啊啊/",
		b: map[string]interface{}{
			"type": "rnewbc",
			"rid":  "0",
			"gid":  "0",
			"bt":   "1",
			"uid":  "172172172",
			"unk":  "啊啊啊",
			"uic":  "avanew/face/201712/19/07/4a1d3c24cf37e54a9ba2c4f3ceea4118",
			"drid": "123123123",
			"donk": "SS水水水水",
			"nl":   "3",
		},
	},
}

func TestDecode(t *testing.T) {
	for i, cc := range decodeCaseList {
		out := Decode(cc.a)
		if !reflect.DeepEqual(out, cc.b) {
			t.Errorf("#%d: sst decode\nhave: %#+v\nwant: %#+v", i, out, cc.b)
			continue
		}
	}
}

func TestEncode(t *testing.T) {
	for i, cc := range encodeCaseList {
		out := Encode(cc.b)
		if out != cc.a {
			t.Errorf("#%d: sst encode\nhave: %#+v\nwant: %#+v", i, out, cc.a)
			continue
		}

		encodeOut := Decode(out)
		if !reflect.DeepEqual(encodeOut, cc.b) {
			t.Errorf("#%d: sst decode\nhave: %#+v\nwant: %#+v", i, encodeOut, cc.b)
			continue
		}
	}
}
