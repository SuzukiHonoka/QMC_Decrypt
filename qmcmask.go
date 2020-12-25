package main

import (
	"errors"
	"fmt"
)

type QmcMask struct {
	superA []byte
	superB []byte
	super58A byte
	super58B byte
	//
	mapping map[float32][]int
	//
	mask128to44 map[int]int
	//
	matrix []byte
	matrix44 []byte
	matrix58 []byte
	matrix128 []byte

}

func (x QmcMask) QMC() QmcMask{
	fmt.Println("MatrixLen:", len(x.matrix))
	if x.superA == nil && x.superB == nil {
		x.matrix128 = x.matrix
		x.generateMask44from128()
		x.generateMask58from128()
	}
	return x
}

func (x QmcMask) Decrypt(data []byte)[]byte{
	index,maskIdx := -1,-1
	for cur := 0; cur < len(data); cur++{
		index++
		maskIdx++
		if index == 0x8000 || (index > 0x8000 && (index + 1) % 0x8000 == 0) {
			index++
			maskIdx++
		}
		if maskIdx >= 128 {
			maskIdx -=128
		}
		data[cur] ^= x.matrix128[maskIdx]
	}
	return data
}

func (x QmcMask) generateMask58from128(){
	if len(x.matrix128) != 128{
		errors.New("incorrect mask128 length")
	}
	superA := x.matrix128[0]
	superB := x.matrix128[8]

	for rowIdx := 0; rowIdx < 8; rowIdx+=1{
		lenStart := 16 * rowIdx
		lenRightStart := 120 - lenStart
		if x.matrix128[lenStart] != superA || x.matrix128[lenStart + 8] != superB {
			errors.New("decode mask-128 to mask-58 failed")
		}
		rowLeft := x.matrix128[lenStart + 1:lenStart + 8]
		rowRight := revers(x.matrix128[lenRightStart + 1:lenRightStart + 8])

		if IsBytesEqual(rowLeft,rowRight){
			x.matrix58 = append(x.matrix58,rowLeft... )
		}else {
			errors.New("decode mask-128 to mask-58 failed")
		}
	}
	fmt.Println("decode mask-128 to mask-58 succeed")
	x.super58A = superA
	x.super58B = superB
}

func (x QmcMask) generateMask44from128(){
	if len(x.matrix) != 128 {
		errors.New("incorrect mask128 matrix length")
	}
	x.getMapping()
	idxI44 := 0
	for k := range x.mapping{
		it256Len := len(x.mapping[k])
		for i := 1; i < it256Len; i++{
			if x.matrix128[x.mapping[k][0]] != x.matrix128[x.mapping[k][i]] {
				errors.New("decode mask-128 to mask-44 failed")
			}
		}
		x.matrix44[idxI44] = x.matrix128[x.mapping[k][0]]
		idxI44++
	}
	fmt.Println("decode mask-128 to mask-44 succeed")
}

func (x QmcMask) getMapping() map[float32][]int {
	x.mapping = map[float32][]int{}
	x.mask128to44 = map[int]int{}
	for i := 0; i < 128; i++{
		realIdx := float32((i*i + 27) % 256)
		if _,exist := x.mapping[realIdx]; exist{
			x.mapping[realIdx] = append(x.mapping[realIdx], i)
		}else {
			x.mapping[realIdx] = []int{i}
		}
	}
	idx44 := 0
	for k := range x.mapping{
		for kk := range x.mapping[k]{
			x.mask128to44[kk] = idx44
		}
		idx44++
	}
	return x.mapping
}