package defines

const (
	REGEXP_IMAGE                    = `(\.png|\.gif|\.jpeg|\.jpg|\.bmp|\.svg)$`
	REGEXP_MULTI_NUMBER_WITH_COMMA  = `^\d+(,\d+)*$`
	REGEXP_ID_LIST                  = `^\d{16}(,\d{16})*$`
	REGEXP_FILE                     = `^\/.+(,\/.+)*$`
	REGEXP_GENERAL_WORD             = "^[a-zA-Z\u4e00-\u9fa5]*$"
	REGEXP_GENERAL_WORD_WITH_SYMBOL = "^[a-zA-Z0-9\\s\u4e00-\u9fa5!@#$%^&*]{1,}$"
	REGEXP_USER_PLATFORM            = `(?i)Android|webOS|iPhone|iPod|BlackBerry|iPad`
	GEGEXP_YEAST_NUMBER             = `(.*)(?=-.*)`
)
