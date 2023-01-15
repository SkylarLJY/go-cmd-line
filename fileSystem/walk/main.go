package main

import (
	"flag"
	"fmt"
	"io"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type config struct {
	ext     string
	size    int64
	list    bool
	del     bool
	wLog    io.Writer
	archive string
}

func run(root string, out io.Writer, cfg config) error {
	delLogger := log.New(cfg.wLog, "DELETED: ", log.LstdFlags)
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

		if cfg.archive != "" {
			if err := archiveFile(cfg.archive, root, path); err != nil {
				return err
			}
		}

		if cfg.del {
			return delFile(path, delLogger)
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
	log := flag.String("log", "", "log deleted files to this file")
	archive := flag.String("archive", "", "archive directory")
	flag.Parse()

	var (
		delLog = os.Stdout
		err    error
	)
	if *log != "" {
		delLog, err = os.OpenFile(*log, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0644)
		if err != nil {
			fmt.Fprintln(os.Stderr, err)
			os.Exit(1)
		}
		defer delLog.Close()

	}
	c := config{*ext, *size, *list, *delete, delLog, *archive}

	if err := run(*root, os.Stdout, c); err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
