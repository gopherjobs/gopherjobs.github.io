package htmlCreator

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gopherjobs/gopherjobs.github.io/internal/entity"
	"github.com/gopherjobs/gopherjobs.github.io/internal/html_creator/templates"
	"github.com/yuin/goldmark"
	"html/template"
	"log/slog"
	"os"
	"path/filepath"
	"time"
)

var (
	ErrJobIsEmpty = errors.New("job is empty")
)

type HTMLCreator interface {
	Generate(jobs []*entity.Job, outputPath string) error
}

type htmlCreator struct {
	logger *slog.Logger
}

func NewHTMLCreator(logger *slog.Logger) HTMLCreator {
	return &htmlCreator{
		logger: logger,
	}
}

func (h *htmlCreator) Generate(jobs []*entity.Job, outputPath string) error {
	if outputPath == "" {
		outputPath = "."
	}
	if jobs == nil {
		return fmt.Errorf("%w: jobs slice is nil", ErrJobIsEmpty)
	}

	htmlDir := filepath.Join(outputPath, "html")
	jobsDir := filepath.Join(htmlDir, "jobs")

	if err := os.MkdirAll(jobsDir, 0o755); err != nil {
		return fmt.Errorf("failed to create html directories %s: %w", jobsDir, err)
	}

	perPage := 20
	total := len(jobs)
	pages := (total + perPage - 1) / perPage
	if pages == 0 {
		pages = 1
	}

	createdAt := time.Now().Unix()

	// generate search index (jobs.json) used by the client-side search on the home page
	if err := h.writeSearchIndex(jobs, htmlDir); err != nil {
		return err
	}

	// generate list pages
	for p := 1; p <= pages; p++ {
		start := (p - 1) * perPage
		end := start + perPage
		if end > total {
			end = total
		}

		pageJobs := jobs[start:end]

		data := map[string]interface{}{
			"Title":      "Gopher Jobs",
			"CreatedAt":  createdAt,
			"Jobs":       pageJobs,
			"Page":       p,
			"TotalPages": pages,
			"MyLink":     "https://github.com/milad-rasouli",
			"Footer":     "Created by Milad Rasouli",
		}

		outFile := filepath.Join(htmlDir, fmt.Sprintf("list%d.html", p))
		f, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("failed to create list file %s: %w", outFile, err)
		}

		if err := templates.JobListTemplate.Execute(f, data); err != nil {
			if cerr := f.Close(); cerr != nil {
				h.logger.Warn("failed to close file after template execution error", slog.String("file", outFile), slog.String("err", cerr.Error()))
			}
			return fmt.Errorf("failed to execute job list template for %s: %w", outFile, err)
		}

		if cerr := f.Close(); cerr != nil {
			h.logger.Warn("failed to close file after writing", slog.String("file", outFile), slog.String("err", cerr.Error()))
		}

		h.logger.Info("wrote job list", slog.String("file", outFile))
	}

	// generate individual job pages
	for _, job := range jobs {
		data := map[string]interface{}{
			"Title":           job.Title,
			"CreatedAt":       createdAt,
			"Job":             job,
			"DescriptionHTML": h.renderMarkdown(job.Description),
			"MyLink":          "https://github.com/milad-rasouli",
			"Footer":          "Created by Milad Rasouli",
		}

		outFile := filepath.Join(jobsDir, fmt.Sprintf("%s.html", job.ID))
		f, err := os.Create(outFile)
		if err != nil {
			return fmt.Errorf("failed to create job file %s: %w", outFile, err)
		}

		if err := templates.JobDescriptionTemplate.Execute(f, data); err != nil {
			if cerr := f.Close(); cerr != nil {
				h.logger.Warn("failed to close file after template execution error", slog.String("file", outFile), slog.String("err", cerr.Error()))
			}
			return fmt.Errorf("failed to execute job description template for %s: %w", outFile, err)
		}

		if cerr := f.Close(); cerr != nil {
			h.logger.Warn("failed to close file after writing", slog.String("file", outFile), slog.String("err", cerr.Error()))
		}

		h.logger.Info("wrote job description", slog.String("file", outFile))
	}

	return nil
}

// searchEntry is the trimmed-down job representation written to jobs.json
type searchEntry struct {
	ID              string `json:"id"`
	Title           string `json:"title"`
	Company         string `json:"company"`
	Location        string `json:"location"`
	Level           string `json:"level"`
	Type            string `json:"type"`
	CompanyIndustry string `json:"company_industry"`
	Remote          bool   `json:"remote"`
	Relocation      bool   `json:"relocation"`
}

func (h *htmlCreator) writeSearchIndex(jobs []*entity.Job, htmlDir string) error {
	entries := make([]searchEntry, 0, len(jobs))
	for _, job := range jobs {
		entries = append(entries, searchEntry{
			ID:              job.ID,
			Title:           job.Title,
			Company:         job.Company,
			Location:        job.Location,
			Level:           job.Level,
			Type:            job.Type,
			CompanyIndustry: job.CompanyIndustry,
			Remote:          job.Remote,
			Relocation:      job.Relocation,
		})
	}

	data, err := json.Marshal(entries)
	if err != nil {
		return fmt.Errorf("failed to marshal search index: %w", err)
	}

	outFile := filepath.Join(htmlDir, "jobs.json")
	if err := os.WriteFile(outFile, data, 0o644); err != nil {
		return fmt.Errorf("failed to write search index %s: %w", outFile, err)
	}

	h.logger.Info("wrote search index", slog.String("file", outFile), slog.Int("jobs", len(entries)))
	return nil
}

func (h *htmlCreator) renderMarkdown(md string) template.HTML {
	var buf bytes.Buffer
	if err := goldmark.Convert([]byte(md), &buf); err != nil {
		return template.HTML(md)
	}
	return template.HTML(buf.String())
}
