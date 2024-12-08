/*
утилита -like

Должна принимать некоторый путь и набор параметров командной строки, чтобы иметь возможность находить различные типы записей
Нас интересуют три типа записей: каталоги, обычные файлы и символические ссылки.

./myFind /foo
выводит все содержимое директории

Пользователь должен иметь возможность указать один, два или все три из них явно,
например ./myFind -f -sl /path/to/dirили ./myFind -d /path/to/other/dir

Вам также следует реализовать еще одну опцию - -ext(работает ТОЛЬКО при указании -f),
чтобы пользователь мог печатать только файлы с определенным расширением
например ./myFind -f -ext 'go' /go

Вам также понадобится разрешить символические ссылки. Так что, если /foo/bar/buzzэто символическая ссылка,
указывающая на какое-то другое место в FS, например /foo/bar/baz, выведите оба пути, разделенные ->
Символические ссылки могут быть повреждены (указывая на несуществующий файловый узел).
В этом случае ваш код должен вывести [broken] вместо пути назначения символической ссылки.
/foo/bar/buzz -> /foo/bar/baz
/foo/bar/broken_sl -> [broken]
*/
package main

import (
	"errors"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"os"
	"path/filepath"
)

type directory struct {
	folder  string
	ext     string
	links   bool
	dirs    bool
	files   bool
	noFlags bool
}

// -ext "j" -f -sl test_folder

func getDirInfo() (*directory, error) {
	extention := flag.String("ext", "", "file extention")
	links := flag.Bool("sl", false, "only links output")
	dirs := flag.Bool("d", false, "only directories output")
	files := flag.Bool("f", false, "only files output")

	flag.Parse()
	folder := flag.Args()

	if len(folder) > 1 || len(folder) == 0 {
		err := errors.New("invalid arguments: add folder")
		return nil, err
	}
	if !*files && *extention != "" {
		err := errors.New("\"ext\" flag can be used with \"f\" flag only")
		return nil, err
	}
	noFlags := false
	if !*links && !*dirs && !*files {
		noFlags = true
	}
	return &directory{folder: folder[0],
		ext:     *extention,
		links:   *links,
		dirs:    *dirs,
		files:   *files,
		noFlags: noFlags,
	}, nil
}

func printContent(d directory) error {
	data, err := os.ReadDir(d.folder)
	if err != nil {
		return err
	}

	for _, val := range data {
		path := d.folder + "/" + val.Name()
		info, err := os.Lstat(path)
		if err != nil {
			return fmt.Errorf("can not  get dir info: %v", err)
		}

		switch mode := info.Mode(); {
		case mode.IsRegular() && (d.files || d.noFlags):
			fileExtention := filepath.Ext(val.Name())
			if d.ext == "" || fileExtention[1:] == d.ext {
				fmt.Printf("/%s\n", path)
			}
		case mode.IsDir() && (d.dirs || d.noFlags):
			newDir := d
			newDir.folder = path
			err = printContent(newDir)
			if err != nil {
				return err
			}
		case mode&fs.ModeSymlink != 0 && (d.links || d.noFlags):
			//errPermission возникает если нет прав к открытию файла - в таком случае его просто пропускаем
			//errInvalid - если передается не символическая ссылка
			link, err := os.Readlink(path)
			if err != nil {
				if errors.Is(err, os.ErrPermission) || errors.Is(err, os.ErrInvalid) || os.IsNotExist(err) {
					continue
				}
				return err
			}
			_, err = os.ReadFile(link)
			if err != nil {
				fmt.Printf("/%s -> [broken]\n", path)
			} else {
				fmt.Printf("/%s -> %s\n", path, link)
			}
		default:
			return fmt.Errorf("unsupported file type")
		}
	}
	return nil
}

func main() {
	d, err := getDirInfo()
	if err != nil {
		log.Fatalln(err)
	}
	err = printContent(*d)
	if err != nil {
		log.Fatalln(err)
	}
}
