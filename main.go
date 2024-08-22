package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"math"
	"net/http"
	"net/url"
	"os"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
)

var db = make(map[string]string)

func setupRouter() *gin.Engine {
	r := gin.Default()

	r.LoadHTMLGlob("static/*")
	r.Static("/static", "./static")

	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusFound, "/static/index.html")
	})

	r.Static("/weather", "./weather_app/build")

	// Weather Comparison API
	r.GET("/api/weather/:location", func(c *gin.Context) {
		locationName := c.Params.ByName("location")
		openWeatherMapKey := os.Getenv("OPEN_WEATHER_MAP")

		locationURL := "http://api.openweathermap.org/geo/1.0/direct?q=" + url.QueryEscape(locationName) + "&limit=1&appid=" + openWeatherMapKey
		res, err := http.Get(locationURL)
		if err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "request to api failed"})
			return
		}
		resData, err := io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "can't parse api response"})
			return
		}
		var resLocationDataJSON []struct {
			Lat float64
			Lon float64
		}
		json.Unmarshal(resData, &resLocationDataJSON)
		if len(resLocationDataJSON) == 0 {
			c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "can't parse api response"})
			return
		}
		latitude := strconv.FormatFloat(resLocationDataJSON[0].Lat, 'f', -1, 64)
		longitude := strconv.FormatFloat(resLocationDataJSON[0].Lon, 'f', -1, 64)

		weatherURL := "https://api.openweathermap.org/data/2.5/weather?lat=" + latitude + "&lon=" + longitude + "&appid=" + openWeatherMapKey
		res, err = http.Get(weatherURL)
		if err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "request to api failed"})
			return
		}
		resData, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "can't parse api response"})
			return
		}
		var resCityDataJSON struct {
			Name string
			Main struct {
				Temp float64
			}
			Sys struct {
				Country string
			}
		}
		json.Unmarshal(resData, &resCityDataJSON)
		cityName := resCityDataJSON.Name
		countryName := resCityDataJSON.Sys.Country
		cityTemp := resCityDataJSON.Main.Temp - 273.15

		forecastURL := "https://api.openweathermap.org/data/2.5/forecast?lat=" + latitude + "&lon=" + longitude + "&appid=" + openWeatherMapKey
		res, err = http.Get(forecastURL)
		if err != nil {
			fmt.Println(err.Error())
			c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "request to api failed"})
			return
		}
		resData, err = io.ReadAll(res.Body)
		if err != nil {
			log.Fatal(err)
			c.JSON(http.StatusServiceUnavailable, gin.H{"msg": "can't parse api response"})
			return
		}
		var resForecaseDataJSON struct {
			List []struct {
				Main struct {
					Temp float64
				}
			}
		}
		json.Unmarshal(resData, &resForecaseDataJSON)
		var forecastTemps []int
		for _, elem := range resForecaseDataJSON.List {
			forecastTemps = append(forecastTemps, int(math.Round(elem.Main.Temp-273.15)))
		}
		resJSON := gin.H{
			"city_name":      cityName,
			"country_name":   countryName,
			"city_temp":      cityTemp,
			"forecast_temps": forecastTemps,
		}
		fmt.Println(resJSON)
		c.JSON(http.StatusOK, resJSON)

	})

	// Ping test
	r.GET("/api/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong")
	})

	// Get user value
	r.GET("/api/user/:name", func(c *gin.Context) {
		user := c.Params.ByName("name")
		value, ok := db[user]
		if ok {
			c.JSON(http.StatusOK, gin.H{"user": user, "value": value})
		} else {
			c.JSON(http.StatusOK, gin.H{"user": user, "status": "no value"})
		}
	})

	// Authorized group (uses gin.BasicAuth() middleware)
	// Same than:
	// authorized := r.Group("/")
	// authorized.Use(gin.BasicAuth(gin.Credentials{
	//	  "foo":  "bar",
	//	  "manu": "123",
	//}))
	// authorized := r.Group("/", gin.BasicAuth(gin.Accounts{
	// 	"foo":  "bar", // user:foo password:bar
	// 	"manu": "123", // user:manu password:123
	// }))

	/* example curl for /admin with basicauth header
	   Zm9vOmJhcg== is base64("foo:bar")

		curl -X POST \
	  	http://localhost:8080/admin \
	  	-H 'authorization: Basic Zm9vOmJhcg==' \
	  	-H 'content-type: application/json' \
	  	-d '{"value":"bar"}'
	*/
	// authorized.POST("admin", func(c *gin.Context) {
	// 	user := c.MustGet(gin.AuthUserKey).(string)

	// 	// Parse JSON
	// 	var json struct {
	// 		Value string `json:"value" binding:"required"`
	// 	}

	// 	if c.Bind(&json) == nil {
	// 		db[user] = json.Value
	// 		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	// 	}
	// })

	return r
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	r := setupRouter()
	// Listen and Server in 0.0.0.0:8080
	r.Run(":8080")
}
