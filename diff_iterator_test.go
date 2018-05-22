package go_three_way_merge

import (
    "testing"

    "github.com/sergi/go-diff/diffmatchpatch"
)

func getIterator(base, newVersion string) diffIterator {
    dmp := diffmatchpatch.New()
    diffs := dmp.DiffMain(base, newVersion, false)
    return newDiffIterator(diffs)
}

func testNextIteration(t *testing.T, iterator *diffIterator, expectedValue string, expectedType int) {
    if iterator.value() != expectedValue {
        t.Errorf("Expected value '%s', got '%s'.", expectedValue, iterator.value())
    }
    if iterator.diffType() != expectedType {
        t.Errorf("Expected diff type %d, got %d for expected value '%s'", expectedType, iterator.diffType(), expectedValue)
    }
    if iterator.isFinished() {
        t.Error("Iterator should not finish.")
    }
    iterator.next()
}

func TestBasicIterator(t *testing.T) {
    iterator := getIterator("abcde fg", "acdefghij")

    testNextIteration(t, &iterator, "a", diffEqual)
    testNextIteration(t, &iterator, "b", diffDelete)
    testNextIteration(t, &iterator, "c", diffEqual)
    testNextIteration(t, &iterator, "d", diffEqual)
    testNextIteration(t, &iterator, "e", diffEqual)
    testNextIteration(t, &iterator, " ", diffDelete)
    testNextIteration(t, &iterator, "f", diffEqual)
    testNextIteration(t, &iterator, "g", diffEqual)
    testNextIteration(t, &iterator, "h", diffInsert)
    testNextIteration(t, &iterator, "i", diffInsert)
    testNextIteration(t, &iterator, "j", diffInsert)

    if !iterator.isFinished() {
        t.Error("Iterator should be finished by now.")
    }
}

func TestUnicodeSupport(t *testing.T) {
    iterator := getIterator("ùéçà$¤", "ùèça€ĸ")

    testNextIteration(t, &iterator, "ù", diffEqual)
    testNextIteration(t, &iterator, "é", diffDelete)
    testNextIteration(t, &iterator, "è", diffInsert)
    testNextIteration(t, &iterator, "ç", diffEqual)
    testNextIteration(t, &iterator, "à", diffDelete)
    testNextIteration(t, &iterator, "$", diffDelete)
    testNextIteration(t, &iterator, "¤", diffDelete)
    testNextIteration(t, &iterator, "a", diffInsert)
    testNextIteration(t, &iterator, "€", diffInsert)
    testNextIteration(t, &iterator, "ĸ", diffInsert)
}
