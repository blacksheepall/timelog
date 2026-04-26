package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"os"
	"strconv"

	"github.com/blacksheepaul/timelog/pkg/temppassword"
)

func main() {
	jsonFlag := flag.Bool("json", false, "output as JSON")
	ttlFlag := flag.Int("ttl", 900, "time to live in seconds (default: 900)")
	countFlag := flag.Int("count", 1, "number of passwords to generate (default: 1)")
	helpFlag := flag.Bool("help", false, "show help")
	helpShortFlag := flag.Bool("h", false, "show help (short)")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, `Standalone temporary password generator.

Usage:
  temp-password [options]
  temp-password [ttl] [options]

Options:
`)
		flag.PrintDefaults()
		fmt.Fprintf(os.Stderr, `
Examples:
  # Generate password with default TTL (900s)
  temp-password

  # Generate password with 1800 second TTL
  temp-password 1800

  # Generate 5 passwords as JSON
  temp-password -count 5 -json

  # Generate password with custom TTL and JSON output
  temp-password 600 -json
`)
	}

	flag.Parse()

	if *helpFlag || *helpShortFlag {
		flag.Usage()
		os.Exit(0)
	}

	// Parse positional argument for TTL if provided
	args := flag.Args()
	if len(args) > 0 {
		if ttl, err := strconv.Atoi(args[0]); err == nil {
			*ttlFlag = ttl
		} else {
			fmt.Fprintf(os.Stderr, "error: invalid ttl value: %v\n", err)
			os.Exit(1)
		}
	}
	if len(args) > 1 {
		fmt.Fprintln(os.Stderr, "error: too many positional arguments")
		os.Exit(1)
	}

	if *countFlag < 1 {
		fmt.Fprintf(os.Stderr, "error: count must be >= 1\n")
		os.Exit(1)
	}

	if *ttlFlag < 0 {
		fmt.Fprintf(os.Stderr, "error: ttl must be >= 0\n")
		os.Exit(1)
	}

	var results []*temppassword.PasswordWithExpiry
	for i := 0; i < *countFlag; i++ {
		pwd, err := temppassword.GeneratePasswordWithExpiry(*ttlFlag)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error: failed to generate password: %v\n", err)
			os.Exit(1)
		}
		results = append(results, pwd)
	}

	if *jsonFlag {
		outputJSON(results)
	} else {
		outputText(results)
	}
}

func outputJSON(results []*temppassword.PasswordWithExpiry) {
	var output interface{}
	if len(results) == 1 {
		output = results[0]
	} else {
		output = results
	}

	data, err := json.MarshalIndent(output, "", "  ")
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: failed to marshal JSON: %v\n", err)
		os.Exit(1)
	}
	fmt.Println(string(data))
}

func outputText(results []*temppassword.PasswordWithExpiry) {
	for i, pwd := range results {
		if len(results) > 1 {
			fmt.Printf("Password %d:\n", i+1)
		}
		fmt.Printf("  Password:  %s\n", pwd.Password)
		fmt.Printf("  Hash:      %s\n", pwd.Hash)
		fmt.Printf("  Expires:   %s\n", pwd.ExpiresAt.Format("2006-01-02 15:04:05"))
		if i < len(results)-1 {
			fmt.Println()
		}
	}
}
