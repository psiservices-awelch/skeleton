package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os"
	"reflect"
	"strings"
	"testing"
)

func Test_copyTemplate(t *testing.T) {
	newName = "test"
	oldName = "skeletor/template"
	os.RemoveAll("./test/")
	defer os.RemoveAll("./test/")

	err := copyTemplate("./template", "./test")
	if err != nil {
		t.Fatal("failed to copy files")
		return
	}

	list1, err := getTemplateFilePaths("./template")
	if err != nil {
		t.Fatal("failed to get list of files to copy")
		return
	}

	list2, err := getTemplateFilePaths("./test")
	if err != nil {
		t.Fatal("failed to get list of files copyied")
		return
	}

	for i := 0; i < len(list1); i++ {
		l1Parts := strings.Split(list1[i], "/")
		l2Parts := strings.Split(list2[i], "/")

		l1Parts = l1Parts[1:]
		l2Parts = l2Parts[1:]

		list1[i] = strings.Join(l1Parts, "/")
		list2[i] = strings.Join(l2Parts, "/")
	}

	if !reflect.DeepEqual(list1, list2) {
		t.Fatal("list of files is different after copy")
		return
	}

}

func Test_copyTemplateAbort(t *testing.T) {
	err := os.Mkdir("./test/", 0777)
	if err != nil {
		t.Fatal(err)
		return
	}

	contents := make([][]byte, 3)
	contents[0] = []byte("no")
	contents[1] = []byte("n")
	contents[2] = []byte("random value")
	tmpfile, err := ioutil.TempFile("", "example")
	if err != nil {
		t.Fatal(err)
		return
	}

	defer os.Remove(tmpfile.Name()) // clean up

	for _, content := range contents {
		if _, err := tmpfile.Seek(0, 0); err != nil {
			t.Fatal(err)
			return
		}

		if _, err := tmpfile.Write(content); err != nil {
			t.Fatal(err)
			return
		}

		if _, err := tmpfile.Seek(0, 0); err != nil {
			t.Fatal(err)
			return
		}

		oldStdin := os.Stdin
		defer func() { os.Stdin = oldStdin }()

		os.Stdin = tmpfile

		err = copyTemplate("./template", "./test")
		if err != nil {
			t.Fatal("failed to copy files", err)
			return
		}
	}

	if err := tmpfile.Close(); err != nil {
		log.Fatal(err)
	}
}

func Test_copyFile(t *testing.T) {
	oldName = ""
	newName = ""
	if _, err := os.Stat("./template/main.go"); os.IsNotExist(err) {
		t.Fatal("./template folder needs to exist for tests")
		return
	}

	if _, err := os.Stat("./test/"); !os.IsNotExist(err) {
		err := os.RemoveAll("./test/")
		if err != nil {
			t.Fatal(err)
			return
		}
	}

	err := os.Mkdir("./test/", 0777)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer os.RemoveAll("./test/")

	err = copyFile("./template/main.go", "./test/main.go")
	if err != nil {
		t.Fatal(err)
		return
	}

	src, err := os.Open("./template/main.go")
	defer src.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	dst, err := os.Open("./test/main.go")
	defer dst.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	srcData, err := ioutil.ReadAll(src)
	if err != nil {
		t.Fatal(err)
		return
	}

	dstData, err := ioutil.ReadAll(dst)
	if err != nil {
		t.Fatal(err)
		return
	}

	if !bytes.Equal(srcData, dstData) {
		t.Fatal("files are not the same")
		return
	}
}

func Test_copyFileAfterReplace(t *testing.T) {
	if _, err := os.Stat("./template/main.go"); os.IsNotExist(err) {
		t.Fatal("./template folder needs to exist for tests")
		return
	}

	if _, err := os.Stat("./test/"); !os.IsNotExist(err) {
		err := os.RemoveAll("./test/")
		if err != nil {
			t.Fatal(err)
			return
		}
	}
	err := os.Mkdir("./test/", 0777)
	if err != nil {
		t.Fatal(err)
		return
	}

	defer os.RemoveAll("./test/")

	oldName = "skeletor/template"
	newName = "test"

	err = copyFile("./template/main.go", "./test/main.go")
	if err != nil {
		t.Fatal(err)
		return
	}

	dst, err := os.Open("./test/main.go")
	defer dst.Close()
	if err != nil {
		t.Fatal(err)
		return
	}

	dstData, err := ioutil.ReadAll(dst)
	if err != nil {
		t.Fatal(err)
		return
	}

	dstParts := bytes.Split(dstData, []byte("\""))

	dstImport := dstParts[1]

	if !bytes.Equal(dstImport, []byte("github.com/awelch/test/print")) {
		t.Fatal("import was improperly changed")
		return
	}
}
