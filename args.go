package mcpvmix

type BaseVMixArguments struct {
	IP   string `json:"ip" jsonschema:"required,description=The IP address of the vMix instance. generally this is 127.0.0.1"`
	Port int    `json:"port" jsonschema:"required,description=The port of the vMix instance. generally this is 8088."`
}

type ConnectVmixArguments struct {
	BaseVMixArguments
}

type VmixInput struct {
	Input string `json:"input" jsonschema:"required,description=The input to cut to. This could be input number or input key(UUID). key would be preferred."`
}

type VmixCutArguments struct {
	BaseVMixArguments
	VmixInput
}

type VmixFadeArguments struct {
	BaseVMixArguments
	VmixInput
	Duration int `json:"duration" jsonschema:"required,description=The duration of the fade. This is the duration of the fade in milliseconds."`
}
type VmixRecordingArguments struct {
	BaseVMixArguments
}

type VmixStreamingArguments struct {
	BaseVMixArguments
	StreamNumber int `json:"streamNumber" jsonschema:"required,description=The stream number to start streaming on. Generally this is 1~4."`
}

type VmixBasicArguments struct {
	BaseVMixArguments
}

type GetShortcutURLArguments struct {
	BaseVMixArguments
	Function string            `json:"function" jsonschema:"required,description=The function to get the shortcut URL for"`
	Queries  map[string]string `json:"queries" jsonschema:"required,description=The Key/Value queries for the function arguments. e.g. {\"Input\": \"1\"} would be for the function Input=1 and the argument Input. {\"Mix\": \"1\"} would be for the function Mix=1 and the argument Mix."`
}

type AddBlankArguments struct {
	BaseVMixArguments
	Numbers       int  `json:"numbers" jsonschema:"required,description=The number of blank inputs to add"`
	IsTransparent bool `json:"isTransparent" jsonschema:"required,description=Whether the blank inputs should be transparent"`
}

type GetCurrentScreenshotArguments struct {
	BaseVMixArguments
	SaveDir string `json:"saveDir" jsonschema:"required,description=The directory to save the screenshot to. This needs to be a valid directory file path. e.g. C:/Users/SPDG/Desktop/test.jpg . the content type depends on the file extension."`
}

type GetCurrentScreenshotInputArguments struct {
	BaseVMixArguments
	VmixInput
	SaveDir string `json:"saveDir" jsonschema:"required,description=The directory to save the screenshot to. This needs to be a valid directory file path. e.g. C:/Users/SPDG/Desktop/test.jpg . the content type depends on the file extension."`
}

type CheckScreenshotArguments struct {
	BaseVMixArguments
}

type CheckScreenshotInputArguments struct {
	BaseVMixArguments
	VmixInput
}

type MakeSceneArguments struct {
	BaseVMixArguments
	VmixInput
	Layers []MakeSceneLayerArguments `json:"layers" jsonschema:"required,description=The layers to make the scene. Up to 10 layers are supported."`
}

type MakeSceneLayerArguments struct {
	VmixInput
	PanX float64 `json:"panX" jsonschema:"required,description=The layer position x1 to make the scene for. 0 means center. Range should be between -2 to 2."`
	PanY float64 `json:"panY" jsonschema:"required,description=The layer position y1 to make the scene for. 0 means center. Range should be between -2 to 2."`
	Zoom float64 `json:"zoom" jsonschema:"required,description=The layer zoom to make the scene for. 1 means 100%(default). Range should be between 0-5."`
}
