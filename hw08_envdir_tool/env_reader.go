package main

import (
	"os"
	"path/filepath"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	env := make(Environment)
	dirEntries, err := os.ReadDir(dir)
	if err != nil {
		return env, err
	}
	var fname string
	for _, entry := range dirEntries {
		if entry.IsDir() {
			continue
		}
		fname = entry.Name()
		envval, err := ReadEnvFile(filepath.Join(dir, fname))
		if err == nil {
			env[fname] = envval
		}
	}
	return env, nil
}

// ReadEnvFile reads file and return EnvValue.
func ReadEnvFile(fname string) (EnvValue, error) {
	var envval EnvValue

	fData, err := os.ReadFile(fname)
	if err != nil {
		return envval, err
	}
	if len(fData) == 0 {
		envval.NeedRemove = true
		return envval, err
	}

	var resBuilder strings.Builder
	for _, r := range fData {
		if r == 0 {
			resBuilder.WriteString("\n")
			continue
		}
		if r == 10 || r == 13 { // LF, CR
			break
		}
		resBuilder.WriteByte(r)
	}
	envval.Value = strings.TrimRight(resBuilder.String(), "\t ")

	return envval, nil
}

// UpdateOsEnv updates os.Environ with Environment values.
func (env *Environment) UpdateOsEnv() error {
	for k, v := range *env {
		var err error
		switch v.NeedRemove {
		case true:
			err = os.Unsetenv(k)
		case false:
			err = os.Setenv(k, v.Value)
		}
		if err != nil {
			return err
		}
	}
	return nil
}
