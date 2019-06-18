package utils

func IsTrue(bln ...bool) bool {
	return len(bln) < 1 || (len(bln) > 0 && bln[0] != false)
}
