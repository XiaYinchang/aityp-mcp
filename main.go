package main

import (
	"context"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

func main() {
	s := server.NewMCPServer(
		"image-mirror",
		"0.1.1",
		server.WithLogging(),
		server.WithRecovery(),
	)
	tool := mcp.NewTool("search",
		mcp.WithDescription("search iamges"),
		mcp.WithString("name",
			mcp.Required(),
			mcp.Description("name of the image to search"),
		),
		mcp.WithString("site",
			mcp.DefaultString("All"),
			mcp.Enum("All", "gcr.io",
				"ghcr.io",
				"quay.io",
				"k8s.gcr.io",
				"docker.io",
				"registry.k8s.io",
				"docker.elastic.co",
				"skywalking.docker.scarf.sh",
				"mcr.microsoft.com")),
		mcp.WithString("platform",
			mcp.DefaultString("linux/amd64"),
			mcp.Enum("All", "linux/386",
				"linux/amd64",
				"linux/arm64",
				"linux/arm",
				"linux/ppc64le",
				"linux/s390x",
				"linux/mips64le",
				"linux/riscv64",
				"linux/loong64")),
	)

	s.AddTool(tool, searchHandler)

	// Start the stdio server
	if err := server.ServeStdio(s); err != nil {
		fmt.Printf("Server error: %v\n", err)
	}
}

func searchHandler(ctx context.Context, request mcp.CallToolRequest) (*mcp.CallToolResult, error) {
	name, ok := request.Params.Arguments["name"].(string)
	if !ok {
		return nil, errors.New("name must be a string")
	}
	site, ok := request.Params.Arguments["site"].(string)
	if !ok {
		return nil, errors.New("site must be a string")
	}
	platform, ok := request.Params.Arguments["platform"].(string)
	if !ok {
		return nil, errors.New("site must be a string")
	}
	url := fmt.Sprintf("https://docker.aityp.com/api/v1/image?search=%s&site=%s&platform=%s", name, site, platform)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	return mcp.NewToolResultText(string(body)), nil
}

// type Image struct {
// 	Source    string `json:"source"`
// 	Mirror    string `json:"mirror"`
// 	Platform  string `json:"platform"`
// 	Size      string `json:"size"`
// 	CreatedAt string `json:"createdAt"`
// }

// type RawResult struct {
// 	Results []Image `json:"results"`
// 	Search  string  `json:"search"`
// 	Count   int     `json:"count"`
// 	Error   bool    `json:"error"`
// }
