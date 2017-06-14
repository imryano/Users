FROM golang:latest

ADD . /go/src/github.com/imryano/Users
#ADD . /go/src/github.com/dgrijalva/jwt-go
#ADD . /go/src/github.com/imryano/Users/userPackage
#ADD . /go/src/github.com/imryano/utils/password
#ADD . /go/src/github.com/imryano/utils/webservice
#ADD . /go/src/gopkg.in/mgo.v2
#ADD . /go/src/gopkg.in/mgo.v2/bson

RUN go get ./src/github.com/imryano/Users/userRead

RUN go build ./src/github.com/imryano/Users/userRead
CMD ["go/src/github.com/imryano/Users/userRead/main"]

EXPOSE 8080