package local

import (
	"small_projects/utility_tool/setup/helper"
	"os"
	"fmt"
)

func PrepareLocalFiles(target_box string) string{
	// create boom-magento-localsetup.yml

	// check if path exists
	path := "../../../../devops/ansible/"
	if _, path_err := os.Stat(path); os.IsNotExist(path_err) {
		// if path does not exists, create it
		fmt.Println("Cannot find devops/ansbile path!")
	}

	// create empty yml file
	helper.CreateFile(path,  "boom-magento-localsetup.yml")

	// read from example text file and write to newly created yml
	helper.ReadFromAndWriteTo("../setup/local/boom-magento-localsetup.txt", path + "boom-magento-localsetup.yml")

	// create boom-magento-localsetup under group_vars/

	path = "../../../../devops/ansible/group_vars/"
	if _, path_err := os.Stat(path); os.IsNotExist(path_err) {
		// if path does not exists, create it
		fmt.Println("Cannot find devops/ansbile/group_vars path!")
	}

	// create empty boom-magento-localsetup file
	helper.CreateFile(path, "boom-magento-localsetup")

	// read from example boom-magento-localsetup file and write to the newly created file with vagrant params
	helper.ReadFromAndWriteTo("../setup/local/boom-magento-localsetup", path + "boom-magento-localsetup")

	// replace box name
	helper.ReplaceStringInFile(path + "boom-magento-localsetup", "replaceme", target_box, path + "boom-magento-localsetup")

	// add [boom-magento-localsetup] in hosts/myhosts and point to desired box
	// first mv original myhosts file to myhosts.tmp
	path = "../../../../devops/ansible/hosts/"
	err := os.Rename(path + "myhosts", path + "myhosts.temp")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// then create a new myhosts, and replace box name in it
	helper.CreateFile(path, "myhosts")

	// read from example myhosts
	helper.ReadFromAndWriteTo("../setup/local/myhosts-localsetup.txt", path + "myhosts")

	// replace box name in file
	helper.ReplaceStringInFile(path + "myhosts", "replaceme", target_box, path + "myhosts")

	// at this state, ready to run ansible using boom-magento-localsetup.yml
	return "../../../../devops/ansible/boom-magento-localsetup.yml"
}

func DeleteLocalFiles() {
	// clean up after running ansible
	// files created: devops/ansible/boom-magento-localsetup.yml,
	// devops/ansible/group_vars/boom-magento-localsetup,
	// devops/ansible/hosts/myhosts,
	// and mv devops/ansible/hosts/myhosts.temp back
	fmt.Println("Start deleting newly creataed files for localsetup. . .")
	path := "../../../../devops/ansible/"
	err := os.Remove(path + "boom-magento-localsetup.yml")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path = "../../../../devops/ansible/group_vars/"
	err = os.Remove(path + "boom-magento-localsetup")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	path = "../../../../devops/ansible/hosts/"
	err = os.Remove(path + "myhosts")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	err = os.Rename(path + "myhosts.temp", path + "myhosts")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}