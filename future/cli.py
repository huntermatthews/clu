''' The main() function implementation for program.
'''

import argparse
import sys

from .__about__ import __version__
from .backend import factory


def main():
    if sys.version_info < (3,6):
        print('ERROR: Must use at least python 3.6')
        sys.exit(1)

    # TODO : check for root and bomb -- give up.

    parser = argparse.ArgumentParser(
        description='Show various information about the host',
        prog='hwid',
        )
    parser.add_argument('--version', action='version', version='%%(prog)s %s' % __version__)
    parser.add_argument('--get', dest='list_mode', action='store_false', default=False)
    parser.add_argument('--list', dest='list_mode', action='store_true')
    parser.add_argument('--all', '-A', dest='list_all', action='store_true')
    parser.add_argument('--recurse', '-R', dest='recurse', action='store_true')   ## What is this for -- unused?
    parser.add_argument('--human', '-H', action='store_true')
    parser.add_argument('trait_paths', nargs='*')
    args = parser.parse_args()

    be = factory()

#        print(f'DEBUG: args = {args}')
#        print(f'DEBUG: list_all = {args.list_all}')
#        print(f'DEBUG: list_mode = {args.list_mode}')
#        print(f'DEBUG: paths = {args.trait_paths}')

    if args.list_all:
        args.trait_paths = be.reg.lookup_top_groups()

    # no traits specified. Dang Humans.
    # Instead of erroring out, give them a bent with top level group names
    if not args.trait_paths:
        args.trait_paths = [
            'id.nodename',
            'os.name', 'os.version',
            'sys.vendor', 'sys.model', 'sys.serial', 'sys.asset',
            'phy.ram', 'phy.cpus',
            'bmc.ip_config', 'bmc.ip_address',
        ]

    trait_names = []
    for path in args.trait_paths:
        trait_names.extend(be.reg.convert_path_to_trait_names(path, args.recurse))

    return _output_traits(be, trait_names, args.list_mode, args.human)


def _output_traits(be, trait_names, list_only=False, human=False):
    ''' BUG: This methods args is a dumpster fire - what makes more sense here?
    '''

    if not trait_names:
        print(f'ERROR: no valid paths given.')
        return 1

    # regardless of mode, we _should_ now have a simple list of trait names
    for name in trait_names:
        trait = be.reg.get_trait(name)
        if trait:
            #print('debug: lookup trait worked')
            if list_only:
                print(f'{trait.name}')
            else:
                if human:
                    value = trait.human
                else:
                    value = trait.value
                print(f'{trait.name} = {value}')
        else:
            # Tbes should NEVER happen, as the output of _path_to_trait_names should
            # only be valid trait names
            print(f'error: trait = {trait} is not valid')
            return 1
    return 0

if __name__ == "__main__":
    main()

## END OF LINE ##
