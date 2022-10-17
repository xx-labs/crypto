////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// This file is compiled for all architectures except WebAssembly.
//go:build !js || !wasm
// +build !js !wasm

package rsa

// DefaultRSABitLen is the RSA key length used in the system, in bits.
var DefaultRSABitLen = 4096
