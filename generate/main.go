package main

import (
	"bufio"
	"bytes"
	"encoding/csv"
	"errors"
	"flag"
	"fmt"
	"go/format"
	"io"
	"net/textproto"
	"os"
	"strings"
	"time"
)

func main() {
	err := run()
	if err != nil {
		_, _ = fmt.Fprintf(os.Stderr, "failed to run: %v", err.Error())
		os.Exit(1)
	}
}

var (
	separator = []byte("%%")
	crlf      = []byte("\r\n")
)

func run() error {
	var inputFile, pairsFile, outputFile string

	flag.StringVar(&inputFile, "i", "", "Input file")
	flag.StringVar(&pairsFile, "p", "", "Language region pair file")
	flag.StringVar(&outputFile, "o", "", "Output file")

	flag.Parse()

	input, err := os.Open(inputFile)
	if err != nil {
		return fmt.Errorf("open input file: %w", err)
	}

	defer input.Close()

	pairsInput, err := os.Open(pairsFile)
	if err != nil {
		return fmt.Errorf("open pairs input file: %w", err)
	}

	defer pairsInput.Close()

	var (
		printErr error
		outBuf   bytes.Buffer
	)

	printf := func(format string, a ...any) {
		if printErr != nil {
			return
		}

		_, printErr = fmt.Fprintf(&outBuf, format, a...)
	}

	scanner := bufio.NewScanner(input)

	var buf bytes.Buffer

	header, err := readSection(scanner, &buf)
	if err != nil {
		return fmt.Errorf("read file header: %w", err)
	}

	const dateFormat = "2006-01-02"

	sourceDate, err := time.Parse(dateFormat, header.Get("File-Date"))
	if err != nil {
		return fmt.Errorf("invalid File-Date: %w", err)
	}

	printf(`//go:generate go run ./generate -o %s -i %s -p %s
`, outputFile, inputFile, pairsFile)
	printf(`//
// Based on language-subtag-registry updated %s, generated %s
package langos

`,
		sourceDate.Format(dateFormat),
		time.Now().Format(dateFormat))

	typeMaps := make(map[string][]textproto.MIMEHeader)

	for {
		entry, err := readSection(scanner, &buf)
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("read entry: %w", err)
		}

		t := entry.Get("Type")
		if t != "region" && t != "language" {
			continue
		}

		typeMaps[t] = append(typeMaps[t], entry)
	}

	for t, declarations := range typeMaps {
		printf("var %sMap = map[string]string{\n", t)

		for _, entry := range declarations {
			identifier := entry.Get("Subtag")

			printf("%q: %q,\n",
				identifier,
				entry.Get("Description"))
		}

		printf("}\n\n")
	}

	pairReader := csv.NewReader(pairsInput)
	pairReader.Comma = '\t'

	printf("var pairMap = map[string]bool{\n")

	for {
		record, err := pairReader.Read()
		if errors.Is(err, io.EOF) {
			break
		} else if err != nil {
			return fmt.Errorf("failed to read pairs from file: %w", err)
		}

		l, r, ok := strings.Cut(record[0], "-")
		if !ok {
			return fmt.Errorf("invalid language-region pair: %q", record[0])
		}

		canon := l + "-" + strings.ToUpper(r)

		printf("%q: true,\n", canon)
	}

	printf("}\n")

	if printErr != nil {
		return fmt.Errorf("failed to write to code buffer: %w", printErr)
	}

	source, err := format.Source(outBuf.Bytes())
	if err != nil {
		return fmt.Errorf("failed to format generated code: %w", err)
	}

	err = os.WriteFile(outputFile, source, 0o600)
	if err != nil {
		return fmt.Errorf("create output file: %w", err)
	}

	return nil
}

func readSection(scanner *bufio.Scanner, buf *bytes.Buffer) (textproto.MIMEHeader, error) {
	var writeErr error

	defer buf.Reset()

	write := func(d []byte) {
		if writeErr != nil {
			return
		}

		_, err := buf.Write(d)
		if err != nil {
			writeErr = err
		}
	}

	var lines int

	for scanner.Scan() {
		lines++

		if bytes.Equal(scanner.Bytes(), separator) {
			write(crlf)

			break
		}

		write(scanner.Bytes())
		write(crlf)
	}

	if writeErr != nil {
		return nil, fmt.Errorf("write to buffer: %w", writeErr)
	}

	if scanner.Err() != nil {
		return nil, fmt.Errorf("read input line: %w", scanner.Err())
	}

	if lines == 0 {
		return nil, io.EOF
	}

	reader := textproto.NewReader(bufio.NewReader(buf))

	section, err := reader.ReadMIMEHeader()
	if err != nil {
		return nil, fmt.Errorf("parse values: %w", err)
	}

	return section, nil
}
