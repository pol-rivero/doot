
codegen: lib/cache/Colfer.go

lib/cache/Colfer.go: lib/cache/cache.colf
	bin/colf -b lib Go lib/cache/cache.colf
