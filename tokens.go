package mustache

import ()

type tokens []token
type token interface{}

type tmplToken struct {
	start int
	end   int
}

type valToken struct {
	key    string
	escape bool
}

type sectionToken struct {
	key string
	t   *Template
}

type invertedSectionToken struct {
	key string
	t   *Template
}
