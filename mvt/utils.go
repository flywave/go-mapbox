package mvt

func appendAll(b ...[]byte) []byte {
	total := 0
	for _, i := range b {
		total += len(i)
	}
	pos := 0
	totalbytes := make([]byte, total)
	for _, i := range b {
		for _, byteval := range i {
			totalbytes[pos] = byteval
			pos += 1
		}
	}
	return totalbytes
}
