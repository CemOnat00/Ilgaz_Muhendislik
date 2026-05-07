package models

import "time"

// Urun (Product) Model
type Urun struct {
	ID           string   `json:"id"`
	Kategori     string   `json:"kategori"`
	Baslik       string   `json:"baslik"`       // Admin and legacy
	Isim         string   `json:"isim"`         // Used in catalog
	Aciklama     string   `json:"aciklama"`
	Resim        string   `json:"resim"`
	Fiyat        float64  `json:"fiyat"`
	StoktaVarMi  bool     `json:"stokta_var_mi"`
	Ozellikler   []string `json:"ozellikler"`
	Durum        string   `json:"durum"`        // Aktif/Pasif
}

// Mesaj (Contact Message) Model
type Mesaj struct {
	ID       string    `json:"id"`
	AdSoyad  string    `json:"ad_soyad"`
	Email    string    `json:"email"`
	Konu     string    `json:"konu"`
	Mesaj    string    `json:"mesaj"`
	Tarih    time.Time `json:"tarih"`
	OkunduMu bool      `json:"okundu_mu"`
}

// Proje (Engineering Project) Model
type Proje struct {
	ID         string    `json:"id"`
	Baslik     string    `json:"baslik"`
	Musteri    string    `json:"musteri"`
	Kategori   string    `json:"kategori"`
	Durum      string    `json:"durum"` // Devam Ediyor, Tamamlandı
	Tarih      time.Time `json:"tarih"`
	Aciklama   string    `json:"aciklama"`
	Resim      string    `json:"resim"`
}

// Admin Login Models
type LoginRequest struct {
	KullaniciAdi string `json:"kullanici_adi"`
	Sifre        string `json:"sifre"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
