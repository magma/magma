/*
Copyright (c) Facebook, Inc. and its affiliates.
All rights reserved.

This source code is licensed under the BSD-style license found in the
LICENSE file in the root directory of this source tree.
*/

package utils

import "fmt"

// LTEBand struct for converting EARFCN to Band
type LTEBand struct {
	ID            uint32
	Mode          DuplexMode
	CountEarfcn   uint32
	StartEarfcnDl uint32
	StartEarfcnUl uint32
}

// DuplexMode of LTE Band
type DuplexMode int32

const (
	// TDDMode
	TDDMode DuplexMode = iota
	// FDDMode
	FDDMode
)

var bands = [...]LTEBand{
	// FDDMode
	{ID: 1, Mode: FDDMode, StartEarfcnDl: 0, StartEarfcnUl: 18000, CountEarfcn: 600},
	{ID: 2, Mode: FDDMode, StartEarfcnDl: 600, StartEarfcnUl: 18600, CountEarfcn: 600},
	{ID: 3, Mode: FDDMode, StartEarfcnDl: 1200, StartEarfcnUl: 19200, CountEarfcn: 750},
	{ID: 4, Mode: FDDMode, StartEarfcnDl: 1950, StartEarfcnUl: 19950, CountEarfcn: 450},
	{ID: 28, Mode: FDDMode, StartEarfcnDl: 9210, StartEarfcnUl: 27210, CountEarfcn: 450},
	// TDDMode
	{ID: 38, Mode: TDDMode, StartEarfcnDl: 37750, CountEarfcn: 500},
	{ID: 39, Mode: TDDMode, StartEarfcnDl: 38250, CountEarfcn: 400},
	{ID: 40, Mode: TDDMode, StartEarfcnDl: 38650, CountEarfcn: 1000},
	{ID: 41, Mode: TDDMode, StartEarfcnDl: 39650, CountEarfcn: 1940},
	{ID: 42, Mode: TDDMode, StartEarfcnDl: 41590, CountEarfcn: 2000},
	{ID: 43, Mode: TDDMode, StartEarfcnDl: 43590, CountEarfcn: 2000},
	{ID: 48, Mode: TDDMode, StartEarfcnDl: 55240, CountEarfcn: 1500},
}

// EarfcnDLInRange checks that an EARFCN-DL belongs to a band
func (band LTEBand) EarfcnDLInRange(earfcndl uint32) bool {
	return band.StartEarfcnDl <= earfcndl && earfcndl < band.StartEarfcnDl+band.CountEarfcn
}

// EarfcnULInRange checks that an EARFCN-UL belongs to a band
func (band LTEBand) EarfcnULInRange(earfcnul uint32) bool {
	if band.Mode == FDDMode {
		return band.StartEarfcnUl <= earfcnul && earfcnul < band.StartEarfcnUl+band.CountEarfcn
	}
	return band.EarfcnDLInRange(earfcnul)
}

// GetBand for a EARFCN-UL
func GetBand(earfcndl uint32) (*LTEBand, error) {
	for _, band := range bands {
		if band.EarfcnDLInRange(earfcndl) {
			return &band, nil
		}
	}
	return nil, fmt.Errorf("Invalid EARFCNDL: no matching band")
}
