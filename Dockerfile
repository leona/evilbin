FROM golang:1.21-bullseye

WORKDIR /app

ENV GOPATH="/root/go"
ENV PATH="$PATH:$GOPATH/bin"
RUN go install -v golang.org/x/tools/gopls@latest
RUN go install -v golang.org/x/tools/cmd/goimports@latest
RUN go install -v github.com/rogpeppe/godef@latest
RUN go install -v github.com/stamblerre/gocode@latest