import sys


HEADER = (
u'''
// func {name}(nr int, xk *uint32, dst, src *byte)
TEXT \u00b7{name}(SB),NOSPLIT,$0
	MOVQ nr+0(FP), CX
	MOVQ xk+8(FP), AX
	MOVQ dst+16(FP), DX
	MOVQ src+24(FP), BX
	MOVUPS 0(AX), {reg_key}''')

def file_header():
    print '#include "textflag.h"'


def func_name(n):
    """
    Name of function for n blocks.
    """
    if n == 1:
        return 'encryptBlockAsm'
    return 'encryptBlocks{}Asm'.format(n)


def generate(n):
    """
    Generate Go assembly for encrypting n blocks at once with one key.
    """
    params = {
        'name': func_name(n),
        'reg_key': 'X{}'.format(n),
        }

    # header
    print HEADER.format(**params)

    # load plain
    for i in xrange(n):
        print '\tMOVUPS {offset}(BX), X{i}'.format(offset=16*i, i=i)

    # initial key add
    print '\tADDQ $16, AX'
    for i in xrange(n):
        print '\tPXOR {reg_key}, X{i}'.format(i=i, **params)

    # num rounds branching
    print '\tSUBQ $12, CX'
    print '\tJE Lenc192'
    print '\tJB Lenc128'

    def enc(ax, inst='AESENC'):
        print '\tMOVUPS {offset}(AX), {reg_key}'.format(offset=16*ax, **params)
        for i in xrange(n):
            print '\t{inst} {reg_key}, X{i}'.format(inst=inst, i=i, **params)

    # 2 extra rounds for 256-bit keys
    print 'Lenc256:'
    enc(0)
    enc(1)
    print '\tADDQ $32, AX'

    # 2 extra rounds for 192-bit keys
    print 'Lenc192:'
    enc(0)
    enc(1)
    print '\tADDQ $32, AX'

    # 10 rounds for 128-bit (with special handling for final)
    print 'Lenc128:'
    for r in xrange(9):
        enc(r)
    enc(9, inst='AESENCLAST')

    # write results to destination
    for i in xrange(n):
        print '\tMOVUPS X{i}, {offset}(DX)'.format(i=i, offset=16*i)

    # return
    print '\tRET'


def generate_file(sizes):
    file_header()
    for size in sizes:
        generate(size)


def main(args):
    sizes = map(int, args[1].split(','))
    generate_file(sizes)


if __name__ == '__main__':
    main(sys.argv)
