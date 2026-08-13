package main

import (
	"fmt"
	"os"
	"math"
)

var data []byte
var out []byte

func zip() {
	HH := byte(len(data) >> 24)
	HL := byte(len(data) >> 16)
	LH := byte(len(data) >> 8)
	LL := byte(len(data) & 0xFF)
	out = append(out, []byte{
		HH, HL, LH, LL,
	}...)

	_max := len(data) - 1

	pos := 0
	exit := false

	for {
		current := data[pos]
		count := uint32(0)

		for j := pos; j < len(data); j++ {
			now := data[j]

			if now == current {
				count++
			}

			if now != current || count == math.MaxUint32 || j == _max {
				HH := byte(count >> 24)
				HL := byte(count >> 16)
				LH := byte(count >> 8)
				LL := byte(count & 0xFF)
				out = append(out, []byte{
					HH, HL, LH, LL, current,
				}...)

				if now != current {
					pos = j
				} else {
					j++
					pos = j
				}

				if j >= _max {
					exit = true
				}
				break
			}
		}

		if exit == true {
			break
		}
	}

	if len(out) >= len(data) {
		fmt.Println("fzip: warning: no net compression")
	}
	os.WriteFile(os.Args[3], out, 0644)
}

func unzip() {
	pos := uint32(0)
	size := uint32(0)
	_max := len(data) - 1

	ReadUint32 := func() uint32 {
		n := uint32(data[pos]) << 24 | uint32(data[pos + 1]) << 16 | uint32(data[pos + 2]) << 8 | uint32(data[pos + 3])
		return n
	}
	
	size = ReadUint32()
	pos += 4

	for {
		if int(pos) >= _max {
			break
		}

		n := ReadUint32()
		pos += 4
		b := data[pos]
		pos++

		for i := 0; i < int(n); i++ {
			out = append(out, b)
		}
	}

	if uint32(len(out)) != size {
		fmt.Println("fzip: warning: unzip fault")
	} 

	os.WriteFile(os.Args[3], out, 0644)
}

func main() {
	UsageMsg := "Usage: fzip <input file> <zip/unzip> <output file>"
	if len(os.Args) < 4 {
		fmt.Println(UsageMsg)
		os.Exit(1)
	}

	filename := os.Args[1]

	var err error
	data, err = os.ReadFile(filename)
	
	if err != nil {
		fmt.Println("fzip: could not open file '" + filename + "':", err)
		os.Exit(1)
	}


	switch os.Args[2] {
	case "zip":
		zip()	
	case "unzip":
		unzip()	
	default:
		fmt.Println(UsageMsg)
		os.Exit(1)	
	}
}
