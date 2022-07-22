package file

import (
	"os"
)

// file_dir +"/country"+"/company"
const FILE_DIR = "data"

const FILE_FLAG_APPEND = os.O_RDWR | os.O_CREATE | os.O_APPEND
const FILE_FLAG_TRUNC = os.O_RDWR | os.O_CREATE | os.O_TRUNC
