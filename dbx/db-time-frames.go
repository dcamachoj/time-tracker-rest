package dbx

// TimeFrames
type TimeFrames struct {
	Daily   TimeFrame `json:"daily,omitempty"`
	Weekly  TimeFrame `json:"weekly,omitempty"`
	Monthly TimeFrame `json:"monthly,omitempty"`
}

// TimeFramesEntity
type TimeFramesEntity struct {
	Daily   TimeFrameEntity
	Weekly  TimeFrameEntity
	Monthly TimeFrameEntity
}

func (e *TimeFramesEntity) PreScan(scanner EntityScanner) error {
	e.Daily.PreScan(scanner)
	e.Weekly.PreScan(scanner)
	e.Monthly.PreScan(scanner)
	return nil
}
func (e *TimeFramesEntity) PosScan(scanner EntityScanner) error {
	return nil
}
func (e *TimeFramesEntity) Result() interface{} {
	var res = &TimeFrames{}
	e.Daily.AsTimeFrame(&res.Daily)
	e.Weekly.AsTimeFrame(&res.Weekly)
	e.Monthly.AsTimeFrame(&res.Monthly)
	return res
}
func (e *TimeFramesEntity) AsTimeFrames(target *TimeFrames) {
	e.Daily.AsTimeFrame(&target.Daily)
	e.Weekly.AsTimeFrame(&target.Weekly)
	e.Monthly.AsTimeFrame(&target.Monthly)
}
func (e *TimeFramesEntity) Select(prefix string) []string {
	var res []string
	res = append(res, e.Daily.Select(prefix+"daily_")...)
	res = append(res, e.Weekly.Select(prefix+"weekly_")...)
	res = append(res, e.Monthly.Select(prefix+"monthly_")...)
	return res
}
