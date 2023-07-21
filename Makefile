NAME=anicat
BUILDIR=build
VERSION := $(patsubst v%,%,$(shell git describe --tags || echo "x.x.x"))
GOBUILD=CGO_ENABLED=0 go build  -ldflags '-X "github.com/NullpointerW/anicat/conf.Ver=$(VERSION)"'
		
all: 
        $(MAKE) windows
        $(MAKE) linux

windows:
        GOARCH=amd64 GOOS=windows $(GOBUILD) -o $(BUILDIR)/$(NAME)-$@-amd64.exe
		zip -j $(BUILDIR)/$(NAME)-$@-amd64.zip $(BUILDIR)/$(NAME)-$@-amd64.exe

linux:
        GOARCH=amd64 GOOS=linux $(GOBUILD) -o $(BUILDIR)/$(NAME)-$@-amd64
		tar czvf $(BUILDIR)/$(NAME)-$@-amd64.tar.gz -C $(BUILDIR) $(NAME)-$@-amd64

docker-build:
        docker build -t wmooon/anict -f Docker/Dockerfile .

docker-push:
        docker push  wmooon/anict
		
clean:
	rm -rf $(BUILDIR)/*