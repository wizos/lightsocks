package lightsocks

import (
	"io"
	"testing"
)

type fakeReadWriterCloser struct {
}

func (f fakeReadWriterCloser) Read(p []byte) (int, error) {
	return 0, io.EOF
}
func (f fakeReadWriterCloser) Write(p []byte) (int, error) {
	return 0, nil
}
func (f fakeReadWriterCloser) Close() error {
	return nil
}

// 从src中源源不断的读取原数据加密后写入到dst，直到src中没有数据可以再读取
func (secureSocket *SecureTCPConn) OldEncodeCopy(dst io.ReadWriteCloser) error {
	buf := make([]byte, bufSize)
	for {
		readCount, errRead := secureSocket.Read(buf)
		if errRead != nil {
			if errRead != io.EOF {
				return errRead
			} else {
				return nil
			}
		}
		if readCount > 0 {
			writeCount, errWrite := (&SecureTCPConn{
				ReadWriteCloser: dst,
				Cipher:          secureSocket.Cipher,
			}).EncodeWrite(buf[0:readCount])
			if errWrite != nil {
				return errWrite
			}
			if readCount != writeCount {
				return io.ErrShortWrite
			}
		}
	}
}
func BenchmarkOldEncodeCopy(b *testing.B) {
	b.Run("make slice", func(b *testing.B) {
		var pswd password
		dst := SecureTCPConn{
			fakeReadWriterCloser{},
			newCipher(&pswd),
		}
		src := SecureTCPConn{
			fakeReadWriterCloser{},
			newCipher(&pswd),
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				dst.OldEncodeCopy(src)
			}
		})
	})
}

func BenchmarkEncodeCopy(b *testing.B) {
	b.Run("use object pool", func(b *testing.B) {
		var pswd password
		dst := SecureTCPConn{
			fakeReadWriterCloser{},
			newCipher(&pswd),
		}
		src := SecureTCPConn{
			fakeReadWriterCloser{},
			newCipher(&pswd),
		}
		b.ResetTimer()
		b.RunParallel(func(pb *testing.PB) {
			for pb.Next() {
				dst.EncodeCopy(src)
			}
		})
	})
}
