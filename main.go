package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
)

var (
	newName string
	oldName string
)

func main() {
	var replace string
	flag.StringVar(&replace, "r", "", "the old and new name to be used to update go files imports separated by a colon")

	flag.Usage = func() {
		fmt.Fprintf(os.Stderr, "Usage: %s [options] <SOURCE PATH> <DEST PATH>:\n", os.Args[0])

		fmt.Fprintln(os.Stderr, "Options:")

		flag.PrintDefaults()

		fmt.Fprintln(os.Stderr, "  -h\n\tviews the skeleton command help")
	}

	flag.Parse()

	if replace != "" {
		names := strings.Split(replace, ":")
		if len(names) == 2 {
			oldName = names[0]
			newName = names[1]
		} else if len(names) > 2 {
			log.Fatal("Too many colons present in the -r flag")
		} else {
			log.Fatal("Too few colons given in the -r flag")
		}
	}

	args := flag.Args()
	if len(args) != 2 {
		log.Fatal("ERR: incorrect number of arguments used, use -h or --help to view intended use")
	}

	err := copyTemplate(args[0], args[1])
	if err != nil {
		log.Fatal("ERR: ", err)
	}
}

//CopyTemplate copies the files and directories in the given path list to the given destination
func copyTemplate(srcPath, dstPath string) error {
	if _, err := os.Stat(dstPath); os.IsNotExist(err) {
		err := os.MkdirAll(dstPath, 0777)
		if err != nil {
			return errors.Wrapf(err, "failed to create %s", dstPath)
		}
	} else {
		reader := bufio.NewScanner(os.Stdin)
		fmt.Printf("%s already exists! Current directory contents will be removed, are you sure you want to continue? (yes/no): ", dstPath)
		var text string
		if reader.Scan() {
			text = reader.Text()
		}
		if err := reader.Err(); err != nil {
			return errors.Wrap(err, "failed to get input from user")
		}
		lw := strings.ToLower(text)
		lw = strings.TrimSpace(lw)
		if !(strings.EqualFold(lw, "yes") || strings.EqualFold(lw, "y")) {
			log.Println("project not copied")
			return nil
		}
		err = os.RemoveAll(dstPath)
		if err != nil {
			return errors.Wrap(err, "failed to remove current files and directories")
		}
		err = os.MkdirAll(dstPath, 0777)
		if err != nil {
			return errors.Wrapf(err, "failed to create %s", dstPath)
		}
	}

	if _, err := os.Stat(srcPath); os.IsNotExist(err) {
		return errors.Errorf("the %s directory does not exist", srcPath)
	}

	gitInit := exec.Command("git", "init")
	gitInit.Dir = dstPath
	_, err := gitInit.Output()
	if err != nil {
		return errors.Wrapf(err, "failed to initalize git at %s", dstPath)
	}

	list, err := getTemplateFilePaths(srcPath)
	if err != nil {
		log.Fatal("ERR: ", err)
	}

	err = copySource(srcPath, dstPath, list)
	if err != nil {
		return errors.Wrap(err, "failed to copy source files")
	}

	return nil
}

func copySource(src, dst string, list []string) error {
	for _, tempPath := range list {
		file, err := os.Stat(tempPath)
		if err != nil {
			return errors.Wrapf(err, "failed of find stats of the file %s", tempPath)
		}
		if strings.Contains(src, "./") || strings.Contains(src, ".\\") {
			src = strings.TrimPrefix(src, "./")
			src = strings.TrimPrefix(src, ".\\")
		}
		newPath := strings.TrimPrefix(tempPath, src)
		switch mode := file.Mode(); {
		case mode.IsDir():
			os.MkdirAll(dst+newPath, 0777)
		case mode.IsRegular():
			err := copyFile(tempPath, dst+newPath)
			if err != nil {
				return errors.Wrap(err, "failed to copy source file")
			}
		default:
			return errors.Errorf("file, %s, is an unrecognized file mode", newPath)
		}
	}
	return nil
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	defer in.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to open %s", src)
	}

	out, err := os.Create(dst)
	defer out.Close()
	if err != nil {
		return errors.Wrapf(err, "failed to create %s", dst)
	}

	var data io.Reader

	if strings.Contains(src, ".go") || strings.Contains(src, "go.mod") || strings.Contains(src, "go.sum") {
		if oldName != "" && newName != "" {
			inData, err := ioutil.ReadAll(in)
			if err != nil {
				return errors.Wrap(err, "failed to read src file")
			}
			newData := bytes.ReplaceAll(inData, []byte(oldName), []byte(newName))
			data = bytes.NewReader(newData)
		} else {
			log.Println("Warn: no names given to update go files, go files may need to be manually updated")
			data = in
		}
	} else {
		data = in
	}

	_, err = io.Copy(out, data)
	if err != nil {
		return errors.Wrap(err, "failed to copy source file into output file")
	}

	err = out.Sync()
	if err != nil {
		return errors.Wrap(err, "failed to flush output file contents")
	}

	return nil
}

func getTemplateFilePaths(path string) ([]string, error) {
	var files []string

	root := path
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if !strings.Contains(path, ".git") {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return nil, errors.Wrapf(err, "failed to retrieve names of files and directories in %s", path)
	}

	files = files[1:]

	return files, nil
}
