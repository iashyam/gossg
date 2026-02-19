package main

import (
	"fmt"
	"strings"
)

type Parser struct {
	tokens []Token
	pos    int
}

func NewParser(tokens []Token) *Parser {
	return &Parser{tokens: tokens, pos: 0}
}

func (p *Parser) Parse() string {
	var sb strings.Builder

	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]
		switch token.Type {
		case HEADING:
			p.parseHeading(&sb)
		case LIST:
			p.parseList(&sb)
		case QUOTE:
			p.parseQuote(&sb)
		case BLANKLINE:
			p.pos++ // Skip blanklines between blocks
		case EOF:
			p.pos++
		default:
			p.parseParagraph(&sb)
		}
	}
	return sb.String()
}

func (p *Parser) parseHeading(sb *strings.Builder) {
	token := p.tokens[p.pos]
	level := strings.Count(token.value, "#")
	sb.WriteString(fmt.Sprintf("<h%d>", level))
	p.pos++

	p.parseInline(sb)

	sb.WriteString(fmt.Sprintf("</h%d>\n", level))
}

func (p *Parser) parseList(sb *strings.Builder) {
	sb.WriteString("<ul>\n")
	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == LIST {
		p.pos++ // consume list marker
		sb.WriteString("<li>")
		p.parseInline(sb)
		sb.WriteString("</li>\n")
	}
	sb.WriteString("</ul>\n")
}

func (p *Parser) parseQuote(sb *strings.Builder) {
	sb.WriteString("<blockquote>\n<p>")
	// Quote might span multiple lines with > markers or just be a block?
	// Lexer emits QUOTE for '> '.
	// We can treat consecutive QUOTE lines as one blockquote.

	for p.pos < len(p.tokens) && p.tokens[p.pos].Type == QUOTE {
		p.pos++ // consume >
		p.parseInline(sb)
		// If parseInline consumed NEWLINE, next might be QUOTE again.
		// We insert a space if joining lines?
		// parseInline converts NEWLINE to space.
	}
	sb.WriteString("</p>\n</blockquote>\n")
}

func (p *Parser) parseParagraph(sb *strings.Builder) {
	sb.WriteString("<p>")
	p.parseInline(sb)
	sb.WriteString("</p>\n")
}

// parseInline consumes tokens until a block separator (BLANKLINE) or block starter (HEADING, LIST, QUOTE) is found.
// It effectively handles "inline" content which can span multiple lines (paragraphs).
func (p *Parser) parseInline(sb *strings.Builder) {
	p.parseUntil(sb, EOF)
}

// parseUntil consumes tokens until endTokenType is found (if not EOF).
// If endTokenType is EOF, it runs until block boundaries.
func (p *Parser) parseUntil(sb *strings.Builder, endTokenType TokenType) {
	for p.pos < len(p.tokens) {
		token := p.tokens[p.pos]

		if endTokenType != EOF && token.Type == endTokenType {
			// Found the closer
			return
		}

		if token.Type == BLANKLINE || token.Type == EOF {
			return // End of block
		}

		// Check for block starters only if we are at start of line (implied by previous NEWLINE consumed?
		// Actually Lexer only emits HEADING/LIST/QUOTE at start of line.
		// So if we see them here, it means we are at start of line.
		if token.Type == HEADING || token.Type == LIST || token.Type == QUOTE {
			return
		}

		if token.Type == NEWLINE {
			sb.WriteString(" ")
			p.pos++
			continue
		}

		switch token.Type {
		case TEXT:
			sb.WriteString(token.value)
			p.pos++
		case BOLD:
			sb.WriteString("<strong>")
			p.pos++
			p.parseUntil(sb, BOLD)
			sb.WriteString("</strong>")
			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == BOLD {
				p.pos++
			}
		case ITALIC:
			sb.WriteString("<em>")
			p.pos++
			p.parseUntil(sb, ITALIC)
			sb.WriteString("</em>")
			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == ITALIC {
				p.pos++
			}
		case INLINE_CODE:
			sb.WriteString("<code>")
			p.pos++
			p.parseUntil(sb, INLINE_CODE)
			sb.WriteString("</code>")
			if p.pos < len(p.tokens) && p.tokens[p.pos].Type == INLINE_CODE {
				p.pos++
			}
		case MATH:
			// Strip $ if present in value? Lexer keeps them now?
			// Lexer MathHandler returns content between $.
			// Wait, prior fix: "return Token{Type: MATH, value: l.input[start : l.pos-1]}"
			// This strips the $.
			sb.WriteString(fmt.Sprintf("<math>%s</math>", token.value))
			p.pos++
		default:
			// Should not happen for handled types.
			p.pos++
		}
	}
}
