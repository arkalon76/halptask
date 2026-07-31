package updater

import (
	"archive/tar"
	"archive/zip"
	"bytes"
	"compress/gzip"
	"testing"
)

func TestCompareVersions(t *testing.T) {
	tests := []struct {
		v1       string
		v2       string
		expected int // 1 if v2 > v1, 0 if v2 == v1, -1 if v2 < v1
	}{
		{"0.1.0", "0.2.0", 1},
		{"v0.1.0", "v0.2.0", 1},
		{"0.1.0", "0.1.0", 0},
		{"v1.2.3", "1.2.3", 0},
		{"0.2.0", "0.1.0", -1},
		{"1.0.0", "2.0.0", 1},
		{"0.1.0", "0.1.1", 1},
	}

	for _, tt := range tests {
		res := CompareVersions(tt.v1, tt.v2)
		if res != tt.expected {
			t.Errorf("CompareVersions(%q, %q) = %d, expected %d", tt.v1, tt.v2, res, tt.expected)
		}
	}
}

func TestMatchAsset(t *testing.T) {
	rel := &ReleaseInfo{
		Assets: []ReleaseAsset{
			{Name: "halptask_Darwin_x86_64.tar.gz", BrowserDownloadURL: "https://example.com/darwin_x86_64"},
			{Name: "halptask_Darwin_arm64.tar.gz", BrowserDownloadURL: "https://example.com/darwin_arm64"},
			{Name: "halptask_Linux_x86_64.tar.gz", BrowserDownloadURL: "https://example.com/linux_x86_64"},
			{Name: "halptask_Windows_x86_64.zip", BrowserDownloadURL: "https://example.com/windows_x86_64"},
		},
	}

	matchAsset(rel)

	if rel.AssetURL == "" {
		t.Errorf("matchAsset failed to find any matching asset for platform")
	}
}

func TestExtractBinaryTarGz(t *testing.T) {
	var buf bytes.Buffer
	gw := gzip.NewWriter(&buf)
	twWriter := tar.NewWriter(gw)
	dummyContent := []byte("dummy binary content")
	hdr := &tar.Header{
		Name: "halptask",
		Mode: 0755,
		Size: int64(len(dummyContent)),
	}
	if err := twWriter.WriteHeader(hdr); err != nil {
		t.Fatalf("Failed to write tar header: %v", err)
	}
	if _, err := twWriter.Write(dummyContent); err != nil {
		t.Fatalf("Failed to write tar content: %v", err)
	}
	twWriter.Close()
	gw.Close()

	extracted, err := extractBinary("halptask_Darwin_arm64.tar.gz", buf.Bytes())
	if err != nil {
		t.Fatalf("extractBinary failed: %v", err)
	}
	if string(extracted) != string(dummyContent) {
		t.Errorf("Extracted content mismatch: got %q, want %q", string(extracted), string(dummyContent))
	}
}

func TestExtractBinaryZip(t *testing.T) {
	var buf bytes.Buffer
	zw := zip.NewWriter(&buf)
	dummyContent := []byte("dummy zip content")

	f, err := zw.Create("halptask.exe")
	if err != nil {
		t.Fatalf("Failed to create zip entry: %v", err)
	}
	if _, err := f.Write(dummyContent); err != nil {
		t.Fatalf("Failed to write zip content: %v", err)
	}
	zw.Close()

	extracted, err := extractBinary("halptask_Windows_x86_64.zip", buf.Bytes())
	if err != nil {
		t.Fatalf("extractBinary failed for zip: %v", err)
	}
	if string(extracted) != string(dummyContent) {
		t.Errorf("Extracted zip content mismatch: got %q, want %q", string(extracted), string(dummyContent))
	}
}

func TestCanUpdate(t *testing.T) {
	canUpdate, realPath, _ := CanUpdate()
	t.Logf("CanUpdate: %v, Path: %s", canUpdate, realPath)
	if realPath == "" {
		t.Errorf("CanUpdate returned empty executable path")
	}
}
