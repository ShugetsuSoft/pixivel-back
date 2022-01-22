package utils

func UniqueArray(m []interface{}) []interface{} {
	d := make([]interface{}, 0)
	tempMap := make(map[interface{}]bool, len(m))
	for _, v := range m {
		if tempMap[v] == false {
			tempMap[v] = true
			d = append(d, v)
		}
	}
	return d
}
