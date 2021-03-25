package iso8583

import (
	"encoding/hex"
	"fmt"
	"strconv"
)

// MtiType is the message type identifier type
type MtiType struct {
	mti string
}

// String returns the mti as a string
func (m *MtiType) String() string {
	return m.mti
}

// ElementsType stores iso8583 elements in a map
type ElementsType struct {
	elements map[int64]string
}

// GetElements returns the available elemts as a map
func (e *ElementsType) GetElements() map[int64]string {
	return e.elements
}

// IsoStruct is an iso8583 container
type IsoStruct struct {
	Spec     Spec
	Mti      MtiType
	Bitmap   []int64
	Elements ElementsType
	Tpdu     []byte
}

// ToString packs the mti, bitmap and elements into a string
func (iso *IsoStruct) ToString() (string, error) {
	var str string
	// get done with the mti and the bitmap
	bitmapString, err := BitMapArrayToHex(iso.Bitmap)
	if err != nil {
		return str, err
	}
	elementsStr, err := iso.packElements()
	if err != nil {
		return str, err
	}

	if iso.Spec.fields[0].HeaderHex {
		mtiByte, _ := hex.DecodeString(iso.Mti.String())

		var bitmapByte []byte
		if iso.Spec.fields[1].HeaderHex {
			bitmapByte, _ = hex.DecodeString(bitmapString)
		}

		var isomsgByte []byte
		if len(iso.Tpdu) > 0 {
			isomsgByte = append(mtiByte, bitmapByte...)
			isomsgByte = append(iso.Tpdu, isomsgByte...)
		} else {
			isomsgByte = append(mtiByte, bitmapByte...)
		}

		return string(isomsgByte) + elementsStr, nil
	}

	fmt.Printf("create message, tpdu %v\n", iso.Tpdu)
	if len(iso.Tpdu) > 0 {
		str = string(iso.Tpdu) + iso.Mti.String() + bitmapString + elementsStr
	} else {
		str = iso.Mti.String() + bitmapString + elementsStr
	}
	return str, nil
}

// AddMTI adds the provided iso8583 MTI into the current struct
// also updates the bitmap in the process
func (iso *IsoStruct) AddMTI(data string) error {
	mti := MtiType{mti: data}
	_, err := MtiValidator(mti)
	if err != nil {
		return err
	}
	iso.Mti = mti
	return nil
}

// AddField adds the provided iso8583 field into the current struct
// also updates the bitmap in the process
func (iso *IsoStruct) AddField(field int64, data string) error {
	if field < 2 || field > int64(len(iso.Bitmap)) {
		return fmt.Errorf("expected field to be between %d and %d found %d instead", 2, len(iso.Bitmap), field)
	}
	iso.Bitmap[field-1] = 1
	iso.Elements.elements[field] = data
	return nil
}

func (iso *IsoStruct) RemoveField(field int64) error {
	if field < 2 || field > int64(len(iso.Bitmap)) {
		return fmt.Errorf("expected field to be between %d and %d found %d instead", 2, len(iso.Bitmap), field)
	}
	iso.Bitmap[field-1] = 0
	iso.Elements.elements[field] = ""
	return nil
}

// Parse parses an iso8583 string
func (iso *IsoStruct) Parse(i string, useTpdu bool) (IsoStruct, error) {
	var q IsoStruct
	spec := iso.Spec
	var msg string
	var tpdu []byte

	if useTpdu {
		var err error
		tpdu, msg, err = extractTpdu(i)
		if err != nil {
			fmt.Println(err.Error())
		}

		iso.Tpdu = tpdu

	} else {
		msg = i
	}
	fmt.Printf("tpdu: %v\n", iso.Tpdu)

	mti, rest := extractMTI(msg, spec.fields[0].HeaderHex)
	bitmap, elementString, err := extractBitmap(rest, spec.fields[1].HeaderHex)

	if err != nil {
		return q, err
	}

	// validat the mti
	_, err = MtiValidator(mti)
	if err != nil {
		return q, err
	}

	elements, err := unpackElements(bitmap, elementString, spec)
	if err != nil {
		return q, err
	}

	q = IsoStruct{Spec: spec, Mti: mti, Bitmap: bitmap, Elements: elements, Tpdu: tpdu}
	return q, nil
}

func (iso *IsoStruct) packElements() (string, error) {
	var str string
	bitmap := iso.Bitmap
	elementsMap := iso.Elements.GetElements()
	elementsSpec := iso.Spec

	for index := 1; index < len(bitmap); index++ { // index 0 of bitmap isn't need here
		if bitmap[index] == 1 { // if the field is present
			field := int64(index + 1)
			fieldDescription := elementsSpec.fields[int(field)]
			if fieldDescription.LenType == "fixed" {

				if fieldDescription.HeaderHex {
					strtemp := elementsMap[field]
					strByte, _ := hex.DecodeString(strtemp)
					str = str + string(strByte)
				} else {
					str = str + elementsMap[field]
				}

			} else {
				lengthType, err := getVariableLengthFromString(fieldDescription.LenType)
				if err != nil {
					return str, err
				}
				actualLength := len(elementsMap[field])
				paddedLength := leftPad(strconv.Itoa(actualLength), int(lengthType), "0")

				if fieldDescription.HeaderHex {
					strtemp := (paddedLength + elementsMap[field])
					strByte, _ := hex.DecodeString(strtemp)
					str = str + string(strByte)
				} else {
					str = str + elementsMap[field]
				}

			}
		}
	}
	return str, nil
}

func extractTpdu(rest string) ([]byte, string, error) {

	if len(rest) < 5 {
		return nil, "", fmt.Errorf("could not slice %d string of %d\n", len(rest), 5)
	}

	frontHex := rest[0:5]
	tpdu := []byte(frontHex)
	msg := rest[5:len(rest)]

	return tpdu, msg, nil
}

// extractMTI extracts the mti from an iso8583 string
func extractMTI(str string, isHex bool) (MtiType, string) {

	if !isHex {

		if len(str) < 4 {
			return MtiType{}, ""
		}

		mti := str[0:4]
		rest := str[4:len(str)]

		return MtiType{mti: mti}, rest
	} else {

		if len(str) < 2 {
			return MtiType{}, ""
		}

		mti := hex.EncodeToString([]byte(str[0:2]))
		rest := str[2:len(str)]

		return MtiType{mti: string(mti)}, rest
	}
}

func extractBitmap(rest string, isHex bool) ([]int64, string, error) {
	var bitmap []int64
	var elementsString string
	var inDec []byte
	var err error

	if len(rest) < 1 {
		return bitmap, elementsString, fmt.Errorf("bitmap length = 0, no bitmap to be processed")
	}

	if !isHex {
		// remove first two characters
		frontHex := rest[0:2]
		//fmt.Println(frontHex)
		inDec, err = hex.DecodeString(frontHex)
		if err != nil {
			return bitmap, elementsString, err
		}
	} else {
		// remove first characters
		frontHex := hex.EncodeToString([]byte(rest[0:1]))
		//fmt.Println(frontHex)
		inDec, err = hex.DecodeString(string(frontHex))
		if err != nil {
			return bitmap, elementsString, err
		}
	}

	inBinary := fmt.Sprintf("%8b", inDec[0])
	compare := "1"
	var bitmapHexLength int

	// if the first bit of the bitmap is 1,
	// it means a secondary bitmap exist hence its a 128 bit bitmap (hex length 32)
	if inBinary[0] == compare[0] { // don't why I did it like this
		// secondary bitmap exists
		bitmapHexLength = 32
	} else {
		// only primary bitmap is there
		// 64 bit bitmap (hex length 16)
		bitmapHexLength = 16
	}

	if isHex {
		bitmapHexLength = bitmapHexLength / 2
	}

	var bitmapHexString string
	if !isHex {
		bitmapHexString = rest[0:bitmapHexLength]
	} else {
		bitmapHexByte := hex.EncodeToString([]byte(rest[0:bitmapHexLength]))
		bitmapHexString = string(bitmapHexByte)
	}

	elementsString = rest[bitmapHexLength:len(rest)]

	bitmap, err = HexToBitmapArray(bitmapHexString)
	if err != nil {
		return bitmap, elementsString, err
	}
	return bitmap, elementsString, nil
}

func getVariableLengthFromString(str string) (int64, error) {
	var num int64
	if str == "llvar" {
		return 2, nil
	}
	if str == "lllvar" {
		return 3, nil
	}
	if str == "llllvar" {
		return 4, nil
	}

	return num, fmt.Errorf("%s is an invalid LenType", str)
}

func extractFieldFromElements(spec Spec, field int, str string) (string, string, error) {
	var extractedField, substr string
	fieldDescription := spec.fields[int(field)]

	if fieldDescription.LenType == "fixed" {

		var err error
		extractedField, substr, err = getFieldValue(fieldDescription.HeaderHex, fieldDescription.MaxLen, fieldDescription.Contain, str)
		if err != nil {
			return extractedField, substr, fmt.Errorf("spec error: field %d: %s", field, err.Error())
		}

	} else {
		// varianle length fields have their lengths embedded into the string
		length, err := getVariableLengthFromString(fieldDescription.LenType)
		if err != nil {
			return extractedField, substr, fmt.Errorf("spec error: field %d: %s", field, err.Error())
		}

		var fieldLength string
		if fieldDescription.HeaderHex {
			if length%2 != 0 {
				length = (length + 1) / 2
			} else {
				length = length / 2
			}

			fieldLength1 := hex.EncodeToString([]byte(str[0:length]))
			fieldLength = string(fieldLength1)

		} else {
			fieldLength = str[0:length] // get the embedded length
		}
		tempSubstr := str[length:len(str)] // get the string with the length removed

		fieldLengthInt, err := strconv.ParseInt(fieldLength, 10, 64)
		if err != nil {
			return extractedField, substr, err
		}

		if fieldDescription.Contain == "string" {
			fieldLengthInt = fieldLengthInt * 2
		}

		extractedField, substr, err = getFieldValue(fieldDescription.HeaderHex, int(fieldLengthInt), fieldDescription.Contain, tempSubstr)
		if err != nil {
			return extractedField, substr, err
		}
	}

	return extractedField, substr, nil
}

func getFieldValue(headerHex bool, maxLen int, contain string, str string) (extractedField string, substr string, err error) {
	if headerHex {
		var length int
		if maxLen%2 != 0 {
			length = (maxLen + 1) / 2
		} else {
			length = maxLen / 2
		}

		//if xxxvar chip, length di x 2
		if contain == "chip-tag" {
			length = length * 2
		}

		if len(str) < length {
			return extractedField, substr, fmt.Errorf("could not slice %d string of %d\n", len(str), length)
		}

		extractedFieldTemp := str[0:length]
		extractedFieldTemp2 := hex.EncodeToString([]byte(extractedFieldTemp))
		extractedField = string(extractedFieldTemp2)

		substr = str[length:len(str)]
	} else {
		if len(str) < maxLen {
			return extractedField, substr, fmt.Errorf("could not slice %d string of %d\n", len(str), maxLen)
		}

		extractedField = str[0:maxLen]
		substr = str[maxLen:len(str)]
	}

	return extractedField, substr, nil
}

func unpackElements(bitmap []int64, elements string, spec Spec) (ElementsType, error) {
	var elem ElementsType
	var m = make(map[int64]string)
	currentString := elements
	// The first (index 0) bit of the bitmap shows the presense(1)/absense(0) of the secondary
	// we therefore start with the second bit (index 1) which is field (2)
	for index := 1; index < len(bitmap); index++ {
		bit := bitmap[index]
		if bit == 1 { // field is present
			field := index + 1 // adjust to account for the fact that arrays start at 0
			extractedField, substr, err := extractFieldFromElements(spec, field, currentString)
			if err == nil {
				m[int64(field)] = extractedField
				currentString = substr
			} else {
				return elem, err
			}
		}
	}

	elem = ElementsType{elements: m}
	return elem, nil
}

// NewISOStruct creates a new IsoStruct
// based on the content of the specfile provided
func NewISOStruct(filename string, secondaryBitmap bool) IsoStruct {
	var iso IsoStruct
	var bitmap []int64
	mti := MtiType{mti: ""}

	if secondaryBitmap == true {
		bitmap = make([]int64, 128)
		bitmap[0] = 1
	} else {
		bitmap = make([]int64, 64)
	}

	emap := make(map[int64]string)
	elements := ElementsType{elements: emap}
	spec, err := SpecFromFile(filename)
	if err != nil {
		panic(err) // we panic because we don't want to do anything without a valid specfile
	}

	var tpdu []byte
	tpdu = make([]byte, 5)
	fmt.Printf("tpdu: %#v", tpdu)

	iso = IsoStruct{Spec: spec, Mti: mti, Bitmap: bitmap, Elements: elements, Tpdu: tpdu}
	return iso
}
