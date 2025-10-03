package api

import (
	"encoding/json"
	"fmt"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/gofiber/fiber/v2"
	"nadhi.dev/sarvar/fun/ai"
	"nadhi.dev/sarvar/fun/auth"
	vela "nadhi.dev/sarvar/fun/bucket"
	"nadhi.dev/sarvar/fun/server"
	sheet "nadhi.dev/sarvar/fun/sheets"
)

// Sheet represents the sheet data structure
type Sheet struct {
	Subject     string   `json:"subject"`
	Course      string   `json:"course"`
	Description string   `json:"description"`
	Tags        []string `json:"tags"`
	Visibility  string   `json:"visibility"`
}

// Last request timestamps for cooldown
var lastRequestTimes = make(map[string]time.Time)

// sheetGen is the global sheet generator instance
var sheetGen *sheet.SheetGenerator

// SheetsIndex registers all sheet related routes
func SheetsIndex() error {
	// Initialize GlobalSheetGenerator if it's nil
	if sheet.GlobalSheetGenerator == nil {
		var err error
		sheet.GlobalSheetGenerator, err = sheet.NewSheetGenerator(nil, "./queue_data", 2)
		if err != nil {
			log.Printf("Failed to initialize GlobalSheetGenerator: %v", err)
			return err
		}
		log.Printf("GlobalSheetGenerator initialized successfully")
	}

	server.Route.Post("/api/v1/sheets/generate-tags", generateTags)
	server.Route.Post("/api/v1/sheets/generate-subject", generateSubject)
	server.Route.Post("/api/v1/sheets/generate-course", generateCourse)
	server.Route.Post("/api/v1/sheets/generate-description", generateDescription)
	server.Route.Post("/api/v1/sheets/queue/:id", func(c *fiber.Ctx) error {
    id := c.Params("id")
    if id == "" {
        return c.Status(400).JSON(fiber.Map{"error": "missing job id"})
    }
    if sheet.GlobalSheetGenerator == nil || sheet.GlobalSheetGenerator.Queue == nil {
        return c.Status(500).JSON(fiber.Map{"error": "Sheet queue not initialized"})
    }
    err := sheet.GlobalSheetGenerator.Queue.DeleteJob(id)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": err.Error()})
    }
    return c.JSON(fiber.Map{"status": "deleted"})
})


	server.Route.Get("/api/v1/sheets/get", func(c *fiber.Ctx) error {
    // Query params
    search := c.Query("search", "")
    latest := c.Query("latest", "true") == "true"
    objNumStr := c.Query("obj_num", "10")
    objNum, err := strconv.Atoi(objNumStr)
    if err != nil || objNum <= 0 {
        objNum = 10
    }

    queuePath := "./queue_data/queue.json"
    items, err := vela.GetQueueItems(queuePath, latest, objNum, search)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to get queue items"})
    }
    return c.JSON(items)
})

	
server.Route.Post("/api/v1/sheets/create", func(c *fiber.Ctx) error {
    var req struct {
        Subject             string `json:"subject"`
        Course              string `json:"course"`
        Description         string `json:"description"`
        Tags                string `json:"tags"`
        Curriculum          string `json:"curriculum"`
        SpecialInstructions string `json:"specialInstructions"`
        Visibility          string `json:"visibility"`
    }
    if err := c.BodyParser(&req); err != nil {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request"})
    }

    // Validate required fields
    if req.Subject == "" || req.Course == "" || req.Description == "" || req.Tags == "" || req.Curriculum == "" || req.SpecialInstructions == "" || req.Visibility == "" {
        return c.Status(400).JSON(fiber.Map{"error": "Invalid request: missing required fields"})
    }

    // Extract and validate session
    authHeader := c.Get("Authorization")
    if len(authHeader) < 8 || !strings.HasPrefix(authHeader, "Bearer ") {
        return c.Status(401).JSON(fiber.Map{"error": "missing or invalid authorization header"})
    }
    sessionID := authHeader[7:]
    valid, err := auth.IsSessionValid(sessionID)
    if err != nil || !valid {
        return c.Status(401).JSON(fiber.Map{"error": "invalid session"})
    }

    // Get user from session
    user, err := auth.GetUserBySession(sessionID)
    if err != nil {
        return c.Status(401).JSON(fiber.Map{"error": "user not found or session invalid"})
    }
    userID := user.Username

    // Create a proper GenerationRequest
    genRequest := &ai.GenerationRequest{
        Subject:             req.Subject,
        Course:              req.Course,
        Description:         req.Description,
        Tags:                strings.Split(req.Tags, ","), // Convert comma-separated string to slice
        Curriculum:          req.Curriculum,
        SpecialInstructions: req.SpecialInstructions,
    }

    // Double-check GlobalSheetGenerator is not nil
    if sheet.GlobalSheetGenerator == nil {
        return c.Status(500).JSON(fiber.Map{"error": "Sheet generator not initialized"})
    }

    jobID, err := sheet.GlobalSheetGenerator.CreateSheet(userID, genRequest)
    if err != nil {
        return c.Status(500).JSON(fiber.Map{"error": "Failed to enqueue sheet"})
    }

    return c.JSON(fiber.Map{"jobId": jobID, "status": "queued"})
})


	server.Route.Get("/api/v1/sheets/queue", func(c *fiber.Ctx) error {
		userID := c.Locals("username")
		if userID == nil {
			userID = "anonymous"
		}
		jobs, err := sheetGen.GetUserJobs(userID.(string))
		if err != nil {
			return c.Status(500).JSON(fiber.Map{"error": "Failed to get queue"})
		}
		return c.JSON(jobs)
	})

	return nil
}

// getGeminiKey retrieves the Gemini API key
func getGeminiKey() string {
	return "AIzaSyAcjAEDWdKqbOJBGfkkwd6TkS6XdqHcdcM"
}

// getGeminiModel returns the model to use
func getGeminiModel() string {
	return "gemini-2.5-flash-lite"
}

// getCooldown returns the cooldown time in seconds
func getCooldown() int {
	return 2
}

// checkCooldown checks if the cooldown period has passed for a given endpoint
func checkCooldown(endpoint string) bool {
	cooldown := getCooldown()
	lastTime, exists := lastRequestTimes[endpoint]
	if !exists {
		lastRequestTimes[endpoint] = time.Now()
		return true
	}

	if time.Since(lastTime).Seconds() < float64(cooldown) {
		return false
	}

	lastRequestTimes[endpoint] = time.Now()
	return true
}

// extractTags extracts tags from a response
func extractTags(response string) ([]string, error) {
	var tags []string

	err := json.Unmarshal([]byte(response), &tags)
	if err != nil {
		// Try to extract JSON array from text
		startIdx := strings.Index(response, "[")
		endIdx := strings.LastIndex(response, "]")
		if startIdx >= 0 && endIdx > startIdx {
			jsonStr := response[startIdx : endIdx+1]
			err = json.Unmarshal([]byte(jsonStr), &tags)
			if err != nil {
				// As a fallback, split by commas and clean up
				cleanResponse := strings.Trim(response, "[]\" \n")
				tags = strings.Split(cleanResponse, ",")
				for i, tag := range tags {
					tags[i] = strings.Trim(tag, "\" ")
				}
			}
		}
	}

	return tags, nil
}

// generateTags handles requests to generate tags using AI
func generateTags(c *fiber.Ctx) error {
	// Check cooldown
	if !checkCooldown("tags") {
		return c.Status(429).JSON(fiber.Map{"error": "Too many requests, please wait"})
	}

	var sheet Sheet
	if err := c.BodyParser(&sheet); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request data"})
	}

	if sheet.Subject == "" && sheet.Course == "" && sheet.Description == "" {
		return c.Status(400).JSON(fiber.Map{"error": "At least one of subject, course, or description is required"})
	}

	systemPrompt := `You are a tag generator for educational content. 
Your task is to generate 3-7 relevant tags based on the subject, course title, and description provided.
Return ONLY a JSON array of strings with the tags, nothing else.
Example response: ["mathematics", "algebra", "equations", "polynomials"]`

	userPrompt := fmt.Sprintf(`Generate tags for the following educational content:
Subject: %s
Course: %s
Description: %s`,
		sheet.Subject,
		sheet.Course,
		sheet.Description)

	apiKey := getGeminiKey()
	model := getGeminiModel()

	response, err := ai.GenerateResponse(apiKey, model, systemPrompt, userPrompt, getCooldown())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to generate tags: %v", err)})
	}

	tags, _ := extractTags(response)
	return c.Status(200).JSON(fiber.Map{"tags": tags})
}

// generateSubject generates a subject based on course and/or description
func generateSubject(c *fiber.Ctx) error {
	// Check cooldown
	if !checkCooldown("subject") {
		return c.Status(429).JSON(fiber.Map{"error": "Too many requests, please wait"})
	}

	var request struct {
		Course       string `json:"course"`
		Description  string `json:"description"`
		GenerateTags bool   `json:"generateTags"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request data"})
	}

	if request.Course == "" && request.Description == "" {
		return c.Status(400).JSON(fiber.Map{"error": "At least course or description is required"})
	}

	systemPrompt := `You are an educational content creator. 
Based on the course title and description provided, generate an appropriate subject field.
Return ONLY the subject name, nothing else. Keep it concise (1-3 words).`

	userPrompt := fmt.Sprintf(`Generate a subject name for the following course:
Course: %s
Description: %s`,
		request.Course,
		request.Description)

	apiKey := getGeminiKey()
	model := getGeminiModel()

	response, err := ai.GenerateResponse(apiKey, model, systemPrompt, userPrompt, getCooldown())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to generate subject: %v", err)})
	}

	// Clean the response
	subject := strings.Trim(response, " \n\"")

	result := fiber.Map{"subject": subject}

	// Generate tags only if requested AND the tags query param is set to true
	if request.GenerateTags && c.Query("tags") == "true" {
		if checkCooldown("subject_tags") {
			tagSystemPrompt := `Generate 3-5 tags for this academic subject. Return only a JSON array of strings.
Example: ["physics", "mechanics", "motion"]`

			tagUserPrompt := fmt.Sprintf("Subject: %s\nCourse: %s\nDescription: %s",
				subject, request.Course, request.Description)

			tagResponse, err := ai.GenerateResponse(apiKey, model, tagSystemPrompt, tagUserPrompt, 0)
			if err == nil {
				tags, _ := extractTags(tagResponse)
				result["tags"] = tags
			}
		}
	}

	return c.Status(200).JSON(result)
}

// generateCourse generates a course title based on subject and/or description
func generateCourse(c *fiber.Ctx) error {
	// Check cooldown
	if !checkCooldown("course") {
		return c.Status(429).JSON(fiber.Map{"error": "Too many requests, please wait"})
	}

	var request struct {
		Subject      string `json:"subject"`
		Description  string `json:"description"`
		GenerateTags bool   `json:"generateTags"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request data"})
	}

	if request.Subject == "" && request.Description == "" {
		return c.Status(400).JSON(fiber.Map{"error": "At least subject or description is required"})
	}

	systemPrompt := `You are an educational content creator. 
Based on the subject and description provided, generate an appropriate course title.
Return ONLY the course title, nothing else. Make it sound like an actual academic course.`

	userPrompt := fmt.Sprintf(`Generate a course title for the following:
Subject: %s
Description: %s`,
		request.Subject,
		request.Description)

	apiKey := getGeminiKey()
	model := getGeminiModel()

	response, err := ai.GenerateResponse(apiKey, model, systemPrompt, userPrompt, getCooldown())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to generate course: %v", err)})
	}

	// Clean the response
	course := strings.Trim(response, " \n\"")

	result := fiber.Map{"course": course}

	// Generate tags only if requested AND the tags query param is set to true
	if request.GenerateTags && c.Query("tags") == "true" {
		if checkCooldown("course_tags") {
			tagSystemPrompt := `Generate 3-5 tags for this academic course. Return only a JSON array of strings.
Example: ["calculus", "mathematics", "derivatives"]`

			tagUserPrompt := fmt.Sprintf("Subject: %s\nCourse: %s\nDescription: %s",
				request.Subject, course, request.Description)

			tagResponse, err := ai.GenerateResponse(apiKey, model, tagSystemPrompt, tagUserPrompt, 0)
			if err == nil {
				tags, _ := extractTags(tagResponse)
				result["tags"] = tags
			}
		}
	}

	return c.Status(200).JSON(result)
}

// generateDescription generates a description based on subject and/or course
func generateDescription(c *fiber.Ctx) error {
	// Check cooldown
	if !checkCooldown("description") {
		return c.Status(429).JSON(fiber.Map{"error": "Too many requests, please wait"})
	}

	var request struct {
		Subject      string `json:"subject"`
		Course       string `json:"course"`
		GenerateTags bool   `json:"generateTags"`
	}

	if err := c.BodyParser(&request); err != nil {
		return c.Status(400).JSON(fiber.Map{"error": "Invalid request data"})
	}

	if request.Subject == "" && request.Course == "" {
		return c.Status(400).JSON(fiber.Map{"error": "At least subject or course is required"})
	}

	systemPrompt := `You are an educational content creator. 
Based on the subject and course title provided, generate an appropriate description.
The description should be 2-3 sentences that explain what the course covers.`

	userPrompt := fmt.Sprintf(`Generate a description for the following course:
Subject: %s
Course: %s
Make an apporiate description with the instructions on how to prepare for the course and create an exam course.`,
		request.Subject,
		request.Course)

	apiKey := getGeminiKey()
	model := getGeminiModel()

	response, err := ai.GenerateResponse(apiKey, model, systemPrompt, userPrompt, getCooldown())
	if err != nil {
		return c.Status(500).JSON(fiber.Map{"error": fmt.Sprintf("Failed to generate description: %v", err)})
	}

	// Clean the response
	description := strings.Trim(response, " \n\"")

	result := fiber.Map{"description": description}

	// Generate tags only if requested AND the tags query param is set to true
	if request.GenerateTags && c.Query("tags") == "true" {
		if checkCooldown("description_tags") {
			tagSystemPrompt := `Generate 3-5 tags for this course description. Return only a JSON array of strings.
Example: ["chemistry", "organic", "synthesis"]`

			tagUserPrompt := fmt.Sprintf("Subject: %s\nCourse: %s\nDescription: %s",
				request.Subject, request.Course, description)

			tagResponse, err := ai.GenerateResponse(apiKey, model, tagSystemPrompt, tagUserPrompt, 0)
			if err == nil {
				tags, _ := extractTags(tagResponse)
				result["tags"] = tags
			}
		}
	}

	return c.Status(200).JSON(result)
}
