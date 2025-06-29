package config

import (
	"time"
)

type JapaneseReading struct {
	SearchTimeout time.Duration `yaml:"search_timeout"`
	MecabDicDir   string        `yaml:"mecab_dic_dir"`
}

type CardsConfig struct {
	Common          Config `yaml:"common"`
	Database        `yaml:"database"`
	JapaneseReading JapaneseReading `yaml:"japanese_reading"`
}
