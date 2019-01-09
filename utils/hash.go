package utils

// APHash算法
//unsigned int APHash(char* str, unsigned int len)
//{
//    unsigned int hash = 0xAAAAAAAA;
//    unsigned int i    = 0;
//    for(i = 0; i < len; str++, i++)
//    {
//        hash ^= ((i & 1) == 0) ? ( (hash << 7) ^ (*str) * (hash >> 3)) :
//            (~((hash << 11) + (*str) ^ (hash >> 5)));
//    }
//    return hash;
//}
func APHash(data []byte) uint32 {
	var hash uint32 = 0xAAAAAAAA

	for i := 0; i < len(data); i++ {
		by := uint32(data[i]) & 0xff
		if (i & 1) == 0 {
			hash ^= (hash << 7) ^ by*(hash>>3)
		} else {
			hash ^= ^((hash << 11) + by ^ (hash >> 5))
		}
	}

	return hash
}
