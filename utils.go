package main

var FlacHeader = []byte{0x66, 0x4C, 0x61, 0x43}

func min(a int,b int)int{
	if a < b{
		return a
	}
	return b
}

func revers(s []byte) []byte {
	a := make([]byte,len(s))
	copy(a,s)
	for i := len(a)/2-1; i >= 0; i-- {
		opp := len(a)-1-i
		a[i], a[opp] = a[opp], a[i]
	}
	return a
}

func IsBytesEqual(first []byte, second []byte)bool{
	for k := range first{
		if first[k] != second[k]{
			return false
		}
	}
	return true
}