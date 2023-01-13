package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"os"
	"path/filepath"
)

type config struct {
	ext  string
	size int64
	list bool
	del  bool
}

func run(root string, out io.Writer, cfg config) error {
	return filepath.Walk(root, func(path string, info fs.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if filterOut(path, cfg.ext, cfg.size, info) {
			return nil
		}

		if cfg.list {
			return listFile(path, out)
		}
		if cfg.del {
			return delFile(path)
		}

		return listFile(path, out)
	})
}

func main() {
	root := flag.String("root", ".", "root directory to start the search")
	ext := flag.String("ext", "", "file extension to filter for")
	size := flag.Int64("size", 0, "the minimum size of the files to list")
	list := flag.Bool("list", false, "list all files")
	delete := flag.Bool("del", false, "delete files")
	flag.Parse()

	c := config{*ext, *size, *list, *delete}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
