package main

import (
	_ "embed"
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/atomicgo/cursor"
	"github.com/eiannone/keyboard"
	"github.com/leaanthony/clir"
	"github.com/pterm/pterm"
)

//go:embed data/words.json
var rawWords []byte

const (
	NotGuessed int = 0
	RightPlace int = 1
	WrongPlace int = 2
	NotFound   int = 3
)

// Keeps track of whether or not letters have been found
var guessedLetters map[string]int

func customBanner(cli *clir.Cli) string {

	return `
 ██     ██  ██████  ██████  ██████  ██      ███████ ██████
 ██     ██ ██    ██ ██   ██ ██   ██ ██      ██      ██   ██
 ██  █  ██ ██    ██ ██████  ██   ██ ██      █████   ██████
 ██ ███ ██ ██    ██ ██   ██ ██   ██ ██      ██      ██   ██
  ███ ███   ██████  ██   ██ ██████  ███████ ███████ ██   ██
                     ` + cli.Version() + " - " + cli.ShortDescription()
}

func Date(year, month, day int) time.Time {
	return time.Date(year, time.Month(month), day, 0, 0, 0, 0, time.UTC)
}

type Words struct {
	Solutions []string `json:"solutions"`
	Herrings  []string `json:"herrings"`
}

func Contains(a []string, x string) bool {
	for _, n := range a {
		if x == n {
			return true
		}
	}
	return false
}

func displayInputWord(word string) {
	var blanks = strings.Repeat("_", 5-len(word))
	t, _ := pterm.DefaultBigText.WithLetters(
		pterm.NewLettersFromStringWithStyle(word, pterm.NewStyle(pterm.FgWhite, pterm.Bold)), pterm.NewLettersFromStringWithStyle(blanks, pterm.NewStyle(pterm.FgWhite, pterm.Bold))).
		Srender()
	pterm.DefaultCenter.Print(t)
	pterm.DefaultCenter.Print()
	displayAlphabet()
	pterm.DefaultCenter.Print()
	pterm.DefaultCenter.Print()
}

func displayAlphabet() {

	var sb strings.Builder

	// lowercase a - z ASCII codes
	for i := 97; i < 123; i++ {

		letter := string(i - 32)

		var style *pterm.Style

		if guessedLetters[string(i)] == RightPlace {
			style = pterm.NewStyle(pterm.FgGreen)
		} else if guessedLetters[string(i)] == WrongPlace {
			style = pterm.NewStyle(pterm.FgYellow)
		} else if guessedLetters[string(i)] == NotFound {
			style = pterm.NewStyle(pterm.FgDarkGray)
		} else {
			style = pterm.NewStyle(pterm.FgWhite, pterm.Bold)
		}

		sb.WriteString(style.Sprint(letter))
		if i < 122 {
			sb.WriteString(" ")
		}

	}

	pterm.DefaultCenter.Print(sb.String())

}

func displayGuess(guess string, answer string) {
	var used = [5]bool{false, false, false, false, false}

	var letters []pterm.Letters
	for i := 0; i < len(guess); i++ {
		if guess[i] == answer[i] {
			letters = append(letters, pterm.NewLettersFromStringWithStyle(string(guess[i]), pterm.NewStyle(pterm.FgGreen)))
			guessedLetters[string(guess[i])] = RightPlace
			used[i] = true
		} else {
			var found = false
			for j := 0; j < len(guess); j++ {
				if i == j {
					continue
				}
				if !used[j] && guess[i] == answer[j] {
					used[j] = true
					found = true
					break
				}
			}
			if found {
				letters = append(letters, pterm.NewLettersFromStringWithStyle(string(guess[i]), pterm.NewStyle(pterm.FgYellow)))
				if guessedLetters[string(guess[i])] != RightPlace {
					guessedLetters[string(guess[i])] = WrongPlace
				}
			} else {
				letters = append(letters, pterm.NewLettersFromStringWithStyle(string(guess[i]), pterm.NewStyle(pterm.FgDarkGray)))
				if guessedLetters[string(guess[i])] != RightPlace && guessedLetters[string(guess[i])] != WrongPlace {
					guessedLetters[string(guess[i])] = NotFound
				}
			}

		}
	}
	t, _ := pterm.DefaultBigText.WithLetters(letters...).
		Srender()
	pterm.DefaultCenter.Print(t)
}

func main() {

	numberOfGuesses := 0
	maxNumberOfGuesses := 6

	guessedLetters = make(map[string]int)

	for i := 97; i < 123; i++ {

		letter := string(i)
		guessedLetters[letter] = NotGuessed

	}

	// Number of lines that displaying the input word will take up
	const lines = 10

	// Create new cli
	cli := clir.NewCli("wordler", "Play WORDLE on the command line", "v1.0.2")

	cli.SetBannerFunction(customBanner)

	// Name
	dateString := ""
	cli.StringFlag("date", "Date to play (e.g. 2021-11-24)", &dateString)

	// Define action for the command
	cli.Action(func() error {

		if err := keyboard.Open(); err != nil {
			panic(err)
		}
		defer func() {
			_ = keyboard.Close()
		}()

		// Read the solutions

		var words Words

		json.Unmarshal(rawWords, &words)

		t1 := Date(2021, 6, 19)
		today := time.Now()
		var t2 time.Time
		if dateString != "" {

			date, err := time.Parse("2006-01-02", dateString)
			if err != nil {
				panic(err)
			}

			t2 = Date(date.Year(), int(date.Month()), date.Day())
		} else {
			t2 = Date(today.Year(), int(today.Month()), today.Day())
		}

		pterm.DefaultHeader.WithFullWidth().Println(t2.Format("Monday, January 2, 2006"))
		pterm.Println()

		days := t2.Sub(t1).Hours() / 24
		var answer = words.Solutions[int(days)%len(words.Solutions)]

		cursor.Hide()

		var exit = false
		var sb strings.Builder

		sb.Reset()
		displayInputWord("")

		for {

			var submit = false

			if exit {
				break
			}

			for {
				char, key, err := keyboard.GetKey()
				if err != nil {
					panic(err)
				}

				if key == keyboard.KeyBackspace || key == keyboard.KeyBackspace2 {
					cursor.ClearLinesUp(lines)
					var backup = sb.String()
					if len(backup) > 0 {
						sb.Reset()
						sb.WriteString(backup[0 : len(backup)-1])
					}
					displayInputWord(sb.String())

				} else if key == keyboard.KeyEnter {
					if len(sb.String()) == 5 {
						submit = true
						break
					}
				} else if key == keyboard.KeyEsc || key == keyboard.KeyCtrlC {
					exit = true
					break
				} else {

					if key == 0 && char >= 'a' && char <= 'z' {
						cursor.ClearLinesUp(lines)
						if len(sb.String()) < 5 {
							sb.WriteRune(char)
						}
						displayInputWord(sb.String())
					}

				}
			}

			if exit && len(sb.String()) == 0 {
				cursor.ClearLinesUp(lines)
			}

			if submit {

				var guess = sb.String()

				if !Contains(words.Solutions, guess) && !Contains(words.Herrings, guess) {
					cursor.Up(1)
					pterm.DefaultCenter.Print("Not in word list")
				} else {

					sb.Reset()

					cursor.ClearLinesUp(lines)
					displayGuess(guess, answer)
					numberOfGuesses++

					if guess == answer {
						exit = true
						switch numberOfGuesses {
						case 1:
							pterm.DefaultCenter.Print("Genius")
						case 2:
							pterm.DefaultCenter.Print("Magnificent")
						case 3:
							pterm.DefaultCenter.Print("Impressive")
						case 4:
							pterm.DefaultCenter.Print("Splendid")
						case 5:
							pterm.DefaultCenter.Print("Great")
						case 6:
							pterm.DefaultCenter.Print("Nice")
						}
					} else if numberOfGuesses >= maxNumberOfGuesses {
						pterm.DefaultCenter.Print("The word was " + strings.ToUpper(answer))
						exit = true
					} else {
						displayInputWord("")
					}

				}

			}

		}

		cursor.Show()
		return nil

	})

	if err := cli.Run(); err != nil {
		fmt.Printf("Error encountered: %v\n", err)
	}

}
