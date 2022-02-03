package dbx

import "database/sql"

// TimeFrame
type TimeFrame struct {
	Current  int64 `json:"current,omitempty"`
	Previous int64 `json:"previous,omitempty"`
}

// TimeFrameEntity
type TimeFrameEntity struct {
	Current  sql.NullInt64
	Previous sql.NullInt64
}

func (e *TimeFrameEntity) PreScan(scanner EntityScanner) error {
	scanner.AddScanners(
		&e.Current,
		&e.Previous,
	)
	return nil
}
func (e *TimeFrameEntity) PosScan(scanner EntityScanner) error {
	return nil
}
func (e *TimeFrameEntity) Result() interface{} {
	return &TimeFrame{
		Current:  e.Current.Int64,
		Previous: e.Previous.Int64,
	}
}
func (e *TimeFrameEntity) AsTimeFrame(target *TimeFrame) {
	target.Current = e.Current.Int64
	target.Previous = e.Previous.Int64
}
func (e *TimeFrameEntity) Select(prefix string) []string {
	return []string{
		prefix + "current",
		prefix + "previous",
	}
}
