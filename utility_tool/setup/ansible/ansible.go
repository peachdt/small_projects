package ansible

import (
	"small_projects/utility_tool/setup/local"
	"small_projects/utility_tool/setup/server"
	"fmt"
	"os"
)

func CreateAnsibleFiles(target_box, env string) *string{
	if env == "local" {
//		rebuilding local vagrant box
		yml_file := local.PrepareLocalFiles(target_box)
		return &yml_file

	} else if env == "server" {
//		rebuilding server box
		yml_file := server.PrepareServerFiles(target_box)
		return &yml_file
	} else {
//		log error
		fmt.Println("env is invalid, should be either local or server")
		os.Exit(1)
		return nil
	}
}

func RunAnsible(yml_file string) {
	// run ansible with correct files
	fmt.Println("Running " + yml_file)
	fmt.Println("Ansible ran successfully!")
}

