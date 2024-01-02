package random

import (
	"encoding/json"
	"fmt"
	"math/rand"
	"os"
	"os/exec"
	"time"

	"github.com/robfig/cron/v3"
)

// LotteryData struct for JSON data
type LotteryData struct {
	RandomNumbers map[string]map[string][]int `json:"random_numbers"`
}

// StartRandomNumberGenerator initializes and starts the lottery number generator
func StartRandomNumberGenerator() {
	c := cron.New()
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	c.AddFunc("CRON_TZ=America/Los_Angeles 0 9 * * *", func() {
		generateAndCommit(location)
	})
	c.Start()
}

// generateAndCommit generates lottery numbers and commits the changes
func generateAndCommit(location *time.Location) {
	fmt.Println("Starting the lottery number generation process...")

	// Generate a random number between 5 and 18
	count := generateCount()

	// Load or Initialize numbers.json
	var data LotteryData
	file, err := os.ReadFile("./internal/random/numbers.json")
	if err != nil || len(file) == 0 || json.Unmarshal(file, &data) != nil {
		fmt.Println("Initializing new numbers.json file")
		data = LotteryData{
			RandomNumbers: make(map[string]map[string][]int),
		}
		// Immediately create or update the file
		if err := saveNumbersJSON(&data); err != nil {
			fmt.Println("Error saving initialized numbers.json:", err)
			return
		}
	}

	for i := 0; i < count; i++ {
		// Wait for a random time between 5 to 9 seconds
		waitTime := rand.Intn(5) + 5 // 5 to 9
		fmt.Printf("Waiting for %d seconds...\n", waitTime)
		time.Sleep(time.Duration(waitTime) * time.Second)

		lotteryNumbers := generateLotteryNumbers()

		now := time.Now().In(location)
		dateStr := now.Format("2006-01-02")
		timeStr := now.Format("15:04:05")
		if data.RandomNumbers == nil {
			data.RandomNumbers = make(map[string]map[string][]int)
		}
		if data.RandomNumbers[dateStr] == nil {
			data.RandomNumbers[dateStr] = make(map[string][]int)
		}
		data.RandomNumbers[dateStr][timeStr] = lotteryNumbers

		fmt.Printf("Generated numbers for %s at %s: %v\n", dateStr, timeStr, lotteryNumbers)
	}

	// Save the updated numbers.json
	if err := saveNumbersJSON(&data); err != nil {
		fmt.Println("Error saving numbers.json:", err)
		return
	}

	commitAndPush()

	fmt.Println("Lottery number generation and commit process completed.")
}

// saveNumbersJSON saves the LotteryData to numbers.json
func saveNumbersJSON(data *LotteryData) error {
	updatedFile, err := json.MarshalIndent(data, "", "  ")
	if err != nil {
		return err
	}
	return os.WriteFile("./internal/random/numbers.json", updatedFile, 0644)
}

// generateLotteryNumbers generates a set of lottery numbers
func generateLotteryNumbers() []int {
	numbers := make([]int, 6)
	for i := 0; i < 5; i++ {
		numbers[i] = rand.Intn(69) + 1 // 1-69
	}
	numbers[5] = rand.Intn(26) + 1 // 1-26 for Powerball
	return numbers
}

// commitAndPush commits and pushes changes to the repository
func commitAndPush() {
	// Define the commit message
	commitMessage := "Update lottery numbers"

	// Execute Git commands
	if err := executeGitCommand("git", "add", "./internal/random/numbers.json"); err != nil {
		fmt.Println("Error staging changes:", err)
		return
	}
	if err := executeGitCommand("git", "commit", "-m", commitMessage); err != nil {
		fmt.Println("Error committing changes:", err)
		return
	}
	if err := executeGitCommand("git", "push"); err != nil {
		fmt.Println("Error pushing changes:", err)
		return
	}

	fmt.Println("Changes committed and pushed successfully.")
}

// executeGitCommand executes a given Git command and returns any errors
func executeGitCommand(command ...string) error {
	cmd := exec.Command(command[0], command[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	return cmd.Run()
}

// testGenerateAndCommit is a manual test function for generateAndCommit
func TestGenerateAndCommit() {
	location, err := time.LoadLocation("America/Los_Angeles")
	if err != nil {
		fmt.Println("Error loading location:", err)
		return
	}

	// Manually set the random count to 1 for testing
	originalGenerateCount := generateCount
	defer func() { generateCount = originalGenerateCount }()
	generateCount = func() int { return 1 }

	// Call generateAndCommit for testing
	generateAndCommit(location)
}

// generateCount returns a random count for lottery number generation
var generateCount = func() int {
	rand.Seed(time.Now().UnixNano())
	return rand.Intn(14) + 5 // 5 to 18
}
