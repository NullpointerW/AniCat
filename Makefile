NAME=anicat
BUILDIR=build
BUILDIR_CLI=build/cli
VERSION := $(patsubst v%,%,$(shell git describe --tags || echo "x.x.x"))
GOBUILD=CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)"'
CLI_FILE=net/client/cli.go

all: 
	$(MAKE) windows
	$(MAKE) linux
	$(MAKE) windows-cli
	$(MAKE) linux-cli

windows:
	GOARCH=amd64 GOOS=windows $(GOBUILD) -ldflags '-H=windowsgui' -o $(BUILDIR)/$(NAME)-$@-amd64.exe
	zip -j $(BUILDIR)/$(NAME)-$@-amd64.zip $(BUILDIR)/$(NAME)-$@-amd64.exe

windows-cli:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -o $(BUILDIR_CLI)/$(NAME)-$@-amd64.exe $(CLI_FILE)

linux:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BUILDIR)/$(NAME)-$@-amd64
	tar czvf $(BUILDIR)/$(NAME)-$@-amd64.tar.gz -C $(BUILDIR) $(NAME)-$@-amd64

linux-cli:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(BUILDIR_CLI)/$(NAME)-$@-amd64 $(CLI_FILE)

docker-build:
	docker build -t wmooon/anict -f Docker/Dockerfile .

docker-push:
	docker push  wmooon/anict
        
clean:
	rm -rf $(BUILDIR)/*
