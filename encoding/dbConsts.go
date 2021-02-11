package encoding

//obfuscating types
const (
	TinyintType   = "tinyint"   //1 bytes
	SmallintType  = "smallint"  //2 bytes
	MediumintType = "mediumint" //3 bytes
	IntType       = "int"       //4 bytes
	BigintType    = "bigint"    //8 bytes

	UTinyintType   = "tinyint unsigned"   //1 bytes
	USmallintType  = "smallint unsigned"  //2 bytes
	UMediumintType = "mediumint unsigned" //3 bytes
	UIntType       = "int unsigned"       //4 bytes
	UBigintType    = "bigint unsigned"    //8 bytes

	FloatType   = "float"   //4 bytes
	DoubleType  = "double"  //other names cast to this automatically by mysql //8 bytes
	DecimalType = "decimal" //default (10) -> (10, 0)//other names cast to this automatically by mysql

	CharType    = "char"    //default (1)
	VarcharType = "varchar" //no default

	TinytextType   = "tinytext"
	TextType       = "text"
	MediumtextType = "mediumtext"
	LongtextType   = "longtext"
)

//obfuscating bounds
const (
	LowerBoundTinyint   = -128
	UpperBoundTinyint   = 127
	LowerBoundSmallint  = -32768
	UpperBoundSmallint  = 32767
	LowerBoundMediumint = -8388608
	UpperBoundMediumint = 8388607
	LowerBoundInt       = -2147483648
	UpperBoundInt       = 2147483647
	LowerBoundBigint    = -9223372036854775808
	UpperBoundBigint    = 9223372036854775807

	UpperBoundUTinyint   = 255
	UpperBoundUSmallint  = 65535
	UpperBoundUMediumint = 16777215
	UpperBoundUInt       = 4294967295
	UpperBoundUBigint    = 18446744073709551615
)
