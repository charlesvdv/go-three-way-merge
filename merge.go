// Package go_three_way_merge implements an intuitive three way merge
// algorithm. See https://en.wikipedia.org/wiki/Merge_(version_control)#Three-way_merge
// for more detail.
package go_three_way_merge

import (
    "strings"
    "errors"

    "github.com/sergi/go-diff/diffmatchpatch"
)

const (
    diffDelete = int(diffmatchpatch.DiffDelete)
    diffInsert = int(diffmatchpatch.DiffInsert)
    diffEqual = int(diffmatchpatch.DiffEqual)
)

var (
    inconsistencyErr error = errors.New("Inconsistency error")
)

// Merge two versions of the same base file. Returns the resulting merged
// file if the boolean is true. If the boolean is false, there is either
// a merge-conflict or an error. In theory, this algorithm should never
// return an error (if so, feel free to open an issue).
func Merge(base, versionA, versionB string) (string, bool, error) {
    result, ok, err := MergeRunes([]rune(base), []rune(versionA), []rune(versionB))
    return result, ok, err
}

// Same as `Merge` with runes instead.
func MergeRunes(base, versionA, versionB []rune) (string, bool, error) {
    dmp := diffmatchpatch.New()

    diffA := dmp.DiffMainRunes(base, versionA, false)
    diffIteratorA := newDiffIterator(diffA)
    diffB := dmp.DiffMainRunes(base, versionB, false)
    diffIteratorB := newDiffIterator(diffB)

    var result strings.Builder
    offsetBase := 0

    for true {
        if diffIteratorA.isFinished() && diffIteratorB.isFinished() {
            break
        } else if diffIteratorA.isFinished() || diffIteratorB.isFinished() {
            var nonFinished *diffIterator
            if diffIteratorA.isFinished() {
                nonFinished = &diffIteratorB
            } else {
                nonFinished = &diffIteratorA
            }

            // TODO: verify this is always correct.
            if nonFinished.diffType() != diffInsert {
                return result.String(), false, errors.New("Unexpected diff.")
            }
            result.WriteRune(nonFinished.valueRune())
            nonFinished.next()
            continue
        }

        aType, bType := diffIteratorA.diffType(), diffIteratorB.diffType()
        aValue, bValue := diffIteratorA.valueRune(), diffIteratorB.valueRune()
        baseValue := base[offsetBase]

        if aType == diffEqual && bType == diffEqual {
            if aValue != baseValue || bValue != baseValue {
                return result.String(), false, inconsistencyErr
            }
            diffIteratorA.next()
            diffIteratorB.next()
            offsetBase += 1

            result.WriteRune(baseValue)
        } else if aType != diffEqual && bType != diffEqual {
            // Possible merge conflict...

            if aValue == bValue && aType == bType {
                // Same diff so doesn't conflict.
                diffIteratorA.next()
                diffIteratorB.next()

                if aType == diffInsert {
                    result.WriteRune(aValue)
                } else if aType == diffDelete {
                    offsetBase += 1
                }
            } else {
                // A merge-conflict has been found.
                return result.String(), false, nil
            }
        } else if aType == diffDelete || bType == diffDelete {
            if aValue != baseValue || bValue != baseValue {
                return result.String(), false, inconsistencyErr
            }
            diffIteratorA.next()
            diffIteratorB.next()
            offsetBase += 1
        } else if aType == diffInsert || bType == diffInsert {
            var insertIterator *diffIterator
            if aType == diffInsert {
                insertIterator = &diffIteratorA
            } else {
                insertIterator = &diffIteratorB
            }

            result.WriteRune(insertIterator.valueRune())
            insertIterator.next()
        }
    }

    return result.String(), true, nil
}
