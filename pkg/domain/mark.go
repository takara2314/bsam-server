package domain

const FirstMarkNo = 1

func CalcNextMarkNo(wantMarkCounts int, passedMarkNo int) int {
	if passedMarkNo+1 > wantMarkCounts {
		return 1
	}
	return passedMarkNo + 1
}
