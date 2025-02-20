package registry

const (
	DEVPROP_MASK_TYPE = 0x00000FFF
	//Constants
	RegSZ                       = 0x0001
	RegExpandSZ                 = 0x0002
	RegBin                      = 0x0003
	RegDWord                    = 0x0004
	RegMultiSZ                  = 0x0007
	RegQWord                    = 0x000B
	RegNone                     = 0x0000
	RegBigEndian                = 0x0005
	RegLink                     = 0x0006
	RegResourceList             = 0x0008
	RegFullResourceDescriptor   = 0x0009
	RegResourceRequirementsList = 0x000A
	RegFileTime                 = 0x0010
	// Following are new types from settings.dat
	RegUint8              = 0x101
	RegInt16              = 0x102
	RegUint16             = 0x103
	RegInt32              = 0x104
	RegUint32             = 0x105
	RegInt64              = 0x106
	RegUint64             = 0x107
	RegFloat              = 0x108
	RegDouble             = 0x109
	RegUnicodeChar        = 0x10A
	RegBoolean            = 0x10B
	RegUnicodeString      = 0x10C
	RegCompositeValue     = 0x10D
	RegDateTimeOffset     = 0x10E
	RegTimeSpan           = 0x10F
	RegGUID               = 0x110
	RegUnk111             = 0x111
	RegUnk112             = 0x112
	RegUnk113             = 0x113
	RegBytesArray         = 0x114
	RegInt16Array         = 0x115
	RegUint16Array        = 0x116
	RegInt32Array         = 0x117
	RegUInt32Array        = 0x118
	RegInt64Array         = 0x119
	RegUInt64Array        = 0x11A
	RegFloatArray         = 0x11B
	RegDoubleArray        = 0x11C
	RegUnicodeCharArray   = 0x11D
	RegBooleanArray       = 0x11E
	RegUnicodeStringArray = 0x11F

	//Constants to support the transaction log files (new format)
	LOG_ENTRY_SIZE_HEADER    = 40
	LOG_ENTRY_SIZE_ALIGNMENT = 0x200
)

var tt = []int{RegUint8, RegInt16, RegUint16, RegInt32, RegUint32,
	RegInt64, RegUint64, RegFloat, RegDouble, RegUnicodeChar,
	RegBoolean, RegUnicodeString, RegCompositeValue, RegDateTimeOffset,
	RegTimeSpan, RegGUID, RegUnk111, RegUnk112, RegUnk113, RegBytesArray,
	RegInt16Array, RegUint16Array, RegInt32Array, RegUInt32Array,
	RegInt64Array, RegUInt64Array, RegFloatArray, RegDoubleArray,
	RegUnicodeCharArray, RegBooleanArray, RegUnicodeStringArray}

// 字符串类型切片
var stringTypes = []int{
	RegSZ,
	RegExpandSZ,
	RegMultiSZ,
	RegUnicodeChar,
	RegUnicodeString,
	RegUnicodeCharArray,
	RegUnicodeStringArray,
}

// 字节数组类型切片
var byteArrayTypes = []int{
	RegBin,
	RegBytesArray,
}

// 32 位整数类型切片
var int32Types = []int{
	RegInt32,
	RegUint32,
	RegDWord,
	RegInt32Array,
	RegUInt32Array,
}

// 64 位整数类型切片
var int64Types = []int{
	RegInt64,
	RegUint64,
	RegQWord,
	RegInt64Array,
	RegUInt64Array,
}
