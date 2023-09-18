NAME=anicat
BUILDIR=build
BUILDIR_CLI=build/cli
SRV_FILE=service-exec/service.go
VERSION := $(patsubst v%,%,$(shell git describe --tags || echo "x.x.x"))
GOBUILD=CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)"'
CLI_FILE=net/client/cli.go

all: 
	$(MAKE) windows
	$(MAKE) linux
	$(MAKE) windows-cli
	$(MAKE) linux-cli

windows:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)" -H=windowsgui' -o $(BUILDIR)/$(NAME)-$@-amd64.exe
	zip -j $(BUILDIR)/$(NAME)-$@-amd64.zip $(BUILDIR)/$(NAME)-$@-amd64.exe

windows-cli:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build -o $(BUILDIR_CLI)/$(NAME)-$@-amd64.exe $(CLI_FILE)

windows-service:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)" ' -o $(BUILDIR)/$(NAME)-$@-amd64.exe $(SRV_FILE)
	zip -j $(BUILDIR)/$(NAME)-$@-amd64.zip $(BUILDIR)/$(NAME)-$@-amd64.exe

linux:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BUILDIR)/$(NAME)-$@-amd64
	tar czvf $(BUILDIR)/$(NAME)-$@-amd64.tar.gz -C $(BUILDIR) $(NAME)-$@-amd64

linux-cli:
	GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go build -o $(BUILDIR_CLI)/$(NAME)-$@-amd64 $(CLI_FILE)

docker-build:
	docker build --build-arg VER=$(VERSION) -t wmooon/anicat -f docker/Dockerfile .

docker-push:
	docker push  wmooon/anicat
        
clean:
	rm -rf $(BUILDIR)/*
