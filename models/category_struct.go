package models

type Category struct {
	Number    int
	Name      string
	Ads_score int
}

type Chosen_Category struct {
	Category_info string
	Chosen        string
}
