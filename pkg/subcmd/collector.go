package subcmd

// Go port of src/clu/cmd/collection.py excluding parse_args. Implements system
// snapshot archiving: gathers required files, program outputs, and metadata
// into a temporary working directory then creates a gzipped tarball.
// No build/test executed per user instruction.

import (
	"archive/tar"
	"compress/gzip"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/huntermatthews/clu/pkg/facts"
	"github.com/huntermatthews/clu/pkg/facts/types"
	"github.com/huntermatthews/clu/pkg/global"
	"github.com/huntermatthews/clu/pkg/input"
)

// CollectorCmd implements the "collector" subcommand (stub only).
type CollectorCmd struct {
	OutDir string `name:"out-dir" help:"Output directory for the tarball." default:"/tmp"`
}

func (c *CollectorCmd) Run(stdout input.Stdout, stderr input.Stderr) error {
	fmt.Fprintln(stdout, "collector: running")

	// Ensure output directory exists (Kong provides default via struct tag)
	if err := os.MkdirAll(c.OutDir, 0o755); err != nil {
		return fmt.Errorf("failed to create out-dir %s: %w", c.OutDir, err)
	}

	hostname, _ := os.Hostname()
	workDir, err := setupWorkdir(hostname)
	if err != nil {
		return err
	}
	defer cleanupWorkdir(workDir)

	osys := facts.OpSysFactory()
	requires := osys.Requires()

	if err := collectMetadata(workDir, hostname); err != nil {
		return fmt.Errorf("collect metadata: %w", err)
	}

	if err := collectFiles(requires, workDir); err != nil {
		return fmt.Errorf("collect files: %w", err)
	}

	if err := collectPrograms(requires, workDir); err != nil {
		return fmt.Errorf("collect programs: %w", err)
	}

	collectionPath, err := createCollection(hostname, workDir, c.OutDir)
	if err != nil {
		return fmt.Errorf("error creating collection: %w", err)
	}

	fmt.Fprintf(stdout, "Collection created at %s\n", collectionPath)
	return nil
}

func setupWorkdir(hostname string) (string, error) {
	dir, err := os.MkdirTemp("", fmt.Sprintf("%s.", hostname))
	if err != nil || dir == "" {
		return "", fmt.Errorf("could not create temp dir: %w", err)
	}
	return dir, nil
}

func cleanupWorkdir(dir string) error { return os.RemoveAll(dir) }

func collectMetadata(workDir, hostname string) error {
	metaDir := filepath.Join(workDir, "_meta")
	if err := os.MkdirAll(metaDir, 0o755); err != nil {
		return fmt.Errorf("failed to create metadata dir %s: %w", metaDir, err)
	}

	path := filepath.Join(metaDir, "clu_version")
	if err := os.WriteFile(path, []byte(global.Version+"\n"), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	path = filepath.Join(metaDir, "go_version")
	if err := os.WriteFile(path, []byte(runtime.Version()+"\n"), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	path = filepath.Join(metaDir, "hostname")
	if err := os.WriteFile(path, []byte(hostname+"\n"), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	path = filepath.Join(metaDir, "path")
	if err := os.WriteFile(path, []byte(os.Getenv("PATH")+"\n"), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	path = filepath.Join(metaDir, "date")
	if err := os.WriteFile(path, []byte(time.Now().Format(time.RFC3339)+"\n"), 0o644); err != nil {
		return fmt.Errorf("write %s: %w", path, err)
	}

	return nil
}

func collectFiles(reqs *types.Requires, workDir string) error {
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
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			return fmt.Errorf("create dir for %s: %w", dest, err)
		}
		copyFile(file, dest)
	}
	return nil
}

func collectPrograms(reqs *types.Requires, workDir string) error {
	progDir := filepath.Join(workDir, "_programs")
	if err := os.MkdirAll(progDir, 0o755); err != nil {
		return fmt.Errorf("create programs dir %s: %w", progDir, err)
	}
	for _, prog := range reqs.Programs {
		if prog == "" {
			continue
		}
		cmdName, rcName := transformCmdlineToFilename(prog)
		dataPath := filepath.Join(progDir, cmdName)
		rcPath := filepath.Join(progDir, rcName)
		stdout, rc, perr := input.CommandRunner(prog)
		if perr != nil {
			// Non-fatal: still record whatever output we have and rc below.
			log.Printf("collector: command runner error for %q: %v", prog, perr)
		}
		if err := os.WriteFile(dataPath, []byte(stdout), 0o644); err != nil {
			return fmt.Errorf("write %s: %w", dataPath, err)
		}
		if rc != 0 {
			if err := os.WriteFile(rcPath, []byte(fmt.Sprintf("%d\n", rc)), 0o644); err != nil {
				return fmt.Errorf("write %s: %w", rcPath, err)
			}
		}
	}
	return nil
}

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
	if _, err := io.Copy(out, in); err != nil {
		log.Printf("collector: copy %s -> %s failed: %v", src, dst, err)
	}
}

func createCollection(hostname, workDir, outDir string) (string, error) {
	outPath := filepath.Join(outDir, fmt.Sprintf("%s_%s.tgz", global.GetAppInfo().Name, hostname))
	f, err := os.Create(outPath)
	if err != nil {
		return "", err
	}
	defer f.Close()

	gw := gzip.NewWriter(f)
	// do not defer; close with error handling below

	tw := tar.NewWriter(gw)
	// do not defer; close with error handling below

	err = filepath.Walk(workDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, relErr := filepath.Rel(workDir, path)
		if relErr != nil {
			return relErr
		}

		// Skip root entry
		if rel == "." {
			return nil
		}

		hdr, err := tar.FileInfoHeader(info, "")
		if err != nil {
			return err
		}
		hdr.Name = rel

		if err := tw.WriteHeader(hdr); err != nil {
			return err
		}

		if info.IsDir() {
			return nil
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
	if cerr := tw.Close(); cerr != nil {
		return "", fmt.Errorf("closing tar writer: %w", cerr)
	}
	if cerr := gw.Close(); cerr != nil {
		return "", fmt.Errorf("closing gzip writer: %w", cerr)
	}
	if err := os.Chmod(outPath, 0o644); err != nil {
		log.Printf("collector: chmod failed for %s: %v", outPath, err)
	}
	return outPath, nil
}
