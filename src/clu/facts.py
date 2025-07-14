
# I can global if I want to
# (remember, this was ported from a couple of scripts that used globals)
facts = {}

def add_fact(key, value):
    facts[key] = value

def get_fact(key):
    return facts.get(key)

def get_all_facts():
    return facts
