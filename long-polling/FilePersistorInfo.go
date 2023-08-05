package long_polling

type FilePersistorInfo struct {
	Use                     bool
	FileName                string
	WriteBufferSize         int
	WriteFlushPeriodSeconds int
}
