from pathlib import Path
from letterboxed import LetterBoxed
from dictionary import get_todays_dictionary, get_dictionary_from_file


# all_words = get_todays_dictionary()
all_words = get_dictionary_from_file(Path(__file__).parent.resolve() / ".." / "dictionary.txt")

num_sols = 0

# Print all solutions with 3 or fewer words
for sol in LetterBoxed(all_words, ["rxn", "aof", "htc", "epi"]).solutions(3):
    print(" -> ".join(sol))
    num_sols += 1

print(f"Found {num_sols} solutions")
