package main

import (
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gobuffalo/plush"
	pagetemplate "github.com/innovation-upstream/codema-hub/internal/page-template"
	"github.com/julienschmidt/httprouter"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// Pattern represents the structure for a code pattern.
type Pattern struct {
	Label                            string    `bson:"label"`
	Description                      string    `bson:"description"`
	IsPublic                         bool      `bson:"is_public"`
	UpdatedAt                        time.Time `bson:"updated_at"`
	CreatedAt                        time.Time `bson:"created_at"`
	FunctionImplementationDefinition string    `bson:"function_implementation_definition"`
}

var (
	landingPageTmpl  *plush.Template
	patternPageTmpl  *plush.Template
	bountiesPageTmpl *plush.Template
	mongoClient      *mongo.Client
	minioClient      *minio.Client
)

func init() {
	var err error
	landingPageTmpl, err = plush.NewTemplate(pagetemplate.IndexPageTemplate)
	if err != nil {
		log.Fatalf("Error parsing landing page template: %v", err)
	}

	patternPageTmpl, err = plush.NewTemplate(pagetemplate.PatternDetailsPageTemplate)
	if err != nil {
		log.Fatalf("Error parsing pattern page template: %v", err)
	}

	bountiesPageTmpl, err = plush.NewTemplate(pagetemplate.BountiesPageTemplate)
	if err != nil {
		log.Fatalf("Error parsing bounties page template: %v", err)
	}

	// Initialize MongoDB client
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	mongoClient, err = mongo.Connect(ctx, options.Client().ApplyURI("mongodb://mongodb:27017"))
	if err != nil {
		log.Fatalf("Error connecting to MongoDB: %v", err)
	}

	// Ping the MongoDB server to verify the connection
	err = mongoClient.Ping(ctx, nil)
	if err != nil {
		log.Fatalf("Error pinging MongoDB: %v", err)
	}

	// Initialize MinIO client
	minioClient, err = minio.New("minio:9000", &minio.Options{
		Creds:  credentials.NewStaticV4("minioadmin", "minioadmin", ""),
		Secure: false,
	})
	if err != nil {
		log.Fatalf("Error initializing MinIO client: %v", err)
	}

	// Ensure the "patterns" bucket exists
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	exists, err := minioClient.BucketExists(ctx, "patterns")
	if err != nil {
		log.Fatalf("Error checking if 'patterns' bucket exists: %v", err)
	}
	if !exists {
		err = minioClient.MakeBucket(ctx, "patterns", minio.MakeBucketOptions{})
		if err != nil {
			log.Fatalf("Error creating 'patterns' bucket: %v", err)
		}
	}
}

func main() {
	router := httprouter.New()
	router.GET("/", LandingPageHandler)
	router.GET("/bounties", BountiesPageHandler)
	router.GET("/health", HealthHandler)
	router.GET("/pattern/:patternLabel", PatternPageHandler)
	router.GET("/api/pattern/pull/*patternLabel", PullPatternByLabelHandler)
	router.POST("/api/pattern/publish/*patternLabelWithVersion", PublishPatternHandler)

	router.ServeFiles("/static/*filepath", http.Dir("static"))

	port := 8090
	fmt.Printf("Server running on http://localhost:%d\n", port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", port), router))
}

// LandingPageHandler handles the landing page.
func LandingPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := plush.NewContext()
	tmpl, err := landingPageTmpl.Exec(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// PatternPageHandler handles the pattern page.
func PatternPageHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	patternLabel := ps.ByName("patternLabel")

	// Fetch the pattern from MongoDB
	collection := mongoClient.Database("codema").Collection("patterns")
	var pattern Pattern
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"label": patternLabel}).Decode(&pattern)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Pattern not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching pattern", http.StatusInternalServerError)
		}
		return
	}

	pctx := plush.NewContext()
	pctx.Set("pattern", pattern)

	tmpl, err := patternPageTmpl.Exec(pctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

// BountiesPageHandler handles the bounties page.
func BountiesPageHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	ctx := plush.NewContext()
	tmpl, err := bountiesPageTmpl.Exec(ctx)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Set("Content-Type", "text/html")
	w.Write([]byte(tmpl))
}

func HealthHandler(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	// Check MongoDB connection
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := mongoClient.Ping(ctx, nil)
	if err != nil {
		http.Error(w, "MongoDB connection failed", http.StatusServiceUnavailable)
		return
	}

	// Check MinIO connection
	_, err = minioClient.ListBuckets(ctx)
	if err != nil {
		http.Error(w, "MinIO connection failed", http.StatusServiceUnavailable)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write([]byte("ok"))
}

func PublishPatternHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	patternLabelWithVersion := strings.TrimPrefix(ps.ByName("patternLabelWithVersion"), "/")

	// Validate input
	parts := strings.Split(patternLabelWithVersion, "/")
	if len(parts) != 3 {
		http.Error(w, fmt.Sprintf("Invalid pattern label and version format %s", patternLabelWithVersion), http.StatusBadRequest)
		return
	}
	patternLabel, version := parts[0]+"/"+parts[1], parts[2]

	// Create a new pattern in MongoDB
	pattern := Pattern{
		Label:     patternLabel,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	collection := mongoClient.Database("codema").Collection("patterns")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_, err := collection.InsertOne(ctx, pattern)
	if err != nil {
		http.Error(w, "Error saving pattern to database", http.StatusInternalServerError)
		return
	}

	// Create a temporary file to store the uploaded content
	tempFile, err := os.CreateTemp("", "pattern-*.zip")
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error creating temporary file", http.StatusInternalServerError)
		return
	}
	defer os.Remove(tempFile.Name())

	// Copy the uploaded content to the temporary file
	_, err = io.Copy(tempFile, r.Body)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error reading uploaded content", http.StatusInternalServerError)
		return
	}

	// Rewind the file for reading
	_, err = tempFile.Seek(0, 0)
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error preparing file for upload", http.StatusInternalServerError)
		return
	}

	// Upload the file to MinIO
	objectName := fmt.Sprintf("%s/%s.zip", patternLabel, version)
	_, err = minioClient.PutObject(ctx, "patterns", objectName, tempFile, -1, minio.PutObjectOptions{ContentType: "application/zip"})
	if err != nil {
		fmt.Printf("%+v\n", err)
		http.Error(w, "Error uploading to storage", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	fmt.Fprintf(w, "Pattern %s version %s published successfully", patternLabel, version)
}

func PullPatternByLabelHandler(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	// Parse the pattern label and version from the URL
	fullPath := strings.TrimPrefix(ps.ByName("patternLabel"), "/")
	parts := strings.Split(fullPath, "@")
	if len(parts) != 3 {
		http.Error(w, "Invalid pattern URL format", http.StatusBadRequest)
		return
	}
	patternLabel := "@" + parts[0] + parts[1]
	version := parts[2]

	// Fetch pattern from MongoDB
	collection := mongoClient.Database("codema").Collection("patterns")
	var pattern Pattern
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	err := collection.FindOne(ctx, bson.M{"label": patternLabel}).Decode(&pattern)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			http.Error(w, "Pattern not found", http.StatusNotFound)
		} else {
			http.Error(w, "Error fetching pattern", http.StatusInternalServerError)
		}
		return
	}

	// Construct the object name for the specific version
	objectName := fmt.Sprintf("%s/%s.zip", patternLabel, version)

	// Get the object from MinIO
	ctx, cancel = context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	obj, err := minioClient.GetObject(ctx, "patterns", objectName, minio.GetObjectOptions{})
	if err != nil {
		http.Error(w, "Error fetching pattern file", http.StatusInternalServerError)
		return
	}
	defer obj.Close()

	// Get object info to set the correct Content-Length header
	objInfo, err := obj.Stat()
	if err != nil {
		http.Error(w, "Error getting object info", http.StatusInternalServerError)
		return
	}

	// Set headers for zip file download
	w.Header().Set("Content-Type", "application/zip")
	w.Header().Set("Content-Disposition", fmt.Sprintf("attachment; filename=%s-%s.zip", patternLabel, version))
	w.Header().Set("Content-Length", fmt.Sprintf("%d", objInfo.Size))

	// Stream the zip file to the response
	_, err = io.Copy(w, obj)
	if err != nil {
		http.Error(w, "Error sending pattern files", http.StatusInternalServerError)
		return
	}
}
