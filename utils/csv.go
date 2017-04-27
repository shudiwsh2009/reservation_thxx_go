package utils

import (
	"encoding/csv"
	re "github.com/shudiwsh2009/reservation_thxx_go/rerror"
	"golang.org/x/text/encoding/simplifiedchinese"
	"golang.org/x/text/transform"
	"os"
	"path/filepath"
)

var EXPORT_FOLDER = filepath.Join("static", "export")

const (
	EXCEL_FILE_SUFFIX = ".xlsx"
	CSV_FILE_SUFFIX   = ".csv"
)

func WriteToCSV(data [][]string, path string) error {
	// 写入文件
	fout, err := os.Create(path)
	if err != nil {
		return re.NewRError("建立文件失败", err)
	}
	defer fout.Close()
	w := csv.NewWriter(transform.NewWriter(fout, simplifiedchinese.GB18030.NewEncoder()))
	w.UseCRLF = true
	if err = w.WriteAll(data); err != nil {
		return re.NewRError("写入表数据失败", err)
	}
	w.Flush()
	return nil
}

func ReadFromCSV(path string) ([][]string, error) {
	fin, err := os.Open(path)
	if err != nil || fin == nil {
		return nil, re.NewRError("打开文件失败: %v", err)
	}
	defer fin.Close()
	w := csv.NewReader(transform.NewReader(fin, simplifiedchinese.GB18030.NewDecoder()))
	data, err := w.ReadAll()
	if err != nil {
		return nil, re.NewRError("读取文件失败: %v", err)
	}
	return data, nil
}
