package processors

import (
	"log"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const mappingRegex = "(\\(([0-9]*)\\))"
const shellRegex = "(`(.*)`)"


// RegexProcessor holds necessary information for processing relavent data.
type RegexProcessor struct {
	regex         *regexp.Regexp
	mappingRegex *regexp.Regexp
	shellRegex   *regexp.Regexp
	mapping       string
    name          string
}

/*NewRegexProcessor Constructor
  @params: parserRegex, mapping
  parserRegex: the parsing information for inputs from configuration file.
  mapping: the parsing information for outputs from configuration file.
  
  returns new RegexProcessor struct containing compiled regex struct instances 
*/
func NewRegexProcessor(name, parserRegex, mapping string) *RegexProcessor {
	regex := regexp.MustCompile(parserRegex)
	mappingRegex := regexp.MustCompile(mappingRegex)
	shellRegex := regexp.MustCompile(shellRegex)

	reg := &RegexProcessor{
		mapping:       mapping,
		regex:         regex,
		mappingRegex:  mappingRegex,
		shellRegex:    shellRegex,
        name:          name,
	}

	return reg
}

/*Matches Simplification of regular expression matching function
  @params: input
  
  returns input matched against the parsing regex in the RegexProcessor instance
*/
func (regexProcessor *RegexProcessor) Matches(input string) bool {
	return regexProcessor.regex.MatchString(input)
}

// ToString returns the string representation of the processor.
func (regexProcessor *RegexProcessor) ToString() string {
    return "PROCESSOR(Type: \"Regex\", Name: \"" + regexProcessor.name + "\")"
}

/*Process Generations output for rocketlog
  @params: input
  
  returns mapping regex with replaced data captured in groups in input in addition to captured commands captured by backticks
*/
func (regexProcessor *RegexProcessor) Process(input string) string {
	result := regexProcessor.regex.FindStringSubmatch(input)
	mappingResult := regexProcessor.mappingRegex.FindAllStringSubmatch(regexProcessor.mapping, -1)
	shellResult := regexProcessor.shellRegex.FindAllStringSubmatch(regexProcessor.mapping, -1)

	output := regexProcessor.mapping

	for _, v := range mappingResult {
		currentMappingToken := v[1]
		resultIndex, err := strconv.Atoi(v[2])
		if err != nil {
			log.Fatal(err)
		}

		currentResultToken := result[resultIndex]

		output = strings.Replace(output, currentMappingToken, currentResultToken, -1)
	}

	for _, v := range shellResult {
		currentMappingToken := v[1]
		cmdSlices := strings.Split(v[2], " ")

		byteStream, err := exec.Command(cmdSlices[0], cmdSlices[1:]...).Output()
		if err != nil {
			log.Fatal(err)
		}
		currentResultToken := strings.Trim(string(byteStream), "\n")

		output = strings.Replace(output, currentMappingToken, currentResultToken, -1)
	}

	return output
}
