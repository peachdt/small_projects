package main

import (
	"flag"
	"fmt"
	"BoomPayments/labs/zhen_projects/utility_tool/ansible"

)

func main() {
	action := flag.String("action", "", "main.go action argument")
	target_box := flag.String("target_box", "", "main.go action argument")
	env := flag.String("env", "", "main.go action argument")

	flag.Parse()

	switch *action {

	case "rebuild":
//		create ansible files
		ansible.CreateAnsibleFiles(*target_box, *env)
//		write to files
//		gives out warning, press y or resume
//		delete files
		fmt.Println(fmt.Sprintf("rebuilding box %s in %s", *target_box, *env))
	default:
		fmt.Println("Arguments missing!")
		fmt.Println("go run main.go -action=rebuild -target_box=test -env=local/server")
	}
}
