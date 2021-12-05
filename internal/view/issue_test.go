package view

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/ankitpokhrel/jira-cli/pkg/adf"
	"github.com/ankitpokhrel/jira-cli/pkg/jira"
	"github.com/ankitpokhrel/jira-cli/pkg/tui"
)

func TestIssueDetailsRenderInPlainView(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer

	data := &jira.Issue{
		Key: "TEST-1",
		Fields: jira.IssueFields{
			Summary: "This is a test",
			Resolution: struct {
				Name string `json:"name"`
			}{Name: "Fixed"},
			Description: &adf.ADF{
				Version: 1,
				DocType: "doc",
				Content: []*adf.Node{
					{
						NodeType: "paragraph",
						Content: []*adf.Node{
							{NodeType: "text", NodeValue: adf.NodeValue{Text: "Test description"}},
						},
					},
				},
			},
			IssueType: struct {
				Name string `json:"name"`
			}{Name: "Bug"},
			Assignee: struct {
				Name string `json:"displayName"`
			}{Name: "Person A"},
			Priority: struct {
				Name string `json:"name"`
			}{Name: "High"},
			Reporter: struct {
				Name string `json:"displayName"`
			}{Name: "Person Z"},
			Status: struct {
				Name string `json:"name"`
			}{Name: "Done"},
			Components: []struct {
				Name string `json:"name"`
			}{{Name: "BE"}, {Name: "FE"}},
			Watches: struct {
				IsWatching bool `json:"isWatching"`
				WatchCount int  `json:"watchCount"`
			}{IsWatching: true, WatchCount: 4},
			Created: "2020-12-13T14:05:20.974+0100",
			Updated: "2020-12-13T14:07:20.974+0100",
		},
	}

	issue := Issue{
		Data:    data,
		Display: DisplayFormat{Plain: true},
	}

	expected := "🐞 Bug  ✅ Done  ⌛ Sun, 13 Dec 20  👷 Person A  🔑️ TEST-1  💭 0 comments  \U0001F9F5 0 linked issues\n# This is a test\n⏱️  Sun, 13 Dec 20  🔎 Person Z  🚀 High  📦 BE, FE  🏷️  None  👀 You + 3 watchers\n\n------------------------ Description ------------------------\n\nTest description\n\n"
	actual := issue.String()

	assert.NoError(t, issue.renderPlain(&b))
	assert.Equal(t, tui.TextData(expected), tui.TextData(actual))
}

func TestIssueDetailsWithV2Description(t *testing.T) {
	t.Parallel()

	var b bytes.Buffer

	data := &jira.Issue{
		Key: "TEST-1",
		Fields: jira.IssueFields{
			Summary: "This is a test",
			Resolution: struct {
				Name string `json:"name"`
			}{Name: "Fixed"},
			Description: "h1. Title\nh2. Subtitle\n\nThis is a *bold* and _italic_ text with [a link|https://ankit.pl] in between.",
			IssueType: struct {
				Name string `json:"name"`
			}{Name: "Bug"},
			Assignee: struct {
				Name string `json:"displayName"`
			}{Name: "Person A"},
			Priority: struct {
				Name string `json:"name"`
			}{Name: "High"},
			Reporter: struct {
				Name string `json:"displayName"`
			}{Name: "Person Z"},
			Status: struct {
				Name string `json:"name"`
			}{Name: "Done"},
			Components: []struct {
				Name string `json:"name"`
			}{{Name: "BE"}, {Name: "FE"}},
			Comment: struct {
				Total int `json:"total"`
			}{Total: 3},
			IssueLinks: []struct {
				LinkType struct {
					Name    string `json:"name"`
					Inward  string `json:"inward"`
					Outward string `json:"outward"`
				} `json:"type"`
				InwardIssue  *jira.Issue `json:"inwardIssue,omitempty"`
				OutwardIssue *jira.Issue `json:"outwardIssue,omitempty"`
			}{
				{
					LinkType: struct {
						Name    string `json:"name"`
						Inward  string `json:"inward"`
						Outward string `json:"outward"`
					}{Name: "blocks", Inward: "blocks", Outward: "is blocked by"},
					InwardIssue: &jira.Issue{
						Key: "TEST-2",
						Fields: jira.IssueFields{
							Summary: "Something is broken", IssueType: struct {
								Name string `json:"name"`
							}{Name: "Bug"}, Priority: struct {
								Name string `json:"name"`
							}{Name: "High"}, Status: struct {
								Name string `json:"name"`
							}{Name: "TO DO"},
						},
					},
				},
				{
					LinkType: struct {
						Name    string `json:"name"`
						Inward  string `json:"inward"`
						Outward string `json:"outward"`
					}{Name: "relates", Inward: "relates", Outward: "relates to"},
					OutwardIssue: &jira.Issue{
						Key: "TEST-3",
						Fields: jira.IssueFields{
							Summary: "Everything is on fire", IssueType: struct {
								Name string `json:"name"`
							}{Name: "Bug"}, Priority: struct {
								Name string `json:"name"`
							}{Name: "Urgent"}, Status: struct {
								Name string `json:"name"`
							}{Name: "Done"},
						},
					},
				},
			},
			Created: "2020-12-13T14:05:20.974+0100",
			Updated: "2020-12-13T14:07:20.974+0100",
		},
	}

	issue := Issue{
		Data:    data,
		Display: DisplayFormat{Plain: true},
	}
	assert.NoError(t, issue.renderPlain(&b))

	expected := "🐞 Bug  ✅ Done  ⌛ Sun, 13 Dec 20  👷 Person A  🔑️ TEST-1  💭 3 comments  \U0001F9F5 2 linked issues\n# This is a test\n⏱️  Sun, 13 Dec 20  🔎 Person Z  🚀 High  📦 BE, FE  🏷️  None  👀 0 watchers\n\n------------------------ Description ------------------------\n\n# Title\n## Subtitle\nThis is a **bold** and _italic_ text with [a link](https://ankit.pl) in between.\n\n\n------------------------ Linked Issues ------------------------\n\n\n  BLOCKS\n\n    TEST-2 Something is broken   • Bug • High   • TO DO\n\n  RELATES TO\n\n    TEST-3 Everything is on fire • Bug • Urgent • Done \n\n"
	actual := issue.String()

	assert.Equal(t, tui.TextData(expected), tui.TextData(actual))
}

func TestSeparator(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		body        string
		plain       bool
		expected    string
		expected256 string
	}{
		{
			name:        "it returns striaght horizontal bar for empty message",
			body:        "",
			expected:    "\x1b[0;90m————————————————————————————————————————————————\x1b[0m",
			expected256: "\x1b[38;5;242m————————————————————————————————————————————————\x1b[m",
		},
		{
			name:        "it returns raw horizontal bar for empty message in plain output",
			body:        "",
			plain:       true,
			expected:    "------------------------------------------------",
			expected256: "------------------------------------------------",
		},
		{
			name:        "it returns greyed out message wrapped in horizontal bar",
			body:        "Some text",
			expected:    "\x1b[0;90m———————————————————————— Some text ————————————————————————\x1b[0m",
			expected256: "\x1b[38;5;242m———————————————————————— Some text ————————————————————————\x1b[m",
		},
		{
			name:        "it returns greyed out message wrapped in raw horizontal bar for plain output",
			body:        "Some text",
			plain:       true,
			expected:    "------------------------ Some text ------------------------",
			expected256: "------------------------ Some text ------------------------",
		},
		{
			name:        "it doesn't trim spaces",
			body:        "  ",
			expected:    "\x1b[0;90m————————————————————————    ————————————————————————\x1b[0m",
			expected256: "\x1b[38;5;242m————————————————————————    ————————————————————————\x1b[m",
		},
		{
			name:        "it doesn't trim spaces for plain output",
			body:        "  ",
			plain:       true,
			expected:    "------------------------    ------------------------",
			expected256: "------------------------    ------------------------",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			issue := Issue{
				Data: &jira.Issue{
					Key: "TEST-1",
				},
				Display: DisplayFormat{Plain: tc.plain},
			}

			if xterm256() {
				assert.Equal(t, tc.expected256, issue.separator(tc.body))
			} else {
				assert.Equal(t, tc.expected, issue.separator(tc.body))
			}
		})
	}
}
