import logging
import os
import shutil
import sys
import tempfile
import tarfile
from datetime import datetime
from pathlib import Path

from clu import __about__, panic
from clu.config import get_config
from clu.opsys.factory import opsys_factory
from clu.input import transform_cmdline_to_filename, text_program


log = logging.getLogger(__name__)
cfg = get_config()


def parse_args(subparsers):
    subp_archive = subparsers.add_parser(
        "archive", help="Create an archive of the current system state"
    )
    subp_archive.set_defaults(func=make_archive)


def make_archive() -> int:
    """Create an archive of the current system state."""

    log.info(f"Running command {cfg.cmd} with cfg={cfg}")

    requires = opsys_factory().requires()

    hostname = os.uname().nodename
    work_dir = setup_workdir(hostname)
    try:
        collect_metadata(work_dir, hostname)
        collect_files(requires, work_dir)
        collect_programs(requires, work_dir)
        create_archive(hostname, work_dir)

    finally:
        cleanup_workdir(work_dir)

    return 0


def cleanup_workdir(work_dir):
    shutil.rmtree(work_dir, ignore_errors=True)


def setup_workdir(hostname):
    work_dir = tempfile.mkdtemp(prefix=f"{hostname}.")
    work_path = Path(work_dir)
    if not work_dir or not work_path.is_dir():
        panic("Could not create temp dir")
    return work_dir


def create_archive(hostname, work_dir):
    archive_path = f"/tmp/{__about__.__title__}_{hostname}.tgz"
    with tarfile.open(archive_path, "w:gz") as tar:
        tar.add(work_dir, arcname=".")
    os.chmod(archive_path, 0o644)  # u=rw,go=r
    print(f"Archive created at {archive_path}")


def collect_files(requires, work_dir):
    work_path = Path(work_dir)
    for file in requires.files:
        file_path = Path(file)
        if file_path.is_file():
            dest = work_path / file.lstrip("/")
            dest.parent.mkdir(parents=True, exist_ok=True)
            try:
                shutil.copy(file, dest)
            except Exception as e:
                log.error(f"Failed to copy {file}: {e}")


def collect_programs(requires, work_dir):
    work_path = Path(work_dir)
    prog_dir = work_path / "_programs"
    prog_dir.mkdir(exist_ok=True)
    for prog in requires.programs:
        cmd_name, rc_name = transform_cmdline_to_filename(prog)
        data_path = prog_dir / cmd_name
        rc_path = prog_dir / rc_name
        stdout, rc = text_program(prog)

        with open(data_path, "w") as f:
            # BUG - cant write None
            f.write(stdout)
        # Save return code if nonzero
        if rc != 0:
            with open(rc_path, "w") as f:
                f.write(str(rc))


def collect_metadata(work_dir, hostname):
    work_path = Path(work_dir)
    meta_dir = work_path / "_meta"
    meta_dir.mkdir(exist_ok=True)
    with open(meta_dir / "clu_version", "w") as f:
        f.write(__about__.__version__ + "\n")
    with open(meta_dir / "python_version", "w") as f:
        f.write(".".join(map(str, sys.version_info[:3])) + "\n")
    with open(meta_dir / "hostname", "w") as f:
        f.write(hostname + "\n")
    with open(meta_dir / "path", "w") as f:
        f.write(os.environ.get("PATH", "") + "\n")
    with open(meta_dir / "date", "w") as f:
        f.write(datetime.now().isoformat() + "\n")
