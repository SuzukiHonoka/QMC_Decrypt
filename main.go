package main

import (
	"encoding/binary"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
)

const (
	root = "C:\\Users\\a2540\\Music"
)

func check(err error)  {
	if err != nil {
		panic(err)
	}
}

func supported(ext string)bool{
	switch ext {
	case ".mflac",".mgg":
		return true
	default:
		return false
	}
}

func main() {
	// scan files
	var files []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if info.IsDir() || !supported(filepath.Ext(path)) {
			return nil
		}
		files = append(files, path)
		return nil
	})
	check(err)
	// scan done
	for _,path := range files{
		fmt.Println("File:",path)
		data,err := ioutil.ReadFile(path)
		check(err)
		keyLen := binary.LittleEndian.Uint32(data[len(data)-4:])
		keyPos := len(data) - 4 - int(keyLen)
		fmt.Println("datalen:",len(data),"keyLen:", keyLen,"keyPos:", keyPos)
		audioData := data[:keyPos]
		//keyData := data[keyPos:keyPos + int(keyLen)]
		seed := QmcMaskDetectMflac(audioData).QMC()
		musicDecoded := seed.Decrypt(audioData)
		//fmt.Printf("%+v\n",seed)
		ioutil.WriteFile("1.flac",musicDecoded,os.ModePerm)
	}
}
