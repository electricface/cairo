package cairo

// internal status
type intStatus int

const (
	intStatusSuccess intStatus = iota
	intStatusNoMemory
	intStatusInvalidRestore
	intStatusInvalidPopGroup
	intStatusNoCurrentPoint
	intStatusInvalidMatrix
	intStatusInvalidintStatus
	intStatusNullPointer
	intStatusInvalidString
	intStatusPathData
	intStatusReadError
	intStatusWriteError
	intStatusSurfaceFinished
	intStatusSurfaceTypeMismatch
	intStatusPatternTypeMismatch
	intStatusInvalidContent
	intStatusInvalidFormat
	intStatusInvalidVisual
	intStatusFileNotFound
	intStatusInvalidDash
	intStatusInvalidDscComment
	intStatusInvalidIndex
	intStatusClipNotRepresentable
	intStatusTempFileError
	intStatusInvalidStride
	intStatusFontTypeMismatch
	intStatusUserFontImmutable
	intStatusUserFontError
	intStatusNegativeCount
	intStatusInvalidClusters
	intStatusInvalidSlant
	intStatusInvalidWeight
	intStatusInvalidSize
	intStatusUserFontNotImplemented
	intStatusDeviceTypeMismatch
	intStatusDeviceError
	intStatusInvalidMeshConstruction
	intStatusDeviceFinished
	intStatusJBIG2GlobalMissing
	intStatusPngError
	intStatusFreeTypeError
	intStatusWin32GdiError
	intStatusTagError
	intStatusLastStatus
)

const (
	intStatusUnsupported intStatus = iota + 100
	intStatusDegenerate
	intStatusNothingToDo
	intStatusFlattenTransparency
	intStatusAnalyzeRecordingSurfacePattern
)

func (s Status) isError() bool {
	return s != StatusSuccess && s < StatusLastStatus
}

func (s intStatus) isError() bool {
	return s != intStatusSuccess && s < intStatusLastStatus
}
