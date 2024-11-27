package entities

type Song struct {
	ID          uint64 `gorm:"primary_key" json:"id,omitempty"`
	Song        string `json:"song"`
	Group       string `json:"group"`
	ReleaseDate string `json:"release_date,omitempty"`
	Text        string `json:"text,omitempty"`
	Link        string `json:"link,omitempty"`
}
