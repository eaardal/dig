package utils

func DeterministicItemForValue[T any](value string, list []T) T {
	hashValue := HashStringToUint32(value)
	colorIndex := int(hashValue) % len(list)
	return list[colorIndex]
}
