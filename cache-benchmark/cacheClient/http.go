package cacheclient

import (
	"io"
	"log"
	"net/http"
	"strings"
)

type httpClient struct {
	*http.Client
	host string
}

func (c *httpClient) get(key string) string {
	resp, err := c.Get(c.host + key)
	if err != nil {
		log.Panicln(key, err)
	}
	if resp.StatusCode == http.StatusNotFound {
		return ""
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
	buf, err := io.ReadAll(resp.Body)
	if err != nil {
		panic(err)
	}
	return string(buf)
}

func (c *httpClient) set(key, value string) {
	req, err := http.NewRequest(http.MethodPut, c.host+key, strings.NewReader(value))
	if err != nil {
		log.Println(key)
		panic(err)
	}
	resp, err := c.Do(req)
	if err != nil {
		log.Println(key)
		panic(err)
	}
	if resp.StatusCode != http.StatusOK {
		panic(resp.Status)
	}
}
func (c *httpClient) Run(cmd *Cmd) {
	if cmd.Name == "get" {
		cmd.Value = c.get(cmd.Key)
		return
	}
	if cmd.Name == "set" {
		c.set(cmd.Key, cmd.Value)
		return
	}
	panic("unknow cmd name " + cmd.Name)
}

func newHTTPClient(host string) *httpClient {
	client := &http.Client{Transport: &http.Transport{MaxIdleConnsPerHost: 1}}
	return &httpClient{client, "http://" + host + ":8080/cache/"}
}

func (c *httpClient) PipelineRun(cmds []*Cmd) {
	panic("not implemented") // TODO: Implement
}
