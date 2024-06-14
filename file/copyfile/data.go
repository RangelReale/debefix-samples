package copyfile

type FileDataList struct {
	Fields map[string]FileData `json:"fields"`
}

type FileData struct {
	Src  string `yaml:"src"`
	Dest string `yaml:"dest"`
}
