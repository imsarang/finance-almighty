package tools

import "github.com/modelcontextprotocol/go-sdk/mcp"

type Registrar interface {
	Register(s *mcp.Server)
}

func ServerInit(name, version string, regsitrars ...Registrar) *mcp.Server {
	s := mcp.NewServer(&mcp.Implementation{
		Name: name,
		Version: version,
	}, nil)

	for _,r := range regsitrars {
		r.Register(s)
	}
	return s
}