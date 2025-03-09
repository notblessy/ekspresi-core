package router

type response struct {
	Success bool        `json:"success"`
	Message string      `json:"message,omitempty"`
	Data    interface{} `json:"data"`
}

const (
	defaultSize int = 10
	defaultPage int = 1
)

type responseWithPaging struct {
	Records     interface{}            `json:"records"`
	PageSummary map[string]interface{} `json:"page_summary"`
}

func withPaging(result interface{}, total int64, page, size int) responseWithPaging {
	offset := (page - 1) * size

	var hasNext bool
	if offset+size < int(total) {
		hasNext = true
	}

	return responseWithPaging{
		Records: result,
		PageSummary: map[string]interface{}{
			"size":    size,
			"page":    page,
			"hasNext": hasNext,
			"total":   total,
		},
	}
}
