
VERSION = 1.0.0
FILES   = modhmm.go $(filter-out %_test.go modhmm.go,$(wildcard *go))
GOBIN   = $(shell echo $${GOPATH}/bin)

# ------------------------------------------------------------------------------

all: modhmm

modhmm: $(wildcard *.go)
	go build -ldflags "\
	   -X main.Version=$(VERSION) \
	   -X main.BuildTime=`TZ=UTC date -u '+%Y-%m-%dT%H:%M:%SZ'` \
	   -X main.GitHash=`git rev-parse HEAD`" \
	   $(FILES)

install: modhmm | $(GOBIN)
ifeq ($(GOBIN),/bin)
	install modhmm $HOME/go/bin
else
	install modhmm $(GOBIN)
endif

$(GOBIN):
	mkdir -p $(GOBIN)

# ------------------------------------------------------------------------------

.PHONY: all install
