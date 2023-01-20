////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// This file is compiled for all architectures except WebAssembly.
//go:build !js || !wasm

package rsa

// defaultRSABitLen is the RSA key length used in the system, in bits.
var defaultRSABitLen = 4096
