package fileresource

// FileResource representation of file resource
type FileResource struct {
	//Original link to file
	Link string

	//File content
	Content []byte

	//Type of the API specification file (json/yaml ...)
	Type AsdFileType
}
