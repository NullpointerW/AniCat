NAME=anicat
BUILDIR=build

SRV_FILE=service-exec/service.go
VERSION := $(patsubst v%,%,$(shell git describe --tags || echo "unknown"))
GOBUILD=CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)"'


all: 
	$(MAKE) windows
	$(MAKE) linux
	$(MAKE) windows-service
	$(MAKE) docker-build
	$(MAKE) docker-push

windows:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)" -H=windowsgui' -o $(BUILDIR)/$(NAME)-$@-amd64.exe
	zip -j $(BUILDIR)/$(NAME)-$@-amd64.zip $(BUILDIR)/$(NAME)-$@-amd64.exe

windows-service:
	GOARCH=amd64 GOOS=windows CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)" ' -o $(BUILDIR)/$(NAME)-$@-amd64.exe $(SRV_FILE)
	zip -j $(BUILDIR)/$(NAME)-$@-amd64.zip $(BUILDIR)/$(NAME)-$@-amd64.exe

linux:
	GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BUILDIR)/$(NAME)-$@-amd64
	tar czvf $(BUILDIR)/$(NAME)-$@-amd64.tar.gz -C $(BUILDIR) $(NAME)-$@-amd64

docker-build:
	docker build --build-arg VER=$(VERSION) -t wmooon/anicat -f docker/Dockerfile .

docker-push:
	docker push  wmooon/anicat
        
clean:
	rm -rf $(BUILDIR)/*
