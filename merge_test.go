package go_three_way_merge

import (
	"testing"
)

func expectSuccessfulMerge(t *testing.T, err error, ok bool) {
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if !ok {
		t.Fatal("Expected no merge conflict")
	}
}

func TestNoMergeRequired(t *testing.T) {
	text := "Hello World!"
	result, ok, err := Merge(text, text, text)

	expectSuccessfulMerge(t, err, ok)
	if result != text {
		t.Errorf("Expected %s, got %s.", text, result)
	}
}

func TestSimpleInsertDeleteMerge(t *testing.T) {
	base := "Hello World!"
	versionA := "Coucou World!"
	versionB := "Hello World and Space!"

	result, ok, err := Merge(base, versionA, versionB)
	expectSuccessfulMerge(t, err, ok)

	expected := "Coucou World and Space!"
	if result != expected {
		t.Errorf("Expected %s, got %s.", expected, result)
	}
}

func TestMergeWithSameDiff(t *testing.T) {
	base := "Hello World!"
	versionA := "Coucou World!"
	versionB := "Coucou World!"

	result, ok, err := Merge(base, versionA, versionB)
	expectSuccessfulMerge(t, err, ok)

	expected := "Coucou World!"
	if result != expected {
		t.Errorf("Expected %s, got %s.", expected, result)
	}
}

func TestMergeConflict(t *testing.T) {
	base := "Hello World!"
	versionA := "Coucou World!"
	versionB := "Hola World!"

	_, ok, err := Merge(base, versionA, versionB)
	if err != nil {
		t.Fatalf("Unexpected error: %s", err)
	}
	if ok {
		t.Error("Should merge-conflict")
	}
}
