package plist

import (
	"bytes"
	"testing"
)

/*
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0', // Magic

		// Object Table
		// Offset Table

		// Trailer
		0x00, 0x00, 0x00, 0x00, 0x00, //  - U8[5] Unused
		0x01,                      //  - U8    Sort Version
		0x01,                      //  - U8    Offset Table Entry Size (#bytes)
		0x01,                      //  - U8    Object Reference Size (#bytes)
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //  - U64   # Objects
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //  - U64   Top Object
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, //  - U64   Offset Table Offset
	},
*/

var InvalidBplists = [][]byte{
	// Too short
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',
		0x00,
	},
	// Bad magic
	[]byte{
		'x', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Bad version
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '3', '0',

		0x00,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Bad version II
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '@', 'A',

		0x00,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Offset table inside trailer
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A,
	},
	// Offset table inside header
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
	},
	// Offset table off end of file
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF, 0x00,
	},
	// Garbage between offset table and trailer
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00,
		0x09,

		0xAB, 0xCD,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A,
	},
	// Top Object out of range
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Object out of range
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00,
		0xFF,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Object references too small (1 byte, but 257 objects)
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00,

		// 257 bytes worth of object table
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Offset references too small (1 byte, but 257 bytes worth of objects)
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		// 257 bytes worth of "objects"
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00,

		0x00,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01, 0x09,
	},
	// Too many objects
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x00,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0xFF,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// String way too long
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x5F, 0x10, 0xFF,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0B,
	},
	// UTF-16 String way too long
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x6F, 0x10, 0xFF,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0B,
	},
	// Data way too long
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x4F, 0x10, 0xFF,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0B,
	},
	// Array way too long
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xAF, 0x10, 0xFF,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0B,
	},
	// Dictionary way too long
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xDF, 0x10, 0xFF,
		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0B,
	},
	// Array self-referential
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xA1, 0x00,

		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A,
	},
	// Dictionary self-referential key
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xD1, 0x00, 0x01,
		0x50, // 0-byte string

		0x08, 0x0B,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C,
	},
	// Dictionary self-referential value
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xD1, 0x01, 0x00,
		0x50, // 0-byte string

		0x08, 0x0B,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C,
	},
	// Dictionary non-string key
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xD1, 0x01, 0x02,
		0x08,
		0x09,

		0x08, 0x0B, 0x0C,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x03,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0D,
	},
	// Array contains invalid reference
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xA1, 0x0F,

		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0A,
	},
	// Dictionary contains invalid reference
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xD1, 0x01, 0x0F,
		0x50, // 0-byte string

		0x08, 0x0B,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C,
	},
	// Invalid float ("7-byte")
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x27,

		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Invalid integer (8^5)
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0x15,

		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},
	// Invalid atom
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xFF,

		0x08,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x09,
	},

	// array refers to self through a second level
	[]byte{
		'b', 'p', 'l', 'i', 's', 't', '0', '0',

		0xA1, 0x01,
		0xA1, 0x00,

		0x08, 0x0A,

		0x00, 0x00, 0x00, 0x00, 0x00,
		0x01,
		0x01,
		0x01,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x02,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x0C,
	},
}

func TestInvalidBinaryPlists(t *testing.T) {
	for _, data := range InvalidBplists {
		buf := bytes.NewReader(data)
		d := newBplistParser(buf)
		_, err := d.parseDocument()
		if err == nil {
			t.Fatal("invalid plist failed to throw error")
		} else {
			t.Log(err)
		}
	}
}