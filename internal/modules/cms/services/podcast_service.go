package cms

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/google/uuid"
	"github.com/kkdai/youtube/v2"
	"github.com/ziyadrw/faslah/internal/base"
	podcastDTOs "github.com/ziyadrw/faslah/internal/modules/cms/dtos"
	podcastModels "github.com/ziyadrw/faslah/internal/modules/cms/models"
	podcast "github.com/ziyadrw/faslah/internal/modules/cms/repositories"
)

type PodcastService struct {
	PodcastRepository *podcast.PodcastRepository
}

func NewPodcastService(podcastRepository *podcast.PodcastRepository) *PodcastService {
	return &PodcastService{PodcastRepository: podcastRepository}
}

// CreateContent creates a new podcast from either a YouTube URL or a direct MP4 upload
func (ps *PodcastService) CreateContent(dto podcastDTOs.CreatePodcastRequest, userID string) base.Response {
	if (dto.SourceURL != "" && dto.File != nil) || (dto.SourceURL == "" && dto.File == nil) {
		return base.SetErrorMessage(
			"خطأ في البيانات المدخلة",
			"يجب توفير إما رابط يوتيوب أو ملف MP4، وليس كلاهما أو لا شيء",
		)
	}

	if dto.File != nil && dto.Title == "" {
		return base.SetErrorMessage(
			"خطأ في البيانات المدخلة",
			"يجب توفير العنوان عند رفع ملف MP4",
		)
	}

	var mediaURL string
	var title string
	var description string
	var tags []string
	var durationSecs int
	var sourceURL string

	if dto.SourceURL != "" {
		sourceURL = dto.SourceURL

		youtubeResponse := ps.FetchYouTube(dto.SourceURL)
		if youtubeResponse.HTTPStatus != http.StatusOK {
			return youtubeResponse
		}

		youtubeData, ok := youtubeResponse.Data.(podcastDTOs.YouTubeMetadataResponse)
		if !ok {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في استخراج بيانات الفيديو")
		}

		uniqueFilename := fmt.Sprintf("%s.mp4", uuid.New().String())
		uploadResponse := ps.UploadMedia(youtubeData.VideoFile, uniqueFilename)
		if uploadResponse.HTTPStatus != http.StatusOK {
			return uploadResponse
		}

		uploadData, ok := uploadResponse.Data.(podcastDTOs.MediaUploadResponse)
		if !ok {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في رفع الفيديو")
		}

		mediaURL = uploadData.MediaURL

		if dto.Title == "" {
			title = youtubeData.Title
		} else {
			title = dto.Title
		}

		if dto.Description == "" {
			description = youtubeData.Description
		} else {
			description = dto.Description
		}

		if dto.Tags == nil || len(dto.Tags) == 0 {
		} else {
			tags = dto.Tags
		}

		durationSecs = youtubeData.DurationSecs
	} else {
		src, err := dto.File.Open()
		if err != nil {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في فتح الملف")
		}
		defer src.Close()

		fileContent := make([]byte, dto.File.Size)
		if _, err = src.Read(fileContent); err != nil {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في قراءة محتوى الملف")
		}

		uploadResponse := ps.UploadMedia(fileContent, dto.File.Filename)
		if uploadResponse.HTTPStatus != http.StatusOK {
			return uploadResponse
		}

		uploadData, ok := uploadResponse.Data.(podcastDTOs.MediaUploadResponse)
		if !ok {
			return base.SetErrorMessage("خطأ في الخادم", "فشل في رفع الملف")
		}

		mediaURL = uploadData.MediaURL
		title = dto.Title
		description = dto.Description
		tags = dto.Tags

		duration, err := ps.GetVideoDuration(fileContent)
		if err != nil {
			log.Printf("Failed to extract video duration: %v", err)
			durationSecs = 0
		} else {
			durationSecs = duration
		}
	}

	podcast := &podcastModels.Podcast{
		UserID:       userID,
		Title:        title,
		Description:  description,
		Tags:         tags,
		MediaURL:     mediaURL,
		SourceURL:    sourceURL,
		DurationSecs: durationSecs,
		PublishedAt:  dto.PublishedAt,
	}

	if err := ps.PodcastRepository.CreatePodcast(podcast); err != nil {
		return base.SetErrorMessage("فشل في إنشاء البودكاست", err.Error())
	}

	response := podcastDTOs.MapPodcastToDTO(podcast)
	return base.SetData(response, "تم إنشاء البودكاست بنجاح")
}

// GetContent returns a podcast
func (ps *PodcastService) GetContent(id string) base.Response {
	podcastID, err := uuid.Parse(id)
	if err != nil {
		return base.SetErrorMessage("معرف البودكاست غير صالح", err.Error())
	}

	podcast, err := ps.PodcastRepository.GetPodcastByID(podcastID)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if podcast == nil {
		return base.SetErrorMessage("البودكاست غير موجود", "لم يتم العثور على البودكاست")
	}

	response := podcastDTOs.MapPodcastToDTO(podcast)

	return base.SetData(response)
}

// UpdateContent updates a podcast
func (ps *PodcastService) UpdateContent(id string, dto podcastDTOs.UpdateContentRequest) base.Response {
	podcastID, err := uuid.Parse(id)
	if err != nil {
		return base.SetErrorMessage("معرف البودكاست غير صالح", err.Error())
	}

	podcast, err := ps.PodcastRepository.GetPodcastByID(podcastID)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if podcast == nil {
		return base.SetErrorMessage("البودكاست غير موجود", "لم يتم العثور على البودكاست")
	}

	if dto.Title != "" {
		podcast.Title = dto.Title
	}
	if dto.Description != "" {
		podcast.Description = dto.Description
	}
	if dto.Tags != nil {
		podcast.Tags = dto.Tags
	}
	if dto.MediaURL != "" {
		podcast.MediaURL = dto.MediaURL
	}
	if dto.SourceURL != "" {
		podcast.SourceURL = dto.SourceURL
	}
	if dto.DurationSecs != 0 {
		podcast.DurationSecs = dto.DurationSecs
	}
	if dto.PublishedAt != nil {
		podcast.PublishedAt = dto.PublishedAt
	}

	if err := ps.PodcastRepository.UpdatePodcast(podcast); err != nil {
		return base.SetErrorMessage("فشل في تحديث البودكاست", err.Error())
	}

	response := podcastDTOs.MapPodcastToDTO(podcast)

	return base.SetData(response, "تم تحديث البودكاست بنجاح")
}

func (ps *PodcastService) DeleteContent(id string) base.Response {
	podcastID, err := uuid.Parse(id)
	if err != nil {
		return base.SetErrorMessage("معرف البودكاست غير صالح", err.Error())
	}

	podcast, err := ps.PodcastRepository.GetPodcastByID(podcastID)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}
	if podcast == nil {
		return base.SetErrorMessage("البودكاست غير موجود", "لم يتم العثور على البودكاست")
	}

	if err := ps.PodcastRepository.DeletePodcast(podcastID); err != nil {
		return base.SetErrorMessage("فشل في حذف البودكاست", err.Error())
	}

	return base.SetSuccessMessage("تم حذف البودكاست بنجاح")
}

// GetMyContent retrieves all podcasts created by the current user
func (ps *PodcastService) GetMyContent(userID string) base.Response {
	podcasts, err := ps.PodcastRepository.GetPodcastsByUserID(userID)
	if err != nil {
		return base.SetErrorMessage("خطأ في الخادم", err.Error())
	}

	var response []podcastDTOs.PodcastResponse
	for _, podcast := range podcasts {
		response = append(response, podcastDTOs.MapPodcastToDTO(&podcast))
	}

	return base.SetData(response, "تم استرجاع المحتوى الخاص بك بنجاح")
}
func (ps *PodcastService) UploadMedia(file []byte, filename string) base.Response {
	uniqueFilename := fmt.Sprintf("%s_%s", uuid.New().String(), filename)

	accessKeyID := os.Getenv("R2_ACCESS_KEY_ID")
	secretAccessKey := os.Getenv("R2_SECRET_ACCESS_KEY")
	bucketName := os.Getenv("R2_BUCKET_NAME")
	bucketCode := os.Getenv("R2_BUCKET_CODE")
	r2Endpoint := os.Getenv("R2_ENDPOINT_URL")

	customResolver := aws.EndpointResolverWithOptionsFunc(func(service, region string, options ...interface{}) (aws.Endpoint, error) {
		return aws.Endpoint{
			URL:           r2Endpoint,
			SigningRegion: "auto",
		}, nil
	})

	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithEndpointResolverWithOptions(customResolver),
		config.WithCredentialsProvider(credentials.NewStaticCredentialsProvider(accessKeyID, secretAccessKey, "")),
		config.WithRegion("auto"),
	)
	if err != nil {
		return base.SetErrorMessage("فشل في تهيئة خدمة التخزين", err.Error())
	}
	s3Client := s3.NewFromConfig(cfg)
	_, err = s3Client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket:      aws.String(bucketName),
		Key:         aws.String(uniqueFilename),
		Body:        bytes.NewReader(file),
		ContentType: aws.String(http.DetectContentType(file)),
	})
	if err != nil {
		return base.SetErrorMessage("فشل في رفع الملف", err.Error())
	}
	mediaURL := fmt.Sprintf("https://pub-%s.r2.dev/%s", bucketCode, uniqueFilename)
	log.Println(mediaURL)
	response := podcastDTOs.MediaUploadResponse{
		MediaURL: mediaURL,
	}

	return base.SetData(response, "تم رفع الملف بنجاح")
}

func (ps *PodcastService) FetchYouTube(youtubeURL string) base.Response {
	ytClient := youtube.Client{}

	video, err := ytClient.GetVideo(youtubeURL)
	if err != nil {
		return base.SetErrorMessage("فشل في الحصول على معلومات الفيديو", err.Error())
	}

	title := video.Title
	description := video.Description
	durationSecs := int(video.Duration.Seconds())

	var videoFormat, audioFormat *youtube.Format

	for _, f := range video.Formats {
		if strings.HasPrefix(f.MimeType, "video/mp4") {
			videoFormat = &f
			break
		}
	}

	for _, f := range video.Formats {
		if f.ItagNo == 140 {
			audioFormat = &f
			break
		}
	}
	if audioFormat == nil {
		for _, f := range video.Formats {
			if strings.HasPrefix(f.MimeType, "audio/") {
				audioFormat = &f
				break
			}
		}
	}

	if videoFormat == nil || audioFormat == nil {
		log.Println("Available formats:", video.Formats)
		return base.SetErrorMessage("لم يتم العثور على تنسيقات مناسبة", "")
	}

	tempDir := os.TempDir()
	videoPath := filepath.Join(tempDir, "video.mp4")
	audioPath := filepath.Join(tempDir, "audio.m4a")
	mergedPath := filepath.Join(tempDir, "merged_output.mp4")

	videoStream, _, err := ytClient.GetStream(video, videoFormat)
	if err != nil {
		return base.SetErrorMessage("فشل تحميل الفيديو", err.Error())
	}
	videoBytes, err := io.ReadAll(videoStream)
	if err != nil {
		return base.SetErrorMessage("فشل قراءة بيانات الفيديو", err.Error())
	}
	if err := os.WriteFile(videoPath, videoBytes, 0644); err != nil {
		return base.SetErrorMessage("فشل حفظ الفيديو", err.Error())
	}

	audioStream, _, err := ytClient.GetStream(video, audioFormat)
	if err != nil {
		return base.SetErrorMessage("فشل تحميل الصوت", err.Error())
	}
	audioBytes, err := io.ReadAll(audioStream)
	if err != nil {
		return base.SetErrorMessage("فشل قراءة بيانات الصوت", err.Error())
	}
	if err := os.WriteFile(audioPath, audioBytes, 0644); err != nil {
		return base.SetErrorMessage("فشل حفظ الصوت", err.Error())
	}

	cmd := exec.Command("ffmpeg", "-y", "-i", videoPath, "-i", audioPath, "-c", "copy", mergedPath)
	if err := cmd.Run(); err != nil {
		return base.SetErrorMessage("فشل دمج الفيديو مع الصوت", err.Error())
	}

	mergedData, err := os.ReadFile(mergedPath)
	if err != nil {
		return base.SetErrorMessage("فشل قراءة الفيديو المدمج", err.Error())
	}

	_ = os.Remove(videoPath)
	_ = os.Remove(audioPath)
	_ = os.Remove(mergedPath)

	response := podcastDTOs.YouTubeMetadataResponse{
		VideoFile:    mergedData,
		Title:        title,
		Description:  description,
		DurationSecs: durationSecs,
	}

	return base.SetData(response, "تم جلب الفيديو بنجاح")
}

// GetVideoDuration extracts the duration in seconds from a video file using ffmpeg.
func (ps *PodcastService) GetVideoDuration(fileContent []byte) (int, error) {
	tempDir := os.TempDir()
	tempFile := filepath.Join(tempDir, fmt.Sprintf("temp_video_%s.mp4", uuid.New().String()))
	if err := os.WriteFile(tempFile, fileContent, 0644); err != nil {
		return 0, fmt.Errorf("failed to write temporary file: %w", err)
	}
	defer os.Remove(tempFile)

	cmd := exec.Command("ffmpeg", "-i", tempFile, "-f", "null", "-")
	var stderr bytes.Buffer
	cmd.Stderr = &stderr

	err := cmd.Run()
	output := stderr.String()

	durationStr := extractDuration(output)
	if durationStr != "" {
		seconds := convertDurationToSeconds(durationStr)
		return seconds, nil
	}

	if err != nil {
		return 0, fmt.Errorf("failed to extract duration: %w", err)
	}

	return 0, fmt.Errorf("duration not found in ffmpeg output")
}

// extractDuration extracts the duration string from ffmpeg output
func extractDuration(output string) string {
	durationIndex := strings.Index(output, "Duration: ")
	if durationIndex == -1 {
		return ""
	}
	durationStr := output[durationIndex+10:]
	endIndex := strings.Index(durationStr, ",")
	if endIndex == -1 {
		return ""
	}

	return strings.TrimSpace(durationStr[:endIndex])
}

// convertDurationToSeconds converts a duration string (HH:MM:SS.MS) to seconds
func convertDurationToSeconds(durationStr string) int {
	parts := strings.Split(durationStr, ":")
	if len(parts) != 3 {
		return 0
	}

	hours, err := strconv.ParseFloat(parts[0], 64)
	if err != nil {
		return 0
	}

	minutes, err := strconv.ParseFloat(parts[1], 64)
	if err != nil {
		return 0
	}

	seconds, err := strconv.ParseFloat(parts[2], 64)
	if err != nil {
		return 0
	}

	totalSeconds := int(hours*3600 + minutes*60 + seconds)
	return totalSeconds
}
