package vfs

import (
	"io"
	"io/fs"
	"net/url"
)

type FileFilter func(file VFile) (bool, error)

//VFile interface provides the basic functions required to interact
type VFile interface {
	//Closer interface included from io package
	io.Closer
	//VFileContent provider interface included
	VFileContent

	//List of this instance identified by name
	List(name string) (VFile, error)
	//ListAll children of this file instance. can be nil in case of file object instead of directory
	ListAll() ([]VFile, error)
	//Delete the file object. If the file type is directory all  files and subdirectories will be deleted
	Delete() error
	//DeleteMatching will delete only the files that match the filter.
	//Throws error if the files is not a dir type
	//If one of the file deletion fails with an error then it stops processing and returns error
	DeleteMatching(filter FileFilter) error
	//Find files based on filter only works if the file.IsDir() is true
	Find(filter FileFilter) ([]VFile, error)
	//Info  Get the file ifo
	Info() (VFileInfo, error)
	//Parent of the file system
	Parent() (VFile, error)
	//Url of the file
	Url() *url.URL
	// AddProperty will add a property to the file
	AddProperty(name string, value string) error
	// GetProperty will add a property to the file
	GetProperty(name string) (string, error)
}

//VFileContent interface providers access to the content
type VFileContent interface {
	io.ReadWriteSeeker
	//AsString content of the file. This should be used very carefully as it is not wise to load a large file in to string
	AsString() (string, error)
	//AsBytes content of the file.This should be used very carefully as it is not wise to load a large file into an array
	AsBytes() ([]byte, error)
	//WriteString method with write the string
	WriteString(s string) (int, error)
	//ContentType of the underlying content. If not set defaults to UTF-8 for text files
	ContentType() string
}

type VFileInfo interface {
	fs.FileInfo
}
