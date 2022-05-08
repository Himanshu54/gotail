package tail

import (
	"fmt"
	"io"
	"os"
	"strings"
	"testing"
)

func TestSeekLineNFromEnd(t *testing.T) {

	f, err := os.OpenFile(fmt.Sprintf("%s/%s", os.TempDir(), "test-read-n-line"), os.O_RDWR|os.O_CREATE|os.O_TRUNC, 0666)
	if err != nil {
		t.Fatal(err)
	}
	defer f.Close()
	defer os.Remove(f.Name())
	msg := "Programming today is a race between software engineers striving to \n" +
		"build bigger and better idiot-proof programs, and the Universe trying \n" +
		"to produce bigger and better idiots. So far, the Universe is winning."

	if _, err := f.Write([]byte(msg)); err != nil {
		t.Fatalf("WriteFile %s: %v", f.Name(), err)
	}
	test := []struct {
		lines   uint
		exected string
	}{
		{
			1,
			"to produce bigger and better idiots. So far, the Universe is winning.",
		},
		{
			2,
			"build bigger and better idiot-proof programs, and the Universe trying \n" +
				"to produce bigger and better idiots. So far, the Universe is winning.",
		},
		{
			10,
			msg,
		},
	}

	data := make([]byte, 1000)
	for _, tt := range test {
		f.Seek(0, 0)
		SeekLineNFromEnd(f, tt.lines)
		n, err := f.Read(data)
		if err != nil && err != io.EOF {
			t.Fatalf("ReadFile %s: %v", f.Name(), err)
		}
		if string(data[:n]) != tt.exected {
			t.Fatalf("ReadFile: wrong data:\nhave --%q--\nwant --%q--", string(data[:n]), tt.exected)
		}

	}
}

func TestTailFile(t *testing.T) {
	filename := "README.md"
	content, err := TailFile(filename, 7)
	if err != nil {
		t.Error(err)
	}
	if len(strings.Split(string(content), "\n")) != 7 {
		t.Fail()
	}
}
