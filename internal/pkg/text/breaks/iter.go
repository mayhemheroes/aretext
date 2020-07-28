package breaks

// BreakIter iterates through valid breakpoints in a string.
type BreakIter interface {
	// NextBreak returns the next valid breakpoint.
	// The breakpoint is an offset from the beginning of the string the iterator has processed.
	// If there are no more breakpoints, NextBreak returns an io.EOF error.
	NextBreak() (uint64, error)
}
