# Go CRUD Işlemleri

Bu uygulama Patika ve Picus Security iş birliğinde gerçekleşen Golang Backend Web Development Bootcamp kapsamındaki ilk ödevi içermektedir.  

 

<br>  
<br>  

# İçindekiler
- [Uygulama İsterleri](#uygulama-i%CC%87sterleri)  
- [Uygulama Öncesi Hazırlık](#uygulama-%C3%B6ncesi-haz%C4%B1rl%C4%B1k)  
- [Uygulama Aşamaları](#uygulama-a%C5%9Famalar%C4%B1)  


<br>  
<br>  

# Uygulama İsterleri
- Kullanıcı kitap oluşturma , arama , silme ve güncelleme işlemlerini yapacaktır.
- Kitap oluşturma işleminde öncellikle bir token almalıdır.
- Yapılan işlemler sırasında düzgün status kodları gösterilecektir. 


<br>  
<br>  


# Uygulama Öncesi Hazırlık

- **Proje içerisinde kullanılan yapılar :** if/else for struct db connection 

- **Proje içeriğinde kullanılan paketler:** `github.com/gorilla/mux` , `github.com/rs/cors` , `github.com/dgrijalva/jwt-go` , `gorm.io/gorm` , `gorm.io/driver/postgres`


<br>  
<br>  


# Uygulama Aşamaları

- Öncellikle kitapları ve yazarları tuttuğumuz birer struct oluşturuyoruz.

```
type Books struct {
	gorm.Model
	StockNumber int     `json:"stock_num"`
	PageNumber  int     `json:"page_num"`
	Price       float64 `json:"price"`
	Name        string  `json:"book_name"`
	StockCode   string  `json:"stock_code"`
	Isbn        string  `json:"isbn"`
	AuthorName  string  `json:"author_name"`
}


type Author struct {
	gorm.Model
	AuthorName string `json:"author_name"`
}


```


- Bu işlemlerden sonra DB bağlantımızı gerçekleştiriyoruz.

```
const (
	host     = "localhost"
	port     = 5432
	user     = "postgres"
	password = "<password>"
	dbname   = "<your_DB>"
)

func InitialMigration() {

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+
		"password=%s dbname=%s sslmode=disable",
		host, port, user, password, dbname)

	DB, err = gorm.Open(postgres.Open(psqlInfo), &gorm.Config{})

	if err != nil {
		panic(err)
	}

	DB.AutoMigrate(&models.Books{})

}

```


- Ardından listeleme , silme ve update işlemlerini gerçekleştiriyoruz.

```

func GetAllBook() ([]models.Books, error) {

	var book []models.Books

	result := migration.DB.Find(&book)

	if result.Error != nil {
		return book, result.Error
	}

	return book, nil

}

func GetBookByName(name string) ([]models.Books, error) {

	var books []models.Books

	result := migration.DB.Where(" Name LIKE ?", "%"+name+"%").Find(&books)
	if result.Error != nil {
		return books, result.Error
	}

	return books, nil
}

func UpdateStock(id, stock int) ([]models.Books, error) {
	var books []models.Books
	var stoc_num int

	migration.DB.Find(&books)

	result := migration.DB.Table("books").Select("stock_number").Where("id = ?", id).Scan(&stoc_num)

	if stoc_num <= 0 || stoc_num < stock {
		return books, result.Error
	}

	if result.Error != nil {
		return books, result.Error
	}

	newStock := stoc_num - stock

	migration.DB.Model(&books).Where("id = ?", id).Update("stock_number", newStock)

	return books, nil
}

func DeleteBookById(id int) ([]models.Books, error) {
	var books []models.Books

	var n []int

	migration.DB.Model(&books).Pluck("id", &n)

	migration.DB.Unscoped().Delete(&books, id)

	isDeleted := helper.CheckSlice(n, id)

	if !isDeleted {

		return books, errors.New("id not found")
	}

	migration.DB.Find(&books)

	return books, nil
}
```


- Kullanıcı bir kitap oluşturmak isterse öncellikle `http://localhost:8080/token` adresinden token almalıdır.
<br>
<br>

```
func GenerateJWT() (string, error) {

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)

	claims["authorized"] = true
	claims["exp"] = time.Now().Add(time.Minute * 15).Unix()

	tokenString, err := token.SignedString(myKey)

	if err != nil {
		fmt.Errorf("something went wrong %s", err.Error())
		return "", err
	}

	return tokenString, nil

}


```

- Yukarıda `github.com/dgrijalva/jwt-go` paketini kullanarak bir jwt üretiyoruz. Bu fonksiyon bu jwt'yi return ediyor.
-  Ardından bu tokenı kullanıcının "GET" methodu ile almasını sağlıyoruz.

```
//? http://localhost:8080/token
func GetToken(w http.ResponseWriter, r *http.Request) {
	validToken, err := GenerateJWT()

	if err != nil {
		fmt.Fprint(w, err.Error())
	}

	fmt.Fprint(w, validToken)
}
```


- Tokenın header kısmına yazılıp doğruluğunu kontrol etmek için bir middleware oluşturuyoruz ve main içinde tanımladığımız `isTokenValid(controller.AddBook)` ile eğer token valid değil ise hata alıyoruz eğer valid ise işlemlerimize devam ediyoruz.
<br>
<br>

```

var myKey = []byte("yourname")

func isTokenValid(endpoint func(http.ResponseWriter, *http.Request)) http.Handler {

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] != nil {

			token, err := jwt.Parse(r.Header["Token"][0], func(token *jwt.Token) (interface{}, error) {

				if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
					return nil, fmt.Errorf("error")
				}

				return myKey, nil

			})

			if err != nil {
				fmt.Fprintf(w, err.Error())
			}

			if token.Valid {
				endpoint(w, r)
			}

		} else {
			http.Error(w, "not authorized", http.StatusUnauthorized)
		}

	})

}

func main() {

	migration.InitialMigration()

	r := mux.NewRouter()

	r.Handle("/create/book", isTokenValid(controller.AddBook)).Methods("POST")
	
}

```

- Son olarak kullanıcı aşağıda belirtilen urlleri kullanarak işlem yapabilir.

1- Create : http://localhost:8080/create/book

2- List All Book : http://localhost:8080/Allbooks

3- List Book by Name : http://localhost/books/{name}

4- Buy Book : http://localhost/books/buy/{id}/{stock}

5- Delete Book : http://localhost/books/delete/{id}

6- Get Token : http://localhost:8080/token




