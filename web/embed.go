package web

import "embed"

var (
	//go:embed templates/* assets/*
	Files embed.FS
)
