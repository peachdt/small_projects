package server

import (
	"small_projects/utility_tool/setup/helper"
	"os"
	"fmt"
	"strings"
)

func PrepareServerFiles(target_box string) string{
	// create boom-magento-serversetup.yml

	// check if path exists
	path := "../../../../devops/ansible/"
	if _, path_err := os.Stat(path); os.IsNotExist(path_err) {
		// if path does not exists, create it
		fmt.Println("Cannot find devops/ansbile path!")
	}

	// create empty yml file
	helper.CreateFile(path, "boom-magento-serversetup.yml")

	// read from example text file and write to newly created yml
	helper.ReadFromAndWriteTo("../setup/server/boom-magento-serversetup.txt", path + "boom-magento-serversetup.yml")

	// create boom-magento-serversetup under group_vars/
	path = "../../../../devops/ansible/group_vars/"
	if _, path_err := os.Stat(path); os.IsNotExist(path_err) {
		// if path does not exists, create it
		fmt.Println("Cannot find devops/ansbile/group_vars path!")
	}

	// create empty boom-magento-serversetup file
	helper.CreateFile(path, "boom-magento-serversetup")

	// read from example boom-magento-localsetup file and write to the newly created file with vagrant params
	helper.ReadFromAndWriteTo("../setup/server/boom-magento-serversetup", path + "boom-magento-serversetup")

	// replace box name
	helper.ReplaceStringInFile(path + "boom-magento-serversetup", "replaceme", target_box, path + "boom-magento-serversetup")
	// take out dash from box name
	box_name_no_dash_list := strings.Split(target_box, "-")
	list_len := len(box_name_no_dash_list)
	box_name_no_dash := ""
	for i := 0; i < list_len; i ++ {
		box_name_no_dash = box_name_no_dash + box_name_no_dash_list[i]
	}

	helper.ReplaceStringInFile(path + "boom-magento-serversetup", "replacemenodash", box_name_no_dash, path + "boom-magento-serversetup")


	// add [boom-magento-serversetup] in hosts/myhosts and point to desired box
	// first mv original myhosts file to myhosts.tmp
	path = "../../../../devops/ansible/hosts/"

	if _, err := os.Stat(path + "myhosts.temp"); os.IsNotExist(err) {
		err := os.Rename(path + "myhosts", path + "myhosts.temp")
		if err != nil {
			fmt.Println(err)
			os.Exit(1)
		}
	} else {
		fmt.Print(path + "myhosts.temp already exists! Please review it first.")
		os.Exit(1)
	}

	// then create a new myhosts, and replace box name in it
	helper.CreateFile(path, "myhosts")

	// read from example myhosts
	helper.ReadFromAndWriteTo("../setup/server/myhosts-serversetup.txt", path + "myhosts")

	// replace box name in file
	helper.ReplaceStringInFile(path + "myhosts", "replaceme", target_box, path + "myhosts")

	// at this state, ready to run ansible using boom-magento-serversetup.yml
	return "../../../../devops/ansible/boom-magento-serversetup.yml"
}

func DeleteServerFiles() {
	// clean up after running ansible
	// files created: devops/ansible/boom-magento-serversetup.yml,
	// devops/ansible/group_vars/boom-magento-serversetup,
	// devops/ansible/hosts/myhosts,
	// and mv devops/ansible/hosts/myhosts.temp back
	fmt.Println("Start deleting newly creataed files for serversetup. . .")

	path := "../../../../devops/ansible/"
	fmt.Println(fmt.Sprintf("Removing boom-magento-serversetup.yml from %s . . .", path))
	err := os.Remove(path + "boom-magento-serversetup.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path = "../../../../devops/ansible/group_vars/"
	fmt.Println(fmt.Sprintf("Removing boom-magento-serversetup from %s . . .", path))
	err = os.Remove(path + "boom-magento-serversetup")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path = "../../../../devops/ansible/hosts/"
	fmt.Println(fmt.Sprintf("Removing myhosts from %s . . .", path))
	err = os.Remove(path + "myhosts")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	fmt.Println("Moving back myhosts.temp to myhosts . . .")
	err = os.Rename(path + "myhosts.temp", path + "myhosts")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
