package fileutils

import (
	"archive/zip"
	"crypto/sha256"
	"fmt"
	"io"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"

	"github.com/eunanhardy/terrapak-action/internal/config"
	"github.com/gofrs/uuid"
)

/*
* REMOTE = 2
* LOCAL = 1
* UNKNOWN = 0
 */
func IdentifyPath(path string) int {
	if !IsLocal(path){
		if IsRemote(path){
			return 2
		}
	}else {
		return 1
	}
	return 0
}

func IsLocal(path string) (bool) {
	fileInfo, err := os.Stat(path)
	if err != nil {
		return false
	}

	return fileInfo.IsDir()
}

func HasChanges(path string) bool {
    target_branch, ext := os.LookupEnv("GITHUB_BASE_REF"); if !ext {
        target_branch = "main"
    }
    git_config := []string{"git","config", "--global","--add","safe.directory","/github/workspace"}
    cmd := []string{"git", "diff","--compact-summary", "HEAD",fmt.Sprintf("origin/%s",target_branch),"--",path}
    err0 := RunCommand(git_config...); if err0 != nil {
        fmt.Println(err0)
    }

    out , err := RunCommandWithOutput(cmd...); if err != nil {
        fmt.Println(err)
    }

    if(len(out) > 0){
        return true
    }else {
        fmt.Println(out)
        return false
    }
}

func HasFileChanges(config *config.ModuleConfig, hash string) (bool,error) {
    
    filepath, local_hash, err := Pack(config); if err != nil {
        return false,nil
    }
    err = os.Remove(filepath); if err != nil {
        return false,fmt.Errorf("Error Removing local file")
    }
    fmt.Printf("Remote: %s - Local: %s",hash, local_hash)

    if local_hash == hash {
        return false,nil
    } else {
        return true,nil
    }
}

func HasPreviousChanges(path string) bool {
    cmd := []string{"git", "diff","--compact-summary", "HEAD","HEAD^","--",path}

    out , err := RunCommandWithOutput(cmd...); if err != nil {
        fmt.Println(err)
    }

    if(len(out) > 0){
        return true
    }else {
        return false
    }
}

func Pack(config *config.ModuleConfig)(string,string,error){
	requestid := uuid.Must(uuid.NewV4())
	localpath := fmt.Sprintf("/tmp/%s/",requestid)
	filepath := fmt.Sprintf("%s/%s.zip",localpath,config.Name)
	err := os.MkdirAll(localpath,os.ModePerm); if err != nil {
		fmt.Println(err)
		return "","",err
	}
	err = ZipDir(config.Path,filepath); if err != nil {
		fmt.Println(err)
		return "","",err
	}
	hash, err := HashFile(filepath); if err != nil {

	}

	return filepath,hash,nil
}

func HashFile(path string) (string, error) {
    f, err := os.Open(path)
    if err != nil {
        return "",err
    }
    defer f.Close()

    h := sha256.New()
    if _, err := io.Copy(h, f); err != nil {
        return "",err
        
    }
    
    return fmt.Sprintf("%s",h.Sum(nil)), nil
}

func IsRemote(path string) (bool) {
	pattern := regexp.MustCompile(`\b(?:\w+:\/\/)?(?:\w+\.)?[a-zA-Z0-9-]+\.[a-zA-Z]{2,}(?:\.[a-zA-Z]{2,})?\b`)
	if pattern.MatchString(path) {
		return true
	}
	return false
}

func FileExists(filename string) bool {
    _, err := os.Stat(filename)
    return !os.IsNotExist(err)
}

func ZipDir(source, target string) error {

    f, err := os.Create(target); if err != nil {
        return err
    }
    defer f.Close()

    writer := zip.NewWriter(f)
    defer writer.Close()

    return filepath.Walk(source, func(path string, info os.FileInfo, err error) error {
        if err != nil {
            return err
        }

        header, err := zip.FileInfoHeader(info); if err != nil {
            return err
        }

        header.Method = zip.Deflate
        header.Name, err = filepath.Rel(source, path); if err != nil {
            return err
        }
        if info.IsDir() {
            header.Name += "/"
        }

        headerWriter, err := writer.CreateHeader(header); if err != nil {
            return err
        }

        if info.IsDir() {
            return nil
        }

        f, err := os.Open(path); if err != nil {
            return err
        }
        defer f.Close()

        _, err = io.Copy(headerWriter, f)
        return err
    })
}

func RunCommand(args ...string) (err error) {
	cmd := exec.Command(args[0], args[1:]...)

	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err1 := cmd.Run()
	if err1 != nil {
		return err1
	}

	return nil
}

func RunCommandWithOutput(args ...string) (string, error) {
	cmd := exec.Command(args[0], args[1:]...)
	out, err := cmd.CombinedOutput()
	if err != nil {
		return "", err
	}

	return string(out), nil
}