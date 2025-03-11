package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	models "github.com/FlowingSPDG/vmix-go"
	vmixhttp "github.com/FlowingSPDG/vmix-go/http"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/samber/lo"

	"github.com/FlowingSPDG/mcp-vmix/logger"
)

var log logger.Logger

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

	// ロガーの初期化
	logPath, err := logger.GetLogFilePath()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to get log file path: %v\n", err)
		return
	}

	log, err = logger.NewFileLogger(logPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to create logger: %v\n", err)
		return
	}
	defer log.Close()

	log.Info("Starting vMix MCP server...")
	server := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	if err := server.RegisterTool("vmix_fetch", "Connect to a vMix instance", func(arguments ConnectVmixArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to connect to vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}
		log.Info(fmt.Sprintf("Successfully connected to vMix instance at %s:%d", arguments.IP, arguments.Port))

		inputs := lo.Map(vmix.Inputs.Input, func(input models.Input, _ int) *mcp_golang.Content {
			return mcp_golang.NewTextContent(fmt.Sprintf("Input Number: %d, Name: %s. State: %s, Position: %d, Duration: %d, Loop: %t", input.Number, input.Name, input.State, input.Position, input.Duration, input.Loop))
		})

		allContents := append([]*mcp_golang.Content{
			mcp_golang.NewTextContent(fmt.Sprintf("Connected to vMix instance %s:%d", arguments.IP, arguments.Port)),
			mcp_golang.NewTextContent(fmt.Sprintf("vMix version is %s, Edition is %s.", vmix.Version, vmix.Edition)),
			mcp_golang.NewTextContent(fmt.Sprintf("vMix is running on %s.", vmix.Preset)),
		}, inputs...)

		log.Info(fmt.Sprintf("Successfully fetched vMix information: Version=%s, Edition=%s, Preset=%s", vmix.Version, vmix.Edition, vmix.Preset))
		return mcp_golang.NewToolResponse(allContents...), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fetch tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_cut", "Perform a cut shortcut on a vMix instance.", func(arguments VmixCutArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to cut to input %s on vMix instance at %s:%d", arguments.Input, arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.Cut(arguments.Input); err != nil {
			errMsg := fmt.Sprintf("Failed to perform cut operation: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info(fmt.Sprintf("Successfully cut to input %s", arguments.Input))
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Cut to input %s", arguments.Input))), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_cut tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_fade", "Perform a Fade shortcut function on a vMix instance", func(arguments VmixFadeArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to fade to input %s with duration %d on vMix instance at %s:%d", arguments.Input, arguments.Duration, arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.Fade(arguments.Input, uint(arguments.Duration)); err != nil {
			errMsg := fmt.Sprintf("Failed to perform fade operation: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info(fmt.Sprintf("Successfully faded to input %s with duration %d", arguments.Input, arguments.Duration))
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Fade to input %s for Duration %d", arguments.Input, arguments.Duration))), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fade tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_recording", "Start recording on a vMix instance", func(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to start recording on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StartRecording", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to start recording: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully started recording")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started recording")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_recording tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_recording", "Stop recording on a vMix instance", func(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to stop recording on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StopRecording", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to stop recording: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully stopped recording")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped recording")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_recording tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_streaming", "Start streaming on a vMix instance", func(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to start streaming on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StartStreaming", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to start streaming: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully started streaming")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started streaming")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_streaming tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_streaming", "Stop streaming on a vMix instance", func(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to stop streaming on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StopStreaming", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to stop streaming: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully stopped streaming")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped streaming")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_streaming tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_fade_to_black", "Perform Fade To Black on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to fade to black on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("FadeToBlack", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to perform fade to black: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully performed fade to black")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Performed Fade To Black")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fade_to_black tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_external", "Start external output on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to start external output on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StartExternal", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to start external output: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully started external output")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started external output")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_external tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_external", "Stop external output on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to stop external output on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StopExternal", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to stop external output: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully stopped external output")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped external output")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_external tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_multicorder", "Start MultiCorder on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to start MultiCorder on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StartMultiCorder", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to start MultiCorder: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully started MultiCorder")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started MultiCorder")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_multicorder tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_multicorder", "Stop MultiCorder on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to stop MultiCorder on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StopMultiCorder", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to stop MultiCorder: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully stopped MultiCorder")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped MultiCorder")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_multicorder tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_playlist", "Start playlist on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to start playlist on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StartPlayList", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to start playlist: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully started playlist")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started playlist")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_playlist tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_playlist", "Stop playlist on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to stop playlist on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("StopPlayList", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to stop playlist: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully stopped playlist")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped playlist")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_playlist tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_fullscreen", "Toggle fullscreen on a vMix instance", func(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
		log.Info(fmt.Sprintf("Attempting to toggle fullscreen on vMix instance at %s:%d", arguments.IP, arguments.Port))

		vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
		if err != nil {
			errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		if err := vmix.SendFunction("FullScreen", nil); err != nil {
			errMsg := fmt.Sprintf("Failed to toggle fullscreen: %v", err)
			log.Error(errMsg)
			return nil, fmt.Errorf(errMsg)
		}

		log.Info("Successfully toggled fullscreen")
		return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Toggled fullscreen")), nil
	}); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fullscreen tool: %v", err))
		return
	}

	log.Info("Starting MCP server...")
	if err := server.Serve(); err != nil {
		log.Error(fmt.Sprintf("Failed to start MCP server: %v", err))
		return
	}

	<-ctx.Done()
}
