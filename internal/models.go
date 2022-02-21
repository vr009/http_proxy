package internal

type Req struct {
	FullMsg []byte
	Secure  bool
	Host    string
	Port    string
}

type Resp struct {
	FullMsg []byte
}
