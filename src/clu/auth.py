"""Simple credential acquisition for a user - username and password.

Does nothing fancy because we need to run on servers without pass or keyring.
"""

import getpass
import logging

from clu.config import get_config

log = logging.getLogger(__name__)
cfg = get_config()


def _get_username(username):
    """Get username either from getpass library or passed override."""
    if not username:
        username = getpass.getuser()

    if username == "" or username is None:
        raise ValueError("no user could be found for authentication.")

    return username


def _get_password(username):
    """Get password for username by prompting user."""

    password = getpass.getpass(prompt=f"{username} Password: ")

    return password


def get_primary_credentials(username=None) -> tuple[str, str]:
    """Get username and password for user running script."""

    username = _get_username(username)

    password = _get_password(username)

    return (username, password)


# import os
# import getpass

# def get_original_user():
#     """
#     Retrieves the username of the user who invoked sudo,
#     or the current user if sudo was not used.
#     """
#     # Check if SUDO_USER environment variable exists (set by sudo)
#     if 'SUDO_USER' in os.environ:
#         return os.environ['SUDO_USER']
#     else:
#         # If not running with sudo, return the current logged-in user
#         return getpass.getuser()

# original_user = get_original_user()
# print(f"The original user is: {original_user}")
