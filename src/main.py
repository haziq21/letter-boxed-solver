import json
from pathlib import Path
from solver import Solver

DICT_FILEPATH = Path(__file__).parent.resolve() / ".." / "dictionary-250502.json"

# Load the dictionary of all words
with open(DICT_FILEPATH, "r") as f:
    all_words: list[str] = json.load(f)
    all_words = [w.lower() for w in all_words]

# Print all solutions with 3 or fewer words
for sol in Solver(all_words).solutions(["rxn", "aof", "htc", "epi"]):
    if len(sol) <= 3:
        print(" -> ".join(sol))
