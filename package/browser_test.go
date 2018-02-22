package jankybrowser

import (
	"reflect"
	"testing"
)

const ex1 = "http://example.com/"
const ex2 = "http://example.com/other"

func TestPopHistory(t *testing.T) {
	b := Browser{
		history: []string{ex1},
	}
	_, err := b.PopHistory()
	expectedErr := "can't go back; already on last page"
	if err == nil || err.Error() != expectedErr {
		t.Fatalf("expected %s; got %v", expectedErr, err)
	}

	b2 := Browser{
		history: []string{ex1, ex2},
	}
	url, err := b2.PopHistory()
	if err != nil {
		t.Fatalf("expected no err; got %v", err)
	}
	if url != ex1 {
		t.Fatalf("expected %s; got %s", ex1, url)
	}
	if !reflect.DeepEqual(b2.history, []string{}) {
		t.Fatalf("expected %v; got %v", []string{}, b2.history)
	}
}
