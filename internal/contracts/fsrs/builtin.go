package fsrs

// BuiltinDefaultParameters is the FSRS Python package built-in default (from GetDefaultFsrsParameters RPC).
type BuiltinDefaultParameters struct {
	Parameters         []float64
	DesiredRetention   float64
	MaximumInterval    int
	FSRSPackageVersion string
}
