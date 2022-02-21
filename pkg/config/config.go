package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/spf13/viper"
)

type Config struct {
	TelegramToken string
	DbUser        string
	DbPassword    string
	Messages
}

type Messages struct {
	Responses
	Errors
}

type Responses struct {
	Start            string `mapstructure:"start"`
	Welcome          string `mapstructure:"welcome"`
	UnknownCommand   string `mapstructure:"unknown_command"`
	UnknownAsk       string `mapstructure:"unknown_ask"`
	TextOnly         string `mapstructure:"text_only"`
	HowMuchCallories string `mapstructure:"how_much_callories"`
	InvalidText      string `mapstructure:"invalid_text"`
	ProductExists    string `mapstructure:"product_exists"`
	ProductAdded     string `mapstructure:"product_added"`
	StartLunch       string `mapstructure:"start_lunch"`
	RestLunch        string `mapstructure:"rest_lunch"`
	StopLunch        string `mapstructure:"stop_lunch"`
	WhatProduct      string `mapstructure:"what_product"`
	ProductNotExists string `mapstructure:"product_not_exists"`
	ProductCallories string `mapstructure:"product_callories"`
	DayReport        string `mapstructure:"day_report"`
	WeekReport       string `mapstructure:"week_report"`
	MonthReport      string `mapstructure:"month_report"`
	AskList          string `mapstructure:"ask_list"`
}

type Errors struct {
	FailedStart          string `mapstructure:"failed_start"`
	UserParamError       string `mapstructure:"user_param_error"`
	InsertFoodError      string `mapstructure:"insert_food_error"`
	CountCalloriesError  string `mapstructure:"count_callories_error"`
	ReportErrorNotFound  string `mapstructure:"report_error_not_found"`
	ReportErrorFailed    string `mapstructure:"report_error_failed"`
	FindCalloriesError   string `mapstructure:"find_callories_error"`
	FindProductNameError string `mapstructure:"find_product_name_error"`
}

func Init() (*Config, error) {
	var cfg Config
	if err := parseEnv(&cfg); err != nil {
		return nil, err
	}
	if err := setUpViper(); err != nil {
		log.Fatal("Config load error")
	}
	err := viper.Unmarshal(&cfg)
	if err != nil {
		log.Printf("Unable to decode into struct, %v", err)
	}
	if err := viper.UnmarshalKey("messages.response", &cfg.Messages.Responses); err != nil {
		return nil, err
	}

	if err := viper.UnmarshalKey("messages.error", &cfg.Messages.Errors); err != nil {
		return nil, err
	}

	return &cfg, nil
}

func parseEnv(cfg *Config) error {
	if err := godotenv.Load(); err != nil {
		log.Print("No .env file found")
	}

	TelegramToken, exists := os.LookupEnv("TOKEN")
	if !exists {
		log.Printf("No %s var found", TelegramToken)
	}
	DbUser, exists := os.LookupEnv("DBUSER")
	if !exists {
		log.Printf("No %s var found", DbUser)
	}
	DbPassword, exists := os.LookupEnv("DBPASSWORD")
	if !exists {
		log.Printf("No %s var found", DbPassword)
	}

	cfg.TelegramToken = TelegramToken
	cfg.DbUser = DbUser
	cfg.DbPassword = DbPassword

	return nil
}

func setUpViper() error {
	viper.AddConfigPath("configs")
	viper.SetConfigName("main")

	return viper.ReadInConfig()
}
