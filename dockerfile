FROM golang:latest 
WORKDIR $GOPATH/src/github.com/crux55/rest_server
COPY . . 
RUN make build
EXPOSE 8000 
RUN go get -d -v ./...
RUN go install -v ./...
CMD ["rest_server"]
