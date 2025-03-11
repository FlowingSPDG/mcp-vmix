package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	mcpvmix "github.com/FlowingSPDG/mcp-vmix"
	"github.com/FlowingSPDG/mcp-vmix/logger"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
)

var log logger.Logger

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

	// MCPvMixインスタンスの作成
	vmixInstance := mcpvmix.NewMCPvMix(log)

	// ツールの登録
	if err := server.RegisterTool("vmix_fetch", "Connect to a vMix instance.", vmixInstance.FetchVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fetch tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_cut", "Perform a cut shortcut on a vMix instance.", vmixInstance.CutVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_cut tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_fade", "Perform a Fade shortcut function on a vMix instance", vmixInstance.FadeVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fade tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_fade_to_black", "Perform Fade To Black on a vMix instance", vmixInstance.FadeToBlackVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fade_to_black tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_recording", "Start recording on a vMix instance", vmixInstance.StartRecordingVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_recording tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_recording", "Stop recording on a vMix instance", vmixInstance.StopRecordingVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_recording tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_streaming", "Start streaming on a vMix instance", vmixInstance.StartStreamingVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_streaming tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_streaming", "Stop streaming on a vMix instance", vmixInstance.StopStreamingVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_streaming tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_external", "Start external output on a vMix instance", vmixInstance.StartExternalVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_external tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_external", "Stop external output on a vMix instance", vmixInstance.StopExternalVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_external tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_multicorder", "Start MultiCorder on a vMix instance", vmixInstance.StartMulticorderVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_multicorder tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_multicorder", "Stop MultiCorder on a vMix instance", vmixInstance.StopMulticorderVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_multicorder tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_start_playlist", "Start playlist on a vMix instance", vmixInstance.StartPlaylistVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_start_playlist tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_stop_playlist", "Stop playlist on a vMix instance", vmixInstance.StopPlaylistVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_stop_playlist tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_fullscreen", "Toggle fullscreen on a vMix instance", vmixInstance.FullscreenVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_fullscreen tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_snapshot", "Take a screenshot of the current vMix instance", vmixInstance.SnapShotVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_snapshot tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_snapshot_input", "Take a screenshot of a specific input on a vMix instance", vmixInstance.SnapShotInputVMix); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_snapshot_input tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_check_screenshot", "Check screenshot of the current vMix instance", vmixInstance.CheckScreenshot); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_check_screenshot tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_check_screenshot_input", "Check screenshot of a specific input on a vMix instance", vmixInstance.CheckScreenshotInput); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_check_screenshot_input tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_get_shortcut_url", "Get shortcut URL for a vMix instance. This is useful for getting the URL of a shortcut function for vMix users.", vmixInstance.GetShortcutURL); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_get_shortcut_url tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_add_blank", "Add blank inputs to a vMix instance. Input will be appended last input slice, so input number must be current number of inputs +1 . This is useful for adding blank inputs to a vMix instance.", vmixInstance.AddBlank); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_add_blank tool: %v", err))
		return
	}

	if err := server.RegisterTool("vmix_make_scene", "Make a complicated composit scene on a vMix instance. This is used to make a new scene with multiple layers. It is always recommended to use this for Blank Input.", vmixInstance.MakeScene); err != nil {
		log.Error(fmt.Sprintf("Failed to register vmix_make_scene tool: %v", err))
		return
	}

	log.Info("Starting MCP server...")
	if err := server.Serve(); err != nil {
		log.Error(fmt.Sprintf("Failed to start MCP server: %v", err))
		return
	}

	<-ctx.Done()
}
