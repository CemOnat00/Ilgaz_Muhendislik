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

// Review (Ürün Değerlendirmesi) Modeli
type Review struct {
	ID          bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	ProductID   bson.ObjectID `json:"product_id" bson:"product_id"`       // hangi ürüne ait
	ProductName string        `json:"product_name" bson:"product_name"`   // denormalized — admin listesi kolaylığı
	Name        string        `json:"ad_soyad" bson:"name"`
	Phone       string        `json:"-" bson:"phone"`                     // JSON'da gizli — sadece admin görür
	Email       string        `json:"-" bson:"email"`                     // JSON'da gizli — sadece admin görür
	Rating      int           `json:"rating" bson:"rating"`               // 1-5 arası
	Content     string        `json:"content" bson:"content"`
	IsApproved  bool          `json:"is_approved" bson:"is_approved"`     // false = beklemede
	AdminReply  string        `json:"admin_reply" bson:"admin_reply"`     // boş = yanıt yok
	HelpfulYes  int           `json:"helpful_yes" bson:"helpful_yes"`
	HelpfulNo   int           `json:"helpful_no" bson:"helpful_no"`
	VoterIPs    []string      `json:"-" bson:"voter_ips"`                 // IP tekrar oy koruması
	IPAddress   string        `json:"-" bson:"ip_address"`               // gönderenin IP'si
	CreatedAt   time.Time     `json:"created_at" bson:"created_at"`
	UpdatedAt   time.Time     `json:"updated_at" bson:"updated_at"`
}

// PublicReview — Müşteri tarafına gönderilecek, hassas alanları maskeli response struct
type PublicReview struct {
	ID         string    `json:"id"`
	MaskedName string    `json:"ad_soyad"`     // "A*** Y***" formatı
	Rating     int       `json:"rating"`
	Content    string    `json:"content"`
	AdminReply string    `json:"admin_reply"`
	HelpfulYes int       `json:"helpful_yes"`
	HelpfulNo  int       `json:"helpful_no"`
	CreatedAt  time.Time `json:"created_at"`
	HasVoted   bool      `json:"has_voted"`    // bu IP oy verdi mi? (handler'da set edilir)
}

// AdminReview — Admin paneline gönderilecek, tüm alanları görünür response struct
type AdminReview struct {
	ID          string    `json:"id"`
	ProductID   string    `json:"product_id"`
	ProductName string    `json:"product_name"`
	Name        string    `json:"ad_soyad"`
	Phone       string    `json:"telefon"`
	Email       string    `json:"email"`
	Rating      int       `json:"rating"`
	Content     string    `json:"content"`
	IsApproved  bool      `json:"is_approved"`
	AdminReply  string    `json:"admin_reply"`
	HelpfulYes  int       `json:"helpful_yes"`
	HelpfulNo   int       `json:"helpful_no"`
	CreatedAt   time.Time `json:"created_at"`
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

// SiteSettings (Site / İletişim Ayarları) — tek dokümanlı koleksiyon
type SiteSettings struct {
	ID           bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	CompanyName  string        `json:"sirket_adi" bson:"company_name"`
	Email        string        `json:"email" bson:"email"`
	Phone        string        `json:"telefon" bson:"phone"`
	SupportEmail string        `json:"destek_email" bson:"support_email"`
	Address      string        `json:"adres" bson:"address"`
	MapLocation  string        `json:"harita_konum" bson:"map_location"` // Harita için adres ya da "enlem,boylam"
	MapZoom      int           `json:"harita_zoom" bson:"map_zoom"`
}

// Lead (Potansiyel Müşteri) — katalog indirme ve yorum formundan toplanan iletişim bilgisi.
type Lead struct {
	ID        bson.ObjectID `json:"id,omitempty" bson:"_id,omitempty"`
	Name      string        `json:"ad_soyad" bson:"name"`
	Phone     string        `json:"telefon" bson:"phone"`
	Email     string        `json:"email" bson:"email"`
	Category  string        `json:"kategori" bson:"category"`  // İndirilen katalog ("Genel" = tüm ürünler)
	Source    string        `json:"kaynak" bson:"source"`      // "catalog_download" | "review_form"
	IPAddress string        `json:"ip_adresi" bson:"ip_address"`
	CreatedAt time.Time     `json:"created_at" bson:"created_at"`
}

// DefaultSettings — veritabanında ayar dokümanı yoksa kullanılan varsayılanlar
func DefaultSettings() SiteSettings {
	return SiteSettings{
		CompanyName:  "Ilgaz Mühendislik A.Ş.",
		Email:        "info@ilgazmuhendislik.com",
		Phone:        "+90 212 123 45 67",
		SupportEmail: "support@ilgazmuhendislik.com",
		Address:      "Mimarsinan, İsmet İnönü Blv. no:140 D:A, 55200 Atakum/Samsun",
		MapLocation:  "41.33043975392915,36.27995416441768",
		MapZoom:      17,
	}
}
