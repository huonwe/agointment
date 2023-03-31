package main

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

var nowDate = time.Now().Format("2006-01-02 15")
var jwtSecret = []byte(fmt.Sprintf("%v%v", nowDate, "dF13ayA"))

type MapClaims map[string]interface{}
type StrStr map[string]string

// Claim是一些实体（通常指的用户）的状态和额外的元数据
type Claims struct {
	Username string `json:"username"`
	UserID   uint   `json:"password"`
	jwt.StandardClaims
}

// 根据用户的用户名和ID产生token
func GenerateToken(userID uint, username string) (string, error) {
	//设置token有效时间
	nowTime := time.Now()
	expireTime := nowTime.Add(36 * time.Hour)

	claims := Claims{
		UserID:   userID,
		Username: username,
		StandardClaims: jwt.StandardClaims{
			// 过期时间
			ExpiresAt: expireTime.Unix(),
			// 指定token发行人
			// Issuer: "gin-blog",
		},
	}

	tokenClaims := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	//该方法内部生成签名字符串，再用于获取完整、已签名的token
	token, err := tokenClaims.SignedString(jwtSecret)
	return token, err
}

// 根据传入的token值获取到Claims对象信息，（进而获取其中的用户名和密码）
func ParseToken(token string) (*Claims, error) {

	//用于解析鉴权的声明，方法内部主要是具体的解码和校验的过程，最终返回*Token
	tokenClaims, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		return jwtSecret, nil
	})

	if tokenClaims != nil {
		// 从tokenClaims中获取到Claims对象，并使用断言，将该对象转换为我们自己定义的Claims
		// 要传入指针，项目中结构体都是用指针传递，节省空间。
		if claims, ok := tokenClaims.Claims.(*Claims); ok && tokenClaims.Valid {
			return claims, nil
		}
	}
	return nil, err
}

// func postReader(req *http.Request) map[string]string {
// 	con, _ := io.ReadAll(req.Body)
// 	data := make(StrStr)
// 	_ = json.Unmarshal(con, &data)
// 	return data
// }

func handle(err error) {
	if err != nil {
		log.Panic(err)
	}
}

func handle_resp(err error, ctx *gin.Context) {
	if err != nil {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "Failed",
			"msg":    "error",
		})
		log.Panicln(err)
	}
}

func str2uint(id string) uint {
	ID, err := strconv.ParseUint(id, 10, 32)
	handle(err)
	return uint(ID)
}

func str2int(i string) int {
	I, err := strconv.ParseInt(i, 10, 32)
	handle(err)
	return int(I)
}

func now() string {
	return time.Now().In(time.FixedZone("CST", 8*3600)).Format("2006-01-02 15:04")
}

func order_desc_createdAt(db *gorm.DB) *gorm.DB {
	return db.Order("created_at desc")
}

func Alphabet2Num(col rune) int {
	return int(col - 'A')
}

func Num2Col(num int) string {
	col := 'A'
	target := col + int32(num)
	return string(target)
}
