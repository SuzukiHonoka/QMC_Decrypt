package main

import "fmt"

func QmcMaskDetectMflac(data []byte) QmcMask{
	var mask QmcMask
	search_len := min(0x8000, len(data))
	for blockIdx := 0; blockIdx < search_len; blockIdx += 128 {
		mask = QmcMask{matrix: data[blockIdx : blockIdx+ 128]}.QMC()
		if IsBytesEqual(FlacHeader, mask.Decrypt(data[0:len(FlacHeader)])){
			fmt.Println("FLAC DETECTED!")
			break
		}
	}
	return mask
}

