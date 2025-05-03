from itertools import pairwise, product, chain
from collections.abc import Collection
from typing import Generator


def to_set[T](col: Collection[T]) -> set[T]:
    return col if isinstance(col, set) else set(col)


class Solver:
    def __init__(self, dictionary: Collection[str]):
        self.dictionary = to_set(dictionary)
        self.prefix_dict = self.get_prefix_dict(self.dictionary)

    def solutions(self, sides: list[str]) -> Generator[list[str], None, None]:
        allowed_letters = set("".join(sides))
        allowed_words = self.get_allowed_words(sides)
        yield from self.sub_solutions(allowed_words, allowed_letters, [], allowed_words)

    def sub_solutions(
        self,
        wordset: Collection[str],
        unused_letters: Collection[str],
        previous_words: Collection[str],
        allowed_words: Collection[str],
    ) -> Generator[list[str], None, None]:
        # Remove words that have already been used since they won't be helpful
        wordset = to_set(wordset) - to_set(previous_words)
        unused_letters = to_set(unused_letters)

        # Consider all possible words, prioritized by the number of new letters they use
        for word in self.sort_words_by_letter_usage(wordset, unused_letters):
            new_unused_letters = unused_letters - to_set(word)

            # If there are no letters left to use, it means we've found a solution
            if len(new_unused_letters) == 0:
                yield [*previous_words, word]
                continue

            # Cap solutions to 5 words (<4 from previous_words + 1 current word + 1 next word)
            if len(previous_words) < 4:
                candidate_next_words = self.prefix_dict[word[-1]] & to_set(allowed_words)
                yield from self.sub_solutions(
                    candidate_next_words,
                    new_unused_letters,
                    [*previous_words, word],
                    allowed_words,
                )

    def get_allowed_words(self, sides: Collection[str]) -> set[str]:
        """
        Get all words from the dictionary that can be formed with the letters in `sides`.
        Following the rules of the game, two letters from the same side are not allowed
        to be used consecutively, and the minimum number of letters in a word is 2.
        """
        allowed_letters = set("".join(sides))
        disallowed_pairs = set(chain.from_iterable(product(s, repeat=2) for s in sides))
        return set(
            w
            for w in self.dictionary
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

    @classmethod
    def sort_words_by_letter_usage(cls, words: Collection[str], counted_letters: Collection[str]) -> list[str]:
        """
        Sort words by the number of distinct letters from `counted_letters` present in each word.
        """
        counted_letters = to_set(counted_letters)
        letter_usages = [(cls.get_letter_usage(w, counted_letters), w) for w in words]
        return [word for (_, word) in sorted(letter_usages, key=lambda x: x[0], reverse=True)]
