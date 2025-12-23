package subcmd

// Go port of src/clu/cmd/archive.py excluding parse_args. Implements system
// snapshot archiving: gathers required files, program outputs, and metadata
// into a temporary working directory then creates a gzipped tarball.
// No build/test executed per user instruction.

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	pkg "github.com/huntermatthews/clu/pkg"
	"github.com/huntermatthews/clu/pkg/facts/types"
)

// CollectorCmd implements the "collector" subcommand (stub only).
type CollectorCmd struct{}

func (c *CollectorCmd) Run() error {
	fmt.Println("collector: running (stub)")
	return nil
}

// ArchiveConfig controls archive creation (currently only output directory override).
type ArchiveConfig struct {
	OutputDir string // where to place final archive (default /tmp)
}

// MakeArchive executes the archive workflow and returns exit code (0 success).
func MakeArchive(cfg *ArchiveConfig) int {
	if cfg == nil {
		cfg = &ArchiveConfig{}
	}
	outDir := cfg.OutputDir
	if outDir == "" {
		outDir = "/tmp"
	}
	hostname := hostName()
	workDir, err := setupWorkdir(hostname)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return 1
	}
	defer cleanupWorkdir(workDir)

	// Fallback if detectOpSys not configured externally.
	var requires *types.Requires
	osysIface := detectOpSys()
	if osysIface != nil {
		requires = osysIface.Requires()
	} else {
		requires = types.NewRequires() // empty fallback
	}

	collectMetadata(workDir, hostname)
	collectFiles(requires, workDir)
	collectPrograms(requires, workDir)
	archivePath, err := createArchive(hostname, workDir, outDir)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error creating archive: %v\n", err)
		return 1
	}
	fmt.Printf("Archive created at %s\n", archivePath)
	return 0
}

// opsysFactory simplified runtime selection.
// detectOpSys delegates to report command's runtime factory (reuse) to avoid duplicate definition.
// Provided by requires.go; if refactoring later, centralize in a shared package.
var detectOpSys = func() interface{ Requires() *types.Requires } { return nil }

func hostName() string { h, _ := os.Hostname(); return h }

func setupWorkdir(hostname string) (string, error) {
	dir, err := os.MkdirTemp("", fmt.Sprintf("%s.", hostname))
	if err != nil || dir == "" {
		return "", fmt.Errorf("could not create temp dir: %w", err)
	}
	return dir, nil
}

func cleanupWorkdir(dir string) { _ = os.RemoveAll(dir) }

func collectMetadata(workDir, hostname string) {
	metaDir := filepath.Join(workDir, "_meta")
	_ = os.MkdirAll(metaDir, 0o755)
	writeFile(metaDir, "clu_version", pkg.Version+"\n")
	writeFile(metaDir, "go_version", runtime.Version()+"\n")
	writeFile(metaDir, "hostname", hostname+"\n")
	writeFile(metaDir, "path", os.Getenv("PATH")+"\n")
	writeFile(metaDir, "date", time.Now().Format(time.RFC3339)+"\n")
}

func collectFiles(reqs *types.Requires, workDir string) {
	for _, file := range reqs.Files {
		if file == "" {
			continue
		}
		fi, err := os.Stat(file)
		if err != nil || fi.IsDir() {
			continue
		}
		// replicate path structure without leading '/'
		rel := strings.TrimPrefix(file, "/")
		dest := filepath.Join(workDir, rel)
		_ = os.MkdirAll(filepath.Dir(dest), 0o755)
		copyFile(file, dest)
	}
}

func collectPrograms(reqs *types.Requires, workDir string) {
	progDir := filepath.Join(workDir, "_programs")
	_ = os.MkdirAll(progDir, 0o755)
	for _, prog := range reqs.Programs {
		if prog == "" {
			continue
		}
		cmdName, rcName := transformCmdlineToFilename(prog)
		dataPath := filepath.Join(progDir, cmdName)
		rcPath := filepath.Join(progDir, rcName)
		stdout, rc := pkg.CommandRunner(prog)
		writeFileRaw(dataPath, stdout)
		if rc != 0 {
			writeFile(rcPath, "", fmt.Sprintf("%d\n", rc))
		}
	}
}

// transformCmdlineToFilename simplified version mirroring Python input.transform_cmdline_to_filename.
func transformCmdlineToFilename(cmd string) (base string, rc string) {
	// Replace spaces and special chars with underscores; keep alnum and dot.
	cleaned := make([]rune, 0, len(cmd))
	for _, r := range cmd {
		if (r >= 'a' && r <= 'z') || (r >= 'A' && r <= 'Z') || (r >= '0' && r <= '9') || r == '.' || r == '-' {
			cleaned = append(cleaned, r)
		} else {
			cleaned = append(cleaned, '_')
		}
	}
	base = string(cleaned)
	rc = base + "_rc"
	return
}

func createArchive(hostname, workDir, outDir string) (string, error) {
	archivePath := filepath.Join(outDir, fmt.Sprintf("%s_%s.tgz", pkg.Title, hostname))
	f, err := os.Create(archivePath)
	if err != nil {
		return "", err
	}
	defer f.Close()
	gw := gzip.NewWriter(f)
	defer gw.Close()
	tw := tar.NewWriter(gw)
	defer tw.Close()
	// Walk workDir and add files.
	err = filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		rel, _ := filepath.Rel(workDir, path)
		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = rel // ensure relative name
		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}
		src, err := os.Open(path)
		if err != nil {
			return err
		}
		defer src.Close()
		if _, err := io.Copy(tw, src); err != nil {
			return err
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	_ = os.Chmod(archivePath, 0o644)
	return archivePath, nil
}

// Helper to copy file contents.
func copyFile(src, dst string) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer out.Close()
	_, _ = io.Copy(out, in)
}

func writeFile(dir, name, content string) {
	_ = os.WriteFile(filepath.Join(dir, name), []byte(content), 0o644)
}

func writeFileRaw(path, content string) { _ = os.WriteFile(path, []byte(content), 0o644) }
