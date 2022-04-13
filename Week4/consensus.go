// Copyright 2017 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

// Package consensus implements different Ethereum consensus engines.
package consensus

import (
	"math/big"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/state"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
)

// ChainHeaderReader mendefinisikan kumpulan kecil metode yang diperlukan untuk mengakses lokal
// blockchain selama verifikasi header.
type ChainHeaderReader interface {
	// Config mengambil konfigurasi rantai blockchain.
	Config() *params.ChainConfig

	// CurrentHeader mengambil header saat ini dari rantai lokal.
	CurrentHeader() *types.Header

	// GetHeader mengambil header blok dari database dengan hash dan nomor.
	GetHeader(hash common.Hash, number uint64) *types.Header

	// GetHeaderByNumber mengambil header blok dari database dengan nomor.
	GetHeaderByNumber(number uint64) *types.Header

	// GetHeaderByHash mengambil header blok dari database dengan hash-nya.
	GetHeaderByHash(hash common.Hash) *types.Header

	// GetTd mengambil total kesulitan dari database dengan hash dan nomor.
	GetTd(hash common.Hash, number uint64) *big.Int
}

// ChainReader mendefinisikan kumpulan kecil metode yang diperlukan untuk mengakses lokal
// blockchain selama verifikasi header
type ChainReader interface {
	ChainHeaderReader

	// GetBlock mengambil blok dari database dengan hash dan nomor.
	GetBlock(hash common.Hash, number uint64) *types.Block
}

// Engine adalah mesin konsensus agnostik algoritma.
type Engine interface {
	// Author retrieves the Ethereum address of the account that minted the given
	// blok, yang mungkin berbeda dari basis koin header jika konsensus.
	// engine didasarkan pada signatura.
	Author(header *types.Header) (common.Address, error)

	// VerifyHeader checks whether a header conforms to the consensus rules of a
	// engine yang diberikan. Memverifikasi segel dapat dilakukan secara opsional di sini, atau secara eksplisit
	// via the VerifySeal method.
	VerifyHeader(chain ChainHeaderReader, header *types.Header, seal bool) error

	// VerifyHeaders is similar to VerifyHeader, but verifies a batch of headers
	// concurrently. The method returns a quit channel to abort the operations and
	// a results channel to retrieve the async verifications (the order is that of
	// the input slice).
	// verify headers akan sama dengan metoda verify header, namun verifikasi header dalam batch secara bersamaan
	VerifyHeaders(chain ChainHeaderReader, headers []*types.Header, seals []bool) (chan<- struct{}, <-chan error)

	// VerifyUncles verifies that the given block's uncles conform to the consensus
	// rules of a given engine.
	// verify uncles akan mengecek apakah block uncles sesuai dengan aturan consensie dari engine tertentu
	VerifyUncles(chain ChainReader, block *types.Block) error

	// Prepare initializes the consensus fields of a block header according to the
	// rules of a particular engine. The changes are executed inline.
	// Untuk mengikuti aturan consencus, state database dan header blok dapat diperbarui yang akan terjadi secara endline.
	Prepare(chain ChainHeaderReader, header *types.Header) error

	// Finalize runs any post-transaction state modifications (e.g. block rewards)
	// but does not assemble the block.
	//
	// Note: The block header and state database might be updated to reflect any
	// consensus rules that happen at finalization (e.g. block rewards).
	Finalize(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header)

	// FinalizeAndAssemble runs any post-transaction state modifications (e.g. block
	// rewards) and assembles the final block.
	//
	// Catatan: Header blok dan basis data status mungkin diperbarui untuk mencerminkan apa pun
	// consensus rules that happen at finalization (e.g. block rewards).
	FinalizeAndAssemble(chain ChainHeaderReader, header *types.Header, state *state.StateDB, txs []*types.Transaction,
		uncles []*types.Header, receipts []*types.Receipt) (*types.Block, error)

	// Seal generates a new sealing request for the given input block and pushes
	// the result into the given channel.
	//
	// Note, the method returns immediately and will send the result async. More
	// than one result may also be returned depending on the consensus algorithm.
	Seal(chain ChainHeaderReader, block *types.Block, results chan<- *types.Block, stop <-chan struct{}) error

	// SealHash returns the hash of a block prior to it being sealed.
	SealHash(header *types.Header) common.Hash

	// CalcDifficulty is the difficulty adjustment algorithm. It returns the difficulty
	// that a new block should have.
	CalcDifficulty(chain ChainHeaderReader, time uint64, parent *types.Header) *big.Int

	// APIs returns the RPC APIs this consensus engine provides.
	APIs(chain ChainHeaderReader) []rpc.API

	// Close terminates any background threads maintained by the consensus engine.
	Close() error
}

// PoW is a consensus engine based on proof-of-work.
type PoW interface {
	Engine

	// Hashrate returns the current mining hashrate of a PoW consensus engine.
	Hashrate() float64
}