package config

import (
	"fmt"

	"github.com/spf13/viper"
)

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start                 string `mapstructure:"start"`
	WordDone              string `mapstructure:"word_done"`
	WordMiss              string `mapstructure:"word_miss"`
	TestDone              string `mapstructure:"test_done"`
	NoPackages            string `mapstructure:"no_packages"`
	ChoosePackageToStart  string `mapstructure:"choose_package_to_start"`
	ChoosePackageToDelete string `mapstructure:"choose_package_to_delete"`
	InsertPackage         string `mapstructure:"insert_package"`
	InsertPackageName     string `mapstructure:"insert_package_name"`
	AddedPackage          string `mapstructure:"package_added"`
}

type Errors struct {
	UnknownCommand string `mapstructure:"unknown_command"`
	StartLesson    string `mapstructure:"start_lesson"`
	SmallPackage   string `mapstructure:"small_package"`
	EmptyStrings   string `mapstructure:"empty_strings"`
	AlredyExist    string `mapstructure:"already_exist"`
	DoesntWork     string `mapstructure:"dosnt_work"`
}

func InitCfg() (*Messages, error) {
	viper.AddConfigPath("configs")
	viper.SetConfigName("config")
	if err := viper.ReadInConfig(); err != nil {
		fmt.Println(err)
		return nil, err
	}

	var messages Messages

	if err := viper.UnmarshalKey("messages.responses", &messages.Responses); err != nil {
		fmt.Println(err)
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.errors", &messages.Errors); err != nil {
		fmt.Println(err)
		return nil, err
	}
	return &messages, nil
}
