package utils

func IndexOf[V comparable](s []V, value V) int {
	for i, e := range s {
		if e == value {
			return i
		}
	}
	return -1
}
