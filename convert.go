package convert

import (
	"io"

	"github.com/pkg/errors"
	"github.com/saintfish/chardet"
	"golang.org/x/text/encoding/charmap"
	"golang.org/x/text/encoding/japanese"
	"golang.org/x/text/encoding/korean"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/encoding/unicode"
	"golang.org/x/text/encoding/unicode/utf32"
	"golang.org/x/text/transform"
)

func transformEncoding(rawReader io.Reader, trans transform.Transformer) io.Reader {
	return transform.NewReader(rawReader, trans)
}

func NewTransformer(r io.ReadSeeker) (io.Reader, error) {
	var detect = make([]byte, 128)
	_, err := r.Read(detect)
	if err != nil {
		return nil, err
	}

	if _, err := r.Seek(0, io.SeekStart);err != nil {
		return nil, err
	}
	charDetector := chardet.NewTextDetector()
	detectResult, err := charDetector.DetectBest(detect)
	if err != nil {
		return nil, err
	}
	switch detectResult.Charset {
	case "GB-18030":
		return transformEncoding(r, simplifiedchinese.GB18030.NewDecoder()), nil
	case "BIG":
		return transformEncoding(r, simplifiedchinese.GBK.NewDecoder()), nil
	case "ISO-2022-CN":
		return nil, errors.Errorf("%s not supported",detectResult.Charset)
	case "ISO-2022-JP":
		return transformEncoding(r,japanese.ISO2022JP.NewDecoder()), nil
	case "ISO-2022-KR":
		return nil, errors.Errorf("%s not supported",detectResult.Charset)
	case "EUC-KR":
		return transformEncoding(r, korean.EUCKR.NewDecoder()), nil
	case "EUC-JP":
		return transformEncoding(r,japanese.EUCJP.NewDecoder()), nil
	case "Shift_JIS":
		return transformEncoding(r,japanese.ShiftJIS.NewDecoder()), nil
	case "ISO-8859-1":
		return transformEncoding(r,charmap.ISO8859_1.NewDecoder()), nil
	case "UTF-16LE":
		return transformEncoding(r,unicode.UTF16(unicode.LittleEndian,unicode.IgnoreBOM).NewDecoder()), nil
	case "UTF-16BE":
		return transformEncoding(r,unicode.UTF16(unicode.BigEndian,unicode.IgnoreBOM).NewDecoder()), nil
	case "UTF-32LE":
		return transformEncoding(r,utf32.UTF32(utf32.LittleEndian,utf32.IgnoreBOM).NewDecoder()), nil
	case "UTF-32BE":
		return transformEncoding(r,utf32.UTF32(utf32.BigEndian,utf32.IgnoreBOM).NewDecoder()), nil

	default:
		return r, nil
	}
}
