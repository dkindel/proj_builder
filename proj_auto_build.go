/////////////////////////////////
// This go script is run by calling go p1_auto_build.go [-x]
// The flag x is used to denote if extraction of files (and copying to root) should occur
// If a student only submitted, for instance, his main and nothing else,
// the project won't build but the provided files can be copied in to the dir with his main
// and the script can be run again without the x flag.
// This will build the executable normally and won't throw errors, trying to extract all the contents again.
// When running that, ALL executables will be rebuilt again, it won't just build what is already there.
/////////////////////////////////

package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//For QMake:
//qmake -project "CONFIG+=c++11" "CONFIG+=app_bundle" "CONFIG+=console" -nopwd *.c* -o myAssembler.pro
//qmake myAssembler.pro
//make

/*//directories
//These can be either relative (from go script dir) or absolute
const project_dir = "Project_1_working"
const binName = "myAssembler"
const workingDirName = "AutoGenDir"
*/

//const ideal_txt_dir = "ideal_dir"
var logfile *os.File

const logfile_name = "proj_build_logfile"

//summary of failures and passes
const sumfile_name = "proj_build_summary"

var failures []string
var successes []string

func build(extract bool) {
	if err := exec.Command("7z").Run(); err != nil {
		fmt.Println("Need to have 7zip installed via comand line.  To test, run '7z' via command line")
		return
	}
	/*var x_bool = flag.Bool("x", false, "Extract the contents of any compressed files before compiling")
	flag.Parse()
	extract := *x_bool*/

	openLog()
	//create the slices failures and success to be filled
	failures = make([]string, 0, 20)
	successes = make([]string, 0, 40)

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
		fmt.Println("Building project for: " + studentName)
		writeLog(0, "Building project for: "+studentName)
		err := os.Chdir(studentName + "/Submission attachment(s)")
		if err != nil {
			fmt.Println(err)
			continue
		}
		submission_files, err := ioutil.ReadDir("./")
		if err != nil {
			addFailure(studentName, "Trouble reading directory.")
			err = os.Chdir("../..")
			continue
		}
		//extract files from compression and copy them to root folder
		if extract {
			if err := extractAnyFiles(submission_files); err != nil {
				addFailure(studentName, err.Error())
				os.Chdir("../..")
				continue
			}
		}
		//if there was a binary placed in the root dir, remove it
		//removeBins() //unnecessary now that we're creating a whole new dir

		//
		err = os.Chdir("../Feedback Attachment(s)/" + workingDirName)
		if err != nil {
			addFailure(studentName, err.Error())
			os.Chdir("../..")
			continue
		}
		if _, err := os.Stat("main.cpp"); os.IsNotExist(err) {
			addFailure(studentName, "No main file.")
			err = os.Chdir("../../..")
			continue
		}

		//build the binary
		if err := buildBinary(); err != nil {
			addFailure(studentName, "Error running g++ and qmake/qt. Try Visual Studio next")
			os.Chdir("../../..")
			continue
		}

		if err := os.Link(binName, "../"+binName); err != nil {
			fmt.Println(err)
		}
		err = os.Chdir("../../..")
		if err != nil {
			fmt.Println(err)
		}

		addSuccess(studentName)
	}
	os.Chdir("../")
	writeSummary()
}

func openLog() {
	templog, err := os.Create(logfile_name)
	logfile = templog
	if err != nil {
		fmt.Println(err)
	}
}

func writeLog(level int, msg string) {
	final_msg := ""
	for i := 0; i < level; i++ {
		final_msg += "\t"
	}
	final_msg += msg + "\n"
	logfile.WriteString(final_msg)
}

func writeSummary() {
	sumfile, err := os.Create(sumfile_name)
	if err != nil {
		fmt.Println(err)
	}
	sumfile.WriteString("***FAILURES****\n")
	for i := 0; i < len(failures); i++ {
		sumfile.WriteString(failures[i] + "\n")
	}
	sumfile.WriteString("***SUCCESSES****\n")
	for i := 0; i < len(successes); i++ {
		sumfile.WriteString(successes[i] + "\n")
	}
}

func addFailure(name string, msg string) {
	fmt.Println(msg)
	writeLog(1, "**FAIL**: "+msg)
	failures = append(failures, name)
}

func addSuccess(name string) {
	writeLog(1, "**SUCCESS**")
	successes = append(successes, name)
}

func isDir(file os.FileInfo) bool {
	switch mode := file.Mode(); {
	case mode.IsDir():
		return true
	}
	return false
}

//This takes in a list of all the files, checks if any files need to be compressed
//if any files do, it uncompressed them
//This will also move all of the .cpp and .hpp files up to the root dir
func extractAnyFiles(file_list []os.FileInfo) error {
	for _, f := range file_list {
		ext := filepath.Ext(f.Name())
		if ext == ".zip" || ext == ".7z" || ext == ".gz" || ext == ".tar" || ext == ".rar" {
			//need syscall isntead of package because people use 7z
			writeLog(1, "unzipping file: "+f.Name())
			cmd := exec.Command("7z", "x", f.Name())
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	//7z will only un gzip files if the ext is .tar.gz - check for a tar
	new_files, err := ioutil.ReadDir("./")
	if err != nil {
		fmt.Println(err)
		return err
	}
	for _, f := range new_files {
		if filepath.Ext(f.Name()) == ".tar" {
			writeLog(1, "unzipping file: "+f.Name())
			cmd := exec.Command("7z", "x", f.Name())
			if err := cmd.Run(); err != nil {
				fmt.Println(err)
				return err
			}
		}
	}
	//make directory in Feedback Attachments (where we put files to compile)
	if err := os.Mkdir("../Feedback Attachment(s)/"+workingDirName, 0777); err != nil {
		return err
	}
	if err := filepath.Walk("./", checkDir); err != nil {
		return err
	}
	return nil
}

func checkDir(path string, info os.FileInfo, err error) error {
	ext := filepath.Ext(path)
	ext = strings.ToLower(ext)
	if ext == ".cpp" || ext == ".c" || ext == ".cc" || ext == ".h" || ext == ".hpp" || ext == ".hh" {
		if err := os.Link(path, "../Feedback Attachment(s)/"+workingDirName+"/"+filepath.Base(path)); err != nil {
			return err
		}
	}
	return nil
}

func buildBinary() error {
	cmd_str := "g++ -std=c++11 *.c* -o " + binName
	cmd := exec.Command("sh", "-c", cmd_str)
	if err := cmd.Run(); err != nil {
		fmt.Println("Error running g++.  Trying qmake/qt.")
		writeLog(1, "Error running g++.  Trying qmake/qt.")
		if err := buildWithQMake(); err != nil {
			return err
		}
	}
	return nil
}

func buildWithQMake() error {
	cmd_str := "qmake -project \"CONFIG+=c++11\" \"CONFIG+=app_bundle\" \"CONFIG+=console\" -nopwd *.c* -o " + binName + ".pro"
	cmd := exec.Command("sh", "-c", cmd_str)
	if err := cmd.Run(); err != nil {
		writeLog(2, "Failed building .pro file")
		return err
	}
	cmd_str = "qmake " + binName + ".pro"
	cmd = exec.Command("sh", "-c", cmd_str)
	if err := cmd.Run(); err != nil {
		writeLog(2, "Failed building Makefile file")
		os.Remove(binName + ".pro")
		return err
	}
	cmd = exec.Command("make")
	if err := cmd.Run(); err != nil {
		writeLog(2, "Failed building binary "+binName)
		os.Remove(binName + ".pro")
		os.Remove("Makefile")
		return err
	}

	return nil
}

func removeBins() {
	os.Remove(binName)
	os.Remove(binName + ".exe")
}
