package database

import (
	"context"
	"log"
	"os"
	"time"

	"github.com/cemonat00/ilgaz-backend/internal/models"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

var MongoClient *mongo.Client
var DB *mongo.Database

const (
	dbName = "IlgazDB"
)

func ConnectDB() {
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("MONGO_URI environment variable is not set")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	clientOptions := options.Client().ApplyURI(uri)
	client, err := mongo.Connect(clientOptions)
	if err != nil {
		log.Fatal("MongoDB Bağlantı Hatası:", err)
	}

	// Ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal("MongoDB Atlas'a erişilemiyor:", err)
	}

	MongoClient = client
	DB = client.Database(dbName)
	log.Println("MongoDB Atlas'a başarıyla bağlanıldı!")
}

func GetCollection(collectionName string) *mongo.Collection {
	return DB.Collection(collectionName)
}

func SeedAdmin() {
	collection := GetCollection("admins")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Default credentials (In production, move to .env)
	defaultUser := "admin"
	defaultPass := "ilgaz2026"

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(defaultPass), bcrypt.DefaultCost)

	// Update or Create admin
	opts := options.UpdateOne().SetUpsert(true)
	filter := bson.M{"username": defaultUser}
	update := bson.M{
		"$set": bson.M{
			"username": defaultUser,
			"password": string(hashedPassword),
			"email":    "admin@ilgazmuhendislik.com",
		},
	}

	_, err := collection.UpdateOne(ctx, filter, update, opts)
	if err != nil {
		log.Println("Admin seed hatası:", err)
	} else {
		log.Println("Admin hesabı güncellendi/oluşturuldu (Güvenli Hashlendi)")
	}
}

func SeedProducts() {
	collection := GetCollection("products")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Check if already seeded
	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Println("Urun sayisi kontrol edilemedi:", err)
		return
	}
	if count > 0 {
		log.Println("Veritabaninda urunler zaten mevcut. Seed islemi atlandi.")
		return
	}

	log.Println("Urunler koleksiyonu bos, ornek urunler yukleniyor...")

	now := time.Now()
	fiveDaysAgo := now.AddDate(0, 0, -5)
	fortyFiveDaysAgo := now.AddDate(0, 0, -45)

	seedProducts := []interface{}{
		// Vana Grubu
		models.Product{
			ID:              bson.NewObjectID(),
			Name:            "Yüksek Basınçlı Buhar Vanası",
			Baslik:          "Yüksek Basınçlı Buhar Vanası - PN40",
			Isim:            "Yüksek Basınçlı Buhar Vanası",
			Category:        "Vana Grubu",
			Kategori:        "Vana Grubu",
			Price:           12500.0,
			StockStatus:     true,
			IsIndustrial:    true,
			IsHomeAppliance: false,
			ImageURL:        "https://images.unsplash.com/photo-1581092160607-ee22621dd758?auto=format&fit=crop&w=600&q=80",
			Description:     "Endüstriyel tesisler, enerji santralleri ve buhar hatlarında yüksek basınç altında çalışmak üzere tasarlanmış, tam sızdırmazlık sağlayan flanşlı buhar vanası.",
			Status:          "Aktif",
			Order:           0,
			CreatedAt:       fiveDaysAgo,
			Images: []string{
				"https://images.unsplash.com/photo-1581092160607-ee22621dd758?auto=format&fit=crop&w=600&q=80",
			},
			TechnicalSpecs: []models.TechnicalSpec{
				{Key: "Gövde Malzemesi", Value: "Döküm Çelik (GP240GH)"},
				{Key: "Bağlantı Tipi", Value: "Flanşlı"},
				{Key: "Maks. Sıcaklık", Value: "+400°C"},
				{Key: "Nominal Basınç", Value: "PN40"},
			},
			FeatureBoxes: []models.FeatureBox{
				{Title: "Yüksek Dayanıklılık", Description: "GP240GH çelik gövde tasarımı ile ekstrem sıcaklık ve basınca karşı direnç."},
				{Title: "Hassas Akış Kontrolü", Description: "Özel klape tasarımı sayesinde hassas debi ayarı ve sıfır sızıntı."},
			},
			ApplicationAreas: "Kimya endüstrisi, buhar hatları, ısıtma ve iklimlendirme sistemleri.",
		},
		models.Product{
			ID:              bson.NewObjectID(),
			Name:            "Kelebek Vana - Flanşlı Tip",
			Baslik:          "Kelebek Vana - Flanşlı Tip DN100",
			Isim:            "Kelebek Vana - Flanşlı Tip",
			Category:        "Vana Grubu",
			Kategori:        "Vana Grubu",
			Price:           45000.0,
			StockStatus:     true,
			IsIndustrial:    true,
			IsHomeAppliance: false,
			ImageURL:        "https://images.unsplash.com/photo-1581092921461-eab62e97a780?auto=format&fit=crop&w=600&q=80",
			Description:     "Kolay montajı ve kompakt yapısıyla bilinen, gaz ve sıvı hatlarında açma-kapama amaçlı kullanılan yüksek kaliteli kelebek vana.",
			Status:          "Aktif",
			Order:           1,
			CreatedAt:       fortyFiveDaysAgo,
			Images: []string{
				"https://images.unsplash.com/photo-1581092921461-eab62e97a780?auto=format&fit=crop&w=600&q=80",
			},
			TechnicalSpecs: []models.TechnicalSpec{
				{Key: "Çap", Value: "DN100"},
				{Key: "Conta", Value: "EPDM / NBR"},
				{Key: "Nominal Basınç", Value: "PN16"},
			},
			FeatureBoxes: []models.FeatureBox{
				{Title: "Kompakt Tasarım", Description: "Hafif ve ince yapısıyla boru hatlarında minimum yer kaplar."},
			},
			ApplicationAreas: "HVAC sistemleri, su arıtma tesisleri, endüstriyel dağıtım hatları.",
		},

		// Pompa Sistemleri
		models.Product{
			ID:              bson.NewObjectID(),
			Name:            "Santrifüj Pompa - Endüstriyel Tip",
			Baslik:          "Santrifüj Pompa - Endüstriyel Tip",
			Isim:            "Santrifüj Pompa - Endüstriyel Tip",
			Category:        "Pompa Sistemleri",
			Kategori:        "Pompa Sistemleri",
			Price:           3500.0,
			StockStatus:     true,
			IsIndustrial:    true,
			IsHomeAppliance: false,
			ImageURL:        "https://images.unsplash.com/photo-1504307651254-35680f356dfd?auto=format&fit=crop&w=600&q=80",
			Description:     "Yüksek debi ihtiyacı olan sanayi tesisleri ve tarımsal sulama sistemleri için tasarlanmış ağır hizmet tipi yatay santrifüj pompa.",
			Status:          "Aktif",
			Order:           2,
			CreatedAt:       fiveDaysAgo,
			Images: []string{
				"https://images.unsplash.com/photo-1504307651254-35680f356dfd?auto=format&fit=crop&w=600&q=80",
			},
			TechnicalSpecs: []models.TechnicalSpec{
				{Key: "Güç", Value: "15 kW (20 HP)"},
				{Key: "Maks. Debi", Value: "120 m³/saat"},
				{Key: "Maks. Basma Yüksekliği", Value: "65 m"},
			},
			FeatureBoxes: []models.FeatureBox{
				{Title: "Yüksek Verim", Description: "Optimize edilmiş çark tasarımı ile maksimum enerji tasarrufu sağlar."},
				{Title: "Ağır Hizmet", Description: "Dökme demir gövde ve paslanmaz çelik mil ile uzun ömürlü kullanım."},
			},
			ApplicationAreas: "Fabrika su temini, basınçlandırma sistemleri, endüstriyel soğutma kuleleri.",
		},
		models.Product{
			ID:              bson.NewObjectID(),
			Name:            "Sirkülasyon Pompası - Frekans Kontrollü",
			Baslik:          "Sirkülasyon Pompası - Frekans Kontrollü",
			Isim:            "Sirkülasyon Pompası - Frekans Kontrollü",
			Category:        "Pompa Sistemleri",
			Kategori:        "Pompa Sistemleri",
			Price:           25000.0,
			StockStatus:     false,
			IsIndustrial:    true,
			IsHomeAppliance: false,
			ImageURL:        "https://images.unsplash.com/photo-1616401784845-180882ba9ba8?auto=format&fit=crop&w=600&q=80",
			Description:     "Isıtma ve soğutma sistemlerinde sirkülasyonu sağlayan, kendinden frekans konvertörlü akıllı sirkülasyon pompası.",
			Status:          "Aktif",
			Order:           3,
			CreatedAt:       fortyFiveDaysAgo,
			Images: []string{
				"https://images.unsplash.com/photo-1616401784845-180882ba9ba8?auto=format&fit=crop&w=600&q=80",
			},
			TechnicalSpecs: []models.TechnicalSpec{
				{Key: "Enerji Sınıfı", Value: "Class A (EEI ≤ 0.20)"},
				{Key: "Güç Kontrolü", Value: "Oransal, Sabit Basınç"},
				{Key: "Maks. Çalışma Basıncı", Value: "10 bar"},
			},
			FeatureBoxes: []models.FeatureBox{
				{Title: "Frekans Konvertörü", Description: "Tesisatın ihtiyacına göre devrini otomatik ayarlayarak elektrik tüketimini azaltır."},
			},
			ApplicationAreas: "Merkezi ısıtma sistemleri, yerden ısıtma sistemleri, klima tesisatları.",
		},

		// Isı Değiştiriciler
		models.Product{
			ID:              bson.NewObjectID(),
			Name:            "Plakalı Eşanjör - Conta Bağlantılı",
			Baslik:          "Plakalı Eşanjör - Conta Bağlantılı",
			Isim:            "Plakalı Eşanjör - Conta Bağlantılı",
			Category:        "Isı Değiştiriciler",
			Kategori:        "Isı Değiştiriciler",
			Price:           15000.0,
			StockStatus:     true,
			IsIndustrial:    false,
			IsHomeAppliance: true,
			ImageURL:        "https://images.unsplash.com/photo-1581094794329-c8112a89af12?auto=format&fit=crop&w=600&q=80",
			Description:     "İki farklı akışkan arasında yüksek verimle ısı transferi sağlayan, plakaları sökülebilir ve kapasitesi artırılabilir conta bağlantılı eşanjör.",
			Status:          "Aktif",
			Order:           4,
			CreatedAt:       fiveDaysAgo,
			Images: []string{
				"https://images.unsplash.com/photo-1581094794329-c8112a89af12?auto=format&fit=crop&w=600&q=80",
			},
			TechnicalSpecs: []models.TechnicalSpec{
				{Key: "Plaka Malzemesi", Value: "AISI 316 Paslanmaz Çelik"},
				{Key: "Conta Malzemesi", Value: "EPDM / Viton"},
				{Key: "Maks. Kapasite", Value: "150 kW"},
			},
			FeatureBoxes: []models.FeatureBox{
				{Title: "Esnek Tasarım", Description: "Gerektiğinde plaka eklenerek kapasitesi kolayca artırılabilir."},
				{Title: "Kolay Bakım", Description: "Saplamalar gevşetilerek plakalar ve contalar kolayca temizlenebilir."},
			},
			ApplicationAreas: "Jeotermal ısıtma, bina altı ısıtma istasyonları, endüstriyel proses soğutma.",
		},

		// Otomasyon
		models.Product{
			ID:              bson.NewObjectID(),
			Name:            "Oransal Vana Motoru - Akıllı Kontrol",
			Baslik:          "Oransal Vana Motoru - Akıllı Kontrol",
			Isim:            "Oransal Vana Motoru - Akıllı Kontrol",
			Category:        "Otomasyon",
			Kategori:        "Otomasyon",
			Price:           3500.0,
			StockStatus:     true,
			IsIndustrial:    true,
			IsHomeAppliance: false,
			ImageURL:        "https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=600&q=80",
			Description:     "Vana açıklık oranını 0-10V veya 4-20mA kontrol sinyali ile hassas şekilde ayarlayan, geri besleme özellikli akıllı motor.",
			Status:          "Aktif",
			Order:           5,
			CreatedAt:       fiveDaysAgo,
			Images: []string{
				"https://images.unsplash.com/photo-1518770660439-4636190af475?auto=format&fit=crop&w=600&q=80",
			},
			TechnicalSpecs: []models.TechnicalSpec{
				{Key: "Kontrol Sinyali", Value: "0-10 V DC, 4-20 mA"},
				{Key: "Besleme Voltajı", Value: "24 V AC/DC"},
				{Key: "Tork", Value: "20 Nm"},
			},
			FeatureBoxes: []models.FeatureBox{
				{Title: "Hassas Pozisyonlama", Description: "Gelişmiş mikroişlemci kontrolü ile vanayı tam açıda konumlandırır."},
			},
			ApplicationAreas: "HVAC otomasyon sistemleri, oransal vana kontrolü, endüstriyel otomasyon.",
		},

		// Yedek Parça
		models.Product{
			ID:          bson.NewObjectID(),
			Name:        "EPDM Sızdırmazlık Contası",
			Baslik:      "EPDM Sızdırmazlık Contası DN100",
			Isim:        "EPDM Sızdırmazlık Contası",
			Category:    "Yedek Parça",
			Kategori:    "Yedek Parça",
			Price:       15.00,
			StockStatus: true,
			ImageURL:    "https://images.unsplash.com/photo-1530124560677-bdaea92b5066?auto=format&fit=crop&w=600&q=80",
			Description: "Kelebek vanalar ve flanş bağlantıları için mükemmel sızdırmazlık ve yüksek elastikiyet sağlayan EPDM conta.",
			Status:      "Aktif",
			Order:       6,
			CreatedAt:   fortyFiveDaysAgo,
			Images: []string{
				"https://images.unsplash.com/photo-1530124560677-bdaea92b5066?auto=format&fit=crop&w=600&q=80",
			},
			TechnicalSpecs: []models.TechnicalSpec{
				{Key: "Sıcaklık Aralığı", Value: "-30°C ila +130°C"},
				{Key: "Malzeme", Value: "EPDM kauçuk"},
				{Key: "Çap Uyumluluğu", Value: "DN100 / PN16"},
			},
			FeatureBoxes: []models.FeatureBox{
				{Title: "Yüksek Direnç", Description: "Sıcak su, buhar, asitler ve alkalilere karşı mükemmel dayanım gösterir."},
			},
			ApplicationAreas: "Su hatları, arıtma tesisleri, genel endüstriyel tesisatlar.",
		},
	}

	_, err = collection.InsertMany(ctx, seedProducts)
	if err != nil {
		log.Println("Urun seed hatası:", err)
	} else {
		log.Println("Toplam", len(seedProducts), "ornek urun basariyla veritabanina yuklendi.")
	}
}

func SeedSettings() {
	collection := GetCollection("settings")
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	count, err := collection.CountDocuments(ctx, bson.M{})
	if err != nil {
		log.Println("Ayar sayisi kontrol edilemedi:", err)
		return
	}
	if count > 0 {
		log.Println("Site ayarlari zaten mevcut. Seed islemi atlandi.")
		return
	}

	def := models.DefaultSettings()
	def.ID = bson.NewObjectID()
	_, err = collection.InsertOne(ctx, def)
	if err != nil {
		log.Println("Ayar seed hatası:", err)
	} else {
		log.Println("Varsayilan site ayarlari olusturuldu.")
	}
}
