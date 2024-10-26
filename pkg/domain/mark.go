package domain

const FirstMarkNo = 1

func CalcNextMarkNo(wantMarkCounts int, passedMarkNo int) int {
	if passedMarkNo+1 > wantMarkCounts {
		return FirstMarkNo
	}
	return passedMarkNo + 1
}

func CalcPreviousMarkNo(wantMarkCounts int, passedMarkNo int) int {
	if passedMarkNo-1 < FirstMarkNo {
		return wantMarkCounts
	}
	return passedMarkNo - 1
}
