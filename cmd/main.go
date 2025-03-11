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
type VmixRecordingArguments struct {
	IP   string `json:"ip" jsonschema:"required,description=The IP address of the vMix instance"`
	Port int    `json:"port" jsonschema:"required,description=The port of the vMix instance"`
}

type VmixStreamingArguments struct {
	IP   string `json:"ip" jsonschema:"required,description=The IP address of the vMix instance"`
	Port int    `json:"port" jsonschema:"required,description=The port of the vMix instance"`
}

type VmixBasicArguments struct {
	IP   string `json:"ip" jsonschema:"required,description=The IP address of the vMix instance"`
	Port int    `json:"port" jsonschema:"required,description=The port of the vMix instance"`
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

	if err := server.RegisterTool("vmix_start_recording", "Start recording on a vMix instance", func(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StartRecording", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started recording")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_stop_recording", "Stop recording on a vMix instance", func(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StopRecording", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped recording")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_start_streaming", "Start streaming on a vMix instance", func(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StartStreaming", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started streaming")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_stop_streaming", "Stop streaming on a vMix instance", func(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StopStreaming", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped streaming")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_fade_to_black", "Perform Fade To Black on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("FadeToBlack", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Performed Fade To Black")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_start_external", "Start external output on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StartExternal", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started external output")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_stop_external", "Stop external output on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StopExternal", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped external output")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_start_multicorder", "Start MultiCorder on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StartMultiCorder", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started MultiCorder")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_stop_multicorder", "Stop MultiCorder on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StopMultiCorder", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped MultiCorder")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_start_playlist", "Start playlist on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StartPlayList", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started playlist")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_stop_playlist", "Stop playlist on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("StopPlayList", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped playlist")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.RegisterTool("vmix_fullscreen", "Toggle fullscreen on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			return nil, err
		}

		if err := vmix.SendFunction("FullScreen", nil); err != nil {
			return nil, err
		}

		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Toggled fullscreen")), nil
	}); err != nil {
		panic(err)
	}

	if err := server.Serve(); err != nil {
		panic(err)
	}

	<-ctx.Done()
}
