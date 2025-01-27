/*
wordle_solver.go is a simple script used to assist/do wordle


preq: go
	  github.com/dlclark/regexp2  (note can use go get as noted by go.mod)


*/

package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"math/rand"
	"os"
	"strings"
	"time"
	"unicode"

	"github.com/dlclark/regexp2"
)

// Helper function to check if a slice of runes contains a specific character
func containsRune(runes []rune, char rune) bool {
	for _, r := range runes {
		if r == char {
			return true
		}
	}
	return false
}

// Helper regex generator function that simply places the characters if they exist into the regex template
func regexGenerator(knownCharacters []rune, knownBadCharacters []rune, word []rune) string {
	regex := "^"
	shouldNotContainStart := "(?!.*["
	shouldNotContainClosing := "])"
	shouldContainStart := "(?=.*["
	shouldContainEnd := "])"

	// See if we have both known bad and known good characters
	if len(knownBadCharacters) != 0 && len(knownCharacters) != 0 {

		regex += shouldNotContainStart + string(knownBadCharacters) + shouldNotContainClosing + shouldContainStart + string(knownCharacters) + shouldContainEnd + string(word) + "$"
		// check if we atleast have known bad characters
	} else if len(knownBadCharacters) != 0 {
		regex += shouldNotContainStart + string(knownBadCharacters) + shouldNotContainClosing + string(word) + "$"
	} else {
		regex += shouldContainStart + string(knownCharacters) + shouldContainEnd + string(word) + "$"
	}

	// DEBUG PRINT
	//fmt.Println(regex)
	return regex

}

func main() {

	// doubly linked list used as generic dynamic array
	l := list.New()
	rand.Seed(time.Now().UnixNano())

	var wordLength int

	// ask the user for the words length (tested with 5 but should work dynamically)
	fmt.Print("Please enter the wordsLength: ")
	_, err := fmt.Scan(&wordLength)
	if err != nil {
		fmt.Println("Entered a non number exiting")
		return
	}

	// FILL IN WITH YOUR WORDLIST for obvious reasons I wont include mine
	file, err := os.Open("")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)

	}
	defer file.Close()

	reader := bufio.NewReader(file)

	// read in the wordlist making sure to trim the newline
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading file: %v", err)
		}
		line = strings.TrimSpace(line)
		if len(line) == wordLength {
			l.PushBack(line)
		}
	}

	// DEBUG
	/*
		for e := l.Front(); e != nil; e = e.Next() {
			fmt.Println(e.Value)

		}
	*/

	var matchBox string
	var knownCharacters []rune
	var knownBadCharacters []rune
	currentGuess := []rune{'t', 'a', 'l', 'e', 's'}
	var regex string
	enteredWordValid := "N"
	var attemptedGuess string
	var wantGuess string

	// set up the word
	word := make([]rune, wordLength)
	for i := range word {
		word[i] = '.'
	}

	// tales is the best first guess technically so also suggest it first but allow for another word
	fmt.Println("Please enter tales for the first guess (statistically the best word for 5 letter word length)")
	for {
		fmt.Println("Enter T for matchbox to exit. Flow of the program is simple. You enter a match box (example .,.f.) with dots indicating the letter doesn't exist, commas indicating it does exist, and a letter indicating it exists that that location.")

		// ask the user if they want to be provided with a guess
		fmt.Println("Would you like to be provided a guess? (Y/N)")
		_, err = fmt.Scan(&wantGuess)
		if err != nil {
			fmt.Println("Failed reading in if you wanted a guess")
			return
		}

		// if they want a guess then get some random element
		if strings.ToUpper(wantGuess) == "Y" {
			for strings.ToUpper(enteredWordValid) != "Y" {
				randomIndex := rand.Intn(l.Len())
				count := 0

				// get the random element making sure to type cast since doubly linked list uses interface
				for element := l.Front(); element != nil; element = element.Next() {
					if count == randomIndex {
						fmt.Println("Do", element.Value)

						// Perform type assertion to convert element.Value to a string
						if str, ok := element.Value.(string); ok {
							// If the type assertion succeeds, assign the string to attemptedGuess
							attemptedGuess = str
						} else {
							// Handle the case where element.Value is not a string
							fmt.Println("Element value is not a string:", element.Value)
						}

						break
					}
					count += 1
				}

				// Since no wordlist will contain all correct guesses (different wordlist or just bad data) check with the user
				fmt.Println("Was the entered word accepted? (Y/N)")
				_, err = fmt.Scan(&enteredWordValid)
				if err != nil {
					fmt.Println("failure reading in if the entered word was valid")
					return
				}

				if strings.ToUpper(enteredWordValid) == "Y" {
					currentGuess = []rune(attemptedGuess)
					enteredWordValid = "N"
					break
				}
			}

			// if the user did their own guess ask them what it was
		} else {
			fmt.Println("What word did you guess? ")
			_, err = fmt.Scan(&attemptedGuess)
			if err != nil {
				fmt.Println("Failed reading in if you wanted a guess")
				return
			}
			currentGuess = []rune(attemptedGuess)
		}

		// ask the user for the matchbox (aka what was good letters, bad letters, and correct letters)
		fmt.Println("Please enter a match box")
		_, err = fmt.Scan(&matchBox)
		if err != nil {
			fmt.Println("failure reading in the match box")
			return
		}
		// check if user wants out
		if matchBox == "T" {
			break
		} else {
			// break down the match box into correct letters, bad letters, and good letters
			for i := 0; i < len(matchBox); i++ {
				r := rune(matchBox[i])

				switch {
				case unicode.IsLetter(r):
					word[i] = r
				case r == '.':
					if !containsRune(knownBadCharacters, currentGuess[i]) {
						knownBadCharacters = append(knownBadCharacters, currentGuess[i])
					}
				case r == ',':
					if !containsRune(knownCharacters, currentGuess[i]) {
						knownCharacters = append(knownCharacters, currentGuess[i])
					}
				default:
					fmt.Println("Bad input at position: %d", i)
				}
			}

		}

		// using regexp2 for lookahead generate the regex
		regex = regexGenerator(knownCharacters, knownBadCharacters, word)
		re, err := regexp2.Compile(regex, 0)
		if err != nil {
			fmt.Println("Invalid regex issue in generating regex: ", err)
			return
		}

		// remove any word that doesn't match the known requirements
		for element := l.Front(); element != nil; {
			nextElement := element.Next()
			if str, ok := element.Value.(string); ok {
				if isMatch, _ := re.MatchString(str); !isMatch {
					l.Remove(element)
				}
			} else {
				// Handle the case where element.Value is not a string
				fmt.Println("Element value is not a string:", element.Value)
			}
			element = nextElement
		}
		// debug statement
		//fmt.Println(l.Len())

	}

}
