import os
import shutil
import subprocess
import tempfile
import tarfile
from datetime import datetime

def do_archive():
    """Create an archive of the current system state."""
    return
    hostname = os.uname().nodename
    work_dir = setup_workdir(hostname)
    try:
        collect_metadata(work_dir, hostname)
        collect_files(work_dir)
        collect_programs(work_dir)
        collect_uname(work_dir)
        create_archive(hostname, work_dir)
    finally:
        cleanup_workdir(work_dir)


def cleanup_workdir(work_dir):
    shutil.rmtree(work_dir, ignore_errors=True)


def setup_workdir(hostname):
    work_dir = tempfile.mkdtemp(prefix=f"{hostname}.", dir="/tmp")
    if not work_dir or not os.path.isdir(work_dir):
        print("ERROR: Could not create temp dir")
        exit(1)
    return work_dir


def create_archive(hostname, work_dir):
    archive_path = f"/tmp/{hostname}.tgz"
    with tarfile.open(archive_path, "w:gz") as tar:
        tar.add(work_dir, arcname=".")
    os.chmod(archive_path, 0o444)  # ugo+r
    print(f"Archive created at {archive_path}")


def collect_files(work_dir):
    for file in FILES:
        if os.path.isfile(file):
            dest = os.path.join(work_dir, file.lstrip("/"))
            os.makedirs(os.path.dirname(dest), exist_ok=True)
            try:
                shutil.copy2(file, dest)
            except Exception as e:
                print(f"Failed to copy {file}: {e}")


def collect_programs(work_dir):
    prog_dir = os.path.join(work_dir, "_programs")
    os.makedirs(prog_dir, exist_ok=True)
    for prog in PROGRAMS:
        if isinstance(prog, list):
            cmd = prog[0]
            args = prog[1:]
        else:
            cmd = prog
            args = []
        prog_path = shutil.which(cmd) if "|" not in " ".join(prog) else None
        out_file = os.path.join(prog_dir, cmd)
        if prog_path and os.access(prog_path, os.X_OK):
            try:
                # Handle pipes by running in shell if needed
                if "|" in " ".join(prog):
                    result = subprocess.run(
                        " ".join(prog), shell=True, capture_output=True, text=True
                    )
                else:
                    result = subprocess.run(
                        [cmd] + args, capture_output=True, text=True
                    )
                with open(out_file, "w") as f:
                    f.write(result.stdout)
                # Save return code if nonzero
                if result.returncode != 0:
                    with open(out_file + "_rc", "w") as f:
                        f.write(str(result.returncode))
            except Exception as e:
                with open(out_file, "w") as f:
                    f.write(f"{cmd} failed: {e}")
        else:
            with open(out_file, "w") as f:
                f.write(f"{cmd} not found")


def collect_uname(work_dir):
    prog_dir = os.path.join(work_dir, "_programs")
    out_file = os.path.join(prog_dir, "uname-opts")
    opts = "a b i K m n o p r s U v".split()
    with open(out_file, "w") as f:
        for opt in opts:
            try:
                result = subprocess.run(
                    ["uname", f"-{opt}"], capture_output=True, text=True
                )
                f.write(f"uname -{opt} {result.stdout.strip()}\n")
            except Exception as e:
                f.write(f"uname -{opt} failed: {e}\n")


def collect_metadata(work_dir, hostname):
    meta_dir = os.path.join(work_dir, "_meta")
    os.makedirs(meta_dir, exist_ok=True)
    with open(os.path.join(meta_dir, "hostname"), "w") as f:
        f.write(hostname + "\n")
    with open(os.path.join(meta_dir, "path"), "w") as f:
        f.write(os.environ.get("PATH", "") + "\n")
    with open(os.path.join(meta_dir, "date"), "w") as f:
        f.write(datetime.now().isoformat() + "\n")


def main():
    hostname = os.uname().nodename
    work_dir = setup_workdir(hostname)
    try:
        collect_metadata(work_dir, hostname)
        collect_files(work_dir)
        collect_programs(work_dir)
        collect_uname(work_dir)
        create_archive(hostname, work_dir)
    finally:
        cleanup_workdir(work_dir)


if __name__ == "__main__":
    main()
