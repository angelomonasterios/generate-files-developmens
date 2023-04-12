package main

import (
	"fmt"
	"github.com/spf13/cobra"
	"html/template"
	"os"
	"path/filepath"
)

const (
	appDirectory          = "app"
	modelsDirectory       = "./app/Models"
	validatorsDirectory   = "./app/Validators"
	storeDirectory        = "./app/Store"
	presentersDirectory   = "./app/Presenters"
	presentersControllers = "./app/Http/Controllers"

	modelTemplate = `<?php
namespace App\Models\{{.Namespace}};

use App\Models\Base\BaseModel;

class {{.Name}} extends BaseModel
{
    //
}
`

	validatorTemplate = `<?php
return [];
`

	storeTemplate = `<?php
namespace App\Store\{{.Namespace}};

class {{.Name}}
{
    //
}
`

	presenterTemplate = `<?php
namespace App\Presenters\{{.Namespace}};

use App\Presenters\Base\CrudPresenter;

class {{.Name}}Presenter extends CrudPresenter
{
    //
}
`

	controllerTemplate = `<?php
namespace App\Http\Controllers\{{.Namespace}};
use App\Http\Controllers\Base\CrudController;

class {{.Name}}Controller extends CrudController
{
    //
}
`
)

func main() {
	var makePHPCmd = &cobra.Command{
		Use:   "make:php",
		Short: "Creates a new PHP file",
		Long:  `Creates a new PHP file and its associated classes and directories.`,
		Run:   createPHPFile,
	}

	makePHPCmd.Flags().StringP("file", "f", "", "The name of the PHP file to create")
	makePHPCmd.Flags().StringP("father_path", "p", "", "The name of preview struture created ")

	var rootCmd = &cobra.Command{Use: "app"}
	rootCmd.AddCommand(makePHPCmd)

	err := rootCmd.Execute()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

// createPHPFile crea un nuevo archivo PHP con su estructura de directorios y clases asociadas.
func createPHPFile(cmd *cobra.Command, args []string) {
	fileName, err := cmd.Flags().GetString("file")
	fatherPath, err := cmd.Flags().GetString("father_path")
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	if fileName == "" {
		fmt.Println("missing argument: --file")
		os.Exit(1)
	}
	directoryName := fileName

	if fatherPath != "" {
		directoryName = fatherPath
		fileName = directoryName + fileName
	}

	directories := []string{
		appDirectory,
		validatorsDirectory,
		storeDirectory,
		presentersDirectory,
		presentersControllers,
		filepath.Join(modelsDirectory, directoryName),
		filepath.Join(validatorsDirectory, directoryName),
		filepath.Join(storeDirectory, directoryName),
		filepath.Join(presentersDirectory, directoryName),
		filepath.Join(presentersControllers, directoryName),
	}

	files := []struct {
		path     string
		template string
	}{
		{
			path:     filepath.Join(modelsDirectory, directoryName, fileName+".php"),
			template: modelTemplate,
		},
		{
			path:     filepath.Join(validatorsDirectory, directoryName, fileName+"Validator.php"),
			template: validatorTemplate,
		},
		{
			path:     filepath.Join(storeDirectory, directoryName, fileName+".php"),
			template: storeTemplate,
		},
		{
			path:     filepath.Join(presentersDirectory, directoryName, fileName+"Presenter.php"),
			template: presenterTemplate,
		},
		{
			path:     filepath.Join(presentersControllers, directoryName, fileName+"Controller.php"),
			template: controllerTemplate,
		},
	}

	namespace := directoryName
	name := filepath.Base(fileName)

	for _, directory := range directories {
		err := os.MkdirAll(directory, 0755)
		if err != nil {
			fmt.Printf("Error creating directory %s: %v\n", directory, err)
		}
	}

	for _, file := range files {
		tmpl, err := template.New("").Parse(file.template)

		if err != nil {
			fmt.Printf("Error parsing template: %v\n", err)
			continue
		}

		filePath := file.path

		if _, err := os.Stat(filePath); err == nil {
			fmt.Printf("%s already exists, skipping...\n", filePath)
			continue
		}

		f, err := os.Create(filePath)
		if err != nil {
			fmt.Printf("Error creating file %s: %v\n", filePath, err)
			continue
		}

		defer f.Close()

		err = tmpl.Execute(f, struct {
			Namespace string
			Name      string
		}{
			Namespace: namespace,
			Name:      name,
		})

		if err != nil {
			fmt.Printf("Error executing template: %v\n", err)
			continue
		}
		fmt.Printf("Created %s\n", filePath)
	}

}
