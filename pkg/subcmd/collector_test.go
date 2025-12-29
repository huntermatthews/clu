package subcmd

import (
    "archive/tar"
    "compress/gzip"
    "os"
    "path/filepath"
    "testing"
)

// TestCreateCollectionIncludesFiles ensures files under workDir are included in the tarball.
func TestCreateCollectionIncludesFiles(t *testing.T) {
    workDir := t.TempDir()
    outDir := t.TempDir()

    // Create directories and files in workDir
    if err := os.MkdirAll(filepath.Join(workDir, "_meta"), 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    if err := os.MkdirAll(filepath.Join(workDir, "etc"), 0o755); err != nil {
        t.Fatalf("mkdir: %v", err)
    }
    if err := os.WriteFile(filepath.Join(workDir, "_meta", "hostname"), []byte("host123\n"), 0o644); err != nil {
        t.Fatalf("write hostname: %v", err)
    }
    if err := os.WriteFile(filepath.Join(workDir, "etc", "os-release"), []byte("NAME=TestOS\n"), 0o644); err != nil {
        t.Fatalf("write os-release: %v", err)
    }

    tarPath, err := createCollection2("host123", workDir, outDir)
    if err != nil {
        t.Fatalf("createCollection: %v", err)
    }

    // Open and inspect tar.gz
    f, err := os.Open(tarPath)
    if err != nil {
        t.Fatalf("open tar: %v", err)
    }
    defer f.Close()
    gr, err := gzip.NewReader(f)
    if err != nil {
        t.Fatalf("gzip reader: %v", err)
    }
    defer gr.Close()
    tr := tar.NewReader(gr)

    entries := map[string]bool{}
    for {
        hdr, err := tr.Next()
        if err != nil {
            if err.Error() == "EOF" {
                break
            }
            t.Fatalf("tar next: %v", err)
        }
        entries[hdr.Name] = true
    }

    if !entries["_meta/hostname"] {
        t.Errorf("missing _meta/hostname in tar; got entries: %v", entries)
    }
    if !entries["etc/os-release"] {
        t.Errorf("missing etc/os-release in tar; got entries: %v", entries)
    }
}
