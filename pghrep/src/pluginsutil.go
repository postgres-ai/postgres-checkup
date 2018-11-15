package main

import (
	"fmt"
    "io/ioutil"
    "log"
    "math/rand"
    "os"
    "os/exec"
    "path/filepath"
    "plugin"
)


type Preparer interface {
	Prepare(map[string]interface{}) map[string]interface{}
}

// loader keeps the context needed to find where plugins and
// objects are stored.
type loader struct {
    pluginsDir string
    objectsDir string
    binDir string
}

func newLoader() (*loader, error) {
    // The directory that will be watched for new plugins.
    wd, err := os.Getwd()
    if err != nil {
        return nil, fmt.Errorf("could not find current directory: %v", err)
    }
    pluginsDir := filepath.Join(wd, "plugins")
    binDir := filepath.Join(wd, "bin")
    //tmp := filepath.Join(wd, "bin")

    // The directory where all .so files will be stored.
    tmp, err := ioutil.TempDir("", "")
    if err != nil {
        return nil, fmt.Errorf("could not create objects dir: %v", err)
    }
    return &loader{pluginsDir: pluginsDir, objectsDir: tmp, binDir: binDir}, nil
}

func (l *loader) destroy() { 
    //os.RemoveAll(l.objectsDir)
}

func (l *loader) compileAndRun(name string, data map[string]interface{}) (map[string]interface{}, error) {
    obj, err := l.compile(name)
    if err != nil {
        return nil, fmt.Errorf("could not compile %s: %v", name, err)
    }
    defer os.Remove(obj)

    result, err := l.call(obj, data)
    if err != nil {
        return nil, fmt.Errorf("could not run %s: %v", obj, err)
    }
    return result, nil
}

// check existance of binary plugin library or compile it from sources
// and returns its path.
func (l *loader) get(name string) (string, error) {
    pluginPath := filepath.Join(l.binDir, name + ".so")
    _, err := os.Stat(filepath.Join(l.binDir, name + ".so"))
    if err != nil && os.IsNotExist(err) {
        log.Println("WARNING: Binary plugin not found. Try compile")
        pluginPath, err = l.compile(name)
    } else {
        log.Println("Binary plugin found", )
        err = nil
    }
    return pluginPath, err
}

// compile compiles the code in the given path, builds a
// plugin, and returns its path.
func (l *loader) compile(name string) (string, error) {
    // Copy the file to the objects directory with a different name
    // each time, to avoid retrieving the cached version.
    // Apparently the cache key is the path of the file compiled and
    // there's no way to invalidate it.

    f, err := ioutil.ReadFile(filepath.Join(l.pluginsDir, name + ".go"))
    if err != nil {
        return "", fmt.Errorf("could not read %s.go: %v", name, err)
    }

    name = fmt.Sprintf("%d.go", rand.Int())
    srcPath := filepath.Join(l.objectsDir, name)
    if err := ioutil.WriteFile(srcPath, f, 0666); err != nil {
        return "", fmt.Errorf("could not write %s: %v", name, err)
    }

    objectPath := srcPath[:len(srcPath)-3] + ".so"

    cmd := exec.Command("go", "build", "-buildmode=plugin", "-o="+objectPath, srcPath)
    cmd.Stderr = os.Stderr
    cmd.Stdout = os.Stdout
    if err := cmd.Run(); err != nil {
        return "", fmt.Errorf("could not compile %s: %v", name, err)
    }

    return objectPath, nil
}

// call loads the plugin object in the given path and runs the Run
// function.
func (l *loader) call(object string, checkData map[string]interface{}) (map[string]interface{}, error) {
    p, err := plugin.Open(object)
    if err != nil {
        return nil, fmt.Errorf("Check plugin not found.")
    }
    symPreparer, err := p.Lookup("Preparer")
	if err != nil {
		fmt.Println(err)
		return nil, fmt.Errorf("plugin failed with error %v", err)
	}
    
    var preparer Preparer
	preparer, ok := symPreparer.(Preparer)
	if !ok {
        return nil, fmt.Errorf("unexpected type from module symbol")
	}

	result := preparer.Prepare(checkData)
    return result, nil
}

// goFiles lists all the files in the plugins
func (l *loader) plugins() []string {
    dir, err := os.Open(l.pluginsDir)
    if err != nil {
        log.Fatal(err)
    }
    defer dir.Close()
    names, err := dir.Readdirnames(-1)
    if err != nil {
        log.Fatal(err)
    }

    var res []string
    for _, name := range names {
        if filepath.Ext(name) == ".go" {
            res = append(res, name)
        }
    }
    return res
}