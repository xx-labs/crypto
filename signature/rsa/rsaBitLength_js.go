////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// This file is only compiled for WebAssembly.

package rsa

// DefaultRSABitLen is the RSA key length used in the system, in bits.
//
// WARNING: This bit size is smaller than the recommended bit size of 4096. Do
// not use this in production. Only use it for testing.
var DefaultRSABitLen = 1024
