from pathlib import Path
import requests
import re
import json


def get_todays_dictionary() -> list[str]:
    """
    Fetch the Letter Boxed homepage (`https://www.nytimes.com/puzzles/letter-boxed`) and extract the dictionary of words.
    """
    url = "https://www.nytimes.com/puzzles/letter-boxed"
    response = requests.get(url)
    response.raise_for_status()

    # Extract the dictionary from the HTML
    match = re.search(r"\"dictionary\":(\[.+?\])", response.text)
    if match is None:
        raise ValueError("Could not find the dictionary in the HTML response.")

    return [w.lower() for w in json.loads(match.group(1))]


def get_dictionary_from_file(filepath: str | Path) -> list[str]:
    """
    Load the dictionary from a local newline-delimited text file.
    """
    with open(filepath, "r") as f:
        all_words = f.read().splitlines()
    return [w.lower() for w in all_words]
