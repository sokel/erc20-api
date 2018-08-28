vet:
	@echo "+ $@"
	@go tool vet $(shell ls -1 . | grep -v -e vendor)

fmt:
	@echo "+ $@"
	@test -z "$$(gofmt -s -l . 2>&1 | grep -v ^vendor/ | tee /dev/stderr)" || \
		(echo >&2 "+ please format Go code with 'gofmt -s'" && false)
