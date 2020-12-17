package format_utils

import "testing"

func TestFloatAmount2Int(t *testing.T) {
	testData := []struct {
		src  string
		dest string
	}{
		{
			src:  "12.13",
			dest: "1213",
		},
		{
			src:  "12.1",
			dest: "1210",
		},
		{
			src:  "12.",
			dest: "1200",
		},
		{
			src:  "12",
			dest: "1200",
		},
		{
			src:  "0.13",
			dest: "13",
		},
		{
			src:  "0.00",
			dest: "0",
		},
		{
			src:  "0",
			dest: "0",
		},
		{
			src:  "",
			dest: "0",
		},
	}

	for _, data := range testData {
		tmp, err := FloatAmount2Int(data.src, 2)
		if err != nil {
			t.Error("covert error->" + err.Error())
		}
		if tmp != data.dest {
			t.Error("fail src:", data.src, " dest:", data.dest, " actually:", tmp)
		}
	}
}
