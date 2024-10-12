package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	gomail "gopkg.in/gomail.v2"
)

// ContactForm represents the structure of the contact form data.
type ContactForm struct {
	Name        string `json:"name"`
	Email       string `json:"email"`
	Mobile      string `json:"mobile"`
	City        string `json:"city"`
	Country     string `json:"country"`
	Message     string `json:"message"`
	CountryCode string `json:"countryCode"`
	UserAgent   string `json:"userAgent"`
	Referrer    string `json:"referrer"`
	Timestamp   string `json:"timestamp"`
}

// MongoDB client and collection
var client *mongo.Client
var collection *mongo.Collection

func init() {
	// MongoDB connection
	var err error
	clientOptions := options.Client().ApplyURI("mongodb+srv://admin:AdMiN@shaik-baleeghuddin-kash.k33e64f.mongodb.net/")
	client, err = mongo.Connect(context.TODO(), clientOptions)
	if err != nil {
		log.Fatal(err)
	}
	collection = client.Database("contact_form").Collection("submissions")
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintln(w, "Hello, World!")
}

func main() {
	http.HandleFunc("/", handler)
	http.HandleFunc("/contactformsubmits", cors(handleContactFormSubmit))

	log.Println("Server is starting on port 8080...")
	if err := http.ListenAndServe(":8080", nil); err != nil {
		log.Fatal(err)
	}
}

// CORS middleware with restricted origins
func cors(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		allowedOrigins := map[string]bool{
			"http://192.168.1.15:5173": true,
			"https://www.sbkashif.com": true,
		}

		origin := r.Header.Get("Origin")
		if allowedOrigins[origin] {
			w.Header().Set("Access-Control-Allow-Origin", origin)           // Allow specific origin
			w.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS") // Allow specific methods
			w.Header().Set("Access-Control-Allow-Headers", "Content-Type")  // Allow specific headers
		}

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next(w, r)
	}
}

// handleContactFormSubmit handles the POST request to /contactformsubmits
func handleContactFormSubmit(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var contactForm ContactForm

	// Parse JSON body
	if err := json.NewDecoder(r.Body).Decode(&contactForm); err != nil {
		log.Println("Error decoding JSON:", err) // Log error
		http.Error(w, "Invalid request format", http.StatusBadRequest)
		return
	}
	defer r.Body.Close()

	// Insert into MongoDB
	_, err := collection.InsertOne(context.TODO(), bson.D{
		{Key: "name", Value: contactForm.Name},
		{Key: "email", Value: contactForm.Email},
		{Key: "mobile", Value: contactForm.Mobile},
		{Key: "city", Value: contactForm.City},
		{Key: "country", Value: contactForm.Country},
		{Key: "message", Value: contactForm.Message},
		{Key: "countryCode", Value: contactForm.CountryCode},
		{Key: "userAgent", Value: contactForm.UserAgent},
		{Key: "referrer", Value: contactForm.Referrer},
		{Key: "timestamp", Value: contactForm.Timestamp},
	})

	if err != nil {
		log.Println("Error inserting into MongoDB:", err) // Log error
		http.Error(w, "Failed to save contact form data", http.StatusInternalServerError)
		return
	}

	// Send email
	if err := sendEmail(contactForm); err != nil {
		log.Println("Error sending email:", err) // Log error
		http.Error(w, "Failed to send email", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Contact form submitted successfully"})
}

// sendEmail sends an email with the contact form message
func sendEmail(contactForm ContactForm) error {
	m := gomail.NewMessage()
	m.SetHeader("From", "shaikbaleeghuddin@gmail.com")
	m.SetHeader("To", "shaikbaleeghuddin@gmail.com")
	m.SetHeader("Cc", contactForm.Email) // CC the user's email
	m.SetHeader("Subject", "New Contact Form Submission from  "+contactForm.Name)
	m.SetBody("text/plain", contactForm.Message)

	d := gomail.NewDialer("smtp.gmail.com", 587, "shaikbaleeghuddin@gmail.com", "fecd lvhs vtrv uotf")

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		return err
	}
	return nil
}
