import logging
import os
import shutil
import sys
import tempfile
import tarfile
from datetime import datetime

from clu import __about__, panic
from clu.os_map import get_os_functions
from clu.readers import transform_cmdline_to_filename, read_program

log = logging.getLogger(__name__)


def do_archive():
    """Create an archive of the current system state."""

    (requires_fn, _) = get_os_functions()
    requires = requires_fn()

    hostname = os.uname().nodename
    work_dir = setup_workdir(hostname)
    try:
        collect_metadata(work_dir, hostname)
        collect_files(requires, work_dir)
        collect_programs(requires, work_dir)
        create_archive(hostname, work_dir)
    finally:
        cleanup_workdir(work_dir)


def cleanup_workdir(work_dir):
    shutil.rmtree(work_dir, ignore_errors=True)


def setup_workdir(hostname):
    work_dir = tempfile.mkdtemp(prefix=f"{hostname}.")
    if not work_dir or not os.path.isdir(work_dir):
        panic("Could not create temp dir")
    return work_dir


def create_archive(hostname, work_dir):
    archive_path = f"/tmp/{hostname}.tgz"
    with tarfile.open(archive_path, "w:gz") as tar:
        tar.add(work_dir, arcname=".")
    os.chmod(archive_path, 0o444)  # ugo+r
    log.info(f"Archive created at {archive_path}")


def collect_files(requires, work_dir):
    for file in requires.files:
        if os.path.isfile(file):
            dest = os.path.join(work_dir, file.lstrip("/"))
            os.makedirs(os.path.dirname(dest), exist_ok=True)
            try:
                shutil.copy(file, dest)
            except Exception as e:
                log.error(f"Failed to copy {file}: {e}")


def collect_programs(requires, work_dir):
    prog_dir = os.path.join(work_dir, "_programs")
    os.makedirs(prog_dir, exist_ok=True)
    for prog in requires.programs:
        cmd_name, rc_name = transform_cmdline_to_filename(prog)
        data_path = os.path.join(prog_dir, cmd_name)
        rc_path = os.path.join(prog_dir, rc_name)
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
