package dom

import "testing"

func TestProcessBackspace(t *testing.T) {
	// TODO: this is kind of silly; hook it up to the real text widget when that exists
	str := "hello world"
	startIdx, endIdx := 3, 10
	newStr := str[:startIdx] + str[endIdx:]
	expectedNewStr := "held"
	if newStr != expectedNewStr {
		t.Fatalf("expected %s got %s", expectedNewStr, newStr)
	}
}
