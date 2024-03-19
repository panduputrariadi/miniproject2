package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"strings"
	"sync"
	"time"

	"github.com/go-pdf/fpdf"
)

type Perpustakaan struct {
	KodeBuku      string
	Judulbuku     string
	Pengarang     string
	Penerbit      string
	JumlahHalaman int
	TahunTerbit   int
}

var DaftarBuku []Perpustakaan

func tambahBuku() {
	jumlahHalaman := 0
	tahunTerbit := 0
	inputUser := bufio.NewReader(os.Stdin)

	draftBuku := []Perpustakaan{}
	for {
		fmt.Println("Tambah Daftar Buku")

		fmt.Print("Silahkan tambah kode buku: ")
		kodeBuku, err := inputUser.ReadString('\n')
		if err != nil {
			fmt.Print("terjadi error: ", err)
			return
		}
		kodeBuku = strings.TrimSpace(kodeBuku)

		fmt.Print("Silahkan tambah judul buku: ")
		judulBuku, err := inputUser.ReadString('\n')
		if err != nil {
			fmt.Print("terjadi error: ", err)
			return
		}
		judulBuku = strings.Replace(judulBuku, "\n", "", 1)

		fmt.Print("Silahkan tambah pengarang buku: ")
		pengarang, err := inputUser.ReadString('\n')
		if err != nil {
			fmt.Print("terjadi error: ", err)
			return
		}
		pengarang = strings.Replace(pengarang, "\n", "", 1)

		fmt.Print("Silahkan tambah penerbit buku: ")
		penerbit, err := inputUser.ReadString('\n')
		if err != nil {
			fmt.Print("terjadi error: ", err)
			return
		}
		penerbit = strings.Replace(penerbit, "\n", "", 1)

		fmt.Print("Silahkan Masukan jumlah halaman buku: ")
		_, err = fmt.Scanln(&jumlahHalaman)
		if err != nil {
			fmt.Println("Terjadi error: ", err)
			return
		}

		fmt.Print("Silahkan Masukan tahun terbit buku: ")
		_, err = fmt.Scanln(&tahunTerbit)
		if err != nil {
			fmt.Println("Terjadi error: ", err)
			return
		}

		draftBuku = append(draftBuku, Perpustakaan{
			KodeBuku:      kodeBuku,
			Judulbuku:     judulBuku,
			Pengarang:     pengarang,
			Penerbit:      penerbit,
			JumlahHalaman: jumlahHalaman,
			TahunTerbit:   tahunTerbit,
		})
		var pilihanOpsi = 0
		fmt.Println("Ketik 1 untuk tambah daftar buku, ketik 0 untuk keluar")
		_, err = fmt.Scanln(&pilihanOpsi)
		if err != nil {
			fmt.Println("Terjadi Error:", err)
			return
		}

		if pilihanOpsi == 0 {
			break
		}
	}

	fmt.Println("Menambah Daftar Buku...")

	_ = os.Mkdir("perpustakaan", 0777)

	ch := make(chan Perpustakaan)

	wg := sync.WaitGroup{}

	jumlahDaftarBuku := 5

	// Mendaftarkan receiver/pemroses data
	for i := 0; i < jumlahDaftarBuku; i++ {
		wg.Add(1)
		go simpanBuku(ch, &wg, i)
	}

	// Mengirimkan data ke channel
	for _, buku := range draftBuku {
		ch <- buku
	}

	close(ch)

	wg.Wait()

	fmt.Println("Berhasil menambahkan daftar buku")

}

func simpanBuku(ch <-chan Perpustakaan, wg *sync.WaitGroup, noStaff int) {

	for buku := range ch {
		dataJson, err := json.Marshal(buku)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		err = os.WriteFile(fmt.Sprintf("perpustakaan/%s.json", buku.KodeBuku), dataJson, 0644)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		fmt.Printf("Staff No %d Memproses buku ID : %s!\n", noStaff, buku.KodeBuku)
	}
	wg.Done()
}

func bacaDaftarBuku(ch <-chan string, chPerpustakaan chan Perpustakaan, wg *sync.WaitGroup) {
	var daftarBuku Perpustakaan
	for kodeDaftarBuku := range ch {
		dataJson, err := os.ReadFile(fmt.Sprintf("perpustakaan/%s", kodeDaftarBuku))

		if err != nil {
			fmt.Println("Terjadi error:", err)
		}

		err = json.Unmarshal(dataJson, &daftarBuku)
		if err != nil {
			fmt.Println("Terjadi error:", err)
		}
		chPerpustakaan <- daftarBuku
	}

	wg.Done()
}

func lihatDaftarBuku() {
	fmt.Println("=================================")
	fmt.Println("Lihat Pesanan")
	fmt.Println("=================================")
	fmt.Println("Memuat data ...")

	DaftarBuku = []Perpustakaan{}

	listJsonBuku, err := os.ReadDir("perpustakaan")
	if err != nil {
		fmt.Println("Terjadi error:", err)
	}

	wg := sync.WaitGroup{}
	ch := make(chan string)
	chBuku := make(chan Perpustakaan, len(listJsonBuku))

	jumlahStaff := 5
	for i := 0; i < jumlahStaff; i++ {
		wg.Add(1)
		go bacaDaftarBuku(ch, chBuku, &wg)
	}

	for _, fileBuku := range listJsonBuku {
		ch <- fileBuku.Name()
	}

	close(ch)
	wg.Wait()
	close(chBuku)

	for dataBuku := range chBuku {
		DaftarBuku = append(DaftarBuku, dataBuku)
	}

	for urutan, buku := range DaftarBuku {
		fmt.Printf("%d. Kode Buku: %s, Judul Buku: %s, Jumlah Halaman: %d, Penerbit: %s, Pengarang: %s, Tahun Terbit: %d\n",
			urutan+1,
			buku.KodeBuku,
			buku.Judulbuku,
			buku.JumlahHalaman,
			buku.Penerbit,
			buku.Pengarang,
			buku.TahunTerbit,
		)
	}

}

func hapusBuku() {
	fmt.Println("=================================")
	fmt.Println("Hapus Buku")
	fmt.Println("=================================")
	lihatDaftarBuku()
	fmt.Println("=================================")
	var urutanBuku int
	fmt.Print("Masukan Urutan Buku : ")
	_, err := fmt.Scanln(&urutanBuku)
	if err != nil {
		fmt.Println("Terjadi error:", err)
	}

	if (urutanBuku-1) < 0 ||
		(urutanBuku-1) > len(DaftarBuku) {
		fmt.Println("Urutan Buku Tidak Sesuai")
		hapusBuku()
		return
	}

	err = os.Remove(fmt.Sprintf("buku/%s.json", DaftarBuku[urutanBuku-1].KodeBuku))
	if err != nil {
		fmt.Println("Terjadi error:", err)
	}

	fmt.Println("Buku Berhasil Dihapus!")

}

func editBuku() {
	var kodeBuku string
	found := false

	fmt.Println("Edit Buku")
	fmt.Print("Masukkan kode buku yang ingin diedit: ")
	fmt.Scanln(&kodeBuku)

	// Membaca data buku dari file JSON
	data, err := os.ReadFile(fmt.Sprintf("perpustakaan/%s.json", kodeBuku))
	if err != nil {
		fmt.Println("Terjadi error:", err)
		return
	}

	var buku Perpustakaan
	err = json.Unmarshal(data, &buku)
	if err != nil {
		fmt.Println("Terjadi error:", err)
		return
	}

	fmt.Printf("Data Buku Ditemukan:\nKode Buku: %s, Judul Buku: %s, Pengarang: %s, Penerbit: %s, Jumlah Halaman: %d, Tahun Terbit: %d\n",
		buku.KodeBuku, buku.Judulbuku, buku.Pengarang, buku.Penerbit, buku.JumlahHalaman, buku.TahunTerbit)

	// Meminta input untuk data baru
	fmt.Println("Masukkan data baru untuk buku dengan kode buku", kodeBuku, ":")
	inputUser := bufio.NewReader(os.Stdin)

	fmt.Print("Judul Buku: ")
	judulBuku, _ := inputUser.ReadString('\n')
	judulBuku = strings.TrimSpace(judulBuku)

	fmt.Print("Pengarang: ")
	pengarang, _ := inputUser.ReadString('\n')
	pengarang = strings.TrimSpace(pengarang)

	fmt.Print("Penerbit: ")
	penerbit, _ := inputUser.ReadString('\n')
	penerbit = strings.TrimSpace(penerbit)

	fmt.Print("Jumlah Halaman: ")
	fmt.Scanln(&buku.JumlahHalaman)

	fmt.Print("Tahun Terbit: ")
	fmt.Scanln(&buku.TahunTerbit)

	// Update data buku
	buku.Judulbuku = judulBuku
	buku.Pengarang = pengarang
	buku.Penerbit = penerbit

	// Simpan data baru ke file JSON
	dataJSON, err := json.Marshal(buku)
	if err != nil {
		fmt.Println("Terjadi error saat marshaling data JSON:", err)
		return
	}

	err = os.WriteFile(fmt.Sprintf("perpustakaan/%s.json", kodeBuku), dataJSON, 0644)
	if err != nil {
		fmt.Println("Terjadi error saat menulis ke file JSON:", err)
		return
	}

	fmt.Println("Buku berhasil diupdate.", found)
}


func GeneratePdfBuku() {
	fmt.Println("=================================")
	fmt.Println("Membuat Daftar Buku ...")
	fmt.Println("=================================")

	pdf := fpdf.New("P", "mm", "A4", "")
	pdf.AddPage()

	pdf.SetFont("Arial", "", 12)
	pdf.SetLeftMargin(10)
	pdf.SetRightMargin(10)

	for i, buku := range DaftarBuku {
		bukuText := fmt.Sprintf(
			"Buku #%d:\nKode Buku : %s\nJudul Buku : %s\nPengarang : %s\nPenerbit : %s\nJumlah Halaman : %d\nTahun Terbit : %d\n",
			i+1, buku.KodeBuku, buku.Judulbuku,
			buku.Pengarang, buku.Penerbit,
			buku.JumlahHalaman, buku.TahunTerbit)

		pdf.MultiCell(0, 10, bukuText, "0", "L", false)
		pdf.Ln(5)
	}

	err := pdf.OutputFileAndClose(
		fmt.Sprintf("daftar_buku_%s.pdf",
			time.Now().Format("2006-01-02-15-04-05")))

	if err != nil {
		fmt.Println("Terjadi error:", err)
	}
}

func main() {
	// pilihanMenu := 0
	var pilihanMenu int
	fmt.Println("=================================")
	fmt.Println("Sistem Manajemen Perpustakaan")
	fmt.Println("=================================")
	fmt.Println("Silahkan Pilih : ")
	fmt.Println("1. Tambah Buku")
	fmt.Println("2. Liat Buku")
	fmt.Println("3. Hapus Buku")
	fmt.Println("4. Edit Buku")
	fmt.Println("5. Generate Daftar Buku")
	fmt.Println("6. Keluar")
	fmt.Println("=================================")
	fmt.Print("Masukan Pilihan : ")
	_, err := fmt.Scanln(&pilihanMenu)
	if err != nil {
		fmt.Println("Terjadi error:", err)
	}

	switch pilihanMenu {
	case 1:
		tambahBuku()
	case 2:
		lihatDaftarBuku()
	case 3:
		hapusBuku()
	case 4:
		editBuku()
	case 5:
		GeneratePdfBuku()
	case 6:
		os.Exit(0)
	}

	main()
}
