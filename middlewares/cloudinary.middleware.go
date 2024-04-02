package middlewares

import (
	"context"
	"os"

	"github.com/cloudinary/cloudinary-go"
	"github.com/cloudinary/cloudinary-go/api/uploader"
	"github.com/joho/godotenv"
)

func setupCloudinary() (*cloudinary.Cloudinary, string, error) {
	envErr := godotenv.Load(".env")
	if envErr != nil {
		return nil, "", envErr
	}
	// cloudinary setup
	cloud_name := os.Getenv("CLOUD_NAME")
	api_key := os.Getenv("CLOUD_API_KEY")
	api_secret := os.Getenv("CLOUD_API_SECRET")
	cld, cldErr := cloudinary.NewFromParams(cloud_name, api_key, api_secret)
	if cldErr != nil {
		return nil, "", cldErr
	}
	cloud_folder := os.Getenv("CLOUD_FOLDER")
	return cld, cloud_folder, nil
}

func UploadFile(file interface{}) (string, string, error) {
	cld, folder, cldErr := setupCloudinary()
	if cldErr != nil {
		return "", "", cldErr
	}
	res, uplErr := cld.Upload.Upload(context.Background(), file, uploader.UploadParams{Folder: folder})
	if uplErr != nil {
		return "", "", uplErr
	}
	return res.URL, res.PublicID, nil
}

func DeleteImageFile(publicID string) error {
	cld, _, cldErr := setupCloudinary()
	if cldErr != nil {
		return cldErr
	}
	_, desErr := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID, ResourceType: "image"})
	return desErr
}

func DeleteVideoFile(publicID string) error {
	cld, _, cldErr := setupCloudinary()
	if cldErr != nil {
		return cldErr
	}
	_, desErr := cld.Upload.Destroy(context.Background(), uploader.DestroyParams{PublicID: publicID, ResourceType: "video"})
	return desErr
}
