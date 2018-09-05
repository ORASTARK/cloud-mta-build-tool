package gen

import (
	"io/ioutil"
	"path/filepath"
	"runtime"
	"text/template"

	"gopkg.in/yaml.v2"

	"cloud-mta-build-tool/cmd/builders"
	"cloud-mta-build-tool/cmd/constants"
	fs "cloud-mta-build-tool/cmd/fsys"
	"cloud-mta-build-tool/cmd/logs"
	"cloud-mta-build-tool/cmd/mta/models"
)

// Generate - Generate mta build file
func Generate(path string) error {

	const mtaScript = "makefile"
	// Using the module context for the template creation
	mta := models.MTA{}
	type API map[string]string
	var data struct {
		File models.MTA
		API  API
	}
	// Get working directory
	projPath := fs.GetPath()
	// Create the init script file

	bashFile, err := fs.CreateFile(projPath + constants.PathSep + mtaScript)
	if err != nil {
		logs.Logger.Error("yamlFile.Get err   #%v ", err)
		return err
	}
	// Read the MTA
	yamlFile, err := ioutil.ReadFile("mta.yaml")
	if err != nil {
		logs.Logger.Error("yamlFile.Get err   #%v ", err)
		return err
	}
	// Parse mta
	err = yaml.Unmarshal([]byte(yamlFile), &mta)
	data.File = mta

	// Create maps of the template method's
	funcMap := template.FuncMap{
		"CommandProvider": builders.CommandProvider,
	}
	// Get the path of the template source code
	_, file, _, _ := runtime.Caller(0)
	container := filepath.Join(filepath.Dir(file), "script.gotmpl")
	// parse the template txt file
	t, err := template.New("script.txt").Funcs(funcMap).ParseFiles(container)
	if err != nil {
		panic(err)
	}
	// Execute the template
	if err := t.Execute(bashFile, data); err != nil {
		panic(err)
	}
	logs.Logger.Info("MTA build script was generated successfully: " + projPath + constants.PathSep + mtaScript)
	return err

}
