package colorextra

func ConvertUint8ToUint32Color(in uint8) uint32 {
	inUint32 := uint32(in)
	return (inUint32 * 256) + inUint32
}
