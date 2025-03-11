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
	PanX float64 `json:"panX" jsonschema:"required,description=The layer position x1 to make the scene for. 0 means center. Minus value means left, Plus value means right. Range should be between -2 to 2."`
	PanY float64 `json:"panY" jsonschema:"required,description=The layer position y1 to make the scene for. 0 means center. Minus value means up, Plus value means down. Range should be between -2 to 2."`
	Zoom float64 `json:"zoom" jsonschema:"required,description=The layer zoom to make the scene for. 1 means 100%(default). Range should be between 0-5."`
}

type AdjustLayersArguments struct {
	BaseVMixArguments
	VmixInput
	Layers []AdjustLayersLayerArguments `json:"layers" jsonschema:"required,description=The layers to adjust the layers. Up to 10 layers are supported."`
}

type AdjustLayersLayerArguments struct {
	VmixInput
	Index  int     `json:"index" jsonschema:"required,description=The index of the layer to adjust the layers. 1~10."`
	PanX   float64 `json:"panX" jsonschema:"required,description=The layer position x1 to make the scene for. 0 means center. Minus value means left, Plus value means right. Range should be between -2 to 2."`
	PanY   float64 `json:"panY" jsonschema:"required,description=The layer position y1 to make the scene for. 0 means center. Minus value means up, Plus value means down. Range should be between -2 to 2."`
	Zoom   float64 `json:"zoom" jsonschema:"required,description=The layer zoom to make the scene for. 1 means 100%(default). Range should be between 0-5."`
	CropX1 float64 `json:"cropX1" jsonschema:"required,description=The layer crop x1(left) to make the scene for. default is 0. 0=No Crop, 1=Full Crop"`
	CropY1 float64 `json:"cropY1" jsonschema:"required,description=The layer crop y1(top) to make the scene for. default is 0. 0=No Crop, 1=Full Crop"`
	CropX2 float64 `json:"cropX2" jsonschema:"required,description=The layer crop x2(right) to make the scene for. default is 1. 1=No Crop, 0=Full Crop"`
	CropY2 float64 `json:"cropY2" jsonschema:"required,description=The layer crop y2(bottom) to make the scene for. default is 1. 1=No Crop, 0=Full Crop"`
	// Rectangle関連は解像度の指定も必要になるため非常に面倒
	// RectangleX      *float64 `json:"rectangleX" jsonschema:"description=The layer rectangle(rotation x) x to make the scene for. default is null."`
	// RectangleY      *float64 `json:"rectangleY" jsonschema:"description=The layer rectangle(rotation y) y to make the scene for. default is null."`
	// RectangleWidth  *float64 `json:"rectangleWidth" jsonschema:"description=The layer rectangle width to make the scene for. This resizes the layer absolute resolution. default is null."`
	// RectangleHeight *float64 `json:"rectangleHeight" jsonschema:"description=The layer rectangle height to make the scene for. This resizes the layer absolute resolution. default is null."`
}
