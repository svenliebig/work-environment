package cli

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
)

func QuestionFree(q string) string {
	fmt.Print(q)
	reader := bufio.NewReader(os.Stdin)
	i, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("an error occured while reading input. Please try again", err)
		return ""
	}

	i = strings.TrimSuffix(i, "\n")
	return i
}

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

func QuestionMulti(q string, options []string) []string {
	fmt.Println()
	fmt.Println(q)
	fmt.Println()

	for i, o := range options {
		fmt.Printf("  (%d): %s\n", i+1, o)
	}

	fmt.Print("\nSelect the options (komma separated): ")
	reader := bufio.NewReader(os.Stdin)
	i, err := reader.ReadString('\n')

	if err != nil {
		fmt.Println("an error occured while reading input. Please try again", err)
		return []string{}
	}

	selected := make([]string, 0, len(options))
	answers := strings.Split(strings.TrimSuffix(i, "\n"), ",")

	for _, a := range answers {
		v, err := strconv.Atoi(a)

		if err != nil {
			fmt.Println("an error occured while reading input. Please try again", err)
			return []string{}
		}

		if v <= len(options) && v > 0 {
			selected = append(selected, options[v-1])
		} else {
			fmt.Println("please only choose answers that are available for choice")
			return QuestionMulti(q, options)
		}
	}

	return selected
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
