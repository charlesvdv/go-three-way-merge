package go_three_way_merge

import (
	"github.com/sergi/go-diff/diffmatchpatch"
)

type diffIterator struct {
	diffs        []diffmatchpatch.Diff
	diffOffset   int
	runeOffset   int
	currentRunes []rune
}

func newDiffIterator(diffs []diffmatchpatch.Diff) diffIterator {
	var currentRunes []rune
	if len(diffs) != 0 {
		currentRunes = []rune(diffs[0].Text)
	}

	return diffIterator{
		diffs:        diffs,
		diffOffset:   0,
		runeOffset:   0,
		currentRunes: currentRunes,
	}
}

func (i *diffIterator) next() {
	i.runeOffset++

	if i.runeOffset == len(i.currentRunes) {
		i.runeOffset = 0
		i.diffOffset++
		if i.diffOffset < len(i.diffs) {
			i.currentRunes = []rune(i.diffs[i.diffOffset].Text)
		}
	}
}

func (i diffIterator) diffType() int {
	return int(i.diffs[i.diffOffset].Type)
}

func (i diffIterator) value() string {
	return string(i.currentRunes[i.runeOffset])
}

func (i diffIterator) valueRune() rune {
	return i.currentRunes[i.runeOffset]
}

func (i diffIterator) isFinished() bool {
	return i.diffOffset >= len(i.diffs)
}
