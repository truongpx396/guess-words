package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"regexp"
	"strings"
)

const (
	baseURL = "https://wordle.votee.dev:8000"
)

// GuessResult represents the structure of the response from the API for a guess
type GuessResult struct {
	Slot   int    `json:"slot"`
	Guess  string `json:"guess"`
	Result string `json:"result"` // "absent", "present", or "correct"
}

// makeGuess sends a guess to the /random endpoint and retrieves feedback
func makeGuess(guess string, size int, seed int) ([]GuessResult, error) {
	apiURL := fmt.Sprintf("%s/random", baseURL)

	// Add query parameters
	params := url.Values{}
	params.Add("guess", guess)
	params.Add("size", fmt.Sprintf("%d", size))
	params.Add("seed", fmt.Sprintf("%d", seed)) // Add the seed parameter

	fullURL := fmt.Sprintf("%s?%s", apiURL, params.Encode())

	// Send GET request
	resp, err := http.Get(fullURL)
	if err != nil {
		return nil, fmt.Errorf("error making guess: %v", err)
	}
	defer resp.Body.Close()

	// Check for HTTP errors
	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	// Decode JSON response
	var results []GuessResult
	if err := json.NewDecoder(resp.Body).Decode(&results); err != nil {
		return nil, fmt.Errorf("error decoding response: %v", err)
	}

	return results, nil
}

// fiter possible guesses based on correct feedback
func filterByCorrectGuesses(possibleGuesses []string, correctFeedback []GuessResult) []string {
	var refinedGuesses []string
	for _, word := range possibleGuesses {
		match := true
		for _, feedbackItem := range correctFeedback {
			char := string(word[feedbackItem.Slot])
			if char != feedbackItem.Guess {
				match = false
				break
			}
		}
		if match {
			refinedGuesses = append(refinedGuesses, word)
		}
	}
	return refinedGuesses
}

func filterGuesses(possibleGuesses []string, feedback []GuessResult) []string {
	var refinedGuesses []string

	// Filter out correct feedback

	var correctFeedback []GuessResult
	for _, feedbackItem := range feedback {
		if feedbackItem.Result == "correct" {
			correctFeedback = append(correctFeedback, feedbackItem)
		}
	}

	if len(correctFeedback) > 0 {
		possibleGuesses = filterByCorrectGuesses(possibleGuesses, correctFeedback)
	}

	for _, word := range possibleGuesses {
		if isSatisfactoryFeedback(word, feedback) {
			refinedGuesses = append(refinedGuesses, word)
		}
	}

	return refinedGuesses
}

func isSatisfactoryFeedback(word string, feedback []GuessResult) bool {
	for _, feedbackItem := range feedback {
		char := feedbackItem.Guess
		switch feedbackItem.Result {
		case "absent":
			if strings.Contains(word, char) {
				return false
			}
		}
	}
	return true
}

func loadWordsFromFile(filepath string, lenght int) ([]string, error) {
	file, err := os.Open(filepath)
	if err != nil {
		return nil, fmt.Errorf("err opening file %v", err)
	}

	var filterdWords []string
	scanner := bufio.NewScanner(file)
	re := regexp.MustCompile("^[a-zA-Z]+$")

	for scanner.Scan() {
		word := strings.TrimSpace(scanner.Text())
		if len(word) == lenght && re.MatchString(word) {
			filterdWords = append(filterdWords, strings.ToLower(word))
		}
	}

	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("err scanner %v", err)
	}
	return filterdWords, nil
}

func main() {

	seed := 1238

	// Specify the file path and desired word length
	filepath := "words.txt"
	wordLength := 8 // Define word length for the Wordle puzzle

	// Load words of the given length from the file
	possibleGuesses, err := loadWordsFromFile(filepath, wordLength)
	if err != nil {
		log.Fatalf("Failed to load word list: %v", err)
	}

	fmt.Println("listSize: ", len(possibleGuesses))

	for {
		// Use the first word from the possible guesses as the guess
		guess := possibleGuesses[0]
		fmt.Printf("Making guess: %s\n", guess)

		// Make a guess
		feedback, err := makeGuess(guess, wordLength, seed)
		if err != nil {
			log.Fatalf("Error during guess: %v", err)
		}

		// Print feedback
		fmt.Printf("Feedback: %+v\n", feedback)

		// Check if the guess is correct
		correct := true
		for _, item := range feedback {
			if item.Result != "correct" {
				correct = false
				break
			}
		}
		if correct {
			fmt.Printf("Guessed the word correctly: %s\n", guess)
			break
		}

		// Refine the list of possible guesses based on feedback
		possibleGuesses = filterGuesses(possibleGuesses[1:], feedback)
		if len(possibleGuesses) == 0 {
			log.Fatalf("No possible guesses left. Something went wrong.")
		}
	}
}
