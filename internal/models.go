package internal

type Req struct {
	FullMsg []byte
	Host    string
	Port    string
}

type Resp struct {
	FullMsg []byte
}
