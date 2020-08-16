package exec

import (
	"fmt"
	"strings"

	"github.com/wedaly/aretext/internal/pkg/text"
)

// Mutator modifies the state of the cursor or text.
type Mutator interface {
	fmt.Stringer
	Mutate(state *State)
}

// CompositeMutator executes a series of mutations.
type CompositeMutator struct {
	subMutators []Mutator
}

func NewCompositeMutator(subMutators []Mutator) Mutator {
	return &CompositeMutator{subMutators}
}

// Mutate executes a series of mutations in order.
func (cm *CompositeMutator) Mutate(state *State) {
	for _, mut := range cm.subMutators {
		mut.Mutate(state)
	}
}

func (cm *CompositeMutator) String() string {
	args := make([]string, 0, len(cm.subMutators))
	for _, mut := range cm.subMutators {
		args = append(args, mut.String())
	}
	return fmt.Sprintf("Composite(%s)", strings.Join(args, ","))
}

// cursorMutator updates the current location of the cursor.
type cursorMutator struct {
	loc CursorLocator
}

// NewCursorMutator returns a mutator that updates the cursor location.
func NewCursorMutator(loc CursorLocator) Mutator {
	return &cursorMutator{loc}
}

func (cpm *cursorMutator) Mutate(state *State) {
	state.cursor = cpm.loc.Locate(state)
}

func (cpm *cursorMutator) String() string {
	return fmt.Sprintf("MutateCursor(%s)", cpm.loc)
}

// scrollToCursorMutator updates the view origin so that the cursor is visible.
type scrollToCursorMutator struct{}

// NewScrollToCursorMutator returns a mutator that updates the view origin so that the cursor is visible.
func NewScrollToCursorMutator() Mutator {
	return &scrollToCursorMutator{}
}

func (sm *scrollToCursorMutator) Mutate(state *State) {
	state.view.origin = ScrollToCursor(state.cursor.position, state.tree, state.view.origin, state.view.width, state.view.height)
}

func (sm *scrollToCursorMutator) String() string {
	return "ScrollToCursor()"
}

// insertRuneMutator inserts a rune at the current cursor location.
type insertRuneMutator struct {
	r rune
}

// NewInsertRuneMutator returns a mutator that inserts a rune at the current cursor location.
func NewInsertRuneMutator(r rune) Mutator {
	return &insertRuneMutator{r}
}

func (irm *insertRuneMutator) Mutate(state *State) {
	startPos := state.cursor.position
	if err := state.tree.InsertAtPosition(startPos, irm.r); err != nil {
		// Invalid UTF-8 character; ignore it.
		return
	}

	state.cursor.position = startPos + 1
}

func (irm *insertRuneMutator) String() string {
	return fmt.Sprintf("InsertRune(%q)", irm.r)
}

// deleteMutator deletes characters from the cursor up to a location.
type deleteMutator struct {
	loc CursorLocator
}

func NewDeleteMutator(loc CursorLocator) Mutator {
	return &deleteMutator{loc}
}

// Mutate deletes characters from the cursor position up to (but not including) the position returned by the locator.
// It can delete either forwards or backwards from the cursor.
// The cursor position will be set to the start of the deleted region,
// which could be on a newline character or past the end of the text.
func (dm *deleteMutator) Mutate(state *State) {
	startPos := state.cursor.position
	deleteToPos := dm.loc.Locate(state).position

	if startPos < deleteToPos {
		dm.deleteCharacters(state.tree, startPos, deleteToPos-startPos)
		state.cursor = cursorState{position: startPos}
	} else if startPos > deleteToPos {
		dm.deleteCharacters(state.tree, deleteToPos, startPos-deleteToPos)
		state.cursor = cursorState{position: deleteToPos}
	}
}

func (dm *deleteMutator) deleteCharacters(tree *text.Tree, pos uint64, count uint64) {
	for i := uint64(0); i < count; i++ {
		tree.DeleteAtPosition(pos)
	}
}

func (dm *deleteMutator) String() string {
	return fmt.Sprintf("Delete(%s)", dm.loc)
}

// viewOriginZeroMutator sets the view origin to position zero.
type viewOriginZeroMutator struct{}

// NewViewOriginZeroMutator returns a mutator that sets the view origin to position zero.
func NewViewOriginZeroMutator() Mutator {
	return &viewOriginZeroMutator{}
}

func (vm *viewOriginZeroMutator) Mutate(state *State) {
	state.view.origin = 0
}

func (vm *viewOriginZeroMutator) String() string {
	return "SetViewOriginZero()"
}

// viewOriginAtCursorMutator sets the view origin to position of the cursor.
type viewOriginAtCursorMutator struct{}

// NewViewOriginAtCursorMutator returns a mutator that sets the view origin to the position of the cursor.
func NewViewOriginAtCursorMutator() Mutator {
	return &viewOriginAtCursorMutator{}
}

func (vm *viewOriginAtCursorMutator) Mutate(state *State) {
	state.view.origin = state.cursor.position
}

func (vm *viewOriginAtCursorMutator) String() string {
	return "SetViewOriginToCursor()"
}
