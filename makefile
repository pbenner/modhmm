
VERSION = 1.0
FILES   = modhmm.go $(filter-out %_test.go modhmm.go,$(wildcard *go))

# ------------------------------------------------------------------------------

all: modhmm

install: modhmm
	go install

modhmm: $(wildcard *.go)
	go build -ldflags "\
	   -X main.Version=$(VERSION) \
	   -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` \
	   -X main.GitHash=`git rev-parse HEAD`" \
	   $(FILES)

# ------------------------------------------------------------------------------

.PHONY: all install
