
def normalize(value):
    ''' Remove useless parts of typical hostinfo related strings.

    1) Remove leading and trailing whitespace
    2) Remove leading and trailing quotation marks (smc can die in a fire)
    3) Lowercase it
    '''
    # 1) remove leading and trailing whitespace
    value = value.strip()

    # 2) Remove leading and trailing quotation marks (smc can die in a fire)
    # Note we only remove them if they are BALANCED (match on both ends)
    if ((value.startswith('"') and value.endswith('"')) or
        (value.startswith("'") and value.endswith("'"))):
         value = value[1:-1]

    # 3) lowercase it
    value = value.lower()

    return value


## END OF LINE ##
