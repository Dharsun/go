package main

import (
	"archive/zip"
	"bytes"
	"fmt"
	"io"
	"os"
	"regexp"
	"strings"
)

const outputFilePath = "output.html"

func searchZipFile(zipFilePath string, keywordRegexMap map[string]string) (map[string]int, error) {
	reader, err := zip.OpenReader(zipFilePath)
	if err != nil {
		return nil, err
	}
	defer reader.Close()

	// Create or open the output HTML file for writing
	outputFile, err := os.Create(outputFilePath)
	if err != nil {
		return nil, err
	}
	defer outputFile.Close()

	// Write the HTML file header with black background and white text color
	_, err = outputFile.WriteString("<html><head><title>Output</title></head><body style=\"background-color:black;color:white;\">\n")
	if err != nil {
		return nil, err
	}

	// Add header
	_, err = outputFile.WriteString("<h1>This is Doctor Strange's magic</h1>\n")
	if err != nil {
		return nil, err
	}

	// Map to store matched keyword counts
	keywordCounts := make(map[string]int)

	// Add buttons for each keyword
	for keyword, buttonLabel := range map[string]string{
		"keyword1": "JDBC",
		"keyword2": "SRV",
	} {
		_, err = outputFile.WriteString(fmt.Sprintf("<button onclick=\"showKeyword('%s')\">%s</button> = <span id=\"%s_count\">0</span>&nbsp;&nbsp;", keyword, buttonLabel, keyword))
		if err != nil {
			return nil, err
		}
		keywordCounts[keyword] = 0
	}
	_, err = outputFile.WriteString("<br><br>") // Add two new lines
	// A map to store matched keywords
	matchedKeywords := make(map[string]bool)

	for _, file := range reader.File {
		fileReader, err := file.Open()
		if err != nil {
			return nil, err
		}
		defer fileReader.Close()

		var buffer bytes.Buffer
		_, err = io.Copy(&buffer, fileReader)
		if err != nil {
			return nil, err
		}
		fileContent := buffer.Bytes()

		lines := strings.Split(string(fileContent), "\n")
		for lineNumber, line := range lines {
			for keyword, regexPattern := range keywordRegexMap {
				regex := regexp.MustCompile(regexPattern)
				matches := regex.FindAllStringIndex(line, -1)
				for _, match := range matches {
					start, end := match[0], match[1]
					keywordMatch := line[start:end]

					// Check if keyword has already been matched
					if !matchedKeywords[keywordMatch] {
						matchedKeywords[keywordMatch] = true

						// Create a span element with red text color for the keyword
						coloredKeyword := fmt.Sprintf("<span style=\"color:red;\">%s</span>", keywordMatch)
						foundOutput := fmt.Sprintf("<br>FOUND = %s<br><br>", keywordMatch)
						modifiedLine := strings.Replace(line, keywordMatch, coloredKeyword, 1)
						outputLine := fmt.Sprintf("<div class=\"result %s\" style=\"display:none;\">Line %d - %s<br> content: %s %s</div>\n", keyword, lineNumber+1, file.Name, modifiedLine, foundOutput)

						_, err := outputFile.WriteString(outputLine)
						if err != nil {
							return nil, err
						}

						// Update the matched keyword count and display
						keywordCounts[keyword]++
						_, err = outputFile.WriteString(fmt.Sprintf("<script>document.getElementById('%s_count').textContent = %d;</script>\n", keyword, keywordCounts[keyword]))
						if err != nil {
							return nil, err
						}
					}
				}
			}
		}
	}

	// JavaScript function to toggle visibility of results
	_, err = outputFile.WriteString(`
		<script>
			function showKeyword(keyword) {
				var elements = document.getElementsByClassName("result");
				for (var i = 0; i < elements.length; i++) {
					elements[i].style.display = "none";
				}
				
				var keywordElements = document.getElementsByClassName(keyword);
				for (var i = 0; i < keywordElements.length; i++) {
					keywordElements[i].style.display = "block";
				}
			}
		</script>
	`)
	if err != nil {
		return nil, err
	}

	// Write the HTML file footer
	_, err = outputFile.WriteString("</body></html>\n")
	if err != nil {
		return nil, err
	}

	return keywordCounts, nil
}

func main() {
	if len(os.Args) != 2 {
		fmt.Println("Usage: go run main.go <zipFilePath>")
		return
	}

	zipFilePath := os.Args[1]

	// Define a map of keywords and their corresponding regular expressions
	keywordRegexMap := map[string]string{
		"keyword1": `keyword`,
		"keyword2": `srv.*AP[A-Za-z0-9]*`,
		// Add more keywords and patterns as needed
	}

	counts, err := searchZipFile(zipFilePath, keywordRegexMap)
	if err != nil {
		fmt.Println("Error:", err)
	} else {
		fmt.Println("Keyword Counts:")
		for keyword, count := range counts {
			fmt.Printf("%s: %d\n", keyword, count)
		}
	}
}