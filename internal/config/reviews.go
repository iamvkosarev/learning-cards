package config

type ReviewsUseCase struct {
	AnswerInfluencePercent         float32 `yaml:"answer_influence_percent"`
	SelectDurationInfluencePercent float32 `yaml:"select_duration_influence_percent"`
}

type CardsServer struct {
	Address string `json:"address"`
}

type ReviewsConfig struct {
	Common         Config `json:"common"`
	Database       `yaml:"database"`
	CardsServer    `yaml:"cards_server"`
	ReviewsUseCase `yaml:"use_case"`
}
