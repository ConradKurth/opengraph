package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"net/url"
	"os"
	"strings"

	"github.com/otiai10/opengraph/v2"
)

// headerFlags allows multiple -H flags
type headerFlags []string

func (h *headerFlags) String() string {
	return strings.Join(*h, ", ")
}

func (h *headerFlags) Set(value string) error {
	*h = append(*h, value)
	return nil
}

func main() {
	flagset := flag.CommandLine
	flagset.Usage = func() {
		fmt.Println("Usage: ogp [OPTIONS] URL")
		fmt.Println("\nFetch URL and extract OpenGraph meta informations.")
		fmt.Println("\nOptions:")
		fmt.Println("  -A           Populate relative URLs to absolute URLs")
		fmt.Println("  -H HEADER    Add custom header (can be used multiple times)")
		fmt.Println("               Example: -H \"User-Agent: MyBot\" -H \"Accept-Language: en-US\"")
	}

	var headers headerFlags
	abs := flagset.Bool("A", false, "populate relative URLs to absolute URLs")
	flagset.Var(&headers, "H", "custom headers (format: \"Header: Value\")")

	flagset.Parse(os.Args[1:])
	if err := run(flagset.Args(), *abs, headers); err != nil {
		fmt.Println("Error:", err.Error())
		os.Exit(1)
	}
}

func run(args []string, absolute bool, headers headerFlags) error {
	if len(args) == 0 {
		return fmt.Errorf("URL must be specified")
	}
	rawurl := args[0]
	u, err := url.Parse(rawurl)
	if err != nil {
		return err
	}
	if u.Scheme == "" {
		u.Scheme = "https"
	}

	og := opengraph.New(u.String())

	// Parse custom headers if provided
	if len(headers) > 0 {
		headerMap := make(map[string]string)
		for _, h := range headers {
			parts := strings.SplitN(h, ":", 2)
			if len(parts) == 2 {
				key := strings.TrimSpace(parts[0])
				value := strings.TrimSpace(parts[1])
				headerMap[key] = value
			}
		}
		og.Intent.Headers = headerMap
	}

	if err := og.Fetch(); err != nil {
		return err
	}
	if absolute {
		if err := og.ToAbs(); err != nil {
			return err
		}
	}
	b, err := json.MarshalIndent(og, "", "\t")
	if err != nil {
		return err
	}
	fmt.Printf("%+v\n", string(b))
	return nil
}
