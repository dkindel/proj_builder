# proj_builder

##About
Project Builder was initially created to build project files download from VT's scholar site for the c++ language.  This was designed with ECE 2500 in mind but *should* work largely for other projects as well.  

This is designed to be run on linux, particularly Ubuntu.  

This builder became larger than I initially expected and wasn't architected at all.  For that reason, it's kind of bad.  It's coding sloppily and there aren't many circumventions in place in case of errors.  It'll simply fail and attempt the next student, which will probably also fail.  It also isn't nearly as modular as I'd like it to be in order to supprt different types of programs, for instance, supporting code that could take in multiple input files.  It's a simplified version because that's all that's been necessary.  In the future, I'd like to make this different and build it up some more but who knows if that day will come.  

##How to install
go here: https://golang.org/doc/install

If using Ubuntu, easiest way is by doing this:
sudo apt-get install golang

###Dependencies
7zip: Must be able to run 7zip from the command line using the call "7z".  We extract student submissions this way and it is mandatory.  
sudo apt-get install p7zip-full

g++: g++ is a standard compiler for linux and should simply work and already be installed.  If it isn't, it;s necessary
sudo apt-get install build-essential

Qt: If compilation fails with gcc, qt will be attempted.  A qt project file will automatically be generated (we ignore anything provided by the user).
sudo apt-get install qt-sdk

##How to use
###Set up file locations
From scholar, all files are downloaded into a bulk zip file named bulk_download.zip.  Pretty straightworward file name.  The user will extract that file into a place of his choosing.  We assume a structure similar to this, after extraction:

```golang
// -.go files
// -idealDir
// -inputsDir
// -project_dir
// --Students
// ---Feedback Attachment(s)
// ---Submission attachments(s)
```

The .go files may be in a different subdirectory, however, as long as they reference the project directory and the inputs and ideal directories appropriatly.  Those can be modified inside the proj_auto_main.go file.

```go
const project_dir = "../Project_1_working"
const binName = "myAssembler"
const workingDirName = "AutoGenDir"
const idealDir = "../proj1_ideal"
const inputsDir = "../proj1_inputs"
const numTestCases = 4
```

Each one of these parameters are customizable to work appropriately with your specific project.  

project_dir: The high level directory of the project currently being run against.
binName: The name of the binary that we'll output.  For instance gcc test.c -o test, "test" is the bin name.
workingDirName: The name of the autmatically generated directory under the feedback submissions.
idealDir: The location of the input files to the project
inputDir: The location of the ideal files to the project
numTestCase: How many input files can be expected.  There is only 1 input file assumed per execution.  

###Running
When the files are all set up, run "go run *.go -h" to see the flag options. Run again with whatever flag options are desired when ready for full compilation.  If all options are set, this program will attempt to extract all student submissions, copy them into an automatically generated directory under FeedbackS Submissions, and compile it there.  It will then copy in the input files and attempt to run the program with each input file and generate an output.  Copying the ideal files are simply for quick comparison and is clearly unnecessary for proper execution.  


##Pitfalls
###Installing
* Not having your GOROOT or GOPATH setup properly.  

###Running
* Not having proper permissions set up to read, write, or execute files
* Having a directroy tree that isn't what's expected
* Not setting up the variables *inside* the go files
* 7zip might not extract .rar files nicely.  I've had issues with them before so make sure this is monitored.
