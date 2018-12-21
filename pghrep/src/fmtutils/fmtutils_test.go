package fmtutils

import (
    "testing"
    "strings"
)

func TestByteFormat(t *testing.T) {
    var value string
    value = ByteFormat(982, 0)
    if strings.Compare(value, "982 bytes") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(1982, 0)
    if strings.Compare(value, "2 KiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(13820, 0)
    if strings.Compare(value, "14 KiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(135820, 0)
    if strings.Compare(value, "133 KiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(1735820, 0)
    if strings.Compare(value, "2 MiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(173583220, 0)
    if strings.Compare(value, "166 MiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(1735823330, 0)
    if strings.Compare(value, "2 GiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(173500823330, 0)
    if strings.Compare(value, "162 GiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = ByteFormat(17350082330230, 0)
    if strings.Compare(value, "16 TiB") != 0 {
        t.Fatal("TestGetFilePathSuccess failed")
    }
}

func TestGetUnit(t *testing.T) {
    var value int64
    value = GetUnit("8kB");
    if value != 8192 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = GetUnit("8MB");
    if value != 8388608 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = GetUnit("8GB");
    if value != 8589934592 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = GetUnit("8TB");
    if value != 8796093022208 {
        t.Fatal("TestGetFilePathSuccess failed")
    }
}

func TestGetUnit2(t *testing.T) {
    var value int64
    value = GetUnit("kB");
    
    if value != 1024 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = GetUnit("MB");
    if value != 1024 * 1024 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = GetUnit("GB");
    if value != 1024 * 1024 * 1024 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

    value = GetUnit("TB");
    if value != 1024 * 1024 * 1024 * 1024 {
        t.Fatal("TestGetFilePathSuccess failed")
    }

}