package iso8583Define

// 解析后域配置
type FieldConfig struct {
	FieldValueType FieldValueType // 域数据类型
	FieldAlignType FieldAlignType // 域数据对齐类型
	IsValueLenFix  bool           // 域数据长度是否固定
	ValueLen       int            // 数据固定长度/最大长度
	LenLen         int            // 可变域长度值的长度
	LenType        FieldValueType // 数据长度类型
	PaddingByte    byte           // 添加数据
	Mask           bool           // 是否打印
}
