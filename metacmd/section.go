package metacmd

import (
	"fmt"
	"io"
	"sort"
	"strings"
)

// Desc holds information about a command or alias description.
type Desc struct {
	Desc   string
	Params string
}

// Section is a meta command section.
type Section string

// Meta command section types.
const (
	SectionGeneral         Section = "General"
	SectionQueryExecute    Section = "Query Execute"
	SectionQueryBuffer     Section = "Query Buffer"
	SectionHelp            Section = "Help"
	SectionTransaction     Section = "Transaction"
	SectionInputOutput     Section = "Input/Output"
	SectionInformational   Section = "Informational"
	SectionFormatting      Section = "Formatting"
	SectionConnection      Section = "Connection"
	SectionOperatingSystem Section = "Operating System"
	SectionVariables       Section = "Variables"
	// SectionLargeObjects    Section = "Large Objects"
)

// String satisfies stringer.
func (s Section) String() string {
	return string(s)
}

// SectionOrder is the order of sections to display via Listing.
var SectionOrder = []Section{
	SectionGeneral, SectionQueryExecute, SectionQueryBuffer, SectionHelp,
	SectionInputOutput, SectionInformational, SectionFormatting,
	SectionTransaction,
	SectionConnection, SectionOperatingSystem, SectionVariables,
}

// Listing writes the formatted command listing to w, separated into different
// sections for all known commands.
func Listing(w io.Writer) {
	sectionDescs := make(map[Section][][]string, len(SectionOrder))
	var plen int
	for _, section := range SectionOrder {
		var descs [][]string
		for _, c := range sectMap[section] {
			cmd := cmds[c]
			s, opts := optText(cmd.Desc)
			descs, plen = add(descs, `  \`+cmd.Name+opts, s, plen)
			// sort aliases
			var aliases []string
			for alias, desc := range cmd.Aliases {
				if desc.Desc == "" && desc.Params == "" {
					continue
				}
				aliases = append(aliases, alias)
			}
			sort.Slice(aliases, func(i, j int) bool {
				return strings.ToLower(aliases[i]) < strings.ToLower(aliases[j])
			})
			for _, alias := range aliases {
				s, opts := optText(cmd.Aliases[alias])
				descs, plen = add(descs, `  \`+strings.TrimSpace(alias)+opts, s, plen)
			}
		}
		sectionDescs[section] = descs
	}
	for _, section := range SectionOrder {
		fmt.Fprintln(w, section)
		for _, line := range sectionDescs[section] {
			fmt.Fprintln(w, rpad(line[0], plen), "", line[1])
		}
		fmt.Fprintln(w)
	}
}

// rpad right pads a string.
func rpad(s string, l int) string {
	return s + strings.Repeat(" ", l-len(s))
}

// add adds b, c to a, returning the max of pad or len(b).
func add(a [][]string, b, c string, pad int) ([][]string, int) {
	return append(a, []string{b, c}), max(pad, len(b))
}

// optText returns a string and the opt text.
func optText(desc Desc) (string, string) {
	if desc.Params != "" {
		return desc.Desc, " " + desc.Params
	}
	return desc.Desc, desc.Params
}

// max returns maximum of a, b.
func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}
