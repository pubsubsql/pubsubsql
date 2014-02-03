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
)

var failCount = 0
var OS = "linux"
var architecture = 64 

func main() {
	start()	
	//
	buildServer()	
	//
	done()
}

// server

func buildServer() {
	print("Building pubsubsql server...")
	bin := "build/stage/bin/"
	cd("..")
	rm(serverFileName())
	execute("go", "build")
	cp(serverFileName(), bin + serverFileName(), true)
	cd("build")
}

func serverFileName() string {
	switch OS {
		case "windows":
			return "pubsubsql.exe"
		default:
			return "pubsubsql"
	}
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

func start() {
	print("BUILD STARTED")
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
}

func done() {
	if failCount > 0 {
		print("BUILD FAILED")
	} else {
		print("BUILD SUCCEEDED")
	}
}

func prepareStagingArea() {
	rm("stage")
	mkdir("./stage/bin")	
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

	err = dstFile.Chmod(os.ModePerm)
    if err != nil {
        return
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
