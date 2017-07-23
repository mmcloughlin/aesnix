all: asm_amd64.s multiblock_amd64.s RESULTS.md

RESULTS.md:
	echo '```' > $@
	go version >> $@
	go test -v -bench . -benchtime 10s | tee | grep -v Users >> $@
	echo '```' >> $@

asm_amd64.s:
	wget 'https://raw.githubusercontent.com/golang/go/b81735924936291303559fd71dabaa1aa88f57c5/src/crypto/aes/asm_amd64.s'

multiblock_amd64.s: gen.py
	PYTHONIOENCODING=utf8 python $< 8 > $@
