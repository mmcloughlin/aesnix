import sys
import random


DECL = '// func {name}(nr int, xk *uint32, dst, src *byte)'
TEXT = u'TEXT \u00b7{name}(SB),NOSPLIT,$0'
HEADER = (
u'''	MOVQ nr+0(FP), CX
	MOVQ xk+8(FP), AX
	MOVQ dst+16(FP), DX
	MOVQ src+24(FP), BX
	MOVUPS 0(AX), {reg_key}''')


def file_header():
    print '#include "textflag.h"'
    print


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
    for tmpl in [DECL, TEXT, HEADER]:
        print tmpl.format(**params)

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
    print


def rand_enc(inst='AESENC'):
    print '\t{inst} X{i}, X{j}'.format(
            inst=inst,
            i=random.randrange(16),
            j=random.randrange(16),
            )


def nomem(size):
    """
    nomem generates a function with the same AES-NI instructions required for
    encrypting size blocks, but without any memory accesses at all.
    """
    rounds = 10
    name = 'nomem{}'.format(size)
    print TEXT.format(name=name)
    for i in xrange(size):
        for r in xrange(rounds-1):
            rand_enc()
        rand_enc(inst='AESENCLAST')
    print '\tRET'
    print


GENERATORS = dict(
        multiblock=generate,
        nomem=nomem,
        )


def generate_file(sizes, method='multiblock'):
    generator = GENERATORS[method]
    file_header()
    for size in sizes:
        generator(size)


def main(args):
    method = args[1]
    sizes = map(int, args[2].split(','))
    generate_file(sizes, method=method)


if __name__ == '__main__':
    main(sys.argv)
