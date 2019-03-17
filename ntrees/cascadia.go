package ntrees

import (
	"bytes"
	"errors"
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Query parses a selector and returns, if successful, a Selector object
// that can be used to match against Node objects.
func Query(sel string) (Selector, error) {
	p := &parser{s: sel}
	compiled, err := p.parseSelectorGroup()
	if err != nil {
		return nil, err
	}

	if p.i < len(sel) {
		return nil, fmt.Errorf("parsing %q: %d bytes left over", sel, len(sel)-p.i)
	}

	return compiled, nil
}

//*******************************************************
// css selector functions
//*******************************************************

// the Selector type, and functions for creating them

// A Selector is a function which tells whether a node matches or not.
type Selector func(*Node) bool

// hasChildMatch returns whether n has any child that matches a.
func hasChildMatch(n *Node, a Selector) bool {
	for c, err := n.FirstChild(); err == nil; c, err = c.NextSibling() {
		if a(c) {
			return true
		}
	}
	return false
}

// hasDescendantMatch performs a depth-first search of n's descendants,
// testing whether any of them match a. It returns true as soon as a match is
// found, or false if no match is found.
func hasDescendantMatch(n *Node, a Selector) bool {
	for c, err := n.FirstChild(); err == nil; c, err = c.NextSibling() {
		if a(c) || (c.Type() == ElementNode && hasDescendantMatch(c, a)) {
			return true
		}
	}
	return false
}

// MatchAll returns a slice of the nodes that match the selector,
// from n and its children.
func (s Selector) MatchAll(n *Node) []*Node {
	return s.matchAllInto(n, nil)
}

func (s Selector) matchAllInto(n *Node, storage []*Node) []*Node {
	if s(n) {
		storage = append(storage, n)
	}

	for c, err := n.FirstChild(); err == nil; c, err = c.NextSibling() {
		storage = s.matchAllInto(c, storage)
	}

	return storage
}

// Match returns true if the node matches the selector.
func (s Selector) Match(n *Node) bool {
	return s(n)
}

// MatchFirst returns the first node that matches s, from n and its children.
func (s Selector) MatchFirst(n *Node) *Node {
	if s.Match(n) {
		return n
	}

	for c, err := n.FirstChild(); err == nil; c, err = c.NextSibling() {
		m := s.MatchFirst(c)
		if m != nil {
			return m
		}
	}
	return nil
}

// Filter returns the nodes in nodes that match the selector.
func (s Selector) Filter(nodes []*Node) (result []*Node) {
	for _, n := range nodes {
		if s(n) {
			result = append(result, n)
		}
	}
	return result
}

// typeSelector returns a Selector that matches elements with a given tag name.
func typeSelector(tag string) Selector {
	tag = toLowerASCII(tag)
	return func(n *Node) bool {
		return n.Type() == ElementNode && n.Name() == tag
	}
}

// toLowerASCII returns s with all ASCII capital letters lowercased.
func toLowerASCII(s string) string {
	var b []byte
	for i := 0; i < len(s); i++ {
		if c := s[i]; 'A' <= c && c <= 'Z' {
			if b == nil {
				b = make([]byte, len(s))
				copy(b, s)
			}
			b[i] = s[i] + ('a' - 'A')
		}
	}

	if b == nil {
		return s
	}

	return string(b)
}

// attributeSelector returns a Selector that matches elements
// where the attribute named key satisfies the function f.
func attributeSelector(key string, f func(string) bool) Selector {
	key = toLowerASCII(key)
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}

		if key == "id" {
			if f(n.ID()) {
				return true
			}
			return false
		}

		if attr, ok := n.Attrs.Attr(key); ok {
			if f(attr.Text()) {
				return true
			}
		}
		return false
	}
}

// attributeExistsSelector returns a Selector that matches elements that have
// an attribute named key.
func attributeExistsSelector(key string) Selector {
	return attributeSelector(key, func(string) bool { return true })
}

// attributeEqualsSelector returns a Selector that matches elements where
// the attribute named key has the value val.
func attributeEqualsSelector(key, val string) Selector {
	return attributeSelector(key,
		func(s string) bool {
			return s == val
		})
}

// attributeNotEqualSelector returns a Selector that matches elements where
// the attribute named key does not have the value val.
func attributeNotEqualSelector(key, val string) Selector {
	key = toLowerASCII(key)
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}

		if key == "id" {
			if n.ID() == val {
				return false
			}
			return true
		}

		if attr, ok := n.Attrs.Attr(key); ok {
			if attr.Text() == val {
				return false
			}
		}
		return true
	}
}

// attributeIncludesSelector returns a Selector that matches elements where
// the attribute named key is a whitespace-separated list that includes val.
func attributeIncludesSelector(key, val string) Selector {
	return attributeSelector(key,
		func(s string) bool {
			for s != "" {
				i := strings.IndexAny(s, " \t\r\n\f")
				if i == -1 {
					return s == val
				}
				if s[:i] == val {
					return true
				}
				s = s[i+1:]
			}
			return false
		})
}

// attributeDashmatchSelector returns a Selector that matches elements where
// the attribute named key equals val or starts with val plus a hyphen.
func attributeDashmatchSelector(key, val string) Selector {
	return attributeSelector(key,
		func(s string) bool {
			if s == val {
				return true
			}
			if len(s) <= len(val) {
				return false
			}
			if s[:len(val)] == val && s[len(val)] == '-' {
				return true
			}
			return false
		})
}

// attributePrefixSelector returns a Selector that matches elements where
// the attribute named key starts with val.
func attributePrefixSelector(key, val string) Selector {
	return attributeSelector(key,
		func(s string) bool {
			if strings.TrimSpace(s) == "" {
				return false
			}
			return strings.HasPrefix(s, val)
		})
}

// attributeSuffixSelector returns a Selector that matches elements where
// the attribute named key ends with val.
func attributeSuffixSelector(key, val string) Selector {
	return attributeSelector(key,
		func(s string) bool {
			if strings.TrimSpace(s) == "" {
				return false
			}
			return strings.HasSuffix(s, val)
		})
}

// attributeSubstringSelector returns a Selector that matches nodes where
// the attribute named key contains val.
func attributeSubstringSelector(key, val string) Selector {
	return attributeSelector(key,
		func(s string) bool {
			if strings.TrimSpace(s) == "" {
				return false
			}
			return strings.Contains(s, val)
		})
}

// attributeRegexSelector returns a Selector that matches nodes where
// the attribute named key matches the regular expression rx
func attributeRegexSelector(key string, rx *regexp.Regexp) Selector {
	return attributeSelector(key,
		func(s string) bool {
			return rx.MatchString(s)
		})
}

// intersectionSelector returns a selector that matches nodes that match
// both a and b.
func intersectionSelector(a, b Selector) Selector {
	return func(n *Node) bool {
		return a(n) && b(n)
	}
}

// unionSelector returns a selector that matches elements that match
// either a or b.
func unionSelector(a, b Selector) Selector {
	return func(n *Node) bool {
		return a(n) || b(n)
	}
}

// negatedSelector returns a selector that matches elements that do not match a.
func negatedSelector(a Selector) Selector {
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}
		return !a(n)
	}
}

// writeNodeText writes the text contained in n and its descendants to b.
func writeNodeText(n *Node, b *bytes.Buffer) {
	switch n.Type() {
	case TextNode:
		b.WriteString(n.Name())
	case ElementNode:
		for c, err := n.FirstChild(); err == nil; c, err = c.NextSibling() {
			writeNodeText(c, b)
		}
	}
}

// nodeText returns the text contained in n and its descendants.
func nodeText(n *Node) string {
	return n.Text()
}

// nodeOwnText returns the contents of the text nodes that are direct
// children of n.
func nodeOwnText(n *Node) string {
	var b bytes.Buffer
	for c, err := n.FirstChild(); err == nil; c, err = c.NextSibling() {
		if c.Type() == TextNode {
			b.WriteString(c.Name())
		}
	}
	return b.String()
}

// textSubstrSelector returns a selector that matches nodes that
// contain the given text.
func textSubstrSelector(val string) Selector {
	return func(n *Node) bool {
		text := strings.ToLower(nodeText(n))
		return strings.Contains(text, val)
	}
}

// ownTextSubstrSelector returns a selector that matches nodes that
// directly contain the given text
func ownTextSubstrSelector(val string) Selector {
	return func(n *Node) bool {
		text := strings.ToLower(nodeOwnText(n))
		return strings.Contains(text, val)
	}
}

// textRegexSelector returns a selector that matches nodes whose text matches
// the specified regular expression
func textRegexSelector(rx *regexp.Regexp) Selector {
	return func(n *Node) bool {
		return rx.MatchString(nodeText(n))
	}
}

// ownTextRegexSelector returns a selector that matches nodes whose text
// directly matches the specified regular expression
func ownTextRegexSelector(rx *regexp.Regexp) Selector {
	return func(n *Node) bool {
		return rx.MatchString(nodeOwnText(n))
	}
}

// hasChildSelector returns a selector that matches elements
// with a child that matches a.
func hasChildSelector(a Selector) Selector {
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}
		return hasChildMatch(n, a)
	}
}

// hasDescendantSelector returns a selector that matches elements
// with any descendant that matches a.
func hasDescendantSelector(a Selector) Selector {
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}
		return hasDescendantMatch(n, a)
	}
}

// nthChildSelector returns a selector that implements :nth-child(an+b).
// If last is true, implements :nth-last-child instead.
// If ofType is true, implements :nth-of-type instead.
func nthChildSelector(a, b int, last, ofType bool) Selector {
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}

		parent := n.Parent()
		if parent == nil {
			return false
		}

		if parent.Type() == DocumentNode {
			return false
		}

		i := -1
		count := 0
		for c, err := parent.FirstChild(); err == nil; c, err = c.NextSibling() {
			if (c.Type() != ElementNode) || (ofType && c.Name() != n.Name()) {
				continue
			}
			count++
			if c == n {
				i = count
				if !last {
					break
				}
			}
		}

		if i == -1 {
			// This shouldn't happen, since n should always be one of its parent's children.
			return false
		}

		if last {
			i = count - i + 1
		}

		i -= b
		if a == 0 {
			return i == 0
		}

		return i%a == 0 && i/a >= 0
	}
}

// simpleNthChildSelector returns a selector that implements :nth-child(b).
// If ofType is true, implements :nth-of-type instead.
func simpleNthChildSelector(b int, ofType bool) Selector {
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}

		parent := n.Parent()
		if parent == nil {
			return false
		}

		if parent.Type() == DocumentNode {
			return false
		}

		count := 0
		for c, err := parent.FirstChild(); err == nil; c, err = c.NextSibling() {
			if c.Type() != ElementNode || (ofType && c.Name() != n.Name()) {
				continue
			}
			count++
			if c == n {
				return count == b
			}
			if count >= b {
				return false
			}
		}
		return false
	}
}

// simpleNthLastChildSelector returns a selector that implements
// :nth-last-child(b). If ofType is true, implements :nth-last-of-type
// instead.
func simpleNthLastChildSelector(b int, ofType bool) Selector {
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}

		parent := n.Parent()
		if parent == nil {
			return false
		}

		if parent.Type() == DocumentNode {
			return false
		}

		count := 0
		for c, err := parent.LastChild(); err == nil; c, err = c.PreviousSibling() {
			if c.Type() != ElementNode || (ofType && c.Name() != n.Name()) {
				continue
			}
			count++
			if c == n {
				return count == b
			}
			if count >= b {
				return false
			}
		}
		return false
	}
}

// onlyChildSelector returns a selector that implements :only-child.
// If ofType is true, it implements :only-of-type instead.
func onlyChildSelector(ofType bool) Selector {
	return func(n *Node) bool {
		if n.Type() != ElementNode {
			return false
		}

		parent := n.Parent()
		if parent == nil {
			return false
		}

		if parent.Type() == DocumentNode {
			return false
		}

		count := 0
		for c, err := parent.FirstChild(); err == nil; c, err = c.NextSibling() {
			if (c.Type() != ElementNode) || (ofType && c.Name() != n.Name()) {
				continue
			}
			count++
			if count > 1 {
				return false
			}
		}

		return count == 1
	}
}

// inputSelector is a Selector that matches input, select, textarea and button elements.
func inputSelector(n *Node) bool {
	return n.Type() == ElementNode && (n.Name() == "input" || n.Name() == "select" || n.Name() == "textarea" || n.Name() == "button")
}

// emptyElementSelector is a Selector that matches empty elements.
func emptyElementSelector(n *Node) bool {
	if n.Type() != ElementNode {
		return false
	}

	for c, err := n.FirstChild(); err == nil; c, err = c.NextSibling() {
		switch c.Type() {
		case ElementNode, TextNode:
			return false
		}
	}

	return true
}

// descendantSelector returns a Selector that matches an element if
// it matches d and has an ancestor that matches a.
func descendantSelector(a, d Selector) Selector {
	return func(n *Node) bool {
		if !d(n) {
			return false
		}

		for p := n.Parent(); p != nil; p = p.Parent() {
			if a(p) {
				return true
			}
		}

		return false
	}
}

// childSelector returns a Selector that matches an element if
// it matches d and its parent matches a.
func childSelector(a, d Selector) Selector {
	return func(n *Node) bool {
		return d(n) && n.Parent() != nil && a(n.Parent())
	}
}

// siblingSelector returns a Selector that matches an element
// if it matches s2 and in is preceded by an element that matches s1.
// If adjacent is true, the sibling must be immediately before the element.
func siblingSelector(s1, s2 Selector, adjacent bool) Selector {
	return func(n *Node) bool {
		if !s2(n) {
			return false
		}

		if adjacent {
			for n, err := n.PreviousSibling(); err == nil; n, err = n.PreviousSibling() {
				if n.Type() == TextNode || n.Type() == CommentNode {
					continue
				}
				return s1(n)
			}
			return false
		}

		// Walk backwards looking for element that matches s1
		for c, err := n.PreviousSibling(); err == nil; c, err = c.PreviousSibling() {
			if s1(c) {
				return true
			}
		}

		return false
	}
}

// rootSelector implements :root
func rootSelector(n *Node) bool {
	if n.Type() != ElementNode {
		return false
	}
	if n.Parent() == nil {
		return false
	}
	return n.Parent().Type() == DocumentNode
}

//*******************************************************
// css parser
//*******************************************************

// a parser for CSS selectors
type parser struct {
	s string // the source text
	i int    // the current position
}

// parseEscape parses a backslash escape.
func (p *parser) parseEscape() (result string, err error) {
	if len(p.s) < p.i+2 || p.s[p.i] != '\\' {
		return "", errors.New("invalid escape sequence")
	}

	start := p.i + 1
	c := p.s[start]
	switch {
	case c == '\r' || c == '\n' || c == '\f':
		return "", errors.New("escaped line ending outside string")
	case hexDigit(c):
		// unicode escape (hex)
		var i int
		for i = start; i < p.i+6 && i < len(p.s) && hexDigit(p.s[i]); i++ {
			// empty
		}
		v, _ := strconv.ParseUint(p.s[start:i], 16, 21)
		if len(p.s) > i {
			switch p.s[i] {
			case '\r':
				i++
				if len(p.s) > i && p.s[i] == '\n' {
					i++
				}
			case ' ', '\t', '\n', '\f':
				i++
			}
		}
		p.i = i
		return string(rune(v)), nil
	}

	// Return the literal character after the backslash.
	result = p.s[start : start+1]
	p.i += 2
	return result, nil
}

func hexDigit(c byte) bool {
	return '0' <= c && c <= '9' || 'a' <= c && c <= 'f' || 'A' <= c && c <= 'F'
}

// nameStart returns whether c can be the first character of an identifier
// (not counting an initial hyphen, or an escape sequence).
func nameStart(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' || c > 127
}

// nameChar returns whether c can be a character within an identifier
// (not counting an escape sequence).
func nameChar(c byte) bool {
	return 'a' <= c && c <= 'z' || 'A' <= c && c <= 'Z' || c == '_' || c > 127 ||
		c == '-' || '0' <= c && c <= '9'
}

// parseIdentifier parses an identifier.
func (p *parser) parseIdentifier() (result string, err error) {
	startingDash := false
	if len(p.s) > p.i && p.s[p.i] == '-' {
		startingDash = true
		p.i++
	}

	if len(p.s) <= p.i {
		return "", errors.New("expected identifier, found EOF instead")
	}

	if c := p.s[p.i]; !(nameStart(c) || c == '\\') {
		return "", fmt.Errorf("expected identifier, found %c instead", c)
	}

	result, err = p.parseName()
	if startingDash && err == nil {
		result = "-" + result
	}
	return
}

// parseName parses a name (which is like an identifier, but doesn't have
// extra restrictions on the first character).
func (p *parser) parseName() (result string, err error) {
	i := p.i
loop:
	for i < len(p.s) {
		c := p.s[i]
		switch {
		case nameChar(c):
			start := i
			for i < len(p.s) && nameChar(p.s[i]) {
				i++
			}
			result += p.s[start:i]
		case c == '\\':
			p.i = i
			val, err := p.parseEscape()
			if err != nil {
				return "", err
			}
			i = p.i
			result += val
		default:
			break loop
		}
	}

	if result == "" {
		return "", errors.New("expected name, found EOF instead")
	}

	p.i = i
	return result, nil
}

// parseString parses a single- or double-quoted string.
func (p *parser) parseString() (result string, err error) {
	i := p.i
	if len(p.s) < i+2 {
		return "", errors.New("expected string, found EOF instead")
	}

	quote := p.s[i]
	i++

loop:
	for i < len(p.s) {
		switch p.s[i] {
		case '\\':
			if len(p.s) > i+1 {
				switch c := p.s[i+1]; c {
				case '\r':
					if len(p.s) > i+2 && p.s[i+2] == '\n' {
						i += 3
						continue loop
					}
					fallthrough
				case '\n', '\f':
					i += 2
					continue loop
				}
			}
			p.i = i
			val, err := p.parseEscape()
			if err != nil {
				return "", err
			}
			i = p.i
			result += val
		case quote:
			break loop
		case '\r', '\n', '\f':
			return "", errors.New("unexpected end of line in string")
		default:
			start := i
			for i < len(p.s) {
				if c := p.s[i]; c == quote || c == '\\' || c == '\r' || c == '\n' || c == '\f' {
					break
				}
				i++
			}
			result += p.s[start:i]
		}
	}

	if i >= len(p.s) {
		return "", errors.New("EOF in string")
	}

	// Consume the final quote.
	i++

	p.i = i
	return result, nil
}

// parseRegex parses a regular expression; the end is defined by encountering an
// unmatched closing ')' or ']' which is not consumed
func (p *parser) parseRegex() (rx *regexp.Regexp, err error) {
	i := p.i
	if len(p.s) < i+2 {
		return nil, errors.New("expected regular expression, found EOF instead")
	}

	// number of open parens or brackets;
	// when it becomes negative, finished parsing regex
	open := 0

loop:
	for i < len(p.s) {
		switch p.s[i] {
		case '(', '[':
			open++
		case ')', ']':
			open--
			if open < 0 {
				break loop
			}
		}
		i++
	}

	if i >= len(p.s) {
		return nil, errors.New("EOF in regular expression")
	}
	rx, err = regexp.Compile(p.s[p.i:i])
	p.i = i
	return rx, err
}

// skipWhitespace consumes whitespace characters and comments.
// It returns true if there was actually anything to skip.
func (p *parser) skipWhitespace() bool {
	i := p.i
	for i < len(p.s) {
		switch p.s[i] {
		case ' ', '\t', '\r', '\n', '\f':
			i++
			continue
		case '/':
			if strings.HasPrefix(p.s[i:], "/*") {
				end := strings.Index(p.s[i+len("/*"):], "*/")
				if end != -1 {
					i += end + len("/**/")
					continue
				}
			}
		}
		break
	}

	if i > p.i {
		p.i = i
		return true
	}

	return false
}

// consumeParenthesis consumes an opening parenthesis and any following
// whitespace. It returns true if there was actually a parenthesis to skip.
func (p *parser) consumeParenthesis() bool {
	if p.i < len(p.s) && p.s[p.i] == '(' {
		p.i++
		p.skipWhitespace()
		return true
	}
	return false
}

// consumeClosingParenthesis consumes a closing parenthesis and any preceding
// whitespace. It returns true if there was actually a parenthesis to skip.
func (p *parser) consumeClosingParenthesis() bool {
	i := p.i
	p.skipWhitespace()
	if p.i < len(p.s) && p.s[p.i] == ')' {
		p.i++
		return true
	}
	p.i = i
	return false
}

// parseTypeSelector parses a type selector (one that matches by tag name).
func (p *parser) parseTypeSelector() (result Selector, err error) {
	tag, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return typeSelector(tag), nil
}

// parseIDSelector parses a selector that matches by id attribute.
func (p *parser) parseIDSelector() (Selector, error) {
	if p.i >= len(p.s) {
		return nil, fmt.Errorf("expected id selector (#id), found EOF instead")
	}
	if p.s[p.i] != '#' {
		return nil, fmt.Errorf("expected id selector (#id), found '%c' instead", p.s[p.i])
	}

	p.i++
	id, err := p.parseName()
	if err != nil {
		return nil, err
	}

	return attributeEqualsSelector("id", id), nil
}

// parseClassSelector parses a selector that matches by class attribute.
func (p *parser) parseClassSelector() (Selector, error) {
	if p.i >= len(p.s) {
		return nil, fmt.Errorf("expected class selector (.class), found EOF instead")
	}
	if p.s[p.i] != '.' {
		return nil, fmt.Errorf("expected class selector (.class), found '%c' instead", p.s[p.i])
	}

	p.i++
	class, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	return attributeIncludesSelector("class", class), nil
}

// parseAttributeSelector parses a selector that matches by attribute value.
func (p *parser) parseAttributeSelector() (Selector, error) {
	if p.i >= len(p.s) {
		return nil, fmt.Errorf("expected attribute selector ([attribute]), found EOF instead")
	}
	if p.s[p.i] != '[' {
		return nil, fmt.Errorf("expected attribute selector ([attribute]), found '%c' instead", p.s[p.i])
	}

	p.i++
	p.skipWhitespace()
	key, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()
	if p.i >= len(p.s) {
		return nil, errors.New("unexpected EOF in attribute selector")
	}

	if p.s[p.i] == ']' {
		p.i++
		return attributeExistsSelector(key), nil
	}

	if p.i+2 >= len(p.s) {
		return nil, errors.New("unexpected EOF in attribute selector")
	}

	op := p.s[p.i : p.i+2]
	if op[0] == '=' {
		op = "="
	} else if op[1] != '=' {
		return nil, fmt.Errorf(`expected equality operator, found "%s" instead`, op)
	}
	p.i += len(op)

	p.skipWhitespace()
	if p.i >= len(p.s) {
		return nil, errors.New("unexpected EOF in attribute selector")
	}
	var val string
	var rx *regexp.Regexp
	if op == "#=" {
		rx, err = p.parseRegex()
	} else {
		switch p.s[p.i] {
		case '\'', '"':
			val, err = p.parseString()
		default:
			val, err = p.parseIdentifier()
		}
	}
	if err != nil {
		return nil, err
	}

	p.skipWhitespace()
	if p.i >= len(p.s) {
		return nil, errors.New("unexpected EOF in attribute selector")
	}
	if p.s[p.i] != ']' {
		return nil, fmt.Errorf("expected ']', found '%c' instead", p.s[p.i])
	}
	p.i++

	switch op {
	case "=":
		return attributeEqualsSelector(key, val), nil
	case "!=":
		return attributeNotEqualSelector(key, val), nil
	case "~=":
		return attributeIncludesSelector(key, val), nil
	case "|=":
		return attributeDashmatchSelector(key, val), nil
	case "^=":
		return attributePrefixSelector(key, val), nil
	case "$=":
		return attributeSuffixSelector(key, val), nil
	case "*=":
		return attributeSubstringSelector(key, val), nil
	case "#=":
		return attributeRegexSelector(key, rx), nil
	}

	return nil, fmt.Errorf("attribute operator %q is not supported", op)
}

var errExpectedParenthesis = errors.New("expected '(' but didn't find it")
var errExpectedClosingParenthesis = errors.New("expected ')' but didn't find it")
var errUnmatchedParenthesis = errors.New("unmatched '('")

// parsePseudoclassSelector parses a pseudoclass selector like :not(p).
func (p *parser) parsePseudoclassSelector() (Selector, error) {
	if p.i >= len(p.s) {
		return nil, fmt.Errorf("expected pseudoclass selector (:pseudoclass), found EOF instead")
	}
	if p.s[p.i] != ':' {
		return nil, fmt.Errorf("expected attribute selector (:pseudoclass), found '%c' instead", p.s[p.i])
	}

	p.i++
	name, err := p.parseIdentifier()
	if err != nil {
		return nil, err
	}
	name = toLowerASCII(name)

	switch name {
	case "not", "has", "haschild":
		if !p.consumeParenthesis() {
			return nil, errExpectedParenthesis
		}
		sel, parseErr := p.parseSelectorGroup()
		if parseErr != nil {
			return nil, parseErr
		}
		if !p.consumeClosingParenthesis() {
			return nil, errExpectedClosingParenthesis
		}

		switch name {
		case "not":
			return negatedSelector(sel), nil
		case "has":
			return hasDescendantSelector(sel), nil
		case "haschild":
			return hasChildSelector(sel), nil
		}

	case "contains", "containsown":
		if !p.consumeParenthesis() {
			return nil, errExpectedParenthesis
		}
		if p.i == len(p.s) {
			return nil, errUnmatchedParenthesis
		}
		var val string
		switch p.s[p.i] {
		case '\'', '"':
			val, err = p.parseString()
		default:
			val, err = p.parseIdentifier()
		}
		if err != nil {
			return nil, err
		}
		val = strings.ToLower(val)
		p.skipWhitespace()
		if p.i >= len(p.s) {
			return nil, errors.New("unexpected EOF in pseudo selector")
		}
		if !p.consumeClosingParenthesis() {
			return nil, errExpectedClosingParenthesis
		}

		switch name {
		case "contains":
			return textSubstrSelector(val), nil
		case "containsown":
			return ownTextSubstrSelector(val), nil
		}

	case "matches", "matchesown":
		if !p.consumeParenthesis() {
			return nil, errExpectedParenthesis
		}
		rx, err := p.parseRegex()
		if err != nil {
			return nil, err
		}
		if p.i >= len(p.s) {
			return nil, errors.New("unexpected EOF in pseudo selector")
		}
		if !p.consumeClosingParenthesis() {
			return nil, errExpectedClosingParenthesis
		}

		switch name {
		case "matches":
			return textRegexSelector(rx), nil
		case "matchesown":
			return ownTextRegexSelector(rx), nil
		}

	case "nth-child", "nth-last-child", "nth-of-type", "nth-last-of-type":
		if !p.consumeParenthesis() {
			return nil, errExpectedParenthesis
		}
		a, b, err := p.parseNth()
		if err != nil {
			return nil, err
		}
		if !p.consumeClosingParenthesis() {
			return nil, errExpectedClosingParenthesis
		}
		if a == 0 {
			switch name {
			case "nth-child":
				return simpleNthChildSelector(b, false), nil
			case "nth-of-type":
				return simpleNthChildSelector(b, true), nil
			case "nth-last-child":
				return simpleNthLastChildSelector(b, false), nil
			case "nth-last-of-type":
				return simpleNthLastChildSelector(b, true), nil
			}
		}
		return nthChildSelector(a, b,
				name == "nth-last-child" || name == "nth-last-of-type",
				name == "nth-of-type" || name == "nth-last-of-type"),
			nil

	case "first-child":
		return simpleNthChildSelector(1, false), nil
	case "last-child":
		return simpleNthLastChildSelector(1, false), nil
	case "first-of-type":
		return simpleNthChildSelector(1, true), nil
	case "last-of-type":
		return simpleNthLastChildSelector(1, true), nil
	case "only-child":
		return onlyChildSelector(false), nil
	case "only-of-type":
		return onlyChildSelector(true), nil
	case "input":
		return inputSelector, nil
	case "empty":
		return emptyElementSelector, nil
	case "root":
		return rootSelector, nil
	}

	return nil, fmt.Errorf("unknown pseudoclass :%s", name)
}

// parseInteger parses a  decimal integer.
func (p *parser) parseInteger() (int, error) {
	i := p.i
	start := i
	for i < len(p.s) && '0' <= p.s[i] && p.s[i] <= '9' {
		i++
	}
	if i == start {
		return 0, errors.New("expected integer, but didn't find it")
	}
	p.i = i

	val, err := strconv.Atoi(p.s[start:i])
	if err != nil {
		return 0, err
	}

	return val, nil
}

// parseNth parses the argument for :nth-child (normally of the form an+b).
func (p *parser) parseNth() (a, b int, err error) {
	// initial state
	if p.i >= len(p.s) {
		goto eof
	}
	switch p.s[p.i] {
	case '-':
		p.i++
		goto negativeA
	case '+':
		p.i++
		goto positiveA
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		goto positiveA
	case 'n', 'N':
		a = 1
		p.i++
		goto readN
	case 'o', 'O', 'e', 'E':
		id, nameErr := p.parseName()
		if nameErr != nil {
			return 0, 0, nameErr
		}
		id = toLowerASCII(id)
		if id == "odd" {
			return 2, 1, nil
		}
		if id == "even" {
			return 2, 0, nil
		}
		return 0, 0, fmt.Errorf("expected 'odd' or 'even', but found '%s' instead", id)
	default:
		goto invalid
	}

positiveA:
	if p.i >= len(p.s) {
		goto eof
	}
	switch p.s[p.i] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		a, err = p.parseInteger()
		if err != nil {
			return 0, 0, err
		}
		goto readA
	case 'n', 'N':
		a = 1
		p.i++
		goto readN
	default:
		goto invalid
	}

negativeA:
	if p.i >= len(p.s) {
		goto eof
	}
	switch p.s[p.i] {
	case '0', '1', '2', '3', '4', '5', '6', '7', '8', '9':
		a, err = p.parseInteger()
		if err != nil {
			return 0, 0, err
		}
		a = -a
		goto readA
	case 'n', 'N':
		a = -1
		p.i++
		goto readN
	default:
		goto invalid
	}

readA:
	if p.i >= len(p.s) {
		goto eof
	}
	switch p.s[p.i] {
	case 'n', 'N':
		p.i++
		goto readN
	default:
		// The number we read as a is actually b.
		return 0, a, nil
	}

readN:
	p.skipWhitespace()
	if p.i >= len(p.s) {
		goto eof
	}
	switch p.s[p.i] {
	case '+':
		p.i++
		p.skipWhitespace()
		b, err = p.parseInteger()
		if err != nil {
			return 0, 0, err
		}
		return a, b, nil
	case '-':
		p.i++
		p.skipWhitespace()
		b, err = p.parseInteger()
		if err != nil {
			return 0, 0, err
		}
		return a, -b, nil
	default:
		return a, 0, nil
	}

eof:
	return 0, 0, errors.New("unexpected EOF while attempting to parse expression of form an+b")

invalid:
	return 0, 0, errors.New("unexpected character while attempting to parse expression of form an+b")
}

// parseSimpleSelectorSequence parses a selector sequence that applies to
// a single element.
func (p *parser) parseSimpleSelectorSequence() (Selector, error) {
	var result Selector

	if p.i >= len(p.s) {
		return nil, errors.New("expected selector, found EOF instead")
	}

	switch p.s[p.i] {
	case '*':
		// It's the universal selector. Just skip over it, since it doesn't affect the meaning.
		p.i++
	case '#', '.', '[', ':':
		// There's no type selector. Wait to process the other till the main loop.
	default:
		r, err := p.parseTypeSelector()
		if err != nil {
			return nil, err
		}
		result = r
	}

loop:
	for p.i < len(p.s) {
		var ns Selector
		var err error
		switch p.s[p.i] {
		case '#':
			ns, err = p.parseIDSelector()
		case '.':
			ns, err = p.parseClassSelector()
		case '[':
			ns, err = p.parseAttributeSelector()
		case ':':
			ns, err = p.parsePseudoclassSelector()
		default:
			break loop
		}
		if err != nil {
			return nil, err
		}
		if result == nil {
			result = ns
		} else {
			result = intersectionSelector(result, ns)
		}
	}

	if result == nil {
		result = func(n *Node) bool {
			return n.Type() == ElementNode
		}
	}

	return result, nil
}

// parseSelector parses a selector that may include combinators.
func (p *parser) parseSelector() (result Selector, err error) {
	p.skipWhitespace()
	result, err = p.parseSimpleSelectorSequence()
	if err != nil {
		return
	}

	for {
		var combinator byte
		if p.skipWhitespace() {
			combinator = ' '
		}
		if p.i >= len(p.s) {
			return
		}

		switch p.s[p.i] {
		case '+', '>', '~':
			combinator = p.s[p.i]
			p.i++
			p.skipWhitespace()
		case ',', ')':
			// These characters can't begin a selector, but they can legally occur after one.
			return
		}

		if combinator == 0 {
			return
		}

		c, err := p.parseSimpleSelectorSequence()
		if err != nil {
			return nil, err
		}

		switch combinator {
		case ' ':
			result = descendantSelector(result, c)
		case '>':
			result = childSelector(result, c)
		case '+':
			result = siblingSelector(result, c, true)
		case '~':
			result = siblingSelector(result, c, false)
		}
	}

	panic("unreachable")
}

// parseSelectorGroup parses a group of selectors, separated by commas.
func (p *parser) parseSelectorGroup() (result Selector, err error) {
	result, err = p.parseSelector()
	if err != nil {
		return
	}

	for p.i < len(p.s) {
		if p.s[p.i] != ',' {
			return result, nil
		}
		p.i++
		c, err := p.parseSelector()
		if err != nil {
			return nil, err
		}
		result = unionSelector(result, c)
	}

	return
}
