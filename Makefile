GO=go
GOBUILD=$(GO) build
GOCLEAN=$(GO) clean
GOTEST=$(GO) test
GOGET=$(GO) get

GODEP=dep

#GOPATH=$(PWD)

CMD=tissueabuser
CMDDIR=./cmd

BINDIR=./bin

all: clean build

clean:
	$(GOCLEAN)
	- rm -f $(BINDIR)/$(CMD)

run:
	@echo "use \"run.sh\" to run commands"

build: init build-cmd

init:
	- mkdir $(BINDIR)
	@ echo $(GOPATH)

build-cmd:
	cd $(CMDDIR) && $(GOBUILD) -o $(CMD)
	mv $(CMDDIR)/$(CMD) $(BINDIR)/$(CMD)

deps:
	$(GODEP) ensure
