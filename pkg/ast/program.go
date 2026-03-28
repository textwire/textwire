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
	Chunks     []Chunk
	Components []*CompDir
	Reserves   map[string]*ReserveDir
	Inserts    map[string]*InsertDir
	Slots      map[string]*SlotDir

	// UseDir is used to reference the use directive in the program.
	// We need it because the final program object must have a field UseDir.
	// After parsing a program we link this pointer to program.UseDir.
	UseDir *UseDir
}

func NewProgram(tok token.Token) *Program {
	return &Program{
		BaseNode:   NewBaseNode(tok),
		Components: []*CompDir{},
		Inserts:    map[string]*InsertDir{},
		Reserves:   map[string]*ReserveDir{},
		Slots:      map[string]*SlotDir{},
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

func (p *Program) HasUseDir() bool {
	return p.UseDir != nil
}

func (p *Program) LinkPassBlocksToSlots(compDir *CompDir, compFileProg *Program) *fail.Error {
	if err := compFileProg.linkPassBlocksToSlots(compDir, p.AbsPath); err != nil {
		return err
	}

	compDir.CompProg = compFileProg
	return nil
}

func (p *Program) linkPassBlocksToSlots(compDir *CompDir, compFileAbsPath string) *fail.Error {
	if compDir.DefaultPass != nil {
		if err := p.linkBlockToDefaultSlot(compDir, compFileAbsPath); err != nil {
			return err
		}
	}

	for _, passDir := range compDir.Passes {
		if err := p.linkBlockToSlot(passDir, compDir.Name.Val, compFileAbsPath); err != nil {
			return err
		}
	}

	return nil
}

func (p *Program) linkBlockToDefaultSlot(compDir *CompDir, compFileAbsPath string) *fail.Error {
	slotDir, ok := p.Slots[""]
	if ok {
		slotDir.Block = compDir.DefaultPass.Block
		return nil
	}

	return fail.New(
		compDir.DefaultPass.Pos(),
		compFileAbsPath,
		fail.OriginLink,
		fail.ErrDefaultSlotNotDefined,
		compDir.Name.Val,
		compDir.Name.Val,
	)
}

func (p *Program) linkBlockToSlot(
	passDir *PassDir,
	compName string,
	compFileAbsPath string,
) *fail.Error {
	name := passDir.Name.Val
	slotDir, ok := p.Slots[name]
	if ok {
		slotDir.Block = passDir.Block
		return nil
	}

	return fail.New(
		passDir.Pos(),
		compFileAbsPath,
		fail.OriginLink,
		fail.ErrSlotNotDefined,
		compName,
		name,
	)
}
