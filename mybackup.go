package main

// executables must use the package main

import (
	"fmt"
	"os"
	"log"
	"bufio"
	"strings"
	"io"
	"io/ioutil"
)

type FileHandlerForLineInFile func(line string)

/*
Input argument: Complete or relative path to a file that specifies the following:
A file where each line specifies a folder that is to be backed up (recursively) or a file to be backed up
 */
func main() {
	checkScriptArgsAndExitIfRequired()
	readFile()
}

func ReadFileLineByLine(filename, targetFolder string, handler FileHandlerForLineInFile) {
	file := getFileHandle(filename)
	handleFile(file, targetFolder, handler)
}

func getFileHandle(filename string) *os.File {
	file, err := os.Open(filename)
	if err != nil {
		log.Fatal(err)
	}
	return file
}

func handleFile(file *os.File, targetFolder string, handler FileHandlerForLineInFile) {
	doHandleFileConcurrently(file, targetFolder, handler)
}

func doHandleFileConcurrently(file *os.File, targetFolder string, handler FileHandlerForLineInFile) {
	reader := bufio.NewReader(file)
	scanner := bufio.NewScanner(reader)
	var syncStructure []chan bool

	defer file.Close()
	for scanner.Scan() {
		c := make(chan bool)
		syncStructure = append(syncStructure, c)
		filenameToBeCopied := scanner.Text()
		go func() {
			doHandle(filenameToBeCopied, targetFolder, handler, c)
		}()
	}

	count := len(syncStructure)
	for i := 0; i < count; i++ {
		fmt.Print(<- syncStructure[i])
	}
}

func doHandle(filename, targetFolder string, handler FileHandlerForLineInFile, c chan bool) {
	handler(filename)
	c <- true
}

func readFile() {
	ReadFileLineByLine(os.Args[1], os.Args[2], func(line string) { doBackup(line, os.Args[2]) })
}

func doBackup(filename, targetFolder string) {
	file := getFileHandle(filename)
	if file == nil {
		return
	}
	fInfo, err := file.Stat()
	if (err != nil) {
		log.Println(err)
		return
	}
	mode := fInfo.Mode()
	if mode.IsDir() {
		copyFolder(filename, targetFolder, file)
	} else {
		doCopyFile(filename, targetFolder, file)
	}
}

func copyFolder(fullPathOfFolderToBeCopied, targetFolder string, folderToCopy *os.File) {
	folderToBeCopiedInto := createTargetFolder(fullPathOfFolderToBeCopied, targetFolder)
	files, _ := ioutil.ReadDir(fullPathOfFolderToBeCopied)
	for _, f := range files {
		doBackup(fullPathOfFolderToBeCopied + string(os.PathSeparator) + f.Name(), folderToBeCopiedInto)
	}
}

func createTargetFolder(fullFilePath, targetFolder string) string {
	filenameWithoutPath := getFileNameWithoutPath(fullFilePath)
	fullPathOfFolderToBeCreated := targetFolder + string(os.PathSeparator) + filenameWithoutPath
	fmt.Println("Creating folder " + fullPathOfFolderToBeCreated)
	err := os.MkdirAll(fullPathOfFolderToBeCreated, 0777)
	if err != nil {
		log.Println(err)
	}
	return fullPathOfFolderToBeCreated
}

func getFileNameWithoutPath(fullPath string) string {
	slices := strings.Split(fullPath, string(os.PathSeparator))
	return slices[len(slices) - 1]
}

func doCopyFile(filename, targetFolder string, sourceFile *os.File) {
	targetFile := createTargetFile(filename, targetFolder)
	copyFile(sourceFile, targetFile)
	fmt.Println(filename + " file copied")
}

func copyFile(sourceFile, targetFile *os.File ) {
	_, err := io.Copy(targetFile, sourceFile)
	defer sourceFile.Close()
	defer targetFile.Close()
	if err != nil {
		log.Println(err)
	}
	targetFile.Sync()
}

func createTargetFile(filename, targetFolder string) *os.File {
	filenameWithoutPath := getFileNameWithoutPath(filename)
	file, err := os.Create(targetFolder + string(os.PathSeparator) + filenameWithoutPath)
	if err != nil {
		log.Println(err)
		return nil
	}
	return file
}

func checkScriptArgsAndExitIfRequired() {
	if len(os.Args) < 3 {
		fmt.Printf(getUsageMessage())
		os.Exit(0)
	}
}

func getUsageMessage() string {
	return "usage: " + os.Args[0] + " <input file> <target folder>\n input file: path to file which holds on each line the file or folders to be recursively backed up\n"
}
