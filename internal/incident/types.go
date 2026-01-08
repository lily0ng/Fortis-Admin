package incident

import "time"

type EvidenceType string

const (
	EvidenceAll       EvidenceType = "all"
	EvidenceMemory    EvidenceType = "memory"
	EvidenceDisk      EvidenceType = "disk"
	EvidenceNetwork   EvidenceType = "network"
	EvidenceLogs      EvidenceType = "logs"
	EvidenceSystem    EvidenceType = "system"
	EvidenceProcesses EvidenceType = "processes"
)

type CaptureOptions struct {
	CaseID      string
	Types       []string
	OutputDir   string
	Compress    bool
	Integrity   bool
	Encrypt     bool
	Chain       bool
	Verbose     bool
	CollectedBy string
}

type FileHash struct {
	Path   string `json:"path" yaml:"path"`
	SHA256 string `json:"sha256" yaml:"sha256"`
	Size   int64  `json:"size" yaml:"size"`
}

type CaptureManifest struct {
	CaseID      string     `json:"case_id" yaml:"case_id"`
	Timestamp   time.Time  `json:"timestamp" yaml:"timestamp"`
	OutputDir   string     `json:"output_dir" yaml:"output_dir"`
	CollectedBy string     `json:"collected_by" yaml:"collected_by"`
	Files       []FileHash `json:"files" yaml:"files"`
	Notes       []string   `json:"notes" yaml:"notes"`
}

type TriageOptions struct {
	Quick       bool
	Full        bool
	Processes   bool
	Network     bool
	Persistence bool
	Output      string
}

type TimelineOptions struct {
	Source    string
	From      string
	To        string
	Visualize bool
	Export    string
}

type AnalyzeOptions struct {
	Input     string
	IOCFile   string
	Timeline  bool
	Correlate bool
	Report    string
}

type ReportOptions struct {
	Template  string
	Format    string
	Executive bool
	Technical bool
	Evidence  bool
	Output    string
}
