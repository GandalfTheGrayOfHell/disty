package main

type Project struct {
	Remote string `json:"remote"`
}

func (p *Project) Default() {
	p.Remote = ""
}
