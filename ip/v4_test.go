package ip_test

import (
	"gommon/ip"
	"testing"
)

func TestV4Search(t *testing.T) {
	v4s, err := ip.NewV4s("./test_data/mgiplib-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load v4 data file, %v", err)
		return
	}

	v4 := v4s.Search("223.242.47.30")
	if v4 == nil {
		t.Errorf("Can not find [%s]", "223.242.32.30")
		return
	}

	expected := "223.242.32.0|223.242.47.255|CN|CT|安徽|芜湖|576074"
	if v4s.StringOf(v4) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, v4s.StringOf(v4))
		return
	}

	if v4 := v4s.Search("223.242.64.289"); v4 != nil {
		t.Errorf("Expected not found, but got [%s]", v4s.StringOf(v4))
		return
	}
}
