package utils

func Max(a, b int) int {
	if a < b {
		return b
	} else {
		return a
	}
}

func Min(a, b int) int {
	if a > b {
		return b
	} else {
		return a
	}
}

func PosMod(a, b int) int {
	return (a%b + b) % b
}
