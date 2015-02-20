//////////////////////////////////
// This script will automatically run the project at
// the working project directory.
//
// The binary is assumed to take in all of the files
// that were copied over as arguments, one at a time
/////////////////////////////////

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

//Takes in a param, c, which lets us know if we want to copy ideal files alongside the binary
func run() {

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
		fmt.Println("Running project for: " + studentName)
		err := os.Chdir(studentName + "/Feedback Attachment(s)")
		if err != nil {
			fmt.Println(err)
			os.Chdir("../..")
			continue
		}

		for _, input := range inputs {
			fmt.Println("Testing: " + input)
			cmd := exec.Command("./"+binName, input)
			if err = runCmd(cmd); err != nil {
				//It's POSSIBLE that they are taking in the file name (without extension.  Try that.
				filename := input
				extension := filepath.Ext(filename)
				name := filename[0 : len(filename)-len(extension)]
				cmd = exec.Command("./"+binName, name)
				if err = runCmd(cmd); err != nil {
					fmt.Println("Test case " + name + " failed.")
				}
			}
		}
		os.Chdir("../../")
	}
	os.Chdir("../")
}

//Runs the command passed
func runCmd(cmd *exec.Cmd) error {
	cmd.Start()
	done := make(chan error, 1)
	go func() {
		done <- cmd.Wait()
	}()
	select {
	case <-time.After(1 * time.Second):
		if err := cmd.Process.Kill(); err != nil {
			fmt.Println("failed to kill: " + err.Error())
		}
		<-done // allow goroutine to exit
		fmt.Println("process killed - had to timeout")
	case err := <-done:
		if err != nil {
			return err
		}
	}
	return nil
}
