package aesasm

func expandKeyAsm(nr int, key *byte, enc *uint32, dec *uint32)
func encryptBlockAsm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks2Asm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks4Asm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks6Asm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks8Asm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks10Asm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks12Asm(nr int, xk *uint32, dst, src *byte)
func encryptBlocks14Asm(nr int, xk *uint32, dst, src *byte)
