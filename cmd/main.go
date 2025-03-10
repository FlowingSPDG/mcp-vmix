package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"strconv"

	vmixhttp "github.com/FlowingSPDG/vmix-go/http"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/mark3labs/mcp-go/server"
)

type ToolName string

const (
	listInputs ToolName = "listInputs"
	cut        ToolName = "cut"
	fade       ToolName = "fade"
)

type PromptName string

const (
	connectVMix PromptName = "connect_vmix"
)

func NewMCPServer() *server.MCPServer {
	mcpServer := server.NewMCPServer(
		"vMix",
		"0.1.0",
		server.WithResourceCapabilities(true, true),
		server.WithPromptCapabilities(true),
		server.WithLogging(),
	)

	// connect vMix prompt...
	mcpServer.AddPrompt(mcp.NewPrompt(string(connectVMix),
		mcp.WithPromptDescription("Connect to vMix"),
		mcp.WithArgument("ip",
			mcp.ArgumentDescription("The IP address of the vMix server"),
			mcp.RequiredArgument(),
		),
		mcp.WithArgument("port",
			mcp.ArgumentDescription("The port of the vMix server"),
			mcp.RequiredArgument(),
		),
	), handleConnectVMixPrompt)

	mcpServer.AddTool(mcp.NewTool(string(cut),
		mcp.WithDescription("Cut the current input"),
		mcp.WithNumber("ip",
			mcp.Description("The IP address of the vMix server"),
			mcp.Required(),
		),
		mcp.WithNumber("port",
			mcp.Description("The port of the vMix server"),
			mcp.Required(),
		),
		mcp.WithNumber("input",
			mcp.Description("The input number to cut"),
			mcp.Required(),
		),
	), handleCutTool)

	return mcpServer
}

func handleConnectVMixPrompt(
	ctx context.Context,
	request mcp.GetPromptRequest,
) (*mcp.GetPromptResult, error) {
	arguments := request.Params.Arguments
	ip, ok := arguments["ip"]
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}
	portStr, ok := arguments["port"]
	if !ok {
		return nil, fmt.Errorf("invalid arguments")
	}
	port, err := strconv.Atoi(portStr)
	if err != nil {
		return nil, fmt.Errorf("invalid port: %v", err)
	}

	vc, err := vmixhttp.NewClient(ip, port)
	if err != nil {
		return nil, fmt.Errorf("failed to create vMix client: %v", err)
	}

	return &mcp.GetPromptResult{
		Description: "Connected to vMix and get data",
		Messages: []mcp.PromptMessage{
			{
				Role: mcp.RoleUser,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("Connect to vMix and fetch data. IP: %s, Port: %d", ip, port),
				},
			},
			{
				Role: mcp.RoleAssistant,
				Content: mcp.TextContent{
					Type: "text",
					Text: fmt.Sprintf("I understand. User provided vMix IP and port. It seems to be connected. %v", vc),
				},
			},
		},
	}, nil
}

func handleCutTool(
	ctx context.Context,
	request mcp.CallToolRequest,
) (*mcp.CallToolResult, error) {
	arguments := request.Params.Arguments
	ip, ok1 := arguments["ip"]
	portStr, ok2 := arguments["port"]
	inputStr, ok3 := arguments["input"]
	if !ok1 || !ok2 || !ok3 {
		return nil, fmt.Errorf("invalid number arguments")
	}
	port, err := strconv.Atoi(portStr.(string))
	if err != nil {
		return nil, fmt.Errorf("invalid port: %v", err)
	}
	input, err := strconv.Atoi(inputStr.(string))
	if err != nil {
		return nil, fmt.Errorf("invalid input: %v", err)
	}

	vc, err := vmixhttp.NewClient(ip.(string), port)
	if err != nil {
		return nil, fmt.Errorf("failed to create vMix client: %v", err)
	}

	err = vc.SendFunction("Cut", map[string]string{
		"input": strconv.Itoa(input),
	})
	if err != nil {
		return nil, fmt.Errorf("failed to send function: %v", err)
	}

	return &mcp.CallToolResult{
		Content: []mcp.Content{
			mcp.TextContent{
				Type: "text",
				Text: fmt.Sprintf("Function result is %v. vMix HTTP API result is %v.", err, vc),
			},
		},
	}, nil
}

func main() {
	var transport string
	flag.StringVar(&transport, "t", "stdio", "Transport type (stdio or sse)")
	flag.StringVar(&transport, "transport", "stdio", "Transport type (stdio or sse)")
	flag.Parse()

	mcpServer := NewMCPServer()

	// Only check for "sse" since stdio is the default
	if transport == "sse" {
		sseServer := server.NewSSEServer(mcpServer, "")
		log.Printf("SSE server listening on :8080")
		if err := sseServer.Start(":8080"); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	} else {
		if err := server.ServeStdio(mcpServer); err != nil {
			log.Fatalf("Server error: %v", err)
		}
	}
}
