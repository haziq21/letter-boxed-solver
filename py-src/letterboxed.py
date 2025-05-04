from itertools import pairwise, product, chain
from collections.abc import Collection
from typing import Generator


def to_set[T](col: Collection[T]) -> set[T]:
    return col if isinstance(col, set) else set(col)


class LetterBoxed:
    def __init__(self, dictionary: Collection[str], sides: Collection[str]):
        self.dictionary = self.get_allowed_words(dictionary, sides)
        self.prefix_dict = self.get_prefix_dict(self.dictionary)

    def solutions(self, max_words: int) -> Generator[list[str], None, None]:
        yield from self.sub_solutions([], max_words)

    def sub_solutions(self, previous_words: list[str], max_words: int) -> Generator[list[str], None, None]:
        wordset: set[str]

        if previous_words:
            starting_letter = previous_words[-1][-1]
            wordset = self.prefix_dict[starting_letter]
        else:
            wordset = self.dictionary

        # Remove words that have already been used since they won't be helpful
        wordset = wordset - set(previous_words)
        unused_letters = self.allowed_letters() - set(chain.from_iterable(previous_words))

        # Consider all possible words, prioritized by the number of new letters they use
        for word in self.sort_words_by_letter_usage(wordset, unused_letters):
            new_unused_letters = unused_letters - to_set(word)

            # If there are no letters left to use, it means we've found a solution
            if len(new_unused_letters) == 0:
                yield [*previous_words, word]
                continue

            if len(previous_words) + 1 < max_words:
                yield from self.sub_solutions([*previous_words, word], max_words)

    @staticmethod
    def get_allowed_words(dictionary: Collection[str], sides: Collection[str]) -> set[str]:
        """
        Get all words from the dictionary that can be formed with the letters in `sides`.
        Following the rules of the game, two letters from the same side are not allowed
        to be used consecutively, and the minimum number of letters in a word is 2.
        """
        allowed_letters = set("".join(sides))
        disallowed_pairs = set(chain.from_iterable(product(s, repeat=2) for s in sides))
        return set(
            w
            for w in dictionary
            # Minimum word length is 3
            if len(w) >= 3
            # The word must only use letters from the allowed letters
            and set(w) <= allowed_letters
            # The word must not contain two consecutive letters from the same side
            and all(p not in disallowed_pairs for p in pairwise(w))
        )

    @staticmethod
    def get_letter_usage(word: str, counted_letters: Collection[str]) -> int:
        """
        Count the number of distinct letters from `counted_letters` present in `word`.
        """
        return len(set(word) & to_set(counted_letters))

    @staticmethod
    def get_prefix_dict(words: Collection[str]) -> dict[str, set[str]]:
        """
        Create a dictionary mapping starting letters to words that start with that letter.
        """
        prefix_dict: dict[str, set[str]] = {}

        for word in words:
            prefix = word[0]

            if prefix not in prefix_dict:
                prefix_dict[prefix] = set()

            prefix_dict[prefix].add(word)

        return prefix_dict

    def sort_words_by_letter_usage(self, words: Collection[str], counted_letters: Collection[str]) -> list[str]:
        """
        Sort words by the number of distinct letters from `counted_letters` present in each word.
        """
        counted_letters = to_set(counted_letters)
        letter_usages = {w: self.get_letter_usage(w, counted_letters) for w in words}
        return sorted(letter_usages, key=lambda w: letter_usages[w], reverse=True)

    def allowed_letters(self) -> set[str]:
        """
        Get the set of letters that can be used in the game.
        """
        return set(self.prefix_dict.keys())
