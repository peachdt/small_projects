package main

import (
	"flag"
	"fmt"
	"small_projects/utility_tool/setup/ansible"
	"small_projects/utility_tool/setup/local"
	"small_projects/utility_tool/setup/server"

)

func main() {
	action := flag.String("action", "", "main.go action argument")
	target_box := flag.String("target_box", "", "main.go action argument")
	env := flag.String("env", "", "main.go action argument")

	flag.Parse()

	switch *action {

	case "rebuild":
//		create ansible files
		yml_file := ansible.CreateAnsibleFiles(*target_box, *env)
		if yml_file != nil {
			ansible.RunAnsible(*yml_file)
		}
//		write to files
//		gives out warning, press y or resume
//		delete files
		fmt.Println(fmt.Sprintf("Ansible rebuild done, please run `go run main.go -action=clear -env=%s` to delete created files.", *env))
	case "clear":
		if *env == "local" {
			local.DeleteLocalFiles()
		} else if *env == "server" {
			server.DeleteServerFiles()
		} else {
			fmt.Println("env missing! Please specify env argument (local/server)")
		}
	default:
		fmt.Println("Arguments missing!")
		fmt.Println("go run main.go -action=rebuild -target_box=test -env=local/server")
		fmt.Println("go run main.go -action=clear -env=local/server")
	}
}
