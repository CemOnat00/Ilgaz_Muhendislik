package models

import (
	"time"

	"go.mongodb.org/mongo-driver/v2/bson"
)

// TechnicalSpec (Teknik Özellik)
type TechnicalSpec struct {
	Key   string `json:"key" bson:"key"`
	Value string `json:"value" bson:"value"`
}

// FeatureBox (Ürün Özellik Kutusu)
type FeatureBox struct {
	Title       string `json:"title" bson:"title"`
	Description string `json:"description" bson:"description"`
}

// Product (Ürün) Model
type Product struct {
	ID               bson.ObjectID   `json:"id,omitempty" bson:"_id,omitempty"`
	Name             string          `json:"name" bson:"name"`
	Baslik           string          `json:"baslik" bson:"baslik"`             // Admin/Legacy title
	Isim             string          `json:"isim" bson:"isim"`                 // Catalog display name
	Category         string          `json:"category" bson:"category"`         // Internal/Admin
	Kategori         string          `json:"kategori" bson:"kategori"`         // Catalog display
	Price            float64         `json:"fiyat" bson:"price"`
	StockStatus      bool            `json:"stokta_var_mi" bson:"stock_status"` // true: Stokta, false: Yok
	ImageURL         string          `json:"resim" bson:"image_url"`
	Description      string          `json:"aciklama" bson:"description"`
	Features         []string        `json:"ozellikler" bson:"features"`
	Status           string          `json:"durum" bson:"status"` // Aktif/Pasif (durum for frontend)
	Order            int             `json:"order" bson:"order"`
	CreatedAt        time.Time       `json:"created_at" bson:"created_at"`
	Images           []string        `json:"images" bson:"images"`
	Videos           []string        `json:"videos" bson:"videos"`
	TechnicalSpecs   []TechnicalSpec `json:"technical_specs" bson:"technical_specs"`
	FeatureBoxes     []FeatureBox    `json:"feature_boxes" bson:"feature_boxes"`
	ApplicationAreas string          `json:"application_areas" bson:"application_areas"`
	IsIndustrial     bool            `json:"is_industrial" bson:"is_industrial"`
	IsHomeAppliance  bool            `json:"is_home_appliance" bson:"is_home_appliance"`
}

// Message (İletişim Mesajı) Model
type Message struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string             `json:"ad_soyad" bson:"name"`
	Email     string             `json:"email" bson:"email"`
	Subject   string             `json:"konu" bson:"subject"`
	Content   string             `json:"mesaj" bson:"content"`
	Status    string             `json:"status" bson:"status"`       // Unread/Replied
	IsRead    bool               `json:"okundu_mu" bson:"is_read"`   // Admin panel compatibility
	CreatedAt time.Time          `json:"created_at" bson:"created_at"`
	Tarih     time.Time          `json:"tarih" bson:"-"`             // JSON duplicate for frontend
}

// Project (Proje) Model
type Project struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Title       string             `json:"title" bson:"title"`
	Customer    string             `json:"customer" bson:"customer"`
	Category    string             `json:"category" bson:"category"`
	SubTitle    string             `json:"subtitle" bson:"subtitle"`
	Progress    int                `json:"progress" bson:"progress"`
	Status      string             `json:"status" bson:"status"` // In Progress, Completed
	Description string             `json:"description" bson:"description"`
	ImageURL    string             `json:"image_url" bson:"image_url"`
	Order       int                `json:"order" bson:"order"`
	Date        time.Time          `json:"date" bson:"date"`
}

// Admin Model
type Admin struct {
	ID         bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Username   string             `json:"username" bson:"username"`
	Password   string             `json:"-" bson:"password"`
	Email      string             `json:"email" bson:"email"`
	ResetToken string             `json:"-" bson:"reset_token,omitempty"`
}

// Admin Login Models
type LoginRequest struct {
	Username string `json:"kullanici_adi"`
	Password string `json:"sifre"`
}

type LoginResponse struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

type ForgotPasswordRequest struct {
	Email string `json:"email"`
}

type ResetPasswordRequest struct {
	Token       string `json:"token"`
	NewPassword string `json:"new_password"`
}
