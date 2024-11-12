package cacheclient

type Cmd struct {
	Name  string
	Key   string
	Value string
	Error error
}

type Client interface {
	Run(*Cmd)
	PipelineRun([]*Cmd)
}

func New(typ, host string) Client {
	if typ == "redis" {
		return newRedisClient(host)
	}
	if typ == "http" {
		return newHTTPClient(host)
	}
	if typ == "tcp" {
		return newTcpClient(host)
	}
	panic("unknow client type " + typ)
}
