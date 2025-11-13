// pdf.go
package main

import (
	"bytes"
	"fmt"
	"strings"
	"time"
)

func GeneratePDF(tasks []*Task) []byte {
	var buf bytes.Buffer

	fmt.Fprintln(&buf, "%PDF-1.4")
	objects := []string{}
	objNum := 1

	objects = append(objects, fmt.Sprintf("%d 0 obj\n<< /Type /Catalog /Pages %d 0 R >>\nendobj\n", objNum, objNum+1))
	objNum++

	objects = append(objects, fmt.Sprintf("%d 0 obj\n<< /Type /Pages /Kids [%d 0 R] /Count 1 >>\nendobj\n", objNum, objNum+1))
	objNum++

	objects = append(objects, fmt.Sprintf("%d 0 obj\n<< /Type /Page /Parent %d 0 R /Contents %d 0 R /Resources << /Font << /F1 %d 0 R >> >> /MediaBox [0 0 595 842] >>\nendobj\n", objNum, objNum-1, objNum+1, objNum+2))
	objNum++

	objects = append(objects, fmt.Sprintf("%d 0 obj\n<< /Type /Font /Subtype /Type1 /BaseFont /Helvetica >>\nendobj\n", objNum))
	objNum++

	var content strings.Builder
	content.WriteString("BT /F1 16 Tf 50 800 Td (Отчет по доступности ссылок) Tj ET\n")
	content.WriteString(fmt.Sprintf("BT /F1 10 Tf 50 780 Td (Сформирован: %s) Tj ET\n", time.Now().Format("02.01.2006 15:04")))
	content.WriteString("BT /F1 12 Tf 50 750 Td (ID | URL | Статус) Tj ET\n")
	content.WriteString("BT /F1 10 Tf 50 735 Td (--------------------------------------------------) Tj ET\n")

	y := 710.0
	for _, task := range tasks {
		task.mu.RLock()
		for _, result := range task.Results {
			status := "Доступен"
			if !result.Available {
				status = "Недоступен"
			}
			line := fmt.Sprintf("%3d | %-40s | %s", task.ID, truncate(result.URL, 40), status)
			content.WriteString(fmt.Sprintf("BT /F1 10 Tf 50 %.1f Td (%s) Tj ET\n", y, escapePDF(line)))
			y -= 18
		}
		task.mu.RUnlock()
	}

	contentObj := fmt.Sprintf("%d 0 obj\n<< /Length %d >>\nstream\n%s\nendstream\nendobj\n", objNum, len(content.String()), content.String())
	objects = append(objects, contentObj)

	offsets := make([]int, len(objects))
	current := buf.Len()
	for i, obj := range objects {
		offsets[i] = current
		fmt.Fprint(&buf, obj)
		current += len(obj)
	}

	xref := buf.Len()
	fmt.Fprintln(&buf, "xref")
	fmt.Fprintf(&buf, "0 %d\n", len(objects)+1)
	fmt.Fprintln(&buf, "0000000000 65535 f ")
	for _, off := range offsets {
		fmt.Fprintf(&buf, "%010d 00000 n \n", off)
	}

	fmt.Fprintln(&buf, "trailer")
	fmt.Fprintf(&buf, "<< /Size %d /Root 1 0 R >>\n", len(objects)+1)
	fmt.Fprintln(&buf, "startxref")
	fmt.Fprintf(&buf, "%d\n", xref)
	fmt.Fprintln(&buf, "%%EOF")

	return buf.Bytes()
}

func truncate(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n-3] + "..."
}

func escapePDF(s string) string {
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "(", "\\(")
	s = strings.ReplaceAll(s, ")", "\\)")
	return s
}
