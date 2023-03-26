package processing

import (
	"encoding/json"
	"fmt"
	"gitlab.com/Skinass/hakaton-2023-1-1/common/storage"
	"os"
	"time"
)

func MessageInitTime(data []byte) ([]byte, error) {
	ms := &storage.Message{}

	if err := json.Unmarshal(data, ms); err != nil {
		return nil, fmt.Errorf("unmarshal err: %w, data: %s", err, data)
	}

	ms.Time = time.Now()

	res, err := json.Marshal(ms)

	if err != nil {
		return nil, fmt.Errorf("marshal err: %w", err)
	}

	return res, nil
}

func MessageWrite(data []byte, file *os.File) error {
	_, err := file.Write([]byte(fmt.Sprintf("%s\n", data)))

	if err != nil {
		return fmt.Errorf("write err: %w", err)
	}

	return nil
}
