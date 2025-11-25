from argparse import Namespace

_config: Namespace = Namespace()


def set_config(cfg: Namespace):
    # we do the same thing as the argparse module itself.
    # I wish for a better/less janky way to do this, but python import aliasing is hard.
    for key, value in vars(cfg).items():
        setattr(_config, key, value)


def get_config() -> Namespace:
    return _config
