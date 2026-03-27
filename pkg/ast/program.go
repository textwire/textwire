package ast

import (
	"strings"

	"github.com/textwire/textwire/v4/pkg/fail"
	"github.com/textwire/textwire/v4/pkg/token"
)

type Program struct {
	BaseNode
	IsLayout   bool
	Name       string
	AbsPath    string
	UseDir     *UseDir
	Chunks     []Chunk
	Components []*CompDir
	Reserves   map[string]*ReserveDir
	Inserts    map[string]*InsertDir
}

func NewProgram(tok token.Token) *Program {
	return &Program{
		BaseNode: NewBaseNode(tok),
	}
}

func (p *Program) chunkNode() {}

func (p *Program) String() string {
	var out strings.Builder
	out.Grow(len(p.Chunks))

	for i := range p.Chunks {
		out.WriteString(p.Chunks[i].String())
	}

	return out.String()
}

func (p *Program) AllChunks() []Chunk {
	chunks := make([]Chunk, 0)
	if p.Chunks == nil {
		return []Chunk{}
	}

	for _, chunk := range p.Chunks {
		if chunk == nil {
			continue
		}

		if s, ok := chunk.(NodeWithChunks); ok {
			chunks = append(chunks, s.(Chunk))
			chunks = append(chunks, s.AllChunks()...)
		}
	}

	return chunks
}

// LinkLayoutToUse adds Layout AST program to UseStmt for the current template
// and resets chunks to UseDir chunk only. Because we don't need anything else
// inside a template. Make sure inserts are added before this is called
// because they will be removed by this function.
func (p *Program) LinkLayoutToUse(layoutProg *Program) {
	p.UseDir.LayoutProg = layoutProg
	p.Chunks = []Chunk{p.UseDir}
}

func (p *Program) LinkCompProg(compName string, prog *Program, absPath string) *fail.Error {
	for _, comp := range p.Components {
		if comp.Name.Val != compName {
			continue
		}

		duplicate, times := findDuplicatePasses(comp.Passes)
		if times > 0 && duplicate != nil {
			return fail.New(
				duplicate.Pos(),
				absPath,
				fail.OriginLink,
				fail.ErrDuplicatePass,
				duplicate.Name.Val,
				times,
				compName,
			)
		}

		for _, passDir := range comp.Passes {
			name := passDir.Name.Val
			idx := findSlotIndex(prog.Chunks, name)
			if idx != -1 {
				prog.Chunks[idx].(*SlotDir).Block = passDir.Block
				continue
			}

			if passDir.Name.Val == "" {
				return fail.New(
					passDir.Pos(),
					absPath,
					fail.OriginLink,
					fail.ErrDefaultSlotNotDefined,
					compName,
					compName,
				)
			}

			return fail.New(
				passDir.Pos(),
				absPath,
				fail.OriginLink,
				fail.ErrSlotNotDefined,
				compName,
				name,
			)
		}

		comp.CompProg = prog
	}

	return nil
}

func (p *Program) HasUseDir() bool {
	return p.UseDir != nil
}
