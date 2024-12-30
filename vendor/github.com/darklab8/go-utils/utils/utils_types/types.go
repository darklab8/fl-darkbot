package utils_types

import (
	"embed"
	"fmt"
	"log"
	"path/filepath"
	"strings"
)

type FilePath string

func (f FilePath) ToString() string { return string(f) }

func (f FilePath) Base() FilePath { return FilePath(filepath.Base(string(f))) }

func (f FilePath) Dir() FilePath { return FilePath(filepath.Dir(string(f))) }

func (f FilePath) Join(paths ...string) FilePath {
	paths = append([]string{string(f)}, paths...)
	return FilePath(filepath.Join(paths...))
}

type RegExp string

type TemplateExpression string

type File struct {
	Relpath   FilePath
	Name      string
	Extension string
	Content   string
}

type GetFilesParams struct {
	EmbeddedFilerName string   // required param
	RootFolder        FilePath // to exclude from RelPath
	IsNotRecursive    bool
	relFolder         FilePath
	AllowedExtensions []string
}

func GetFiles(fs embed.FS, params GetFilesParams) []File {
	if len(params.AllowedExtensions) == 0 {
		params.AllowedExtensions = []string{"js", "css", "png", "jpeg"}
	}
	if params.RootFolder == "" {
		params.RootFolder = "."
	}

	files, err := fs.ReadDir(params.RootFolder.ToString())
	var result []File
	if err != nil {
		log.Fatal(err)
	}
	for _, f := range files {
		if f.IsDir() {
			if params.IsNotRecursive {
				continue
			}
			params.relFolder = params.relFolder.Join(f.Name())
			params.RootFolder = params.RootFolder.Join(f.Name())
			result = append(result,
				GetFiles(fs, params)...,
			)
		} else {
			splitted := strings.Split(f.Name(), ".")
			var extension string
			if len(splitted) > 0 {
				extension = splitted[len(splitted)-1]
			}

			path := params.RootFolder.Join(f.Name())
			content, err := fs.ReadFile(path.ToString())
			if err != nil {
				panic(fmt.Sprintln("failed to read file from embeded fs of ", path))
			}

			is_allowed_extension := false
			for _, allowed_extension := range params.AllowedExtensions {
				if allowed_extension == extension {
					is_allowed_extension = true
				}
			}
			if !is_allowed_extension {
				continue
			}

			result = append(result, File{
				Relpath:   params.relFolder.Join(f.Name()),
				Name:      f.Name(),
				Extension: extension,
				Content:   string(content),
			})

		}

	}
	return result
}
