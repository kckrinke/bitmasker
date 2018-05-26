// Copyright 2018 Kevin C. Krinke. All rights reserved.
// Use of this source code is governed by the LGPLv3
// license that can be found in the LICENSE.md file.
// Copyright 2014 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// This file contains simple golden tests for various examples.
// Besides validating the results when the implementation changes,
// it provides a way to look at the generated code without having
// to execute the print statements in one's head.

package main

import (
	"bytes"
	"fmt"
	godmp "github.com/sergi/go-diff/diffmatchpatch"
	"strings"
	"testing"
)

// Golden represents a test case.
type Golden struct {
	name        string
	trimPrefix  string
	lineComment bool
	input       string // input; the package clause is provided when running the test.
	output      string // exected output.
}

var golden = []Golden{
	{"unum", "", false, unum_in, unum_out},
}

// Each example starts with "type XXX uint", with a single space separating them.

// Unsigned integers spanning zero.
const unum_in = `type Unum uint
const (
	m_2 Unum = iota + 253
	m_1
)

const (
	m0 Unum = iota
	m1
	m2
)
`

const unum_out = `
type IUnum interface {
	Has(m Unum) bool
	Set(m Unum) Unum
	Clear(m Unum) Unum
	Toggle(m Unum) Unum
	String() string
}

func (i Unum) Has(m Unum) bool {
	return i&m != 0
}

func (i Unum) Set(m Unum) Unum {
	return i | m
}

func (i Unum) Clear(m Unum) Unum {
	return i &^ m
}

func (i Unum) Toggle(m Unum) Unum {
	return i ^ m
}

const (
	_Unum_name_0 = "m0m1m2"
	_Unum_name_1 = "m_2m_1"
)

var (
	_Unum_index_0 = [...]uint8{0, 2, 4, 6}
	_Unum_index_1 = [...]uint8{0, 3, 6}
)

func (i Unum) String() string {
	switch {
	case 0 <= i && i <= 2:
		return _Unum_name_0[_Unum_index_0[i]:_Unum_index_0[i+1]]
	case 253 <= i && i <= 254:
		i -= 253
		return _Unum_name_1[_Unum_index_1[i]:_Unum_index_1[i+1]]
	default:
		return "Unum(" + strconv.FormatInt(int64(i), 10) + ")"
	}
}
`

func TestGolden(t *testing.T) {
	for _, test := range golden {
		g := Generator{
			trimPrefix:  test.trimPrefix,
			lineComment: test.lineComment,
		}
		input := "package test\n" + test.input
		file := test.name + ".go"
		g.parsePackage(".", []string{file}, input)
		// Extract the name and type of the constant from the first line.
		tokens := strings.SplitN(test.input, " ", 3)
		if len(tokens) != 3 {
			t.Fatalf("%s: need type declaration on first line", test.name)
		}
		g.generate(tokens[1])
		got := string(g.format())
		if got != test.output {
			dmp := godmp.New()
			diffs := dmp.DiffMain(test.output, got, false)
			var buf bytes.Buffer
			_, _ = buf.WriteString(fmt.Sprintf("diff wanted/%[1]s.go received/%[1]s.go\n", test.name))
			_, _ = buf.WriteString(fmt.Sprintf("--- wanted/%s.go\n", test.name))
			_, _ = buf.WriteString(fmt.Sprintf("+++ received/%s.go\n", test.name))
			_, _ = buf.WriteString("@@ -0,0 +0,0 @@\n")
			for _, diff := range diffs {
				text := diff.Text
				switch diff.Type {
				case godmp.DiffInsert:
					_, _ = buf.WriteString("+ ")
				case godmp.DiffDelete:
					_, _ = buf.WriteString("- ")
				case godmp.DiffEqual:
					_, _ = buf.WriteString("  ")

				}
				_, _ = buf.WriteString(text)
			}
			t.Errorf(
				"%s: diff of got vs expected\n=======\n%s\n=======\n",
				test.name,
				buf.String(),
			)
		}
	}
}
