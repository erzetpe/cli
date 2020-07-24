/*
 * Copyright © 2020 Mateusz Kyc
 */

package environment

import (
	"fmt"
	"github.com/google/uuid"
	"github.com/mkyc/epiphany-wrapper-poc/pkg/util"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path"
)

var (
	UsedEnvironmentDirectory string
)

type Environment struct {
	Name string    `yaml:"name"`
	Uuid uuid.UUID `yaml:"uuid"`
}

func init() {
	UsedEnvironmentDirectory = path.Join(util.GetHomeDirectory(), util.DefaultConfigurationDirectory, util.DefaultEnvironmentsSubdirectory)
}

func (e *Environment) Save() error {
	data, err := yaml.Marshal(e)
	if err != nil {
		return err
	}
	err = ioutil.WriteFile(path.Join(UsedEnvironmentDirectory, e.Uuid.String(), util.DefaultEnvironmentConfigFileName), data, 0644)
	if err != nil {
		return err
	}
	return nil
}

func (e *Environment) String() string {
	return fmt.Sprintf(" Name: %s\n UUID: %s\n", e.Name, e.Uuid.String())
}

func Create(name string) (*Environment, error) {
	environment := &Environment{
		Name: name,
		Uuid: uuid.New(),
	}

	newEnvironmentDirectory := path.Join(UsedEnvironmentDirectory, environment.Uuid.String())
	util.EnsureDirectory(newEnvironmentDirectory)
	err := environment.Save()
	if err != nil {
		panic("I wansnt able to save: " + environment.Uuid.String())
	}
	return environment, nil
}

func GetAll() ([]*Environment, error) {
	items, err := ioutil.ReadDir(UsedEnvironmentDirectory)
	if err != nil {
		return nil, err
	}
	var environments []*Environment
	for _, i := range items {
		if i.IsDir() {
			e, err := Get(uuid.MustParse(i.Name()))
			if err == nil {
				environments = append(environments, e)
			}
		}
	}
	return environments, nil
}

func Get(uuid uuid.UUID) (*Environment, error) {
	expectedFile := path.Join(UsedEnvironmentDirectory, uuid.String(), util.DefaultEnvironmentConfigFileName)
	if _, err := os.Stat(expectedFile); os.IsNotExist(err) {
		fmt.Println("file " + expectedFile + " does not exist!") //TODO err?
		return nil, err
	} else {
		e, err := loadEnvironmentFromConfigFile(expectedFile)
		if err != nil {
			fmt.Println("incorrect file?") //TODO warn?
		}
		return e, nil
	}
}

func loadEnvironmentFromConfigFile(configPath string) (*Environment, error) {
	e := &Environment{}
	file, err := os.Open(configPath)
	if err != nil {
		return nil, err
	}
	defer file.Close()
	d := yaml.NewDecoder(file)
	if err := d.Decode(&e); err != nil {
		return nil, err
	}
	return e, nil
}
