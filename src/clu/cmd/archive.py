from clu import config


def setup_args(subparsers):
    subp_archive = subparsers.add_parser("archive")
    subp_archive.set_defaults(func=run)

    subp_archive.add_argument("speed", type=int)


def run():
    print(f"Running command {config.cmd} with args={config}")
