//////////////////////
// This is the main function for unpacking, building, and
// running a project.
// The other files, proj_auto_build.go and proj_auto_run.go
// do all of the legwork and parameters are passed into them.
// Global variables are seen in them
// This project assumes that the main purpose is to RUN
// the project. It also supports building and even extracting
// from zips via cmdline options
// It can be chosen to NOT run the project if provided the -r switch
//
// ****IMPORTANT****
// This project assumes that the directory structre resembles the following
// -.go files
// -idealDir
// -inputsDir
// -project_dir
// --Students
// ---Feedback Attachment(s)
// ---Submission attachments(s)
// *****************
//
// It will create a folder in FeedbackAttahchment(s) called AutoGenDir
// when it extracts the contents.
// If it tries to build the project, that's where it'll look for all of the
// files it needs.
//////////////////////

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
)

//directories
//These can be either relative (from go script dir) or absolute
const project_dir = "Project_3_working"
const binName = "simulate"
const workingDirName = "AutoGenDir"
const idealDir = "proj3_ideal"
const inputsDir = "proj3_inputs"
const numTestCases = 2

var test_cases [numTestCases]string
var inputs [numTestCases]string

func main() {
	var x_bool = flag.Bool("x", false, "Extract the contents of any compressed files before compiling.")
	var cideal_bool = flag.Bool("copy_ideal", false, "Copy the ideal files to the working directory.")
	var cinput_bool = flag.Bool("copy_inputs", false, "Copy the inputs files to the working directory")
	var b_bool = flag.Bool("b", false, "Build the project.")
	var r_bool = flag.Bool("r", true, "Run the project.")
	flag.Parse()
	extract := *x_bool
	cideal := *cideal_bool
	cinput := *cinput_bool
	b := *b_bool
	r := *r_bool

	files, err := ioutil.ReadDir(idealDir)
	if err != nil {
		fmt.Println(err)
	}
	for i, file := range files {
		if isDir(file) {
			continue
		}
		test_cases[i] = file.Name()
		fmt.Println(file.Name())
	}

	files, err = ioutil.ReadDir(inputsDir)
	if err != nil {
		fmt.Println(err)
	}
	for i, file := range files {
		if isDir(file) {
			continue
		}
		inputs[i] = file.Name()
		fmt.Println(file.Name())
	}

	/*test_cases[0] = "test_case_ideal.txt"
	test_cases[1] = "test_case1_ideal.txt"
	test_cases[2] = "test_case2_ideal.txt"
	test_cases[3] = "test_case3_ideal.txt"
	inputs[0] = "test_case.asm"
	inputs[1] = "test_case1.asm"
	inputs[2] = "test_case2.asm"
	inputs[3] = "test_case3.asm"
	*/
	if b {
		fmt.Println("Building project for all students")
		if extract {
			fmt.Println("Before building, extracting files from archives for all students")
		}
		build(extract)
	}
	if cideal {
		fmt.Println("Copying ideal files (to compare tests against) from " + idealDir + " for all students")
		copy_ideal()
	}
	if cinput {
		fmt.Println("Copying inputs " + inputsDir + " for all students")
		copy_inputs()
	}
	if r {
		fmt.Println("Running the project for all students")
		run()
	}
}
