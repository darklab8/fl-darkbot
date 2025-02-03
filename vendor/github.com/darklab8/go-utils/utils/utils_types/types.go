package utils_types

import (
	"crypto/md5"
	"embed"
	"fmt"
	"io"
	"io/fs"
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

func GetFiles(filesystem embed.FS, params GetFilesParams) []File {
	if len(params.AllowedExtensions) == 0 {
		params.AllowedExtensions = []string{"js", "css", "png", "jpeg"}
	}
	if params.RootFolder == "" {
		params.RootFolder = "."
	}

	files, err := filesystem.ReadDir(params.RootFolder.ToString())
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
				GetFiles(filesystem, params)...,
			)
		} else {
			splitted := strings.Split(f.Name(), ".")
			var extension string
			if len(splitted) > 0 {
				extension = splitted[len(splitted)-1]
			}

			path := params.RootFolder.Join(f.Name())
			requested := strings.ReplaceAll(path.ToString(), "\\", "/") // fix for windows
			content, err := filesystem.ReadFile(requested)
			if err != nil {
				PrintFilesForDebug(filesystem)
				fmt.Println(err.Error(), "failed to read file from embeded fs of",
					"path=", path,
					"requested", requested,
				)
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

func PrintFilesForDebug(filesystem embed.FS) {
	fs.WalkDir(filesystem, ".", func(p string, d fs.DirEntry, err error) error {
		if !d.IsDir() {
			st, _ := fs.Stat(filesystem, p)
			r, _ := filesystem.Open(p)
			defer r.Close()

			// Read prefix
			var buf [md5.Size]byte
			n, _ := io.ReadFull(r, buf[:])

			// Hash remainder
			h := md5.New()
			_, _ = io.Copy(h, r)
			s := h.Sum(nil)

			fmt.Printf("%s %d %x %x\n", p, st.Size(), buf[:n], s)
		}
		return nil
	})
}
