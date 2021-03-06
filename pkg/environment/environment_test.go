package environment

import (
	"bytes"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"path"
	"reflect"
	"regexp"
	"testing"

	"github.com/epiphany-platform/cli/pkg/util"
	"github.com/google/uuid"
	"github.com/rs/zerolog"
)

func setup(t *testing.T, suffix string) (string, string, string) {
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	parentDir := os.TempDir()
	mainDirectory, err := ioutil.TempDir(parentDir, fmt.Sprintf("*-e-environment-%s", suffix))
	if err != nil {
		t.Fatal(err)
	}
	envsDirectory, err := ioutil.TempDir(mainDirectory, "environments-*")
	if err != nil {
		t.Fatal(err)
	}

	tempFile, err := ioutil.TempFile(mainDirectory, "config-*.yaml")
	if err != nil {
		t.Fatal(err)
	}
	return tempFile.Name(), mainDirectory, envsDirectory
}

func TestGet(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "get")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	type args struct {
		uuid uuid.UUID
	}
	tests := []struct {
		name    string
		args    args
		mocked  []byte
		want    *Environment
		wantErr error
	}{
		{
			name: "correct",
			args: args{
				uuid: uuid.MustParse("fccf6810-32c4-4500-9414-2de45d2c4097"),
			},
			mocked: []byte(`name: e1
uuid: fccf6810-32c4-4500-9414-2de45d2c4097
installed: []`),
			want: &Environment{
				Name:      "e1",
				Uuid:      uuid.MustParse("fccf6810-32c4-4500-9414-2de45d2c4097"),
				Installed: []InstalledComponentVersion{},
			},
			wantErr: nil,
		},
		{
			name: "missing file",
			args: args{
				uuid: uuid.MustParse("816789aa-7839-4f2b-ac74-b66344e4fbe8"),
			},
			wantErr: errors.New("stat .*-e-environment-get/environments-.*/816789aa-7839-4f2b-ac74-b66344e4fbe8/config.yaml: no such file or directory"),
		},
		{
			name: "incorrect file",
			args: args{
				uuid: uuid.MustParse("66d4cd70-4375-4737-b6ce-7e13f3cc93f9"),
			},
			mocked:  []byte(`incorrect file`),
			wantErr: errors.New("yaml: unmarshal errors:\n  line 1: cannot unmarshal !!str `incorre...` into environment.Environment"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if len(tt.mocked) > 0 {
				envDir := path.Join(util.UsedEnvironmentDirectory, tt.args.uuid.String())
				err := os.MkdirAll(envDir, 0755)
				if err != nil {
					t.Fatal(err)
				}
				envConfigFile := path.Join(envDir, util.DefaultEnvironmentConfigFileName)
				err = ioutil.WriteFile(envConfigFile, tt.mocked, 0644)
				if err != nil {
					t.Fatal(err)
				}
			}
			got, err := Get(tt.args.uuid)
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func TestGetAll(t *testing.T) {
	type mocked struct {
		subdirectory  string
		configName    string
		configContent []byte
	}
	tests := []struct {
		name        string
		mocked      []mocked
		want        []*Environment
		wantErr     error
		shouldPanic bool
	}{
		{
			name: "correct",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory: "4af1705c-c48f-4ca2-be08-53c673da835c",
					configName:   "config.yaml",
					configContent: []byte(`name: e1
uuid: 4af1705c-c48f-4ca2-be08-53c673da835c
installed: []`),
				},
			},
			want: []*Environment{
				{
					Name:      "e2",
					Uuid:      uuid.MustParse("45764648-162a-4526-bdd0-71a438fd6ceb"),
					Installed: []InstalledComponentVersion{},
				},
				{
					Name:      "e1",
					Uuid:      uuid.MustParse("4af1705c-c48f-4ca2-be08-53c673da835c"),
					Installed: []InstalledComponentVersion{},
				},
			},
			wantErr:     nil,
			shouldPanic: false,
		},
		{
			name: "subdirectory name not uuid",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory: "incorrect-directory-name",
					configName:   "config.yaml",
					configContent: []byte(`name: e1
uuid: 4af1705c-c48f-4ca2-be08-53c673da835c
installed: []`),
				},
			},
			wantErr:     errors.New("uuid: Parse\\(incorrect-directory-name\\): invalid UUID length: 24"),
			shouldPanic: true,
		},
		{
			name: "incorrect config file name",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory: "4af1705c-c48f-4ca2-be08-53c673da835c",
					configName:   "incorrect-config.yaml",
					configContent: []byte(`name: e1
uuid: 4af1705c-c48f-4ca2-be08-53c673da835c
installed: []`),
				},
			},
			want: []*Environment{
				{
					Name:      "e2",
					Uuid:      uuid.MustParse("45764648-162a-4526-bdd0-71a438fd6ceb"),
					Installed: []InstalledComponentVersion{},
				},
			},
			wantErr:     nil,
			shouldPanic: false,
		},
		{
			name: "incorrect config file content",
			mocked: []mocked{
				{
					subdirectory: "45764648-162a-4526-bdd0-71a438fd6ceb",
					configName:   "config.yaml",
					configContent: []byte(`name: e2
uuid: 45764648-162a-4526-bdd0-71a438fd6ceb
installed: []`),
				},
				{
					subdirectory:  "4af1705c-c48f-4ca2-be08-53c673da835c",
					configName:    "config.yaml",
					configContent: []byte(`incorrect content`),
				},
			},
			want: []*Environment{
				{
					Name:      "e2",
					Uuid:      uuid.MustParse("45764648-162a-4526-bdd0-71a438fd6ceb"),
					Installed: []InstalledComponentVersion{},
				},
			},
			wantErr:     nil,
			shouldPanic: false,
		},
	}
	for _, tt := range tests {
		util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "get-all")
		t.Run(tt.name, func(t *testing.T) {
			if tt.mocked != nil && len(tt.mocked) > 0 {
				for _, m := range tt.mocked {
					envDir := path.Join(util.UsedEnvironmentDirectory, m.subdirectory)
					err := os.MkdirAll(envDir, 0755)
					if err != nil {
						t.Fatal(err)
					}
					envConfigFile := path.Join(envDir, m.configName)
					err = ioutil.WriteFile(envConfigFile, m.configContent, 0644)
					if err != nil {
						t.Fatal(err)
					}
				}
			}

			f := func() {
				got, err := GetAll()
				if isWrongResult(t, err, tt.wantErr) {
					return
				}
				if !reflect.DeepEqual(got, tt.want) {
					t.Errorf("got = %#v, want %#v", got, tt.want)
				}
			}

			if tt.shouldPanic {
				assertPanic(t, f, tt.wantErr)
			} else {
				f()
			}
		})
		os.RemoveAll(util.UsedConfigurationDirectory)
	}
}

func Test_create(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "create")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	type args struct {
		name string
		uuid string
	}
	tests := []struct {
		name    string
		args    args
		want    *Environment
		wantErr error
	}{
		{
			name: "correct",
			args: args{
				name: "e1",
				uuid: "b03bb900-5d49-4421-a45e-eeeb40e0a5d5",
			},
			want: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("b03bb900-5d49-4421-a45e-eeeb40e0a5d5"),
			},
			wantErr: nil,
		},
		{
			name: "empty name",
			args: args{
				name: "",
				uuid: "66d4cd70-4375-4737-b6ce-7e13f3cc93f9",
			},
			want: &Environment{
				Uuid: uuid.MustParse("66d4cd70-4375-4737-b6ce-7e13f3cc93f9"),
			},
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := create(tt.args.name, uuid.MustParse(tt.args.uuid))
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %+v, want %+v", got, tt.want)
				return
			}
			expectedConfigFile := path.Join(util.UsedEnvironmentDirectory, got.Uuid.String(), util.DefaultEnvironmentConfigFileName)
			if _, err := os.Stat(expectedConfigFile); os.IsNotExist(err) {
				t.Errorf("expected to find file %s but didn't find", expectedConfigFile)
			}
		})
	}
}

func TestEnvironment_Save(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "env-save")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name        string
		environment *Environment
		wantContent []byte
		wantErr     error
	}{
		{
			name: "correct",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
			},
			wantContent: []byte(`name: e1
uuid: 10d52c05-029e-4794-a790-79d6c2af40b6
installed: []
`),
			wantErr: nil,
		},
		{
			name: "missing uuid",
			environment: &Environment{
				Name: "e1",
			},
			wantErr: errors.New("unexpected UUID on Save: 00000000-0000-0000-0000-000000000000"),
		},
		{
			name: "missing name",
			environment: &Environment{
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
			},
			wantContent: []byte(`name: ""
uuid: 10d52c05-029e-4794-a790-79d6c2af40b6
installed: []
`),
			wantErr: nil,
		},
		{
			name: "with installed components",
			environment: &Environment{
				Name: "x",
				Uuid: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("3e5b7269-1b3d-4003-9454-9f472857633a"),
						Name:           "x",
						Type:           "x",
						Version:        "x",
						Image:          "x",
						WorkDirectory:  "x",
						Mounts:         []string{"x"},
						Commands:       []InstalledComponentCommand{},
					},
				},
			},
			wantContent: []byte(`name: x
uuid: 3e5b7269-1b3d-4003-9454-9f472857633a
installed:
- environment_ref: 3e5b7269-1b3d-4003-9454-9f472857633a
  name: x
  type: x
  version: x
  image: x
  workdir: x
  mounts:
  - x
  commands: []
`),
			wantErr: nil,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			dir := path.Join(util.UsedEnvironmentDirectory, tt.environment.Uuid.String())
			err := os.MkdirAll(dir, 0755)
			if err != nil {
				t.Fatal(err)
			}
			e := &Environment{
				Name:      tt.environment.Name,
				Uuid:      tt.environment.Uuid,
				Installed: tt.environment.Installed,
			}
			err = e.Save()
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if len(tt.wantContent) > 0 {
				expectedConfigFile := path.Join(dir, util.DefaultEnvironmentConfigFileName)
				if _, err := os.Stat(expectedConfigFile); os.IsNotExist(err) {
					t.Errorf("expected to find file %s but didn't find", expectedConfigFile)
					return
				}
				savedBytes, _ := ioutil.ReadFile(expectedConfigFile)
				if !bytes.Equal(tt.wantContent, savedBytes) {
					t.Errorf("saved file is \n%s\n\n but expected is \n\n%s\n", string(savedBytes), string(tt.wantContent))
				}
			}
		})
	}
}

func TestEnvironment_GetComponentByName(t *testing.T) {
	util.UsedConfigFile, util.UsedConfigurationDirectory, util.UsedEnvironmentDirectory = setup(t, "env-get-by-name")
	defer os.RemoveAll(util.UsedConfigurationDirectory)

	tests := []struct {
		name          string
		environment   *Environment
		componentName string
		want          *InstalledComponentVersion
		wantErr       error
	}{
		{
			name: "correct single",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c1",
						Type:           "docker",
						Version:        "v1",
						Image:          "i1",
					},
				},
			},
			componentName: "c1",
			want: &InstalledComponentVersion{
				EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Name:           "c1",
				Type:           "docker",
				Version:        "v1",
				Image:          "i1",
			},
			wantErr: nil,
		},
		{
			name: "correct multiple",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c1",
						Type:           "d2",
						Version:        "v1",
						Image:          "i1",
					},
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c2",
						Type:           "d2",
						Version:        "v2",
						Image:          "i3",
					},
				},
			},
			componentName: "c1",
			want: &InstalledComponentVersion{
				EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Name:           "c1",
				Type:           "d2",
				Version:        "v1",
				Image:          "i1",
			},
			wantErr: nil,
		},
		{
			name: "missing",
			environment: &Environment{
				Name: "e1",
				Uuid: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{
					{
						EnvironmentRef: uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
						Name:           "c2",
						Type:           "d2",
						Version:        "v2",
						Image:          "i3",
					},
				},
			},
			componentName: "c1",
			wantErr:       errors.New("no such component installed"),
		},
		{
			name: "empty",
			environment: &Environment{
				Name:      "e1",
				Uuid:      uuid.MustParse("10d52c05-029e-4794-a790-79d6c2af40b6"),
				Installed: []InstalledComponentVersion{},
			},
			componentName: "c1",
			wantErr:       errors.New("no such component installed"),
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := &Environment{
				Name:      tt.environment.Name,
				Uuid:      tt.environment.Uuid,
				Installed: tt.environment.Installed,
			}
			got, err := e.GetComponentByName(tt.componentName)
			if isWrongResult(t, err, tt.wantErr) {
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("got = %#v, want %#v", got, tt.want)
			}
		})
	}
}

func isWrongResult(t *testing.T, err error, wantErr error) bool {
	if err != nil && wantErr != nil {
		re := regexp.MustCompile(wantErr.Error())
		if !re.MatchString(err.Error()) {
			t.Errorf("got \n%v\n, want \n%v\n", err, wantErr)
			return true
		}
	} else if err == nil && wantErr != nil {
		t.Errorf("didn't got error but want: %v", wantErr)
		return true
	} else if err != nil && wantErr == nil {
		t.Errorf("didnt want error but got: %v", err)
		return true
	}
	return false
}

func assertPanic(t *testing.T, f func(), wantErr error) {
	defer func() {
		if r := recover(); r == nil {
			t.Errorf("The code did not panic")
		} else {
			if rc, ok := r.(string); ok {
				if isWrongResult(t, errors.New(rc), wantErr) {
					return
				}
			} else {
				t.Errorf("cannot cast recover resutl to string. r = %#v", r)
				return
			}
		}
	}()
	f()
}
