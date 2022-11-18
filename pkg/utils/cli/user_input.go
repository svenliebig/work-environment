package cli

import (
	"bufio"
	"fmt"
	"os"
	"strings"
)

func Question(q string, allowedAnswers []string) string {
	fmt.Print(q)
	reader := bufio.NewReader(os.Stdin)
	i, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("an error occured while reading input. Please try again", err)
		return ""
	}

	i = strings.TrimSuffix(i, "\n")

	if len(allowedAnswers) != 0 {
		for _, aa := range allowedAnswers {
			if aa == i {
				return i
			}
		}
		fmt.Printf("please choose one of the allowed answers: [%s]\n", strings.Join(allowedAnswers, ", "))
		return Question(q, allowedAnswers)
	}

	return i
}

func QuestionScanner(q string, allowedAnswers []string) string {
	scanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(q)
		// reads user input until \n by default
		scanner.Scan()
		// Holds the string that was scanned
		text := scanner.Text()
		if len(text) != 0 {
			fmt.Println(text)
		} else {
			// exit if user entered an empty string
			break
		}

	}

	// handle error
	if scanner.Err() != nil {
		fmt.Println("Error: ", scanner.Err())
	}
	return ""
}
