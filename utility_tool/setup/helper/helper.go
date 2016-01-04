package helper

import (
	"fmt"
	"os"
	"io/ioutil"
	"bytes"
)

func CreateFile(path, file_name string) {
	f, create_err := os.Create(path + file_name)
	fmt.Println(fmt.Sprintf("Creating %s . . .", file_name))
	if create_err != nil {
		fmt.Println(create_err)
		os.Exit(1)
	}
	defer f.Close()

	fmt.Println(file_name + " successfully created!")
}

func ReadFromAndWriteTo(read_from, write_to string) {
	fmt.Println(fmt.Sprintf("Copying example code %s to %s", read_from, write_to))
	contents, read_err := ioutil.ReadFile(read_from)
	if read_err != nil {
		fmt.Println(read_err)
		os.Exit(1)
	}
	ioutil.WriteFile(write_to, contents, 0644)
}

func ReplaceStringInFile(file_name, replace_from, replace_to, new_file string) {
	fmt.Println(fmt.Sprintf("Replacing '%s' with '%s' in %s",replace_from, replace_to, new_file))
	input, err := ioutil.ReadFile(file_name)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte(replace_from), []byte(replace_to), -1)

	if err = ioutil.WriteFile(new_file, output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}