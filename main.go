package main

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"log"
)

func expandKeyAsm(nr int, key *byte, enc *uint32, dec *uint32)
func encryptBlockAsm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks8Asm(nr int, xk *uint32, dst, src *byte)

type TestVector struct {
	Rounds int
	Key    []byte
	Plain  []byte
	Cipher []byte
}

func VerifySingle(v TestVector) {
	cipher := make([]byte, 16)
	enc := make([]uint32, 4*(v.Rounds+1))
	dec := make([]uint32, 4*(v.Rounds+1))
	expandKeyAsm(v.Rounds, &v.Key[0], &enc[0], &dec[0])
	encryptBlockAsm(v.Rounds, &enc[0], &cipher[0], &v.Plain[0])

	if !bytes.Equal(cipher, v.Cipher) {
		log.Fatal("FAIL")
	}
	log.Print("pass single")
}

func VerifyMulti(v TestVector) {
	const Blocks = 8

	enc := make([]uint32, 4*(v.Rounds+1))
	dec := make([]uint32, 4*(v.Rounds+1))
	expandKeyAsm(v.Rounds, &v.Key[0], &enc[0], &dec[0])

	plain := make([]byte, 16*Blocks)
	for i := 0; i < Blocks; i++ {
		copy(plain[16*i:], v.Plain)
	}
	fmt.Print(hex.Dump(plain))

	cipher := make([]byte, 16*Blocks)
	encryptBlocks8Asm(v.Rounds, &enc[0], &cipher[0], &plain[0])

	fmt.Print(hex.Dump(cipher))
}

func main() {
	v := TestVector{
		Rounds: 10,
		Key:    []byte{0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6, 0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c},
		Plain:  []byte{0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96, 0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a},
		Cipher: []byte{0x3a, 0xd7, 0x7b, 0xb4, 0x0d, 0x7a, 0x36, 0x60, 0xa8, 0x9e, 0xca, 0xf3, 0x24, 0x66, 0xef, 0x97},
	}

	VerifySingle(v)
	VerifyMulti(v)
}
