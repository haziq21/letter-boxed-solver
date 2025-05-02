# from pathlib import Path
from solver import Solver
from dictionary import get_todays_dictionary  # , get_dictionary_from_file


# Alternatively, use get_dictionary_from_file(Path(__file__).parent.resolve() / ".." / ".." / "dictionary.txt")
all_words = get_todays_dictionary()

# Print all solutions with 3 or fewer words
for sol in Solver(all_words).solutions(["rxn", "aof", "htc", "epi"]):
    if len(sol) <= 3:
        print(" -> ".join(sol))
