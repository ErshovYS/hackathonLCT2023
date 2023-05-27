package string_generator

type PageMarkup struct {
	HeaderStart int
	BodyStart   int

	PageStart int
	PageEnd   int

	RowNumber int

	ColumnStart int

	HeaderCurrent   int
	BodyCurrent     int
	BodyCurrentLeft int

	HeaderStep  int
	BodyStep    int
	BodyRowStep int

	BodyStringNum int

	CurNumStep int
	OldNumStep int
}

func initPageMarkup(bodyStart int, rowNumber int, pageEnd int) *PageMarkup {
	val := PageMarkup{
		HeaderStart:     200,
		BodyStart:       bodyStart,
		PageEnd:         pageEnd,
		RowNumber:       rowNumber,
		ColumnStart:     505,
		HeaderCurrent:   200,
		BodyCurrent:     bodyStart,
		BodyCurrentLeft: bodyStart,
		HeaderStep:      20,
		BodyStep:        20,
		BodyRowStep:     50,
		BodyStringNum:   3,
		CurNumStep:      0,
		OldNumStep:      0,
	}

	val.BodyRowStep = (val.PageEnd - val.BodyStart) / val.RowNumber
	if val.BodyRowStep > 70 {
		val.BodyRowStep = 70
	}
	return &val
}
func GetHeaderNextPositionString(val *PageMarkup) string {
	val.HeaderCurrent = val.HeaderCurrent + val.HeaderStep
	return getString(val.HeaderCurrent)
}

func getBodyNextPositionRow(val *PageMarkup) string {
	val.BodyCurrent = val.BodyCurrent + val.BodyRowStep
	val.OldNumStep = val.CurNumStep
	val.CurNumStep = 0
	return getString(val.BodyCurrent)
}

func getImageLinePosition(val *PageMarkup, delta int) string {
	return getString(val.BodyCurrent + (val.BodyRowStep+delta)/2)
}

func getBodyNextPosition(val *PageMarkup) string {
	val.BodyCurrent = val.BodyCurrent + val.BodyStep
	val.CurNumStep++
	return getString(val.BodyCurrent)
}

func getBodyNextLeftPositionRow(val *PageMarkup) string {
	val.BodyCurrentLeft = val.BodyCurrentLeft + val.BodyRowStep + ((val.CurNumStep+val.OldNumStep)*val.BodyStep)/2
	return getString(val.BodyCurrentLeft)
}

type TextMarkup struct {
	StartPos      int
	CurPos        int
	Step          int
	ParagraphStep int
}

func initTextMarkup(startPos int) *TextMarkup {
	return &TextMarkup{
		StartPos:      startPos,
		CurPos:        startPos,
		Step:          30,
		ParagraphStep: 50,
	}
}

func initTextMarkup2(startPos, step int) *TextMarkup {
	return &TextMarkup{
		StartPos:      startPos,
		CurPos:        startPos,
		Step:          step,
		ParagraphStep: 50,
	}
}

func getNextRow(val *TextMarkup) string {
	val.CurPos = val.CurPos + val.Step
	return getString(val.CurPos)
}

func getNextParagraph(val *TextMarkup) string {
	val.CurPos = val.CurPos + val.ParagraphStep
	return getString(val.CurPos)
}
