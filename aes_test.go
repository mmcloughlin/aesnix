package aesnix

import (
	"bytes"
	"strconv"
	"testing"
)

type Encryptor func(nr int, xk *uint32, dst, src *byte)

var v = struct {
	Rounds int
	Key    []byte
	Plain  []byte
	Cipher []byte
}{
	Rounds: 10,
	Key: []byte{
		0x2b, 0x7e, 0x15, 0x16, 0x28, 0xae, 0xd2, 0xa6,
		0xab, 0xf7, 0x15, 0x88, 0x09, 0xcf, 0x4f, 0x3c,
	},
	Plain: []byte{
		0x6b, 0xc1, 0xbe, 0xe2, 0x2e, 0x40, 0x9f, 0x96,
		0xe9, 0x3d, 0x7e, 0x11, 0x73, 0x93, 0x17, 0x2a,
	},
	Cipher: []byte{
		0x3a, 0xd7, 0x7b, 0xb4, 0x0d, 0x7a, 0x36, 0x60,
		0xa8, 0x9e, 0xca, 0xf3, 0x24, 0x66, 0xef, 0x97,
	},
}

func TestSingle(t *testing.T) {
	cipher := make([]byte, 16)
	enc := make([]uint32, 4*(v.Rounds+1))
	dec := make([]uint32, 4*(v.Rounds+1))
	expandKeyAsm(v.Rounds, &v.Key[0], &enc[0], &dec[0])
	encryptBlockAsm(v.Rounds, &enc[0], &cipher[0], &v.Plain[0])
	if !bytes.Equal(cipher, v.Cipher) {
		t.Fatal("encryptBlockAsm failed")
	}
}

func TestMulti(t *testing.T) {
	const Blocks = 8

	enc := make([]uint32, 4*(v.Rounds+1))
	dec := make([]uint32, 4*(v.Rounds+1))
	expandKeyAsm(v.Rounds, &v.Key[0], &enc[0], &dec[0])

	plain := make([]byte, 16*Blocks)
	for i := 0; i < Blocks; i++ {
		copy(plain[16*i:], v.Plain)
	}

	cipher := make([]byte, 16*Blocks)
	encryptBlocks8Asm(v.Rounds, &enc[0], &cipher[0], &plain[0])

	for i := 0; i < Blocks; i++ {
		if !bytes.Equal(cipher[16*i:16*i+16], v.Cipher) {
			t.Errorf("error on block %d", i)
		}
	}
}

func BenchmarkSingle(b *testing.B) {
	cipher := make([]byte, 16)
	enc := make([]uint32, 4*(v.Rounds+1))
	dec := make([]uint32, 4*(v.Rounds+1))
	expandKeyAsm(v.Rounds, &v.Key[0], &enc[0], &dec[0])

	b.SetBytes(16)
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		encryptBlockAsm(v.Rounds, &enc[0], &cipher[0], &v.Plain[0])
	}
}

func BenchmarkMulti(b *testing.B) {
	cases := []struct {
		Encryptor Encryptor
		Blocks    int
	}{
		{encryptBlocks2Asm, 2},
		{encryptBlocks4Asm, 4},
		{encryptBlocks6Asm, 6},
		{encryptBlocks8Asm, 8},
		{encryptBlocks10Asm, 10},
		{encryptBlocks12Asm, 12},
		{encryptBlocks14Asm, 14},
	}
	for _, c := range cases {
		b.Run(strconv.Itoa(c.Blocks), func(b *testing.B) {
			EncryptorBenchmark(b, c.Encryptor, c.Blocks)
		})
	}
}

func EncryptorBenchmark(b *testing.B, f Encryptor, blocks int) {
	enc := make([]uint32, 4*(v.Rounds+1))
	dec := make([]uint32, 4*(v.Rounds+1))
	expandKeyAsm(v.Rounds, &v.Key[0], &enc[0], &dec[0])

	plain := make([]byte, 16*blocks)
	for i := 0; i < blocks; i++ {
		copy(plain[16*i:], v.Plain)
	}
	cipher := make([]byte, 16*blocks)

	b.SetBytes(16 * int64(blocks))
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		f(v.Rounds, &enc[0], &cipher[0], &plain[0])
	}
}

func BenchmarkNomem(b *testing.B) {
	cases := []struct {
		F      func()
		Blocks int
	}{
		{nomem2, 2},
		{nomem4, 4},
		{nomem6, 6},
		{nomem8, 8},
		{nomem10, 10},
		{nomem12, 12},
		{nomem14, 14},
	}
	for _, c := range cases {
		b.Run(strconv.Itoa(c.Blocks), func(b *testing.B) {
			b.SetBytes(16 * int64(c.Blocks))
			for i := 0; i < b.N; i++ {
				c.F()
			}
		})
	}
}
