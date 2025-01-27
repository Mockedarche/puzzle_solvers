/*
spelling_bee.go is  a very very simple approach to solving the spelling bee game on the nytimes website

preq: go
*/
package main

import (
	"bufio"
	"container/list"
	"fmt"
	"io"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {

	// Create a doubly linked list (used as a generic dynamically sized array)
	l := list.New()

	// variable to hold the given letters
	var letters string

	// holds the end of the regex
	regexEnd := "]+$"

	// ask the user for the given letters with the yellow required letter being first
	fmt.Println("Enter the given letters(first being required then the rest): ")
	_, err := fmt.Scan(&letters)
	if err != nil {
		fmt.Println("Please enter only the letters shown for the spelling bee")
		return
	}

	// FILL IN WITH YOUR WORDLIST for obvious reasons I won't include mine
	file, err := os.Open("")
	if err != nil {
		log.Fatalf("Failed to open file: %v", err)

	}
	defer file.Close()
	reader := bufio.NewReader(file)

	// compile the regex that will be used to filter the word list to only correct possible guesses
	re, err := regexp.Compile(string("^" + letters[:1] + "[" + letters[1:] + regexEnd))
	if err != nil {
		fmt.Println("Invalid regex issue in generating regex: ", err)
		return
	}

	// read in the ensuring each word is of correct size, and characteristics, breaking when we finish reading the file
	for {
		line, err := reader.ReadString('\n')
		if err != nil {
			if err == io.EOF {
				break
			}
			log.Fatalf("Error reading file: %v", err)
		}
		// trim the new line
		line = strings.TrimSpace(line)
		// make sure word is valid
		if len(line) >= 4 && re.MatchString(line) {
			l.PushBack(line)
		}
	}

	fmt.Println("Printing out each possible word from the word list")
	// print out only the valid words
	for element := l.Front(); element != nil; element = element.Next() {
		fmt.Println(element.Value)
	}

}
