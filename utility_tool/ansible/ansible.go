package ansible

import (
	"fmt"
	"os"
	"io/ioutil"
	"bytes"
)

func CreateAnsibleFiles(target_box, env string) {
	if env == "local" {
//		rebuilding local vagrant box
		PrepareLocalFiles(target_box)

	} else if env == "server" {
//		rebuilding server box
		PrepareServerFiles()
	} else {
//		log error
		fmt.Println("env is invalid, should be either local or server")
	}
}

func PrepareLocalFiles(target_box string) {
	// create boom-magento-localsetup.yml

	// check if path exists
	path := "../../../../devops/ansible/"
	if _, path_err := os.Stat(path); os.IsNotExist(path_err) {
		// if path does not exists, create it
		fmt.Println("Cannot find devops/ansbile path!")
	}

	// create empty yml file
	f, create_err := os.Create(path + "boom-magento-localsetup.yml")
	fmt.Println(fmt.Sprintf("Creating %s", path + "boom-magento-localsetup.yml"))
	if create_err != nil {
		fmt.Println(create_err)
	}
	defer f.Close()

	// read from example text file and write to newly created yml
	contents, read_err := ioutil.ReadFile("../ansible/boom-magento-localsetup.txt")
	if read_err != nil {
		fmt.Println(read_err)
	}
	ioutil.WriteFile(path + "boom-magento-localsetup.yml", contents, 0644)

	// create boom-magento-localsetup under group_vars/

	path = "../../../../devops/ansible/group_vars/"
	if _, path_err := os.Stat(path); os.IsNotExist(path_err) {
		// if path does not exists, create it
		fmt.Println("Cannot find devops/ansbile/group_vars path!")
	}

	// create empty boom-magento-localsetup file
	f, create_err = os.Create(path + "boom-magento-localsetup.tmp")
	fmt.Println(fmt.Sprintf("Creating %s", path + "boom-magento-localsetup.tmp"))
	if create_err != nil {
		fmt.Println(create_err)
	}
	defer f.Close()

	// read from example boom-magento-localsetup file and write to the newly created file with vagrant params
	contents, read_err = ioutil.ReadFile("../ansible/boom-magento-localsetup")
	if read_err != nil {
		fmt.Println(read_err)
	}
	ioutil.WriteFile(path + "boom-magento-localsetup.tmp", contents, 0644)

	// replace box name
	input, err := ioutil.ReadFile(fmt.Sprintf("%s", path + "boom-magento-localsetup.tmp"))
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	output := bytes.Replace(input, []byte("replaceme"), []byte(target_box), -1)

	if err = ioutil.WriteFile(fmt.Sprintf(path + "boom-magento-localsetup"), output, 0666); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// delete tmp file
	err = os.Remove(fmt.Sprintf("%s", path + "boom-magento-localsetup.tmp"))

	if err != nil {
		fmt.Println(err)
		return
	}

	// add [boom-magento-localsetup] in hosts/myhosts and point to desired box
	

}

func PrepareServerFiles() {
	// create boom-magento-serversetup.yml
	// create boom-magento-serversetup under group_vars/
	// add [boom-magento-serversetup under hosts/myhosts and point to desired box
}


func DeleteLocalFiles() {

}

func DeleteServerFiles() {

}


func RunAnsible() {
	// run ansible with correct files
}