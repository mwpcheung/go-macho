package types

import (
	"encoding/binary"
	"fmt"
	"strings"
)

// An Nlist is a Mach-O generic symbol table entry.
type Nlist struct {
	Name uint32
	Type NType
	Sect uint8
	Desc NDescType
}

// An Nlist32 is a Mach-O 32-bit symbol table entry.
type Nlist32 struct {
	Nlist
	Value uint32
}

func (n *Nlist32) Put32(b []byte, o binary.ByteOrder) uint32 {
	o.PutUint32(b[0:], n.Name)
	b[4] = byte(n.Type)
	b[5] = byte(n.Sect)
	o.PutUint16(b[6:], uint16(n.Desc))
	o.PutUint32(b[8:], uint32(n.Value))
	return 8 + 4
}

// An Nlist64 is a Mach-O 64-bit symbol table entry.
type Nlist64 struct {
	Nlist
	Value uint64
}

func (n *Nlist64) Put64(b []byte, o binary.ByteOrder) uint32 {
	o.PutUint32(b[0:], n.Name)
	b[4] = byte(n.Type)
	b[5] = byte(n.Sect)
	o.PutUint16(b[6:], uint16(n.Desc))
	o.PutUint64(b[8:], n.Value)
	return 8 + 8
}

type NType uint8

/*
 * The n_type field really contains four fields:
 *	unsigned char N_STAB:3,
 *		      N_PEXT:1,
 *		      N_TYPE:3,
 *		      N_EXT:1;
 * which are used via the following masks.
 */
const (
	N_STAB NType = 0xe0 /* if any of these bits set, a symbolic debugging entry */
	N_PEXT NType = 0x10 /* private external symbol bit */
	N_TYPE NType = 0x0e /* mask for the type bits */
	N_EXT  NType = 0x01 /* external symbol bit, set for external symbols */
)

/*
 * Values for N_TYPE bits of the n_type field.
 */
const (
	N_UNDF NType = 0x0 /* undefined, n_sect == NO_SECT */
	N_ABS  NType = 0x2 /* absolute, n_sect == NO_SECT */
	N_SECT NType = 0xe /* defined in section number n_sect */
	N_PBUD NType = 0xc /* prebound undefined (defined in a dylib) */
	N_INDR NType = 0xa /* indirect */
)

func (t NType) IsDebugSym() bool {
	return (t & N_STAB) != 0
}

func (t NType) IsPrivateExternalSym() bool {
	return (t & N_PEXT) != 0
}

func (t NType) IsExternalSym() bool {
	return (t & N_EXT) != 0
}

func (t NType) IsUndefinedSym() bool {
	return (t & N_TYPE) == N_UNDF
}
func (t NType) IsAbsoluteSym() bool {
	return (t & N_TYPE) == N_ABS
}
func (t NType) IsDefinedInSection() bool {
	return (t & N_TYPE) == N_SECT
}
func (t NType) IsPreboundUndefinedSym() bool {
	return (t & N_TYPE) == N_PBUD
}
func (t NType) IsIndirectSym() bool {
	return (t & N_TYPE) == N_INDR
}

func (t NType) String(secName string) string {
	var tStr string
	if t.IsDebugSym() {
		tStr += "debug|"
	}
	if t.IsPrivateExternalSym() {
		tStr += "priv_ext|"
	}
	if t.IsExternalSym() {
		tStr += "ext|"
	}
	if t.IsUndefinedSym() {
		tStr += "undef|"
	}
	if t.IsAbsoluteSym() {
		tStr += "abs|"
	}
	if t.IsDefinedInSection() {
		tStr += fmt.Sprintf("%s|", secName)
	}
	if t.IsPreboundUndefinedSym() {
		tStr += "prebound_undef|"
	}
	if t.IsIndirectSym() {
		tStr += "indir|"
	}
	return strings.TrimSuffix(tStr, "|")
}

type NDescType uint16

func (d NDescType) GetCommAlign() NDescType {
	return (d >> 8) & 0x0f
}

const REFERENCE_TYPE NDescType = 0x7

const (
	/* types of references */
	REFERENCE_FLAG_UNDEFINED_NON_LAZY         NDescType = 0
	REFERENCE_FLAG_UNDEFINED_LAZY             NDescType = 1
	REFERENCE_FLAG_DEFINED                    NDescType = 2
	REFERENCE_FLAG_PRIVATE_DEFINED            NDescType = 3
	REFERENCE_FLAG_PRIVATE_UNDEFINED_NON_LAZY NDescType = 4
	REFERENCE_FLAG_PRIVATE_UNDEFINED_LAZY     NDescType = 5
)

func (d NDescType) IsUndefinedNonLazy() bool {
	return (d & REFERENCE_TYPE) == REFERENCE_FLAG_UNDEFINED_NON_LAZY
}
func (d NDescType) IsUndefinedLazy() bool {
	return (d & REFERENCE_TYPE) == REFERENCE_FLAG_UNDEFINED_LAZY
}
func (d NDescType) IsDefined() bool {
	return (d & REFERENCE_TYPE) == REFERENCE_FLAG_DEFINED
}
func (d NDescType) IsPrivateDefined() bool {
	return (d & REFERENCE_TYPE) == REFERENCE_FLAG_PRIVATE_DEFINED
}
func (d NDescType) IsPrivateUndefinedNonLazy() bool {
	return (d & REFERENCE_TYPE) == REFERENCE_FLAG_PRIVATE_UNDEFINED_NON_LAZY
}
func (d NDescType) IsPrivateUndefinedLazy() bool {
	return (d & REFERENCE_TYPE) == REFERENCE_FLAG_PRIVATE_UNDEFINED_LAZY
}
func (d NDescType) GetLibraryOrdinal() NDescType {
	return (d >> 8) & 0xff
}

func (t NDescType) String() string {
	var tStr string
	if t.IsUndefinedNonLazy() {
		tStr += "undef_nonlazy|"
	}
	if t.IsUndefinedLazy() {
		tStr += "undef_lazy|"
	}
	if t.IsDefined() {
		tStr += "def|"
	}
	if t.IsPrivateDefined() {
		tStr += "priv_def|"
	}
	if t.IsPrivateUndefinedNonLazy() {
		tStr += "pri_undef_nonlazy|"
	}
	if t.IsPrivateUndefinedLazy() {
		tStr += "priv_undef_lazy|"
	}
	// tStr += fmt.Sprintf("libord=%d", t.GetLibraryOrdinal())
	return strings.TrimSuffix(tStr, "|")
}

const (
	SELF_LIBRARY_ORDINAL   NDescType = 0x0
	MAX_LIBRARY_ORDINAL    NDescType = 0xfd
	DYNAMIC_LOOKUP_ORDINAL NDescType = 0xfe
	EXECUTABLE_ORDINAL     NDescType = 0xff
)

// TODO: add these flags to the NDescType String output

const (
	/*
	 * The N_NO_DEAD_STRIP bit of the n_desc field only ever appears in a
	 * relocatable .o file (MH_OBJECT filetype). And is used to indicate to the
	 * static link editor it is never to dead strip the symbol.
	 */
	NO_DEAD_STRIP NDescType = 0x0020 /* symbol is not to be dead stripped */

	/*
	 * The N_DESC_DISCARDED bit of the n_desc field never appears in linked image.
	 * But is used in very rare cases by the dynamic link editor to mark an in
	 * memory symbol as discared and longer used for linking.
	 */
	DESC_DISCARDED NDescType = 0x0020 /* symbol is discarded */

	/*
	 * The N_WEAK_REF bit of the n_desc field indicates to the dynamic linker that
	 * the undefined symbol is allowed to be missing and is to have the address of
	 * zero when missing.
	 */
	WEAK_REF NDescType = 0x0040 /* symbol is weak referenced */

	/*
	 * The N_WEAK_DEF bit of the n_desc field indicates to the static and dynamic
	 * linkers that the symbol definition is weak, allowing a non-weak symbol to
	 * also be used which causes the weak definition to be discared.  Currently this
	 * is only supported for symbols in coalesed sections.
	 */
	WEAK_DEF NDescType = 0x0080 /* coalesed symbol is a weak definition */

	/*
	 * The N_REF_TO_WEAK bit of the n_desc field indicates to the dynamic linker
	 * that the undefined symbol should be resolved using flat namespace searching.
	 */
	REF_TO_WEAK NDescType = 0x0080 /* reference to a weak symbol */

	/*
	 * The N_ARM_THUMB_DEF bit of the n_desc field indicates that the symbol is
	 * a defintion of a Thumb function.
	 */
	ARM_THUMB_DEF NDescType = 0x0008 /* symbol is a Thumb function (ARM) */

	/*
	 * The N_SYMBOL_RESOLVER bit of the n_desc field indicates that the
	 * that the function is actually a resolver function and should
	 * be called to get the address of the real function to use.
	 * This bit is only available in .o files (MH_OBJECT filetype)
	 */
	SYMBOL_RESOLVER NDescType = 0x0100

	/*
	 * The N_ALT_ENTRY bit of the n_desc field indicates that the
	 * symbol is pinned to the previous content.
	 */
	ALT_ENTRY NDescType = 0x0200

	/*
	 * The N_COLD_FUNC bit of the n_desc field indicates that the symbol is used
	 * infrequently and the linker should order it towards the end of the section.
	 */
	N_COLD_FUNC NDescType = 0x0400
)

/*
 * Symbolic debugger symbols.
 */
const (
	N_GSYM  = 0x20 /* global symbol: name,,NO_SECT,type,0 */
	N_FNAME = 0x22 /* procedure name (f77 kludge): name,,NO_SECT,0,0 */
	N_FUN   = 0x24 /* procedure: name,,n_sect,linenumber,address */
	N_STSYM = 0x26 /* static symbol: name,,n_sect,type,address */
	N_LCSYM = 0x28 /* .lcomm symbol: name,,n_sect,type,address */
	N_BNSYM = 0x2e /* begin nsect sym: 0,,n_sect,0,address */
	N_AST   = 0x32 /* AST file path: name,,NO_SECT,0,0 */
	N_OPT   = 0x3c /* emitted with gcc2_compiled and in gcc source */
	N_RSYM  = 0x40 /* register sym: name,,NO_SECT,type,register */
	N_SLINE = 0x44 /* src line: 0,,n_sect,linenumber,address */
	N_ENSYM = 0x4e /* end nsect sym: 0,,n_sect,0,address */
	N_SSYM  = 0x60 /* structure elt: name,,NO_SECT,type,struct_offset */
	N_SO    = 0x64 /* source file name: name,,n_sect,0,address */
	N_OSO   = 0x66 /* object file name: name,,(see below),0,st_mtime */
	/*   historically N_OSO set n_sect to 0. The N_OSO
	 *   n_sect may instead hold the low byte of the
	 *   cpusubtype value from the Mach-O header. */
	N_LSYM    = 0x80 /* local sym: name,,NO_SECT,type,offset */
	N_BINCL   = 0x82 /* include file beginning: name,,NO_SECT,0,sum */
	N_SOL     = 0x84 /* #included file name: name,,n_sect,0,address */
	N_PARAMS  = 0x86 /* compiler parameters: name,,NO_SECT,0,0 */
	N_VERSION = 0x88 /* compiler version: name,,NO_SECT,0,0 */
	N_OLEVEL  = 0x8A /* compiler -O level: name,,NO_SECT,0,0 */
	N_PSYM    = 0xa0 /* parameter: name,,NO_SECT,type,offset */
	N_EINCL   = 0xa2 /* include file end: name,,NO_SECT,0,0 */
	N_ENTRY   = 0xa4 /* alternate entry: name,,n_sect,linenumber,address */
	N_LBRAC   = 0xc0 /* left bracket: 0,,NO_SECT,nesting level,address */
	N_EXCL    = 0xc2 /* deleted include file: name,,NO_SECT,0,sum */
	N_RBRAC   = 0xe0 /* right bracket: 0,,NO_SECT,nesting level,address */
	N_BCOMM   = 0xe2 /* begin common: name,,NO_SECT,0,0 */
	N_ECOMM   = 0xe4 /* end common: name,,n_sect,0,0 */
	N_ECOML   = 0xe8 /* end common (local name): 0,,n_sect,0,address */
	N_LENG    = 0xfe /* second stab entry with length information */
	/*
	 * for the berkeley pascal compiler, pc(1):
	 */
	N_PC = 0x30 /* global pascal symbol: name,,NO_SECT,subtype,line */
)
