# Running GitHub Actions Locally with `act` on ARM64 Mac

To get `act` working on your ARM64 Mac, here's what you'll need:

## Installation

```bash
# Using Homebrew (easiest)
brew install act

# Or using curl
curl https://raw.githubusercontent.com/nektos/act/master/install.sh | sudo bash
```

## Key ARM64 Considerations

**1. Docker Desktop Required**

- You'll need Docker Desktop for Mac with ARM64 support
- Make sure it's configured to use the Apple Silicon optimizations

**2. Platform Emulation Issues**

- Most GitHub Actions use `ubuntu-latest` which defaults to x86_64 images
- These will run under emulation (slower but functional)
- You can override with ARM64 images where available

**3. Configuration Options**

Create `~/.actrc` to specify ARM64-compatible images:

```
--platform linux/arm64
--container-architecture linux/arm64
```

Or use a `.actrc` in your project:

```
-P ubuntu-latest=catthehacker/ubuntu:act-latest-arm64
-P ubuntu-20.04=catthehacker/ubuntu:act-20.04-arm64
```

**4. Testing Your Workflow**

```bash
# Run a specific job
act -j build

# Run with specific platform
act --platform linux/arm64

# Dry run to see what would happen
act --dry-run

# Use different runner image
act -P ubuntu-latest=catthehacker/ubuntu:act-latest-arm64
```

**5. Limitations on ARM64**

- Some actions may not have ARM64 variants
- Container builds (like your CentOS 7 job) may be slower due to emulation
- Go builds should work fine since Go has excellent cross-compilation

**6. For Your Specific Workflow**
Your workflow should mostly work, but the `centos:7` container will run under emulation. You could test specific jobs:

```bash
# Test just the build job
act -j build

# Test RPM build (will be slow due to emulation)
act -j build-rpm-el7
```

## Recommended Setup for This Project

1. Install act via Homebrew
2. Create a project-specific `.actrc` file with ARM64 optimizations
3. Test individual jobs first before running the full workflow
4. Use `--dry-run` to validate workflow parsing without execution
