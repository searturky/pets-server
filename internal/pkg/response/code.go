package response

type CustomCode int

// 业务响应码
const (
	CodeSuccess            CustomCode = 0
	CodeBadRequest         CustomCode = 40000
	CodeUnauthorized       CustomCode = 40100
	CodeForbidden          CustomCode = 40300
	CodeNotFound           CustomCode = 40400
	CodeConflict           CustomCode = 40900
	CodeInternalError      CustomCode = 50000
	CodeServiceUnavailable CustomCode = 50300
	CodePetNotFound        CustomCode = 300100
)

func (c CustomCode) Name() string {
	var names = map[CustomCode]string{
		CodeSuccess:            "success",
		CodeBadRequest:         "bad request",
		CodeUnauthorized:       "unauthorized",
		CodeForbidden:          "forbidden",
		CodeNotFound:           "not found",
		CodeConflict:           "conflict",
		CodeInternalError:      "internal error",
		CodeServiceUnavailable: "service unavailable",
		CodePetNotFound:        "pet not found",
	}
	if name, ok := names[c]; ok {
		return name
	}
	return "unknown"
}
