package main

import (
	"fmt"
	"os"
	"math"
)

var data []byte

func zip(granularity int) (uint32, *[]byte) {
	var out []byte

	HH := byte(len(data) >> 24)
	HL := byte(len(data) >> 16)
	LH := byte(len(data) >> 8)
	LL := byte(len(data) & 0xFF)
	out = append(out, []byte{
		byte('F'), byte('Z'), byte('I'), byte('P'), byte(granularity), HH, HL, LH, LL,
	}...)

	_max := len(data) - 1
	countmax := 0
	switch granularity {
	case 1:
		countmax = math.MaxUint8
	default:
		countmax = math.MaxUint32
	}

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

			if now != current || int(count) == countmax || j == _max {
				switch granularity {
				case 1:
					C := byte(count & 0xFF)
					out = append(out, []byte{
						C, current,
					}...)
				default:
					HH := byte(count >> 24)
					HL := byte(count >> 16)
					LH := byte(count >> 8)
					LL := byte(count & 0xFF)
					out = append(out, []byte{
						HH, HL, LH, LL, current,
					}...)
				}
				

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

	return uint32(len(out)), &out
}

func unzip() {
	pos := uint32(0)
	granularity := uint32(4)
	_max := len(data) - 1
	var out []byte

	ReadUint32 := func() uint32 {
		n := uint32(data[pos]) << 24 | uint32(data[pos + 1]) << 16 | uint32(data[pos + 2]) << 8 | uint32(data[pos + 3])
		return n
	}

	ReadUint8 := func() uint32 {
		n := uint32(data[pos])
		return n
	}

	if (ReadUint32() != 0x465A4950) {
		fmt.Println("fzip: invalid fzip file")
		os.Exit(1)
	}
	pos += 4 // skip over header

	granularity = ReadUint8()
	pos++

	pos += 4 // skip over size

	for {
		if int(pos) > _max {
			break
		}
	
		var n uint32

		switch granularity {
		case 1:
			n = ReadUint8()
			pos++	
		default:
			n = ReadUint32()
			pos += 4
		}
		
		b := data[pos]
		pos++

		for i := 0; i < int(n); i++ {
			out = append(out, b)
		}
	}

	fmt.Printf("fzip: inflated from %d bytes to %d bytes\n", len(data), len(out)) 

	os.WriteFile(os.Args[3], out, 0644)
}

func nn_average() {
	for i := 0; i < len(data); i++ {
		var Occurrences = make(map[byte]int)

		smin := i - 16
		smax := i + 15

		if smin < 0 {
			smin = 0
		}
		if smax > len(data) - 1 {
			smax = len(data) - 1
		}

		subset := data[smin:smax]
		for _, b := range subset {
			Occurrences[b]++
		}

		b := data[i]
		around := []byte {}
		ok := false

		if b <= 253 {
			around = append(around, b + 2)
			ok = true
		}
		if b <= 254 {
			around = append(around, b + 1)
			ok = true
		}
		if b >= 2 {
			around = append(around, b - 2)
			ok = true
		}
		if b >= 1 {
			around = append(around, b - 1)
			ok = true
		}
		around = append(around, b)

		if ok != true {
			break
		}

		_max := 0
		_byte := byte(0)

		for _, a := range around {
			if Occurrences[a] > _max {
				_max = Occurrences[a]
				_byte = a
			}
		}

		if _max == 0 {
			continue
		}

		data[i] = _byte
	}
}

func main() {
	UsageMsg := "Usage: fzip <input file> <zip/unzip/noop> <output file> [-nn (enables nearest-neighbor compression, lossy)]"
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

	for i := 4; i < len(os.Args); i++ {
		arg := os.Args[i]
		switch arg {
		case "-nn":
			nn_average()
		}
	}

	switch os.Args[2] {
	case "zip":	
		size, buf := zip(1)
		size2, buf2 := zip(4)

		var c_size uint32
		var c_buf *[]byte
	
		if size < size2 {
			c_size = size
			c_buf = buf
		} else {
			c_size = size2
			c_buf = buf2	
		}

		fmt.Printf("fzip: deflated from %d bytes to %d bytes\n", len(data), c_size)

		os.WriteFile(os.Args[3], *c_buf, 0644)
	case "unzip":
		unzip()
	case "noop":
		os.WriteFile(os.Args[3], data, 0644)
	default:
		fmt.Println(UsageMsg)
		os.Exit(1)	
	}
}
