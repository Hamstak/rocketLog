package processors

import (
	"regexp"
	"log"
	"strings"
	"strconv"
	"os/exec"
)

const MAPPING_REGEX = "(\\(([0-9]*)\\))"
const SHELL_REGEX = "(`(.*)`)"

type RegexProcessor struct {
	regex         *regexp.Regexp
	mapping_regex *regexp.Regexp
	shell_regex   *regexp.Regexp
	mapping       string
}

func NewRegexProcessor(parser_regex, mapping string) *RegexProcessor{
	regex := regexp.MustCompile(parser_regex)
	mapping_regex := regexp.MustCompile(MAPPING_REGEX)
	shell_regex := regexp.MustCompile(SHELL_REGEX)

	reg := &RegexProcessor{
		mapping: mapping,
		regex: regex,
		mapping_regex: mapping_regex,
		shell_regex: shell_regex,
	}

	return reg
}

func (self *RegexProcessor) Process(input string) string {

	result := self.regex.FindStringSubmatch(input)
	mapping_result := self.mapping_regex.FindAllStringSubmatch(self.mapping, -1)
	shell_result := self.shell_regex.FindAllStringSubmatch(self.mapping, -1)

	output := self.mapping

	for _, v := range mapping_result {
		log.Print(v)
		current_mapping_token := v[1]
		result_index, err := strconv.Atoi(v[2])
		if(err != nil){
			log.Fatal(err)
		}

		current_result_token := result[result_index]

		output = strings.Replace(output, current_mapping_token, current_result_token, -1)
	}

	for _, v:= range shell_result{
		log.Print(v)
		current_mapping_token := v[0]
		cmd_slices := strings.Split(v[2], " ")

		byte_stream, err := exec.Command(cmd_slices[0], cmd_slices[1:]...).Output()
		if(err != nil){
			log.Fatal(err)
		}
		current_result_token := strings.Trim(string(byte_stream), "\n")
		output = strings.Replace(output, current_mapping_token, current_result_token, -1)
	}

	return output
}