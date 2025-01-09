# Guess Words

Guess Words is a simple Go application that interacts with a `Wordle-like API` to guess words based on feedback. The application reads a list of words from a file, makes guesses, and refines the list of possible guesses based on the feedback received from the API.

## Project Structure

```bash
guess-words/ 
├── .gitignore 
├── go.mod 
├── main.go 
└── words.txt
```

## Prerequisites

- Go 1.18 or later
- A `words.txt` file containing a list of words

## Installation

1. Clone the repository:

    ```sh
    git clone https://github.com/truongpx396/guess-words.git
    cd guess-words
    ```

2. Initialize the Go module:

    ```sh
    go mod tidy
    ```

## Usage
1. Ensure you have a `words.txt` file in the project directory. This file should contain a list of words, one word per line.

2. Update the `main.go` file to specify the desired word length and seed value:
    ```go
    seed := 1238
    filepath := "words.txt"
    wordLength := 8 // Define word length for the Wordle puzzle
    ```

3. Run the application:
    ```go
    go run main.go
    ```

## Project Details

`main.go`
The main application logic is contained in the `main.go` file. It includes the following key functions:

- `loadWordsFromFile`: Reads words from `words.txt` and filters them by length, ensuring no special characters and converting to lowercase.
- `makeGuess`: Sends a guess to the `/random` endpoint and retrieves feedback.
- `filterByCorrectGuesses`: Filters possible guesses based on correct feedback.
- `filterGuesses`: Refines the list of possible guesses based on satisfied feedback.
- `isSatisfactoryFeedback`: Checks if a word satifies the feedback (not containing any absent word).

## Example Output
```sh
listSize: 1000
Making guess: example
Feedback: [{Slot:0 Guess:e Result:correct} {Slot:1 Guess:x Result:absent} ...]
Guessed the word correctly: example
```

## Contributing
Contributions are welcome! Please open an issue or submit a pull request for any improvements or bug fixes.

## License
This project is licensed under the MIT License. See the LICENSE file for details.

## Acknowledgements
- [Wordle API](https://wordle.votee.dev:8000) for providing the word guessing endpoint.
- [Go](https://golang.org) for the programming language.

