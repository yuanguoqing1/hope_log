package controllers

import "hope_blog/config"

type AI struct {
	api_key  string `yaml:"api_key"`
	base_url string `yaml:"base_url"`
}

func (a *AI) InitAi() {
	a.api_key = config.Global.AI.API_KEY
	a.base_url = config.Global.AI.BASE_URL
}
