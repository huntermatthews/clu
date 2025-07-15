import os
import shutil
import sys
import tempfile
import tarfile
from datetime import datetime

from clu import requires, __about__
from clu.debug import debug_var, debug
from clu.report_requires import get_os_requirements
from clu.readers import transform_cmdline_to_filename, read_program
from clu.requires import get_all_requires


def do_archive():
    """Create an archive of the current system state."""

    get_os_requirements()
    debug_var("requires", requires)

    hostname = os.uname().nodename
    work_dir = setup_workdir(hostname)
    try:
        collect_metadata(work_dir, hostname)
        collect_files(work_dir)
        collect_programs(work_dir)
        create_archive(hostname, work_dir)
    finally:
        cleanup_workdir(work_dir)


def cleanup_workdir(work_dir):
    shutil.rmtree(work_dir, ignore_errors=True)


def setup_workdir(hostname):
    work_dir = tempfile.mkdtemp(prefix=f"{hostname}.")
    if not work_dir or not os.path.isdir(work_dir):
        print("ERROR: Could not create temp dir")
        sys.exit(1)
    return work_dir


def create_archive(hostname, work_dir):
    archive_path = f"/tmp/{hostname}.tgz"
    with tarfile.open(archive_path, "w:gz") as tar:
        tar.add(work_dir, arcname=".")
    os.chmod(archive_path, 0o444)  # ugo+r
    print(f"Archive created at {archive_path}")


def collect_files(work_dir):
    for file in get_all_requires()["files"]:
        if os.path.isfile(file):
            dest = os.path.join(work_dir, file.lstrip("/"))
            os.makedirs(os.path.dirname(dest), exist_ok=True)
            try:
                shutil.copy(file, dest)
            except Exception as e:
                print(f"Failed to copy {file}: {e}")


def collect_programs(work_dir):
    prog_dir = os.path.join(work_dir, "_programs")
    os.makedirs(prog_dir, exist_ok=True)
    for prog in get_all_requires()["programs"]:
        debug_var("collect_programs:prog", prog)
        cmd_name, rc_name = transform_cmdline_to_filename(prog)
        data_path = os.path.join(prog_dir, cmd_name)
        rc_path = os.path.join(prog_dir, rc_name)
        debug("data_file", data_path)
        debug("rc_file", rc_path)
        stdout, rc = read_program(prog)

        with open(data_path, "w") as f:
            f.write(stdout)
        # Save return code if nonzero
        if rc != 0:
            with open(rc_path, "w") as f:
                f.write(str(rc))


def collect_metadata(work_dir, hostname):
    meta_dir = os.path.join(work_dir, "_meta")
    os.makedirs(meta_dir, exist_ok=True)
    with open(os.path.join(meta_dir, "clu_version"), "w") as f:
        f.write(__about__.__version__ + "\n")
    with open(os.path.join(meta_dir, "python_version"), "w") as f:
        f.write(".".join(map(str, sys.version_info[:3])) + "\n")
    with open(os.path.join(meta_dir, "hostname"), "w") as f:
        f.write(hostname + "\n")
    with open(os.path.join(meta_dir, "path"), "w") as f:
        f.write(os.environ.get("PATH", "") + "\n")
    with open(os.path.join(meta_dir, "date"), "w") as f:
        f.write(datetime.now().isoformat() + "\n")
