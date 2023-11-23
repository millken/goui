package yoga

import "unsafe"

type YGExperimentalFeatureSet uint32

type YGConfigFlags struct {
	useWebDefaults       bool `bitfield:"0"`
	printTree            bool `bitfield:"1"`
	cloneNodeUsesContext bool `bitfield:"2"`
	loggerUsesContext    bool `bitfield:"3"`
}

func (flags *YGConfigFlags) SetUseWebDefaults(useWebDefaults bool) YGConfigFlags {
	flags.useWebDefaults = useWebDefaults
	return *flags
}

func (flags *YGConfigFlags) UseWebDefaults() bool {
	return flags.useWebDefaults
}

func (flags *YGConfigFlags) SetPrintTree(printTree bool) YGConfigFlags {
	flags.printTree = printTree
	return *flags
}

func (flags *YGConfigFlags) PrintTree() bool {
	return flags.printTree
}

func (features *YGExperimentalFeatureSet) SetEnabled(feature YGExperimentalFeature, enabled bool) YGExperimentalFeatureSet {
	if enabled {
		*features |= YGExperimentalFeatureSet(feature)
	} else {
		*features &^= YGExperimentalFeatureSet(feature)
	}
	return *features
}

func (features *YGExperimentalFeatureSet) IsEnabled(feature YGExperimentalFeature) bool {
	return *features&YGExperimentalFeatureSet(feature) != 0
}

type YGConfig struct {
	cloneNodeCallback_ struct {
		withContext YGCloneNodeWithContextFn
		noContext   YGCloneNodeFunc
	}
	logger_ struct {
		withContext LogWithContextFn
		noContext   YGLogger
	}
	flags_               YGConfigFlags
	experimentalFeatures YGExperimentalFeatureSet
	errata_              YGErrata
	pointScaleFactor_    float32
	context_             unsafe.Pointer
}

func NewYGConfig(logger YGLogger) *YGConfig {
	return &YGConfig{
		logger_: struct {
			withContext LogWithContextFn
			noContext   YGLogger
		}{
			noContext: logger,
		},
	}
}

func (config *YGConfig) SetUseWebDefaults(useWebDefaults bool) {
	config.flags_ = config.flags_.SetUseWebDefaults(useWebDefaults)
}

func (config *YGConfig) UseWebDefaults() bool {
	return config.flags_.UseWebDefaults()
}

func (config *YGConfig) SetShouldPrintTree(printTree bool) {
	config.flags_ = config.flags_.SetPrintTree(printTree)
}

func (config *YGConfig) ShouldPrintTree() bool {
	return config.flags_.PrintTree()
}

func (config *YGConfig) SetExperimentalFeatureEnabled(feature YGExperimentalFeature, enabled bool) {
	config.experimentalFeatures = config.experimentalFeatures.SetEnabled(feature, enabled)
}

func (config *YGConfig) IsExperimentalFeatureEnabled(feature YGExperimentalFeature) bool {
	return config.experimentalFeatures.IsEnabled(feature)
}

func (config *YGConfig) GetEnabledExperiments() YGExperimentalFeatureSet {
	return config.experimentalFeatures
}

func (config *YGConfig) SetErrata(errata YGErrata) {
	config.errata_ = errata
}

func (config *YGConfig) AddErrata(errata YGErrata) {
	config.errata_ |= errata
}

func (config *YGConfig) RemoveErrata(errata YGErrata) {
	config.errata_ &^= errata
}

func (config *YGConfig) GetErrata() YGErrata {
	return config.errata_
}

func (config *YGConfig) HasErrata(errata YGErrata) bool {
	return config.errata_&errata != 0
}

func (config *YGConfig) SetPointScaleFactor(pointScaleFactor float32) {
	config.pointScaleFactor_ = pointScaleFactor
}

func (config *YGConfig) GetPointScaleFactor() float32 {
	return config.pointScaleFactor_
}

func (config *YGConfig) SetContext(context unsafe.Pointer) {
	config.context_ = context
}

func (config *YGConfig) GetContext() unsafe.Pointer {
	return config.context_
}

func (config *YGConfig) SetLogger(logger YGLogger) {
	config.logger_.noContext = logger
}

func (config *YGConfig) SetLoggerWithContext(logger LogWithContextFn) {
	config.logger_.withContext = logger
}

func (config *YGConfig) SetLoggerToNull() {
	config.logger_.noContext = nil
}

func (config *YGConfig) Log(node *YGNode, level YGLogLevel, context unsafe.Pointer, format string, args ...interface{}) {
	if config.logger_.noContext != nil {
		config.logger_.noContext(config, node, level, format, args...)
	} else if config.logger_.withContext != nil {
		config.logger_.withContext(config, node, level, context, format, args...)
	}
}

func (config *YGConfig) SetCloneNodeCallback(cloneNode YGCloneNodeFunc) {
	config.cloneNodeCallback_.noContext = cloneNode
}

func (config *YGConfig) SetCloneNodeCallbackWithContext(cloneNode YGCloneNodeWithContextFn) {
	config.cloneNodeCallback_.withContext = cloneNode
}

func (config *YGConfig) SetCloneNodeCallbackToNull() {
	config.cloneNodeCallback_.noContext = nil
}

func (config *YGConfig) CloneNode(node *YGNode, owner *YGNode, childIndex int, cloneContext unsafe.Pointer) *YGNode {
	if config.cloneNodeCallback_.noContext != nil {
		return config.cloneNodeCallback_.noContext(node, owner, childIndex)
	} else if config.cloneNodeCallback_.withContext != nil {
		return config.cloneNodeCallback_.withContext(node, owner, childIndex, cloneContext)
	}
	return nil
}
