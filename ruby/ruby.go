package ruby

import (
	"bufio"
	"strings"
)

func Parse(src string) *RubyClass {
	p := &parser{Src: src}
	return p.Parse()
}

type parser struct {
	Src     string
	rc      *RubyClass
	mode    int
	level   int
	scratch string
}

const (
	base = iota
	class
	function
	selfClass
	classFunction
)

func (p *parser) Parse() *RubyClass {
	p.rc = &RubyClass{File: p.Src}
	s := bufio.NewScanner(strings.NewReader(p.Src))
	for s.Scan() {
		p.ParseLine(s.Text())
	}
	return p.rc
}

func (p *parser) ParseLine(line string) {
	trimmed := strings.TrimSpace(line)
	switch {
	case trimmed == "end":
		p.scratch += "\n" + line
		p.level--
		switch p.level {
		case 0:
			p.rc.Src = p.scratch
		case 1:
			switch p.mode {
			case function:
			case classFunction:
			}
		}
	case strings.HasPrefix(trimmed, "require"):
		p.rc.Libs = append(p.rc.Libs, strings.SplitN(trimmed, " ", 2)[1])
	case strings.HasPrefix(trimmed, "class "):
		className := strings.SplitN(trimmed, " ", 2)[1]
		classParts := strings.Split(className, "::")
		if len(classParts) == 1 {
			p.rc.Name = className
		} else {
			p.rc.Name = classParts[len(classParts)-1]
			p.rc.Modules = classParts[:len(classParts)-1]
		}
		p.mode = class
		p.level++
	}
}

type RubyClass struct {
	File string
	Src  string

	// from class ... < ...
	Modules []string
	Name    string
	Parent  string

	// attr_accessor, attr_writer, etc
	Attrs []Attr

	// from requires
	Libs []string

	// include ...
	Includes []string

	// unrecognized class functions
	Prelude string

	// def self. or class << self
	ClassFuncs []Func

	// Instance functions def ...
	InstFuncs []Func
}

type Attr struct {
	Doc  string
	Name string
	Type string
	Src  string
}

type Func struct {
	Name string
	Args []string
	Code string
	Src  string
}
