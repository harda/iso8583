package iso8583

import (
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/harda/iso8583"
)

func TestISOParseByte(t *testing.T) {
	// MTI = 0200
	// Field (3) = 000010
	// Field (4) = 1500
	// Field (7) = 1206041200
	// Field (11) = 000001
	// Field (41) = 12340001
	// Field (49) = 840
	isobyte := []byte{
		0x60, 0x00, 0x32, 0x00, 0x00, 0x02, 0x00, 0x30, 0x20, 0x07, 0x80, 0x20, 0xC0, 0x12, 0x64, 0x00,
		0x00, 0x00, 0x00, 0x00, 0x00, 0x10, 0x02, 0x00, 0x00, 0x00, 0x92, 0x00, 0x51, 0x00, 0x00, 0x00,
		0x32, 0x00, 0x37, 0x53, 0x36, 0x19, 0x00, 0x02, 0x41, 0x87, 0x65, 0xD2, 0x50, 0x32, 0x01, 0x94,
		0x71, 0x82, 0x12, 0x00, 0x00, 0x00, 0x37, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x36, 0x30, 0x30,
		0x30, 0x30, 0x30, 0x38, 0x37, 0x37, 0x30, 0x30, 0x30, 0x30, 0x30, 0x30, 0x36, 0x5B, 0x6B, 0x19,
		0xED, 0x5B, 0xE9, 0x40, 0x95, 0x01, 0x47, 0x5F, 0x2A, 0x02, 0x03, 0x60, 0x82, 0x02, 0x18, 0x00,
		0x84, 0x07, 0xA0, 0x00, 0x00, 0x00, 0x04, 0x10, 0x10, 0x95, 0x05, 0x80, 0x00, 0x04, 0x08, 0x00,
		0x9A, 0x03, 0x21, 0x03, 0x02, 0x9C, 0x01, 0x00, 0x9F, 0x02, 0x06, 0x00, 0x00, 0x00, 0x10, 0x02,
		0x00, 0x9F, 0x03, 0x06, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x9F, 0x09, 0x02, 0x00, 0x02, 0x9F,
		0x10, 0x12, 0x01, 0x10, 0xA0, 0x00, 0x01, 0x22, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00, 0x00,
		0x00, 0x00, 0x00, 0xFF, 0x9F, 0x1A, 0x02, 0x03, 0x60, 0x9F, 0x1E, 0x08, 0x35, 0x31, 0x34, 0x36,
		0x33, 0x33, 0x39, 0x35, 0x9F, 0x26, 0x08, 0xA3, 0xE1, 0x6C, 0x4A, 0xB9, 0x3C, 0xFC, 0xB1, 0x9F,
		0x27, 0x01, 0x80, 0x9F, 0x33, 0x03, 0xE0, 0xF8, 0xC8, 0x9F, 0x34, 0x03, 0x42, 0x03, 0x00, 0x9F,
		0x35, 0x01, 0x22, 0x9F, 0x36, 0x02, 0x00, 0x19, 0x9F, 0x37, 0x04, 0x8F, 0x1D, 0xD4, 0x95, 0x9F,
		0x41, 0x04, 0x00, 0x00, 0x00, 0x92, 0x9F, 0x53, 0x01, 0x52, 0x00, 0x11, 0xDF, 0x01, 0x08, 0x35,
		0x31, 0x34, 0x36, 0x33, 0x33, 0x39, 0x35, 0x01, 0x60, 0xB4, 0x03, 0xDC, 0x3A, 0xD3, 0x37, 0x6E,
		0xC9, 0x6D, 0x00, 0x43, 0xB6, 0x15, 0x58, 0xE4, 0x8F, 0x99, 0x9B, 0x4C, 0x79, 0x12, 0x8C, 0xDD,
		0xB6, 0x74, 0x5E, 0x0B, 0xE7, 0xFA, 0x8C, 0x72, 0x03, 0xD4, 0xA1, 0x88, 0xA1, 0xAB, 0xF5, 0xA6,
		0xD8, 0x09, 0x67, 0x1F, 0xA0, 0x2D, 0xAA, 0xE9, 0xC5, 0x19, 0x0B, 0x25, 0x35, 0xD8, 0xC2, 0x30,
		0x44, 0x89, 0x6F, 0x8C, 0xEF, 0xD3, 0xF3, 0x31, 0x75, 0xD9, 0xA8, 0x34, 0x58, 0x1E, 0x98, 0x3A,
		0x73, 0x3F, 0xA2, 0x1B, 0x5C, 0x4F, 0x6C, 0x34, 0xC5, 0x90, 0x0C, 0x8B, 0x79, 0x2E, 0xF2, 0xCA,
		0x9F, 0x8D, 0x14, 0xDB, 0x11, 0x31, 0x6D, 0x75, 0x20, 0xF0, 0xE1, 0x15, 0x71, 0xA0, 0xA4, 0xB7,
		0x60, 0x91, 0xB0, 0x0D, 0xDA, 0xA4, 0x45, 0x0A, 0x30, 0x65, 0xC2, 0x9C, 0x31, 0x51, 0xBB, 0x4F,
		0xC2, 0xEC, 0x41, 0xEF, 0x5E, 0xE6, 0xD4, 0xFF, 0xC1, 0x09, 0x3E, 0x80, 0x69, 0x77, 0x0C, 0x9A,
		0x8A, 0xDE, 0x65, 0x9C, 0x37, 0xC1, 0xD8, 0x76, 0xF2, 0x07, 0xA7, 0x43, 0x38, 0x71, 0xE1, 0x79,
		0xD0, 0x4F, 0xA5, 0x4F, 0xC7, 0xC8, 0x97, 0x14, 0x59, 0x00, 0x06, 0x34, 0x30, 0x30, 0x30, 0x30,
		0x31}

	isomsg := string(isobyte)
	isostruct := NewISOStruct("spec1987pos.yml", true)
	parsed, err := isostruct.Parse(isomsg, true)
	if err != nil {
		fmt.Println(err)
		t.Errorf("parse iso message failed")
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		t.Errorf("failed to unpack valid isomsg")
	}
	fmt.Println(isomsgUnpacked)
	// if isomsgUnpacked != isomsg {
	// 	t.Errorf("%s should be %s", isomsgUnpacked, isomsg)
	// }
	fmt.Printf("%#v, %#v\n%#v", parsed.Mti, parsed.Bitmap, parsed.Elements)
}

func TestISOParseInt(t *testing.T) {

	isobyte := []byte{
		96, 0, 50, 0, 0, 4, 0, 112, 36, 7, 128, 0, 192, 2, 100, 22, 83, 4, 135, 32, 0, 0, 8, 72, 0, 0, 0, 0, 0, 0, 2, 1, 0, 0, 1, 32, 35, 6, 0, 81, 0, 1, 0, 50, 0, 55, 55,
		48, 48, 48, 48, 48, 54, 48, 48, 48, 48, 48, 56, 55, 55, 48, 48, 48, 48, 48, 48, 54, 1, 87, 95, 42, 2, 3, 96, 130, 2, 116, 0, 132, 7, 160, 0, 0, 6, 2, 16, 16,
		149, 5, 8, 0, 4, 136, 0, 154, 3, 33, 3, 22, 156, 1, 0, 159, 2, 6, 0, 0, 0, 2, 1, 0, 159, 3, 6, 0, 0, 0, 0, 0, 0, 159, 9, 2, 1, 0, 159, 16, 28, 159, 1, 160, 0, 128,
		0, 0, 145, 243, 17, 8, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 0, 159, 26, 2, 3, 96, 159, 30, 8, 53, 49, 52, 54, 51, 51, 57, 53, 159, 38, 8, 235, 72, 167,
		52, 131, 92, 159, 216, 159, 39, 1, 128, 159, 51, 3, 224, 248, 200, 159, 52, 3, 2, 0, 0, 159, 53, 1, 34, 159, 54, 2, 3, 179, 159, 55, 4, 253, 168, 58, 116,
		159, 65, 4, 0, 0, 1, 32, 159, 83, 1, 82, 0, 17, 223, 1, 8, 53, 49, 52, 54, 51, 51, 57, 53, 1, 96, 180, 3, 220, 58, 211, 55, 110, 201, 198, 159, 188, 92, 177,
		104, 150, 181, 7, 86, 112, 228, 200, 27, 244, 118, 252, 59, 125, 47, 215, 14, 125, 202, 49, 36, 240, 155, 38, 232, 2, 88, 128, 49, 253, 147, 102, 189, 55,
		148, 50, 71, 79, 223, 3, 37, 119, 62, 178, 98, 90, 163, 211, 181, 42, 32, 79, 139, 96, 10, 66, 115, 42, 231, 190, 76, 66, 68, 225, 96, 127, 29, 4, 224, 83, 151,
		57, 40, 246, 53, 27, 27, 124, 91, 221, 13, 110, 13, 3, 46, 55, 127, 198, 230, 71, 215, 105, 129, 198, 121, 216, 234, 55, 187, 227, 241, 187, 188, 246, 112, 8,
		203, 40, 214, 90, 94, 229, 75, 185, 143, 93, 73, 69, 151, 158, 122, 57, 169, 227, 38, 88, 231, 185, 171, 231, 2, 239, 112, 114, 67, 130, 59, 176, 202, 13, 249,
		81, 30, 22, 38, 236, 68, 0, 6, 52, 48, 48, 48, 48, 56}

	isomsg := string(isobyte)
	isostruct := NewISOStruct("spec1987pos.yml", true)
	parsed, err := isostruct.Parse(isomsg, true)
	if err != nil {
		fmt.Println(err)
		t.Errorf("parse iso message failed")
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		t.Errorf("failed to unpack valid isomsg")
	}
	fmt.Println(isomsgUnpacked)
	// if isomsgUnpacked != isomsg {
	// 	t.Errorf("%s should be %s", isomsgUnpacked, isomsg)
	// }
	fmt.Printf("%#v, %#v\n%#v", parsed.Mti, parsed.Bitmap, parsed.Elements)
}

func TestISOParse(t *testing.T) {
	// MTI = 0200
	// Field (3) = 000010
	// Field (4) = 1500
	// Field (7) = 1206041200
	// Field (11) = 000001
	// Field (41) = 12340001
	// Field (49) = 840
	isomsg := "02003220000000808000000010000000001500120604120000000112340001840"
	isostruct := NewISOStruct("spec1987.yml", true)
	parsed, err := isostruct.Parse(isomsg, false)
	if err != nil {
		fmt.Println(err)
		t.Errorf("parse iso message failed")
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		t.Errorf("failed to unpack valid isomsg")
	}
	if isomsgUnpacked != isomsg {
		t.Errorf("%s should be %s", isomsgUnpacked, isomsg)
	}
	// fmt.Printf("%#v, %#v\n%#v", parsed.Mti, parsed.Bitmap, parsed.Elements)
}

func TestEmpty(t *testing.T) {
	one := NewISOStruct("spec1987.yml", false)

	if one.Mti.String() != "" {
		t.Errorf("Empty generates invalid MTI")
	}
	one.AddMTI("0200")
	one.AddField(3, "000010")
	one.AddField(4, "000000001500")
	one.AddField(7, "1206041200")
	one.AddField(11, "000001")
	one.AddField(41, "12340001")
	one.AddField(49, "840")

	expected := "02003220000000808000000010000000001500120604120000000112340001840"
	unpacked, _ := one.ToString()
	if unpacked != expected {
		t.Errorf("Manually constructed isostruct produced %s not %s", unpacked, expected)
	}
}

func TestEmptyPos(t *testing.T) {
	one := NewISOStruct("spec1987pos.yml", false)

	if one.Mti.String() != "" {
		t.Errorf("Empty generates invalid MTI")
	}
	one.AddMTI("0200")
	one.Tpdu = []byte{96, 0, 50, 0, 0}
	one.AddField(3, "000010")
	one.AddField(4, "000000001500")
	one.AddField(7, "1206041200")
	one.AddField(11, "000001")
	one.AddField(41, "12340001")
	one.AddField(49, "840")

	dataByte, _ := hex.DecodeString("0200322000000080800000001000000000150031323036303431323030000001")
	expected := "12340001840"
	expected = string(dataByte) + expected

	unpacked, _ := one.ToString()
	if unpacked != expected {
		t.Errorf("Manually constructed isostruct produced %x not %x", []byte(unpacked), []byte(expected))
	}
}

func TestMessageFromSample1(t *testing.T) {

	isobyte, _ := hex.DecodeString("60001800000800202001000080000492000000029900183737303030303333003748544c45303331303031303031373730303030333330303030303030378ca64de98ca64de9")

	isomsg := string(isobyte)
	isostruct := NewISOStruct("spec1987pos.yml", true)
	parsed, err := isostruct.Parse(isomsg, true)
	if err != nil {
		fmt.Println(err)
		t.Errorf("parse iso message failed")
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		t.Errorf("failed to unpack valid isomsg")
	}
	fmt.Println(isomsgUnpacked)

	one := iso8583.NewISOStruct("spec1987pos.yml", false)

	one.AddMTI("0800")
	one.Tpdu = []byte{96, 0, 24, 0, 0}
	one.AddField(3, "920000")
	one.AddField(11, "000299")
	one.AddField(24, "0018")
	one.AddField(41, "77000033")
	one.AddField(62, "48544c45303331303031303031373730303030333330303030303030378ca64de98ca64de9")

	oneString, _ := one.ToString()

	fmt.Println("-------------")
	fmt.Println("-------------", oneString)
	if isomsgUnpacked != oneString {
		t.Errorf("%s should be %s", isomsgUnpacked, oneString)
	}
	fmt.Printf("visionet sample 1: %#v, %#v\n%#v", parsed.Mti, parsed.Bitmap, parsed.Elements)
	fmt.Printf("length msgbyte: %d , msgstring %d", cap(isobyte), len(isomsg))
	fmt.Println("-------------")
}

func TestMessageFromSample2(t *testing.T) {

	// one := iso8583.NewISOStruct("spec1987pos.yml", false)

	// one.AddMTI("0800")
	// one.Tpdu = []byte{96, 0, 24, 0, 0}
	// one.AddField(1, "2020010000800004")
	// one.AddField(3, "920000")
	// one.AddField(11, "000299")
	// one.AddField(24, "0018")
	// one.AddField(41, "77000006")
	// one.AddField(62, "48544c45303331303031303031373730303030333330303030303030378ca64de98ca64de9")

	isobyte, _ := hex.DecodeString("600047000002003038078020C01224000000000000000100200162202928031900510000004700275178632590094319D2207221800F3132303031353930303030313030303132303030303135AD068BB4DB53A96901615F2A0203605F340100820274008407A0000006021010950508000418009A032103199C01009F02060000000001009F03060000000000009F090201009F101C0101A0000000000050D16B00000000000000000000000000000000009F1A0203609F1E0835313838343138349F2608921771745F33F43B9F2701809F3303E0F8C89F34030200009F3501229F360202C19F3704C8C765CA9F4104002001629F5301520160878A33A2D14751188791CD1AE0C1EBB9EF9284CB378D8A516CB73AA5F055A87E6B43CD74C7ABBFB6161BD628BD1B4AF6AA4DC195A6385D909B0BEC55AF93A3C354CABDCE83F6F818AE81DC41D6D0A52410C79B00236B530CBAC17778CDA57FEB5FAD0C2AAE5C2642A49A0DA77DD919172719DA3CFA3C2CB8D3011F7E3B0EECFABE689933DA576C2D62905CD8A100D98BA66F163FAD98BDEF46456CA00D8235280006323030303033")

	// if one.ToString() != string(hexraw) {
	// 	t.Errorf("%s should be %s", one.ToString(), string(hexraw))
	// }
	// fmt.Printf("%#v, %#v\n%#v", one.Mti, one.Bitmap, one.Elements)

	isomsg := string(isobyte)
	isostruct := NewISOStruct("spec1987pos2.yml", true)
	parsed, err := isostruct.Parse(isomsg, true)
	if err != nil {
		fmt.Println(err)
		t.Errorf("parse iso message failed")
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		t.Errorf("failed to unpack valid isomsg")
	}
	fmt.Println(isomsgUnpacked)

	lenbyte := make([]byte, 2)
	lenbyte[0] = byte(len(isomsg) / 256)
	lenbyte[1] = byte(len(isomsg))
	fmt.Printf("len of sample 2: %#v\n", lenbyte)

	// if isomsgUnpacked != isomsg {
	// 	t.Errorf("%s should be %s", isomsgUnpacked, isomsg)
	// }
	fmt.Printf("visionet sample 2: %#v, %#v\n%#v", parsed.Mti, parsed.Bitmap, parsed.Elements)
}

func TestMessageFromSample3(t *testing.T) {

	isobyte, _ := hex.DecodeString("60001800000800202001000080000492000000014100183737303030303036003748544c45303331303031303031373730303030303630303030303030378ca64de98ca64de9")

	isomsg := string(isobyte)
	isostruct := NewISOStruct("spec1987pos.yml", true)
	parsed, err := isostruct.Parse(isomsg, true)
	if err != nil {
		fmt.Println(err)
		t.Errorf("parse iso message failed")
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		t.Errorf("failed to unpack valid isomsg")
	}
	fmt.Println(isomsgUnpacked)

	lenbyte := make([]byte, 2)
	lenbyte[0] = byte(len(isomsg) / 256)
	lenbyte[1] = byte(len(isomsg))
	fmt.Printf("len of sample 3a: %#v\n", lenbyte)

	// if isomsgUnpacked != isomsg {
	// 	t.Errorf("%s should be %s", isomsgUnpacked, isomsg)
	// }
	fmt.Printf("visionet sample 3a: %#v, %#v\n%#v", parsed.Mti, parsed.Bitmap, parsed.Elements)
}
func TestMessageFromSample4(t *testing.T) {

	isobyte, _ := hex.DecodeString("600009000002003020078020C0124500000000000000030000035900510001000800375304872000000848D2306226000000362000003737303030303333303030303038373730303030303333F9FF7FA34D1778A001575F2A020360820274008407A0000006021010950508000488009A032103039C01009F02060000000003009F03060000000000009F090201009F101C9F01A00000000088692C8C00000000000000000000000000000000009F1A0203609F1E0835313838343138349F26089839C8F4F17310739F2701809F3303E0F8C89F34030200009F3501229F360203A19F37046669A26B9F4104000003599F5301520011DF0108353138383431383400063430303032300000000000000000")

	isomsg := string(isobyte)
	isostruct := NewISOStruct("spec1987pos.yml", true)
	parsed, err := isostruct.Parse(isomsg, true)
	if err != nil {
		fmt.Println(err)
		t.Errorf("parse iso message failed")
	}

	isomsgUnpacked, err := parsed.ToString()
	if err != nil {
		fmt.Println(err)
		t.Errorf("failed to unpack valid isomsg")
	}
	fmt.Println(isomsgUnpacked)

	one := iso8583.NewISOStruct("spec1987pos.yml", false)

	one.AddMTI("0800")
	one.Tpdu = []byte{96, 0, 24, 0, 0}
	one.AddField(3, "000000")
	one.AddField(4, "000000000300")
	one.AddField(11, "000359")
	one.AddField(22, "0051")
	one.AddField(23, "0001")
	one.AddField(24, "0008")
	one.AddField(25, "00")
	one.AddField(35, "5304872000000848d230622600000036200000")
	one.AddField(41, "77000033")
	one.AddField(42, "000008770000033")
	one.AddField(52, "f9ff7fa34d1778a0")
	one.AddField(55, "5f2a020360820274008407a0000006021010950508000488009a032103039c01009f02060000000003009f03060000000000009f090201009f101c9f01a00000000088692c8c00000000000000000000000000000000009f1a0203609f1e0835313838343138349f26089839c8f4f17310739f2701809f3303e0f8c89f34030200009f3501229f360203a19f37046669a26b9f4104000003599f53015200")
	one.AddField(58, "df01083531383834313834")
	one.AddField(62, "48544c45303331303031303031373730303030333330303030303030378ca64de98ca64de9")
	one.AddField(64, "0000000000000000")

	oneString, _ := one.ToString()

	fmt.Println("-------------")
	fmt.Printf("length msgbyte: %d , msgstring %d", cap(isobyte), len(isomsg))

	lenbyte := make([]byte, 2)
	lenbyte[0] = byte(len(isomsg) / 256)
	lenbyte[1] = byte(len(isomsg))
	fmt.Printf("len of sample 4: %#v\n", lenbyte)

	if isomsgUnpacked != oneString {
		t.Errorf("%s should be %s", isomsgUnpacked, oneString)
	}
	fmt.Printf("visionet sample 4: %#v, %#v\n%#v", parsed.Mti, parsed.Bitmap, parsed.Elements)
	// fmt.Println("-------------")
}
