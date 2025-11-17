package errors

import "errors"

var (
	// Transaction errors
	ErrDebitKreditTidakBalance         = errors.New("total debit tidak sama dengan total kredit")
	ErrDebitKreditKeduanya            = errors.New("satu baris tidak boleh memiliki debit dan kredit bersamaan")
	ErrTransaksiTidakDitemukan        = errors.New("transaksi tidak ditemukan")
	ErrTransaksiSudahDiPost           = errors.New("transaksi sudah di-post, tidak bisa diubah")
	ErrTransaksiSudahDiPostTidakBisaHapus = errors.New("transaksi sudah di-post, tidak bisa dihapus")
	ErrTidakAdaBarisTransaksi         = errors.New("tidak ada baris transaksi")
)
