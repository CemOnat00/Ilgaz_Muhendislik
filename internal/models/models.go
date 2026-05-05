package models

import "time"

// Urun (Product) Model
type Urun struct {
	ID         string   `json:"id"`
	Kategori   string   `json:"kategori"`
	Baslik     string   `json:"baslik"`
	Aciklama   string   `json:"aciklama"`
	Resim      string   `json:"resim"`
	Ozellikler []string `json:"ozellikler"`
	Durum      string   `json:"durum"`
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

// Admin Login Models
type LoginRequest struct {
	KullaniciAdi string `json:"kullanici_adi"`
	Sifre        string `json:"sifre"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}
