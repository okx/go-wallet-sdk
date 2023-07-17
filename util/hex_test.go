package util

import (
	"testing"
)

func TestDecodeHexString(t *testing.T) {
	bytes, err := DecodeHexString("0xe33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bytes)
}

func TestDecodeHexString2(t *testing.T) {
	bytes, err := DecodeHexString("e33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
	if err != nil {
		t.Fatal(err)
	}
	t.Log(bytes)
}

func BenchmarkDecodeHexStringWithout0x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes, err := DecodeHexString("e33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
		if err != nil {
			b.Fatal(err)
		}
		b.Log(bytes)
	}
}

func BenchmarkDecodeHexStringWith0x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes, err := DecodeHexString("0xe33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
		if err != nil {
			b.Fatal(err)
		}
		b.Log(bytes)
	}
}

func BenchmarkDecodeHexStringBackupWithout0x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes, err := DecodeHexStringBackup("e33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
		if err != nil {
			b.Fatal(err)
		}
		b.Log(bytes)
	}
}

func BenchmarkDecodeHexStringBackupWith0x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes, err := DecodeHexStringBackup("0xe33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
		if err != nil {
			b.Fatal(err)
		}
		b.Log(bytes)
	}
}

func BenchmarkRemoveZeroHexWithout0x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes := RemoveZeroHex("e33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
		b.Log(bytes)
	}
}

func BenchmarkRemoveZeroHexWith0x(b *testing.B) {
	for i := 0; i < b.N; i++ {
		bytes := RemoveZeroHex("0xe33ef3d7883cd3f6b9c2a72b916c36066cca8443c718fb53bc0a0607de9e4d9a")
		b.Log(bytes)
	}
}
