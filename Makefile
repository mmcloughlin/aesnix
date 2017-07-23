all: asm_amd64.s multiblock_amd64.s nomem_amd64.s

RESULTS.md: aes_test.go
	echo '```' > $@
	go version >> $@
	go test -v -bench . -benchtime 5s | tee | grep -v Users >> $@
	echo '```' >> $@

asm_amd64.s:
	wget 'https://raw.githubusercontent.com/golang/go/b81735924936291303559fd71dabaa1aa88f57c5/src/crypto/aes/asm_amd64.s'

%_amd64.s: gen.py
	PYTHONIOENCODING=utf8 python $< $* 2,4,6,8,10,12,14 > $@
