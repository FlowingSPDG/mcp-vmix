package main

import (
	"context"
	"fmt"
	"os/signal"
	"syscall"

	models "github.com/FlowingSPDG/vmix-go"
	vmixhttp "github.com/FlowingSPDG/vmix-go/http"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/samber/lo"
)

type ConnectVmixArguments struct {
	IP   string `json:"ip" jsonschema:"required,description=The IP address of the vMix instance"`
	Port int    `json:"port" jsonschema:"required,description=The port of the vMix instance"`
}

type VmixCutArguments struct {
	IP    string `json:"ip" jsonschema:"required,description=The IP address of the vMix instance"`
	Port  int    `json:"port" jsonschema:"required,description=The port of the vMix instance"`
	Input string `json:"input" jsonschema:"required,description=The input to cut to"`
}

type VmixFadeArguments struct {
	IP       string `json:"ip" jsonschema:"required,description=The IP address of the vMix instance"`
	Port     int    `json:"port" jsonschema:"required,description=The port of the vMix instance"`
	Input    string `json:"input" jsonschema:"required,description=The input to fade to"`
	Duration int    `json:"duration" jsonschema:"required,description=The duration of the fade"`
}

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
	defer cancel()

	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	if err := server.RegisterTool("vmix_fetch", "Connect to a vMix instance", func(arguments ConnectVmixArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		inputs := lo.Map(vmix.Inputs.Input, func(input models.Input, _ int) *mcp_golang.Content {
			return mcp_golang.NewTextContent(fmt.Sprintf("Input Number: %d, Name: %s. State: %s, Position: %d, Duration: %d, Loop: %t", input.Number, input.Name, input.State, input.Position, input.Duration, input.Loop))
		})

		allContents := append([]*mcp_golang.Content{
			mcp_golang.NewTextContent(fmt.Sprintf("Connected to vMix instance %s:%d", arguments.IP, arguments.Port)),
			mcp_golang.NewTextContent(fmt.Sprintf("vMix version is %s, Edition is %s.", vmix.Version, vmix.Edition)),
			mcp_golang.NewTextContent(fmt.Sprintf("vMix is running on %s.", vmix.Preset)),
		}, inputs...)

		return mcp_golang.NewToolResponse(allContents...), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_cut", "Perform a cut shortcut on a vMix instance.", func(arguments VmixCutArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.Cut(arguments.Input); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Cut to input %s", arguments.Input))), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_fade", "Perform a Fade shortcut function on a vMix instance", func(arguments VmixFadeArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.Fade(arguments.Input, uint(arguments.Duration)); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Fade to input %s for Duration %d", arguments.Input, arguments.Duration))), nil
	}); err != nil {
		panic(err)
	}

	if err := server.Serve(); err != nil {
		panic(err)
	}

	<-ctx.Done()
}
