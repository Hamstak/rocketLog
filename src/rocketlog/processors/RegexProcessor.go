package processors

import (
	"regexp"
	"log"
	"strings"
	"strconv"
)

const MAPPING_REGEX = "(\\(([0-9]*)\\))"

type RegexProcessor struct {
	regex         *regexp.Regexp
	mapping_regex *regexp.Regexp
	mapping       string
}

func NewRegexProcessor(parser_regex, mapping string) *RegexProcessor{
	regex := regexp.MustCompile(parser_regex)
	mapping_regex := regexp.MustCompile(MAPPING_REGEX)

	reg := &RegexProcessor{
		mapping: mapping,
		regex: regex,
		mapping_regex: mapping_regex,
	}

	return reg
}

func (self *RegexProcessor) Process(input string) string {
	result := self.regex.FindStringSubmatch(input)
	mapping_result := self.mapping_regex.FindAllStringSubmatch(self.mapping, -1)

	output := self.mapping

	for _, v := range mapping_result {
		result_index, err := strconv.Atoi(v[2])
		current_mapping_token := v[1]
		current_result_token := result[result_index]
		if(err != nil){
			log.Fatal(err)
		}

		output = strings.Replace(output, current_mapping_token, current_result_token, -1)
	}

	return output
}