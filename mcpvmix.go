package mcpvmix

import (
	"bytes"
	"encoding/base64"
	"fmt"
	"image"
	"image/jpeg"
	"net/url"
	"os"
	"path"
	"strings"
	"sync"
	"time"

	models "github.com/FlowingSPDG/vmix-go"
	vmixhttp "github.com/FlowingSPDG/vmix-go/http"
	mcp_golang "github.com/metoro-io/mcp-golang"
	"github.com/metoro-io/mcp-golang/transport/stdio"
	"github.com/samber/lo"
	"golang.org/x/image/draw"
	"golang.org/x/sync/errgroup"
	"golang.org/x/xerrors"

	"github.com/FlowingSPDG/mcp-vmix/logger"
)

type MCPvMix interface {
	// general
	FetchVMix(arguments ConnectVmixArguments) (*mcp_golang.ToolResponse, error)

	// shortcut functions
	CutVMix(arguments VmixCutArguments) (*mcp_golang.ToolResponse, error)
	FadeVMix(arguments VmixFadeArguments) (*mcp_golang.ToolResponse, error)
	FadeToBlackVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)

	// recording functions
	StartRecordingVMix(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error)
	StopRecordingVMix(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error)

	// streaming functions
	StartStreamingVMix(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error)
	StopStreamingVMix(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error)

	// external output functions
	StartExternalVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)
	StopExternalVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)

	// multicorder functions
	StartMulticorderVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)
	StopMulticorderVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)

	// playlist functions
	StartPlaylistVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)
	StopPlaylistVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)

	// fullscreen function
	FullscreenVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error)

	// Snapshot functions
	SnapShotVMix(arguments GetCurrentScreenshotArguments) (*mcp_golang.ToolResponse, error)
	SnapShotInputVMix(arguments GetCurrentScreenshotInputArguments) (*mcp_golang.ToolResponse, error)

	// support functions
	GetShortcutURL(arguments GetShortcutURLArguments) (*mcp_golang.ToolResponse, error)
	AddBlank(arguments AddBlankArguments) (*mcp_golang.ToolResponse, error)
	CheckScreenshot(arguments CheckScreenshotArguments) (*mcp_golang.ToolResponse, error)
	CheckScreenshotInput(arguments CheckScreenshotInputArguments) (*mcp_golang.ToolResponse, error)
	MakeScene(arguments MakeSceneArguments) (*mcp_golang.ToolResponse, error)
	AdjustLayers(arguments AdjustLayersArguments) (*mcp_golang.ToolResponse, error)
}

type mcpVmix struct {
	logger logger.Logger
	srv    *mcp_golang.Server
}

// FetchVMix implements MCPvMix.
func (m *mcpVmix) FetchVMix(arguments ConnectVmixArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to connect to vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}
	m.logger.Info(fmt.Sprintf("Successfully connected to vMix instance at %s:%d", arguments.IP, arguments.Port))

	inputs := lo.Map(vmix.Inputs.Input, func(input models.Input, _ int) *mcp_golang.Content {
		overlays := lo.Map(input.Overlay, func(overlay models.InputOverlay, _ int) string {
			return fmt.Sprintf("Overlay: %d: Text: %s Key: %s Positions:%v", overlay.Index, overlay.Text, overlay.Key, overlay.Position)
		})
		overlaysStr := strings.Join(overlays, "\n")
		return mcp_golang.NewTextContent(
			fmt.Sprintf("Input: %d: Key:%s, Name: %s. State: %s, Position: %d, Duration: %d, Loop: %t Overlays:%v", input.Number, input.Key, input.Name, input.State, input.Position, input.Duration, input.Loop, overlaysStr),
		)
	})

	allContents := append([]*mcp_golang.Content{
		mcp_golang.NewTextContent(fmt.Sprintf("Connected to vMix instance %s:%d", arguments.IP, arguments.Port)),
		mcp_golang.NewTextContent(fmt.Sprintf("vMix version is %s, Edition is %s.", vmix.Version, vmix.Edition)),
		mcp_golang.NewTextContent(fmt.Sprintf("vMix is running on %s.", vmix.Preset)),
	}, inputs...)

	m.logger.Info(fmt.Sprintf("Successfully fetched vMix information: Version=%s, Edition=%s, Preset=%s", vmix.Version, vmix.Edition, vmix.Preset))
	return mcp_golang.NewToolResponse(allContents...), nil
}

// CutVMix implements MCPvMix.
func (m *mcpVmix) CutVMix(arguments VmixCutArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to cut to input %s on vMix instance at %s:%d", arguments.Input, arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.Cut(arguments.Input); err != nil {
		errMsg := fmt.Sprintf("Failed to perform cut operation: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info(fmt.Sprintf("Successfully cut to input %s", arguments.Input))
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Cut to input %s", arguments.Input))), nil
}

// FadeVMix implements MCPvMix.
func (m *mcpVmix) FadeVMix(arguments VmixFadeArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to fade to input %s with duration %d on vMix instance at %s:%d", arguments.Input, arguments.Duration, arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.Fade(arguments.Input, uint(arguments.Duration)); err != nil {
		errMsg := fmt.Sprintf("Failed to perform fade operation: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info(fmt.Sprintf("Successfully faded to input %s with duration %d", arguments.Input, arguments.Duration))
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("Fade to input %s for Duration %d", arguments.Input, arguments.Duration))), nil
}

// FadeToBlackVMix implements MCPvMix.
func (m *mcpVmix) FadeToBlackVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to fade to black on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.FadeToBlack(); err != nil {
		errMsg := fmt.Sprintf("Failed to perform fade to black: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully performed fade to black")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Performed Fade To Black")), nil
}

// StartRecordingVMix implements MCPvMix.
func (m *mcpVmix) StartRecordingVMix(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to start recording on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StartRecording(); err != nil {
		errMsg := fmt.Sprintf("Failed to start recording: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully started recording")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started recording")), nil
}

// StopRecordingVMix implements MCPvMix.
func (m *mcpVmix) StopRecordingVMix(arguments VmixRecordingArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to stop recording on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StopRecording(); err != nil {
		errMsg := fmt.Sprintf("Failed to stop recording: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully stopped recording")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped recording")), nil
}

// StartStreamingVMix implements MCPvMix.
func (m *mcpVmix) StartStreamingVMix(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to start streaming on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StartStreaming(uint(arguments.StreamNumber)); err != nil {
		errMsg := fmt.Sprintf("Failed to start streaming: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully started streaming")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started streaming")), nil
}

// StopStreamingVMix implements MCPvMix.
func (m *mcpVmix) StopStreamingVMix(arguments VmixStreamingArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to stop streaming on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StopStreaming(uint(arguments.StreamNumber)); err != nil {
		errMsg := fmt.Sprintf("Failed to stop streaming: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully stopped streaming")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped streaming")), nil
}

// StartExternalVMix implements MCPvMix.
func (m *mcpVmix) StartExternalVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to start external output on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StartExternal(); err != nil {
		errMsg := fmt.Sprintf("Failed to start external output: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully started external output")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started external output")), nil
}

// StopExternalVMix implements MCPvMix.
func (m *mcpVmix) StopExternalVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to stop external output on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StopExternal(); err != nil {
		errMsg := fmt.Sprintf("Failed to stop external output: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully stopped external output")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped external output")), nil
}

// StartMulticorderVMix implements MCPvMix.
func (m *mcpVmix) StartMulticorderVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to start MultiCorder on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StartMultiCorder(); err != nil {
		errMsg := fmt.Sprintf("Failed to start MultiCorder: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully started MultiCorder")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started MultiCorder")), nil
}

// StopMulticorderVMix implements MCPvMix.
func (m *mcpVmix) StopMulticorderVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to stop MultiCorder on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.StopMultiCorder(); err != nil {
		errMsg := fmt.Sprintf("Failed to stop MultiCorder: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully stopped MultiCorder")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped MultiCorder")), nil
}

// StartPlaylistVMix implements MCPvMix.
func (m *mcpVmix) StartPlaylistVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to start playlist on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.SendFunction("StartPlayList", nil); err != nil {
		errMsg := fmt.Sprintf("Failed to start playlist: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully started playlist")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Started playlist")), nil
}

// StopPlaylistVMix implements MCPvMix.
func (m *mcpVmix) StopPlaylistVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to stop playlist on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.SendFunction("StopPlayList", nil); err != nil {
		errMsg := fmt.Sprintf("Failed to stop playlist: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully stopped playlist")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Stopped playlist")), nil
}

// FullscreenVMix implements MCPvMix.
func (m *mcpVmix) FullscreenVMix(arguments VmixBasicArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to toggle fullscreen on vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.Fullscreen(); err != nil {
		errMsg := fmt.Sprintf("Failed to toggle fullscreen: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully toggled fullscreen")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Toggled fullscreen")), nil
}

// GetShortcutURL implements MCPvMix.
func (m *mcpVmix) GetShortcutURL(arguments GetShortcutURLArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("vMixインスタンス %s:%d のショートカットURLを生成します。関数: %s", arguments.IP, arguments.Port, arguments.Function))

	// URLを構築
	u := &url.URL{
		Scheme: "http",
		Host:   fmt.Sprintf("%s:%d", arguments.IP, arguments.Port),
		Path:   "/api",
	}

	// クエリパラメータを設定
	q := u.Query()
	q.Set("Function", arguments.Function)
	for key, value := range arguments.Queries {
		q.Set(key, value)
	}
	u.RawQuery = q.Encode()

	shortcutURL := u.String()

	m.logger.Info(fmt.Sprintf("ショートカットURLを生成しました: %s", shortcutURL))
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(shortcutURL)), nil
}

func (m *mcpVmix) addBlank(client *vmixhttp.Client, isTransparent bool) error {
	value := "Black"
	if isTransparent {
		value = "Transparent"
	}

	if err := client.AddInput("Colour", value); err != nil {
		errMsg := fmt.Sprintf("Failed to add blank: %v", err)
		m.logger.Error(errMsg)
		return fmt.Errorf(errMsg)
	}

	return nil
}

// AddBlank implements MCPvMix.
func (m *mcpVmix) AddBlank(arguments AddBlankArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to add %d blank inputs on vMix instance at %s:%d", arguments.Numbers, arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	var wg sync.WaitGroup
	for range arguments.Numbers {
		wg.Add(1)
		go func() {
			defer wg.Done()
			if err := m.addBlank(vmix, arguments.IsTransparent); err != nil {
				errMsg := fmt.Sprintf("Failed to add blank input: %v", err)
				m.logger.Error(errMsg)
			}
		}()
	}
	wg.Wait()

	m.logger.Info("Successfully added blank inputs")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Added blank inputs")), nil
}

// SnapShotVMix implements MCPvMix.
func (m *mcpVmix) SnapShotVMix(arguments GetCurrentScreenshotArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to get screenshot from vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.Snapshot(arguments.SaveDir); err != nil {
		errMsg := fmt.Sprintf("Failed to take screenshot: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully took screenshot")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Took screenshot")), nil
}

// SnapShotInputVMix implements MCPvMix.
func (m *mcpVmix) SnapShotInputVMix(arguments GetCurrentScreenshotInputArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Attempting to get input screenshot from vMix instance at %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	if err := vmix.SnapshotInput(arguments.Input, arguments.SaveDir); err != nil {
		errMsg := fmt.Sprintf("Failed to take input screenshot: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully took input screenshot")
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("Took input screenshot")), nil
}

// CheckScreenshot implements MCPvMix.
func (m *mcpVmix) CheckScreenshot(arguments CheckScreenshotArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("スクリーンショットを確認します。vMixインスタンス %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 一時ディレクトリにスクリーンショットを保存
	now := time.Now().Format("20060102_150405.jpg")
	tmpDir := os.TempDir()
	filePath := path.Join(tmpDir, now)
	if err := vmix.Snapshot(filePath); err != nil {
		errMsg := fmt.Sprintf("Failed to take screenshot: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully checked screenshot")
	time.Sleep(5 * time.Second) // wait until screen shot is saved

	// 保存したスクリーンショットを取得してBase64にエンコード
	snapShotFileBase64, err := retryReadScreenshot(filePath, 30, 200*time.Millisecond)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read screenshot: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	content := mcp_golang.NewImageContent(snapShotFileBase64, "image/jpeg")
	return mcp_golang.NewToolResponse(content), nil
}

// CheckScreenshotInput implements MCPvMix.
func (m *mcpVmix) CheckScreenshotInput(arguments CheckScreenshotInputArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("Checking screenshot for input %s on vMix instance at %s:%d", arguments.Input, arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to connect to vMix instance: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 一時ディレクトリにスクリーンショットを保存
	now := time.Now().Format("20060102_150405.jpg")
	tmpDir := os.TempDir()
	filePath := path.Join(tmpDir, now)
	if err := vmix.SnapshotInput(arguments.Input, filePath); err != nil {
		errMsg := fmt.Sprintf("Failed to take input screenshot: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info("Successfully checked input screenshot")
	time.Sleep(5 * time.Second) // wait until screen shot is saved
	// 保存したスクリーンショットを取得してBase64にエンコード
	snapShotFileBase64, err := retryReadScreenshot(filePath, 30, 200*time.Millisecond)
	if err != nil {
		errMsg := fmt.Sprintf("Failed to read screenshot: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	content := mcp_golang.NewImageContent(snapShotFileBase64, "image/jpeg")
	return mcp_golang.NewToolResponse(content), nil
}

// retryReadScreenshot read a file from filePath, if it fails, it will retry up to maxRetries times with a retryInterval.
// image is resized to 1/2 of the original size due to the performance.
// returned string is base64-encoded string of the file.
func retryReadScreenshot(filePath string, maxRetries int, retryInterval time.Duration) (string, error) {
	for range maxRetries {
		snapShotFile, err := os.Open(filePath)
		if err != nil {
			time.Sleep(retryInterval)
			continue
		}
		defer snapShotFile.Close()
		img, err := jpeg.Decode(snapShotFile)
		if err != nil {
			return "", err
		}

		// 画像を1/2にリサイズ
		bounds := img.Bounds()
		newWidth := bounds.Dx() / 2
		newHeight := bounds.Dy() / 2
		resizedImg := image.NewRGBA(image.Rect(0, 0, newWidth, newHeight))

		draw.NearestNeighbor.Scale(resizedImg, resizedImg.Bounds(), img, img.Bounds(), draw.Over, nil)

		// JPEGにエンコード
		var buf bytes.Buffer
		if err := jpeg.Encode(&buf, resizedImg, &jpeg.Options{Quality: 80}); err != nil {
			return "", err
		}

		// Base64エンコード
		snapShotFileBase64 := base64.StdEncoding.EncodeToString(buf.Bytes())
		return snapShotFileBase64, nil
	}
	return "", fmt.Errorf("Failed to read screenshot")
}

// MakeScene implements MCPvMix.
func (m *mcpVmix) MakeScene(arguments MakeSceneArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("シーン %s を作成します。vMixインスタンス %s:%d", arguments.Input, arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("vMixインスタンスへの接続に失敗しました: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 各レイヤーを設定
	eg := errgroup.Group{}
	for index, layer := range arguments.Layers {
		eg.Go(func() error {
			if err := vmix.SetLayer(arguments.Input, uint8(index+1), layer.Input); err != nil {
				errMsg := fmt.Sprintf("failed to set input layer for index:%d input: %s: error: %v", index+1, layer.Input, err)
				m.logger.Error(errMsg)
				return xerrors.Errorf(errMsg)
			}

			if err := vmix.SetLayerPanX(arguments.Input, uint8(index+1), layer.PanX); err != nil {
				errMsg := fmt.Sprintf("failed to set input layer position for index:%d input: %s: error: %v", index+1, layer.Input, err)
				m.logger.Error(errMsg)
				return xerrors.Errorf(errMsg)
			}

			if err := vmix.SetLayerPanY(arguments.Input, uint8(index+1), layer.PanY); err != nil {
				errMsg := fmt.Sprintf("failed to set input layer position for index:%d input: %s: error: %v", index+1, layer.Input, err)
				m.logger.Error(errMsg)
				return xerrors.Errorf(errMsg)
			}

			if err := vmix.SetLayerZoom(arguments.Input, uint8(index+1), layer.Zoom); err != nil {
				errMsg := fmt.Sprintf("failed to set input layer zoom for index:%d input: %s: error: %v", index+1, layer.Input, err)
				m.logger.Error(errMsg)
				return xerrors.Errorf(errMsg)
			}

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		errMsg := fmt.Sprintf("failed to set input layer: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info(fmt.Sprintf("シーン %s の作成に成功しました", arguments.Input))

	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent(fmt.Sprintf("シーン %s を作成しました", arguments.Input))), nil
}

// AdjustLayers implements MCPvMix.
func (m *mcpVmix) AdjustLayers(arguments AdjustLayersArguments) (*mcp_golang.ToolResponse, error) {
	m.logger.Info(fmt.Sprintf("レイヤーを調整します。vMixインスタンス %s:%d", arguments.IP, arguments.Port))

	vmix, err := vmixhttp.NewClient(arguments.IP, arguments.Port)
	if err != nil {
		errMsg := fmt.Sprintf("vMixインスタンスへの接続に失敗しました: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	// 各レイヤーを設定
	eg := errgroup.Group{}
	for index, layer := range arguments.Layers {
		eg.Go(func() error {
			eg.Go(func() error {
				if err := vmix.SetLayer(arguments.Input, uint8(layer.Index), layer.Input); err != nil {
					errMsg := fmt.Sprintf("failed to set input layer for index:%d input: %s: error: %v", index+1, arguments.Input, err)
					m.logger.Error(errMsg)
					return xerrors.Errorf(errMsg)
				}
				return nil
			})

			eg.Go(func() error {
				if err := vmix.SetLayerPanX(arguments.Input, uint8(layer.Index), layer.PanX); err != nil {
					errMsg := fmt.Sprintf("failed to set input layer for index:%d input: %s: error: %v", index+1, arguments.Input, err)
					m.logger.Error(errMsg)
					return xerrors.Errorf(errMsg)
				}
				return nil
			})

			eg.Go(func() error {
				if err := vmix.SetLayerPanY(arguments.Input, uint8(layer.Index), layer.PanY); err != nil {
					errMsg := fmt.Sprintf("failed to set input layer position for index:%d input: %s: error: %v", index+1, arguments.Input, err)
					m.logger.Error(errMsg)
					return xerrors.Errorf(errMsg)
				}
				return nil
			})

			eg.Go(func() error {
				if err := vmix.SetLayerZoom(arguments.Input, uint8(layer.Index), layer.Zoom); err != nil {
					errMsg := fmt.Sprintf("failed to set input layer zoom for index:%d input: %s: error: %v", index+1, arguments.Input, err)
					m.logger.Error(errMsg)
					return xerrors.Errorf(errMsg)
				}
				return nil
			})

			eg.Go(func() error {
				if err := vmix.SetLayerCrop(arguments.Input, uint8(layer.Index), layer.CropX1, layer.CropY1, layer.CropX2, layer.CropY2); err != nil {
					errMsg := fmt.Sprintf("failed to set input layer crop for index:%d input: %s: error: %v", index+1, arguments.Input, err)
					m.logger.Error(errMsg)
					return xerrors.Errorf(errMsg)
				}
				return nil
			})

			return nil
		})
	}

	if err := eg.Wait(); err != nil {
		errMsg := fmt.Sprintf("failed to set input layer: %v", err)
		m.logger.Error(errMsg)
		return nil, fmt.Errorf(errMsg)
	}

	m.logger.Info(fmt.Sprintf("レイヤーを調整しました"))
	return mcp_golang.NewToolResponse(mcp_golang.NewTextContent("レイヤーを調整しました")), nil
}

func NewMCPvMix(logger logger.Logger) MCPvMix {
	srv := mcp_golang.NewServer(stdio.NewStdioServerTransport())

	return &mcpVmix{
		logger: logger,
		srv:    srv,
	}
}
