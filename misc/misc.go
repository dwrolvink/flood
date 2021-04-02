package misc

func AddUint8(a, b uint8) uint8 {
	var max = 255 - a
	if b > max {
		return 255
	}
	return a + b
}

func Max255Int (a, b int) int {
	s := a + b 
	if s > 255 { return 255 }
	return s
}

func Normalize(a, max, min int) int{
	if a > max {
		return max
	}
	if a < min {
		return min
	}
	return a
}

func ConvertIntToUint8(a int) uint8 {
	if a < 0 { a = 0}
	if a > 255 { a = 255}
	return uint8(a)
}

func SubtractUint8(a, b uint8) uint8 {
	if b > a {
		return 0
	}
	return a - b
}

func AbsInt(a, b int) int {
	if a < b {
		return b - a
	}
	return a - b
}