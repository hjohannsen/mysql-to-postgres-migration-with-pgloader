package main

import (
	"bufio"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

const replaceAll = -1
const intro = `
PGLoader Config Template Expansion

This process will expand placeholders in a PGLoader configuration template for you.
It will prompt you interactively for the MySQL and Postgres connection settings.
Plus for memory limits, timeouts and a new sequence value.

Here we go...
`

type interactiveParam struct {
	placeholder string
	prompt      string
	value       string
}

func (param *interactiveParam) setValue(v string) {
	param.value = v
}

func main() {

	fmt.Println(intro)
	template, outputPath := parseCmdLine()
	templateContent := readContent(template)
	params := readParametersInteractive()
	expandedContent := expandTemplate(params, templateContent)
	write(outputPath, &expandedContent)
	fmt.Println("\n You can run the PGLoader now:")
	fmt.Println("pgloader -v", *outputPath, "\n")

}

func parseCmdLine() (*string, *string) {
	template := flag.String("template", "./src/templates/cfg.load", "The template to use for generation.")
	outputPath := flag.String("outputFile", "./cfg.load", "Output file.")
	flag.Parse()
	fmt.Println("Will use template `" + *template + "`.")
	fmt.Println("Will write to `" + *outputPath + "`.\n")
	return template, outputPath
}

func write(outputPath *string, content *string) {
	f, err := os.Create(*outputPath)
	if err != nil {
		fmt.Println("Couldn't create output file `"+*outputPath+"`. Reason: ", err)
		panic("Couldn't create output file.")
	}
	writer := bufio.NewWriter(f)
	_, err2 := writer.WriteString(*content)
	if err2 != nil {
		fmt.Println("Sorry, something went wrong on writig to `"+*outputPath+"`. Reason: ", err2)
		panic("Couldn't write file.")
	}
	err3 := writer.Flush()
	if err3 != nil {
		fmt.Println("Sorry, something went wrong on flushing to `"+*outputPath+"`. Reason: ", err2)
		panic("Couldn't write file.")
	}
	fmt.Println("Expanded template to `" + *outputPath + "`.")
}

func expandTemplate(params []*interactiveParam, content string) string {
	for _, param := range params {
		content = strings.Replace(content, templatePlaceholder(param.placeholder), param.value, replaceAll)
	}
	return content
}

func templatePlaceholder(placeholder string) string {
	return "{{ACROLINX_MYSQL_PG__" + placeholder + "}}"
}

func readParametersInteractive() []*interactiveParam {
	interactiveParams := getInteractiveParameters()
	reader := bufio.NewReader(os.Stdin)
	length := len(interactiveParams)
	for i, param := range interactiveParams {
		fmt.Print(i+1, "/"+
			"", length, ": ", param.prompt, "\n")
		input, err := reader.ReadString('\n')
		if err != nil {
			fmt.Println("Couldn't read cmd-line input. Reason:", err)
			panic("Couldn't read cmd-line input.")
		}
		input = strings.Trim(input, "\n")
		if input != "" {
			param.setValue(input)
		}
	}
	return interactiveParams
}

func getInteractiveParameters() []*interactiveParam {

	return []*interactiveParam{
		{
			placeholder: "SRC_USER",
			prompt:      "MySQL database user:",
			value:       "",
		},
		{
			placeholder: "SRC_PW",
			prompt:      "MySQL database password:",
			value:       "",
		},
		{
			placeholder: "SRC_HOST",
			prompt:      "MySQL database host (defaults to `localhost`):",
			value:       "localhost",
		},
		{
			placeholder: "SRC_PORT",
			prompt:      "MySQL database port (defaults to " + ("3306") + "):",
			value:       "3306",
		},
		{
			placeholder: "SRC_DB",
			prompt:      "MySQL database name:",
			value:       "",
		},
		{
			placeholder: "TARGET_USER",
			prompt:      "Postgres database user:",
			value:       "",
		},
		{
			placeholder: "TARGET_PW",
			prompt:      "Postgres database password:",
			value:       "",
		},
		{
			placeholder: "TARGET_HOST",
			prompt:      "Postgres database host (defaults to `localhost`):",
			value:       "localhost",
		},
		{
			placeholder: "TARGET_PORT",
			prompt:      "Postgres database port (defaults to " + ("5432") + "):",
			value:       "5432",
		},
		{
			placeholder: "TARGET_DB",
			prompt:      "Postgres database name:",
			value:       "",
		},
		{
			placeholder: "MYSQL_TIMEOUT",
			prompt:      "MySQL net read/write timeout (in seconds, defaults to " + ("6000") + "):",
			value:       "6000",
		},
		{
			placeholder: "PG_WORK_MEM",
			prompt:      "Postgres memory (defaults to " + ("6000") + "):",
			value:       "5000",
		},
		{
			placeholder: "MAX_SEQ",
			prompt:      "Next value of target sequence (integer larger than the biggest PK in the source DB plus a generous buffer):",
			value:       "",
		},
	}
}

func readContent(file *string) string {
	data, err := ioutil.ReadFile(*file)
	if err != nil {
		fmt.Println("Couldn't read file `"+*file+"`. Reason:", err)
		panic("Can't proceed without template file.")
	}
	return string(data)
}
