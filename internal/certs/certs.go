package certs

import _ "embed"

//go:embed server.key
var ServerKey []byte

//go:embed server.crt
var ServerCrt []byte

//go:embed client.key
var ClientKey []byte

//go:embed client.crt
var ClientCrt []byte

//go:embed ca.crt
var CaCrt []byte
