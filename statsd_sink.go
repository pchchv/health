package health

import "bytes"

type StatsDSinkSanitizationFunc func(*bytes.Buffer, string)
