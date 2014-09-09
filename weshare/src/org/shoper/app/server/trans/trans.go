package trans

/**
 *工厂
 */
type TransFactory struct {
	decode ProtocolDecoder
	encode ProtocolEncoder
}

/**
 *解码器接口
 */
type ProtocolDecoder interface {
	decode() error
	finishDecode() error
	dispose() error
}

/**
 *编码器接口
 */
type ProtocolEncoder interface {
	decode() error
	finishDecode() error
	dispose() error
}
