multiblock_amd64.s: gen.py
	PYTHONIOENCODING=utf8 python $< 8 > $@
