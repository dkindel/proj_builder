//////////////////////////////////
// This script will automatically copy all the
// files from an ideal folder or input folder to
// the same directory as the binary
/////////////////////////////////

package main

import (
	"fmt"
	"io/ioutil"
	"os"
)

//Takes in a param, c, which lets us know if we want to copy ideal files alongside the binary
func copy_ideal() {
	if _, err := os.Stat(idealDir); os.IsNotExist(err) {
		fmt.Println(idealDir + " doesn't exist.  Aborting...")
		return
	}

	err := os.Chdir(project_dir)
	if err != nil {
		fmt.Println(err)
	}
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
	}
	//loop through everybody's directory
	for _, file := range files {
		if !isDir(file) {
			break
		}
		studentName := file.Name()
		fmt.Println("Copying ideal files for: " + studentName)
		err := os.Chdir(studentName + "/Feedback Attachment(s)")
		if err != nil {
			fmt.Println(err)
			os.Chdir("../..")
			continue
		}
		for _, f := range test_cases {
			os.Link("../../../"+idealDir+"/"+f, "./"+f)
		}
		os.Chdir("../../")
	}
	os.Chdir("../")
}

func copy_inputs() {
	if _, err := os.Stat(inputsDir); os.IsNotExist(err) {
		fmt.Println(inputsDir + " doesn't exist.  Aborting...")
		return
	}

	err := os.Chdir(project_dir)
	if err != nil {
		fmt.Println(err)
	}
	files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
	}
	//loop through everybody's directory
	for _, file := range files {
		if !isDir(file) {
			break
		}
		studentName := file.Name()
		fmt.Println("Copying inputs for: " + studentName)
		err := os.Chdir(studentName + "/Feedback Attachment(s)")
		if err != nil {
			fmt.Println(err)
			os.Chdir("../..")
			continue
		}
		for _, f := range inputs {
			os.Link("../../../"+inputsDir+"/"+f, "./"+f)
		}
		os.Chdir("../../")
	}
	os.Chdir("../")
}
