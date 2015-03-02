package merkletree

import (
	"hash"
	"io"
)

// ReaderRoot returns the Merkle root of the data read from the reader, where
// each leaf is 'segmentSize' long and 'h' is used as the hashing function. All
// leaves will be 'segmentSize' bytes, the last leaf may have extra zeros.
func ReaderRoot(r io.Reader, h hash.Hash, segmentSize int) (root []byte, err error) {
	tree := New(h)
	for {
		segment := make([]byte, segmentSize)
		_, readErr := io.ReadFull(r, segment)
		if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
			err = readErr
			return
		}
		if readErr == io.EOF {
			break
		}
		tree.Push(segment)
	}
	root = tree.Root()
	return
}

// BuildReaderProof returns a proof that certain data is in the merkle tree
// created by the data in the reader. The merkle root, set of proofs, and the
// number of leaves in the Merkle tree are all returned.
func BuildReaderProof(r io.Reader, h hash.Hash, segmentSize int, index uint64) (root []byte, proofSet [][]byte, numLeaves uint64, err error) {
	tree := New(h)
	tree.SetIndex(index)
	for {
		segment := make([]byte, segmentSize)
		_, readErr := io.ReadFull(r, segment)
		if readErr != nil && readErr != io.EOF && readErr != io.ErrUnexpectedEOF {
			err = readErr
			return
		}
		if readErr == io.EOF {
			break
		}
		tree.Push(segment)
	}
	root, proofSet, _, numLeaves = tree.Prove()
	return
}
