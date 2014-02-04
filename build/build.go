/* Copyright (C) 2013 CompleteDB LLC.
 *
 * This program is free software: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * This program is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with PubSubSQL.  If not, see <http://www.gnu.org/licenses/>.
 */

package main

import (
	"fmt"
	"os"
	"os/exec"
	"io"
	"archive/tar"
	"archive/zip"
	"compress/gzip"
	"path/filepath"
)

var failCount = 0
var OS = "windows"
var architecture = "Win32"

func main() {
	start()	
	//
	buildServer()	
	buildService()
	copyRootFiles()
	createArchive()
	//
	done()
}

// server

func buildServer() {
	emptyln();
	print("Building pubsubsql server...")
	bin := "build/pubsubsql/bin/"
	cd("..")
	rm(serverFileName())
	execute("go", "build")
	cp(serverFileName(), bin + serverFileName(), true)
	cd("build")
	success()
}

func serverFileName() string {
	switch OS {
		case "windows":
			return "pubsubsql.exe"
		default:
			return "pubsubsql"
	}
}

// service installer

func buildService() {
	emptyln();
	print("Building service/installer...") 
	cd("../service/" + OS)		
	switch OS {
		case "linux":
			buildServiceLinux()
		default:
			buildServiceWindows()
	}
	cd("../../build")		
	success()
}

func buildServiceWindows() {
	bin := "../../build/pubsubsql/bin/"
	execute("msbuild.exe", "/t:Clean,Build",  "/p:Configuration=Release", "/p:Platform=" + architecture)
	svc := "pubsubsqlsvc.exe"
	cp(svc, bin + svc, false)
}

func buildServiceLinux() {
	bin := "../../build/pubsubsql/bin/"
	execute("make", "clean")
	execute("make")
	svc := "pubsubsqlsvc"
	cp(svc, bin + svc, true)
}

// copy README LICENSE etc..

func copyRootFiles() {
	emptyln();
	print("Coping root files...") 
	cp("../LICENSE", "./pubsubsql/LICENSE", false)		
	success()
}

// create archive

func createArchive() {
	emptyln();	
	print("Archiving files...") 
	switch OS {
		case "linux":
			targz("temp.tar.gz", "./pubsubsql")						
			dozip("temp.zip", "./pubsubsql")
		case "windows":
			dozip("temp.zip", "./pubsubsql")
	}	
	success()
}

// helpers

func print(str string, v ...interface{}) {
	fmt.Printf(str, v...)
	fmt.Println("")
}

func fail(str string, v ...interface{}) {
	failCount++	
	print("ERROR: " + str, v...)
	os.Exit(1)
}

func emptyln() {
	fmt.Println("")
}

func success() {
	print("SUCCESS")	
}

func start() {
	print("BUILD STARTED")
	emptyln();
	// check OS 
	switch OS {
		case "windows":
			;
		case "linux":
			;
		default:
			fail("Unkown os %v", OS)
	}
	print("Preparing staging area...")
	prepareStagingArea();	
	success()
}

func done() {
	emptyln();
	if failCount > 0 {
		print("BUILD FAILED")
	} else {
		print("BUILD SUCCEEDED")
	}
}

func prepareStagingArea() {
	rm("pubsubsql")
	mkdir("./pubsubsql/bin")	
}

func mkdir(path string) {
	err := os.MkdirAll(path, os.ModeDir | os.ModePerm) 
	if err != nil {
		fail("Failed to create directory: %v error: %v", path, err)
	}
}

func cd(path string) {
	err := os.Chdir(path)  				
	if err != nil {
		fail("Failed to change directory: %v error: %v", path, err)
	}
}

func rm(path string) {
	err := os.RemoveAll(path)
	if err != nil {
		fail("Fialed to remove path: %s error: %v", path, err)
	}
}

func execute(name string, arg ...string) {
	cmd := exec.Command(name, arg ...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Stdin = os.Stdin
	err := cmd.Run()
	if err != nil {
		fail("Failed to execute command %v", err)	
	}	
}

func copyFile(src string, dst string, execute bool)  (err error) {
    srcFile, err := os.Open(src)
    if err != nil {
		return err
    }
    defer srcFile.Close()

    dstFile, err := os.Create(dst)
    if err != nil {
        return
    }
    defer dstFile.Close()

	if OS == "linux" {
		err = dstFile.Chmod(os.ModePerm)
		if err != nil {
			return
		}
	}

    _, err = io.Copy(dstFile, srcFile)
	return err
}

func cp(src string, dst string, execute bool) {
	err := copyFile(src, dst, execute)
	if err != nil {
		fail("Failed to copy file %v", err)	
	}	
}

func open(path string) *os.File {
	file, err := os.Open(path)	
	if err != nil {
		fail("Failed to open file %v error %v", path, err)
	}
	return file
}

func create(path string) *os.File {
	file, err := os.Create(path)
	if err != nil {
		fail("Failed to create file %v error %v", path, err)
	}
	return file
}

func targz(archiveFile string, dir string) {
	// file
	file := create(archiveFile)	
	defer file.Close()
	// gzip
	gzipWriter := gzip.NewWriter(file)	
	defer gzipWriter.Close()
	// tar 
	tarWriter := tar.NewWriter(gzipWriter)
	defer tarWriter.Close()
	//	
	walk := func (path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.Mode().IsDir() {
			return nil
		}
		if err != nil {
			fail("Failed to traverse directory structure %v", err)
		}
		print(path)
		fileToWrite := open(path)
		defer fileToWrite.Close()
		header, err := tar.FileInfoHeader(fileInfo, path)
		header.Name = path
		if err != nil {
			fail("Failed to create tar header from file info %v", err)
		}
		err = tarWriter.WriteHeader(header)
		if err != nil {
			fail("Failed to write tar header %v", err)
		}
		_, err = io.Copy(tarWriter, fileToWrite)		
		return nil
	}
	//
	err := filepath.Walk(dir, walk)
	if err != nil {
		fail("Failed to traverse directory %v %v", dir, err)
	}
}

func dozip(archiveFile string, dir string) {
	// file	
	file := create(archiveFile)
	defer file.Close()
	// zip
	zipWriter := zip.NewWriter(file)
	defer zipWriter.Close()
	//
	walk := func (path string, fileInfo os.FileInfo, err error) error {
		if fileInfo.Mode().IsDir() {
			return nil
		}
		if err != nil {
			fail("Failed to traverse directory structure %v", err)
		}
		print(path)
		fileToWrite := open(path)
		var w io.Writer 
		w, err = zipWriter.Create(path)	
		if err != nil {
			fail("Failed to create zip writer %v", err)
		}
		_, err = io.Copy(w, fileToWrite)		
		return nil
	}
	//
	err := filepath.Walk(dir, walk)
	if err != nil {
		fail("Failed to traverse directory %v %v", dir, err)
	}
}

