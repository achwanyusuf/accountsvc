package svcerr

import (
	"net/http"

	"github.com/achwanyusuf/carrent-lib/pkg/errormsg"
)

const (
	CodeBadRequest = iota + 40000
	CodePSQLErrorTransaction
	CodePSQLErrorCommit
	CodePSQLErrorRollback
	CodePSQLErrorInsert
	CodePSQLErrorUpdate
	CodePSQLErrorDelete
	CodePSQLErrorGet
	CodeInvalidEmptyName
	CodeInvalidEmptyEmail
	CodeInvalidEmailFormat
	CodeInvalidEmptyPassword
	CodeInvalidMinimumPassword
	CodeInvalidMaximumPassword
	CodeInvalidPasswordConfirmation
	CodeInvalidPasswordNotMatch
	CodeInvalidScope
	CodeInvalidClientIDClientSecret

	CodeNotAuthorized = 401000
	CodeNotFound      = 404000
)

var (
	AccountSVCPSQLErrorTransaction        = ErrMsg[CodePSQLErrorTransaction]
	AccountSVCPSQLErrorCommit             = ErrMsg[CodePSQLErrorCommit]
	AccountSVCPSQLErrorRollback           = ErrMsg[CodePSQLErrorRollback]
	AccountSVCPSQLErrorInsert             = ErrMsg[CodePSQLErrorInsert]
	AccountSVCPSQLErrorUpdate             = ErrMsg[CodePSQLErrorUpdate]
	AccountSVCPSQLErrorDelete             = ErrMsg[CodePSQLErrorDelete]
	AccountSVCPSQLErrorGet                = ErrMsg[CodePSQLErrorGet]
	AccountSVCNotAuthorized               = ErrMsg[CodeNotAuthorized]
	AccountSVCNotFound                    = ErrMsg[CodeNotFound]
	AccountSVCBadRequest                  = ErrMsg[CodeBadRequest]
	AccountSVCInvalidEmptyName            = ErrMsg[CodeInvalidEmptyName]
	AccountSVCInvalidEmptyEmail           = ErrMsg[CodeInvalidEmptyEmail]
	AccountSVCInvalidEmailFormat          = ErrMsg[CodeInvalidEmailFormat]
	AccountSVCInvalidEmptyPassword        = ErrMsg[CodeInvalidEmptyPassword]
	AccountSVCInvalidMinimumPassword      = ErrMsg[CodeInvalidMinimumPassword]
	AccountSVCInvalidMaximumPassword      = ErrMsg[CodeInvalidMaximumPassword]
	AccountSVCInvalidPasswordConfirmation = ErrMsg[CodeInvalidPasswordConfirmation]
	AccountSVCInvalidPasswordNotMatch     = ErrMsg[CodeInvalidPasswordNotMatch]
	AccountSVCInvalidScope                = ErrMsg[CodeInvalidScope]
	AccountSVCInvalidClientIDClientSecret = ErrMsg[CodeInvalidClientIDClientSecret]
)

var ErrMsg = map[int]errormsg.Message{
	CodePSQLErrorCommit: {
		Code:       CodePSQLErrorCommit,
		StatusCode: http.StatusBadRequest,
		Message:    "Terdapat kesalahan dalam pembuatan data!",
		Translation: errormsg.Translation{
			EN: "There was an error in creating the data ",
		},
	},
	CodePSQLErrorRollback: {
		Code:       CodePSQLErrorRollback,
		StatusCode: http.StatusBadRequest,
		Message:    "Terdapat kesalahan dalam pembuatan data!",
		Translation: errormsg.Translation{
			EN: "There was an error in creating the data ",
		},
	},
	CodePSQLErrorInsert: {
		Code:       CodePSQLErrorInsert,
		StatusCode: http.StatusBadRequest,
		Message:    "Terdapat kesalahan dalam pembuatan data!",
		Translation: errormsg.Translation{
			EN: "There was an error in creating the data ",
		},
	},
	CodePSQLErrorUpdate: {
		Code:       CodePSQLErrorUpdate,
		StatusCode: http.StatusBadRequest,
		Message:    "Terdapat kesalahan dalam mengubah data!",
		Translation: errormsg.Translation{
			EN: "There was an error in updating the data ",
		},
	},
	CodePSQLErrorDelete: {
		Code:       CodePSQLErrorDelete,
		StatusCode: http.StatusBadRequest,
		Message:    "Terdapat kesalahan dalam menghapus data!",
		Translation: errormsg.Translation{
			EN: "There was an error in deleting the data ",
		},
	},
	CodePSQLErrorGet: {
		Code:       CodePSQLErrorGet,
		StatusCode: http.StatusBadRequest,
		Message:    "Terdapat kesalahan dalam pengambilan data!",
		Translation: errormsg.Translation{
			EN: "There was an error in get data!",
		},
	},
	CodeNotFound: {
		Code:       CodeNotFound,
		StatusCode: http.StatusNotFound,
		Message:    "Data tidak ditemukan!",
		Translation: errormsg.Translation{
			EN: "Data not found!",
		},
	},
	CodeNotAuthorized: {
		Code:       CodeNotAuthorized,
		StatusCode: http.StatusUnauthorized,
		Message:    "Akses tidak diijinkan! Silakan login kembali!",
		Translation: errormsg.Translation{
			EN: "Access not authorized! Please login again!",
		},
	},
	CodeBadRequest: {
		Code:       CodeBadRequest,
		StatusCode: http.StatusBadRequest,
		Message:    "Kesalahan input. Silakan cek kembali masukan anda!",
		Translation: errormsg.Translation{
			EN: "Invalid input. Please validate your input!",
		},
	},
	CodePSQLErrorTransaction: {
		Code:       CodePSQLErrorTransaction,
		StatusCode: http.StatusBadRequest,
		Message:    "Terdapat kesalahan dalam pembuatan data!",
		Translation: errormsg.Translation{
			EN: "There was an error in creating the data ",
		},
	},
	CodeInvalidEmptyName: {
		Code:       CodeInvalidEmptyName,
		StatusCode: http.StatusBadRequest,
		Message:    "Nama tidak boleh kosong!",
		Translation: errormsg.Translation{
			EN: "Name should not be empty!",
		},
	},
	CodeInvalidEmptyEmail: {
		Code:       CodeInvalidEmptyEmail,
		StatusCode: http.StatusBadRequest,
		Message:    "Email tidak boleh kosong!",
		Translation: errormsg.Translation{
			EN: "Email should not be empty!",
		},
	},
	CodeInvalidEmailFormat: {
		Code:       CodeInvalidEmailFormat,
		StatusCode: http.StatusBadRequest,
		Message:    "Format email salah!",
		Translation: errormsg.Translation{
			EN: "Wrong email format!",
		},
	},
	CodeInvalidEmptyPassword: {
		Code:       CodeInvalidEmptyPassword,
		StatusCode: http.StatusBadRequest,
		Message:    "Kata sandi tidak boleh kosong!",
		Translation: errormsg.Translation{
			EN: "Password should not be empty!",
		},
	},
	CodeInvalidMinimumPassword: {
		Code:       CodeInvalidMinimumPassword,
		StatusCode: http.StatusBadRequest,
		Message:    "Kata sandi minimal 5 karakter!",
		Translation: errormsg.Translation{
			EN: "Minimum password is 5 character!",
		},
	},
	CodeInvalidMaximumPassword: {
		Code:       CodeInvalidMaximumPassword,
		StatusCode: http.StatusBadRequest,
		Message:    "Kata sandi maksimal 8 karakter!",
		Translation: errormsg.Translation{
			EN: "Maximum password is 8 character!",
		},
	},
	CodeInvalidPasswordConfirmation: {
		Code:       CodeInvalidPasswordConfirmation,
		StatusCode: http.StatusBadRequest,
		Message:    "Kata sandi dan konfirmasi kata sandi tidak sama!",
		Translation: errormsg.Translation{
			EN: "Password and password confirmation doesn't match!",
		},
	},
	CodeInvalidPasswordNotMatch: {
		Code:       CodeInvalidPasswordNotMatch,
		StatusCode: http.StatusUnauthorized,
		Message:    "Kata sandi salah!",
		Translation: errormsg.Translation{
			EN: "Wrong password!",
		},
	},
	CodeInvalidScope: {
		Code:       CodeInvalidScope,
		StatusCode: http.StatusBadRequest,
		Message:    "Scope tidak valid",
		Translation: errormsg.Translation{
			EN: "Invalid scope!",
		},
	},
	CodeInvalidClientIDClientSecret: {
		Code:       CodeInvalidClientIDClientSecret,
		StatusCode: http.StatusBadRequest,
		Message:    "Client ID/Client Secret harus diisi!",
		Translation: errormsg.Translation{
			EN: "Client ID/Client Secret should not be empty",
		},
	},
}
