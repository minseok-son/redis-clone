package resp

import (
	"bufio"
	"fmt"
	"io"
	"strconv"
	"strings"
)

func extract(line string) (int, error) {
	countStr := strings.TrimSpace(line[1:])
	return strconv.Atoi(countStr)
}

func Parse(reader *bufio.Reader) ([]string, error) {
	// Read the Array Header (e.g., "*3\r\n")
	line, err := reader.ReadString('\n')
	if err != nil {
		return nil, err
	}

	if line[0] != '*' {
		return nil, fmt.Errorf("expected '*' at start of RESP array")
	}

	arrayLen, err := extract(line)
	if err != nil {
		return nil, fmt.Errorf("invalid array length: %v", err)
	}

	command := make([]string, 0, arrayLen)

	for i := 0; i < arrayLen; i++ {
		// Read the Bulk String Header (e.g., "$4\r\n")
		bulkHeader, err := reader.ReadString('\n')
		if err != nil {
			return nil, err
		}

		if bulkHeader[0] != '$' {
			return nil, fmt.Errorf("expected '$' for bulk string, got %q", bulkHeader[0])
		}

		strLen, err := extract(bulkHeader)
		if err != nil {
			return nil, err
		}

		// Read strLen bytes
		data := make([]byte, strLen)
		_, err = io.ReadFull(reader, data)
		if err != nil {
			return nil, err
		}

		_, err = reader.Discard(2)
		if err != nil {
			return nil, err
		}

		command = append(command, string(data))
	}

	return command, nil
}