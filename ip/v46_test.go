package ip_test

import (
	"gommon/ip"
	"testing"
)

func TestV46SearchForV4(t *testing.T) {
	v46s, err := ip.NewV46s("./test_data/mgiplib-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load v4 data file, %v", err)
		return
	}

	v46 := v46s.Search("223.242.47.30")
	if v46 == nil {
		t.Errorf("Can not find [%s]", "223.242.32.30")
		return
	}

	expected := "223.242.32.0|223.242.47.255|CN|CT|安徽|芜湖|576074"
	if v46s.StringOf(v46) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, v46s.StringOf(v46))
		return
	}

	if v46 := v46s.Search("223.242.64.289"); v46 != nil {
		t.Errorf("Expected not found, but got [%s]", v46s.StringOf(v46))
		return
	}
}

func TestV46SearchForV6(t *testing.T) {
	v46s, err := ip.NewV46s("./test_data/mgiplib-v6-std.txt.latest")
	if err != nil {
		t.Errorf("Can not load v6 data file, %v", err)
		return
	}

	v46 := v46s.Search("240e:6af:4700:1111::")
	if v46 == nil {
		t.Errorf("Can not find [%s]", "240e:6af:4700:1111::")
		return
	}

	expected := "240e:06af:4700:0000:0000:0000:0000:0000|240e:6af:47ff:ffff:ffff:ffff:ffff:ffff|CN|CT|江苏|淮安|95566"
	if v46s.StringOf(v46) != expected {
		t.Errorf("Expected[%s], but[%s]", expected, v46s.StringOf(v46))
		return
	}

}
