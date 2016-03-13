package processors

import (
	"regexp"
	"log"
	"strings"
	"strconv"
	"os/exec"
)

const MAPPING_REGEX = "(`.*`|\\([0-9]*\\))"

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
	var current_mapping_token string
	var current_result_token string

	result := self.regex.FindStringSubmatch(input)
	mapping_result := self.mapping_regex.FindAllStringSubmatch(self.mapping, -1)
	log.Print(mapping_result)
	log.Print(result)
	output := self.mapping

	for _, v := range mapping_result {
		log.Print(v[0])
		current_mapping_token = v[0]
		current_result_token = self.generateResultToken(v, result)
		output = strings.Replace(output, current_mapping_token, current_result_token, -1)
	}

	return output
}

func (self *RegexProcessor) generateResultToken(input[]string, result regexp.Regexp) string{
	var current_result_token string
	if bool, err :=regexp.MatchString("`.*`", input[0]); bool {
		if (err != nil){
			log.Fatal(err)
		}
		cmd_slices := strings.Split(input[1][1:len(input)], " ")
		cmd_result, err := exec.Command(cmd_slices[0], cmd_slices[1:]...).Output()
		if(err != nil){
			log.Fatal(err)
		}
		current_result_token = string(cmd_result)

	}else{
		result_index, err := strconv.Atoi(input[1][1:(len(input))])
		if(err != nil){
			log.Fatal(err)
		}
		current_result_token = result[result_index]
	}
	return current_result_token

}