# fileutil
golang programs (GO programming language) dealing with file functions

This is a first attempt at programming in the GO programming language. The program mybackup.go is a simple program to get a feel of the GO language while also providing something of use. To this end, this should be seen as a means of learning the language, hence the code is not perfect.

mybackup.go is GO command. It takes 2 arguments. The first is the (fully qualified or relative) path to a file that contains on each line the full path of a folder or a file that is to be backed up. The second argument is the target folder to which the specified files and folders will be backed up.

Folders specified to be backedup are recursively backedup. There are no checks present to sanitize the input, e.g. there is no check to ensure that the target folder is itself not part of the list of files or folders to be backed up. 
