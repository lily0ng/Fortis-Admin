package backup

import "time"

type Compression string

const (
	CompressionGzip Compression = "gzip"
	CompressionNone Compression = "none"
	CompressionZstd Compression = "zstd"
	CompressionLZ4  Compression = "lz4"
)

type BackupType string

const (
	BackupFull         BackupType = "full"
	BackupIncremental  BackupType = "incremental"
	BackupDifferential BackupType = "differential"
)

type CreateOptions struct {
	TargetDir string
	Sources   []string
	Type      BackupType
	Exclude   []string
	Encrypt   bool
	Compress  Compression
}

type BackupMeta struct {
	ID             string     `json:"id" yaml:"id"`
	CreatedAt      time.Time  `json:"created_at" yaml:"created_at"`
	Type           BackupType `json:"type" yaml:"type"`
	Sources        []string   `json:"sources" yaml:"sources"`
	ArchivePath    string     `json:"archive_path" yaml:"archive_path"`
	SizeBytes      int64      `json:"size_bytes" yaml:"size_bytes"`
	ChecksumSHA256 string     `json:"sha256" yaml:"sha256"`
	Encrypted      bool       `json:"encrypted" yaml:"encrypted"`
	Compression    string     `json:"compression" yaml:"compression"`
	Notes          []string   `json:"notes" yaml:"notes"`
}

type ListOptions struct {
	TargetDir string
	Detailed  bool
	SortBy    string
	Filter    string
	JSON      bool
}

type VerifyOptions struct {
	BackupPath string
	Quick      bool
	Full       bool
	Repair     bool
}

type VerifyResult struct {
	BackupPath string `json:"backup_path"`
	OK         bool   `json:"ok"`
	Reason     string `json:"reason"`
	SHA256     string `json:"sha256"`
}

type RestoreOptions struct {
	BackupPath string
	TargetDir  string
	Items      []string
	DryRun     bool
}

type CatalogOptions struct {
	BackupPath string
	Search     string
	Tree       bool
	Stats      bool
	Extract    string
}
