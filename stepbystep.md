1. Utilizando Ubuntu

2. Instalando Golang
> sudo apt-get update && sudo apt-get -y install golang-go
> wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
> sudo tar -xvf go1.21.0.linux-amd64.tar.gz
> sudo mv go /usr/local
> export GOROOT=/usr/local/go
> export GOPATH=$HOME/go
> export PATH=$GOPATH/bin:$GOROOT/bin:$PATH
> source ~/.profile
> go version
> go mod tidy
> go run .

