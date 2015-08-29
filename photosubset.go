//go:generate genqrc qml

package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"gopkg.in/qml.v1"
)

var dir string
var search string

func main() {
	flag.Usage = func() {
		fmt.Println("photosubset should be run from the directory where all the photos you want to make subsets of are located.")
		fmt.Fprintf(os.Stderr, "Usage of %s:\n", os.Args[0])
		flag.PrintDefaults()
	}

	flag.StringVar(&dir, "tag", ".", "which tag to browse, leave empty for all photos")
	flag.StringVar(&search, "photo", "", "substring of photo filename to start with, leave empty to start at the beginning")
	flag.Parse()

	if err := qml.Run(run); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}

func run() error {
	controller := NewPhotoController(".")

	engine := qml.NewEngine()
	engine.Context().SetVar("controller", controller)

	component, err := engine.LoadFile("qrc:///qml/photosubset.qml")
	if err != nil {
		return err
	}

	win := component.CreateWindow(nil)
	// win.Set("visibility", 5) // fullscreen
	win.Show()
	win.Wait()

	return nil
}

// NewPhotoController sets up a PhotoController by building a list of photos
// to browse and finding the right photo to start at
func NewPhotoController(path string) *PhotoController {
	// read files in the current directory
	files, err := ioutil.ReadDir(dir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// build a list of images
	images := make([]string, 0, len(files))
	for _, f := range files {
		if isImage(f.Name()) {
			absPath, _ := filepath.Abs(f.Name())
			images = append(images, "file://"+absPath)
		}
	}

	// find the correct photo, if user gave us something to search for
	index := 0
	if len(search) > 0 {
		for i, name := range images {
			if strings.Contains(filepath.Base(name), search) {
				index = i
				break
			}
		}
	}

	pc := &PhotoController{files: images, index: index}
	pc.changePhoto()
	return pc
}

// PhotoController takes care of image browsing and tagging
type PhotoController struct {
	files []string
	index int

	CurrentPhotoPath string
	TagsText         string
	currentTagNames  [10]string
	currentTags      [10]bool
}

func isImage(path string) bool {
	f, err := os.Open(path)
	if err != nil {
		return false
	}
	b := make([]byte, 512)
	_, err = f.Read(b)
	if err != nil {
		return false
	}
	contentType := http.DetectContentType(b)
	return strings.HasPrefix(contentType, "image/")
}

func (c *PhotoController) changePhoto() {
	c.CurrentPhotoPath = c.files[c.index]
	qml.Changed(c, &c.CurrentPhotoPath)
	c.getTags()
}

func (c *PhotoController) getTags() {
	c.currentTagNames = [10]string{}
	c.currentTags = [10]bool{}
	files, err := ioutil.ReadDir("./tags")
	if err != nil {
		return
	}

	imageName := filepath.Base(c.CurrentPhotoPath)
	for _, f := range files {
		if f.IsDir() && strings.HasPrefix(f.Name(), "tag-") {
			i, err := strconv.Atoi(string(f.Name()[4]))
			if err != nil {
				return
			}
			c.currentTagNames[i] = f.Name()

			taggedPath := filepath.Join("./tags", f.Name(), imageName)
			if _, err := os.Stat(taggedPath); err == nil {
				c.currentTags[i] = true
			}
		}
	}

	c.TagsText = ""
	for i, tagged := range c.currentTags {
		if tagged {
			c.TagsText = c.TagsText + fmt.Sprintln(c.currentTagNames[i][4:])
		}
	}
	qml.Changed(c, &c.TagsText)
}

// Next moves the UI to next image
func (c *PhotoController) Next() {
	if c.index < len(c.files)-1 {
		c.index++
		c.changePhoto()
	}
}

// Prev moves the UI to previous image
func (c *PhotoController) Prev() {
	if c.index > 0 {
		c.index--
		c.changePhoto()
	}
}

// Tag tags or untags the photo
func (c *PhotoController) Tag(num int) {
	imageName := filepath.Base(c.CurrentPhotoPath)
	tagName := c.currentTagNames[num]
	if len(tagName) == 0 {
		tagName = fmt.Sprintf("tag-%d", num)
		os.Mkdir(filepath.Join("./tags", tagName), 0777)
	}
	taggedPath := filepath.Join("./tags", tagName, imageName)
	if c.currentTags[num] {
		err := os.Remove(taggedPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	} else {
		err := os.Link(imageName, taggedPath)
		if err != nil {
			fmt.Println(err)
			return
		}
	}
	c.getTags()
}
