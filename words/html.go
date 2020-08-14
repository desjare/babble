package words

import (
	"bytes"
	"html"
	"regexp"
	"strings"
)

func HTMLGetTagName(tag []byte) (name string, end bool) {
	end = tag[1] == '/' || tag[len(tag)-2] == '/' || tag[1] == '!'
	startPos := 1
	if tag[1] == '/' {
		startPos = 2
	}
	for i := 0; i < len(tag); i++ {
		if tag[i] == ' ' {
			name = string(tag[startPos:i])
			return
		}
	}
	name = string(tag[startPos : len(tag)-1])
	return
}

func HTMLGetLink(t []byte) string {
	tag := string(t)
	href := `href="`
	start := strings.Index(tag, href)
	if start >= 0 {
		end := strings.Index(tag[start+len(href):], `"`)
		if end >= 0 {
			end = start + len(href) + end
			start = start + len(href)
			return string(tag[start:end])
		}
	}
	return ""
}

func HTMLUnescapeString(data string) (s string) {
	switch data {
	case "&nbsp;":
		s = " "
	case "\n", "\r":
		s = ""
	default:
		s = html.UnescapeString(data)
	}
	return
}

func HTMLProcessTag(body []byte, buffer *bytes.Buffer, i int, tags [][]int, resc *regexp.Regexp) {
	if i < len(tags)-1 && tags[i+1][0] > tags[i][1] {
		data := body[tags[i][1]:tags[i+1][0]]
		escaped := resc.ReplaceAllStringFunc(string(data), HTMLUnescapeString)
		buffer.WriteString(escaped)
	}
	return
}

func HTMLProcessTags(body []byte, tags [][]int) (content string, links []string) {

	var buffer bytes.Buffer

	resc, err := regexp.Compile("(&(.)+;)|(\n)|(\r)")
	if err != nil {
		panic(err)
	}

	for i := 0; i < len(tags)-1; i++ {

		tag := body[tags[i][0]:tags[i][1]]
		name, closing := HTMLGetTagName(tag)

		switch name {
		case "a":
			link := HTMLGetLink(tag)
			if len(link) > 0 {
				links = append(links, link)
			}
		case "style", "img", "ins":
			if closing == false {
				continue
			}
		case "script", "object":
			if closing == false {
				// look for end script
				for i = i + 1; i < len(tags)-1; i++ {
					tag = body[tags[i][0]:tags[i][1]]
					tagName, closing := HTMLGetTagName(tag)
					if tagName == name && closing == true {
						break
					}
				}
			}
		case "p", "h1", "h2", "h3", "h4", "h5", "option", "div", "span", "br", "li", "ul":
			buffer.WriteString("\t")
		}
		HTMLProcessTag(body, &buffer, i, tags, resc)
	}
	content = buffer.String()
	return
}

func HTMLParse(body []byte) (content string, links []string) {
	tagr, err := regexp.Compile("<(.)+?>")
	if err != nil {
		panic(err)
	}
	tags := tagr.FindAllIndex(body, -1)
	content, links = HTMLProcessTags(body, tags)
	return
}
