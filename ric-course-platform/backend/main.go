package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/joho/godotenv"
)

type Server struct {
	pool *pgxpool.Pool
}

type Course struct {
	ID    int    `json:"id"`
	Code  string `json:"code"`
	Title string `json:"title"`
}

type Student struct {
	ID   int    `json:"id"`
	Name string `json:"name"`
}

func mustGetEnv(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func connectDB(ctx context.Context) *pgxpool.Pool {
	// Load .env if present
	_ = godotenv.Load()

	dsn := mustGetEnv("DATABASE_URL", "postgres://postgres:postgres@localhost:5432/ric?sslmode=disable")
	cfg, err := pgxpool.ParseConfig(dsn)
	if err != nil {
		log.Fatalf("failed to parse DATABASE_URL: %v", err)
	}
	cfg.MaxConns = 5
	pool, err := pgxpool.NewWithConfig(ctx, cfg)
	if err != nil {
		log.Fatalf("failed to connect database: %v", err)
	}
	// quick ping
	ctxPing, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()
	if err := pool.Ping(ctxPing); err != nil {
		log.Fatalf("database ping failed: %v", err)
	}
	return pool
}

func main() {
	ctx := context.Background()
	pool := connectDB(ctx)
	defer pool.Close()

	s := &Server{pool: pool}

	r := gin.Default()

	// CORS
	corsOrigins := mustGetEnv("CORS_ORIGINS", "*")
	cfg := cors.DefaultConfig()
	if corsOrigins == "*" {
		cfg.AllowAllOrigins = true
	} else {
		cfg.AllowOrigins = strings.Split(corsOrigins, ",")
	}
	cfg.AllowHeaders = []string{"Origin", "Content-Type", "Accept"}
	r.Use(cors.New(cfg))

	// Health
	r.GET("/healthz", func(c *gin.Context) { c.JSON(http.StatusOK, gin.H{"status": "ok"}) })

	api := r.Group("/api")
	{
		api.GET("/courses", s.listAllCourses)
		api.GET("/students", s.listAllStudents)
		api.GET("/students/:id/courses", s.listCoursesByStudent)
		api.PUT("/students/:id", s.updateStudent)
		api.PUT("/courses/:id", s.updateCourse)
		api.POST("/courses", s.createCourse)
		api.POST("/students", s.createStudent)
		api.POST("/students/:id/courses", s.addCourseForStudent)
		api.DELETE("/students/:id/courses/:courseId", s.removeCourseForStudent)

	}

	port := mustGetEnv("PORT", "8080")
	addr := fmt.Sprintf("0.0.0.0:%s", port)
	log.Printf("Server listening on %s", addr)
	if err := r.Run(addr); err != nil {
		log.Fatal(err)
	}
}

func (s *Server) listAllCourses(c *gin.Context) {
	ctx := c.Request.Context()
	rows, err := s.pool.Query(ctx, `SELECT id, code, title FROM courses ORDER BY code ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []Course
	for rows.Next() {
		var cr Course
		if err := rows.Scan(&cr.ID, &cr.Code, &cr.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, cr)
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "items": items})
}

func (s *Server) listAllStudents(c *gin.Context) {
	ctx := c.Request.Context()
	rows, err := s.pool.Query(ctx, `SELECT id, name FROM students ORDER BY name ASC`)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []Student
	for rows.Next() {
		var st Student
		if err := rows.Scan(&st.ID, &st.Name); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, st)
	}
	c.JSON(http.StatusOK, gin.H{"items": items})
}

func (s *Server) listCoursesByStudent(c *gin.Context) {
	id := c.Param("id")
	ctx := c.Request.Context()
	rows, err := s.pool.Query(ctx, `
		SELECT c.id, c.code, c.title
		FROM enrollments e
		JOIN courses c ON c.id = e.course_id
		WHERE e.student_id = $1
		ORDER BY c.code ASC`, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	defer rows.Close()

	var items []Course
	for rows.Next() {
		var cr Course
		if err := rows.Scan(&cr.ID, &cr.Code, &cr.Title); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		items = append(items, cr)
	}
	c.JSON(http.StatusOK, gin.H{"count": len(items), "items": items})
}

type updateStudentReq struct {
    Name string `json:"name"`
}

func (s *Server) updateStudent(c *gin.Context) {
    id := c.Param("id")
    var req updateStudentReq
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
        return
    }
    name := strings.TrimSpace(req.Name)
    if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
        return
    }

    ctx := c.Request.Context()
    cmd, err := s.pool.Exec(ctx, `UPDATE students SET name=$1 WHERE id=$2`, name, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if cmd.RowsAffected() == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "student not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"id": id, "name": name})
}

type updateCourseReq struct {
    Code  string `json:"code"`
    Title string `json:"title"`
}

func (s *Server) updateCourse(c *gin.Context) {
    id := c.Param("id")
    var req updateCourseReq
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
        return
    }

    code := strings.TrimSpace(req.Code)
    title := strings.TrimSpace(req.Title)
    if code == "" || title == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "code and title are required"})
        return
    }

    ctx := c.Request.Context()
    cmd, err := s.pool.Exec(ctx, `UPDATE courses SET code=$1, title=$2 WHERE id=$3`, code, title, id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }
    if cmd.RowsAffected() == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "course not found"})
        return
    }
    c.JSON(http.StatusOK, gin.H{"id": id, "code": code, "title": title})
}

type createCourseReq struct {
    Code  string `json:"code"`
    Title string `json:"title"`
}

func (s *Server) createCourse(c *gin.Context) {
    var req createCourseReq
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
        return
    }

    code := strings.TrimSpace(req.Code)
    title := strings.TrimSpace(req.Title)
    if code == "" || title == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "code and title are required"})
        return
    }

    ctx := c.Request.Context()
    var id int
    err := s.pool.QueryRow(ctx, `INSERT INTO courses (code, title) VALUES ($1, $2) RETURNING id`, code, title).Scan(&id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": id, "code": code, "title": title})
}

type createStudentReq struct {
    Name string `json:"name"`
}

func (s *Server) createStudent(c *gin.Context) {
    var req createStudentReq
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
        return
    }

    name := strings.TrimSpace(req.Name)
    if name == "" {
        c.JSON(http.StatusBadRequest, gin.H{"error": "name is required"})
        return
    }

    ctx := c.Request.Context()
    var id int
    err := s.pool.QueryRow(ctx, `INSERT INTO students (name) VALUES ($1) RETURNING id`, name).Scan(&id)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"id": id, "name": name})
}

type addCourseReq struct {
    CourseID int `json:"course_id"`
}

func (s *Server) addCourseForStudent(c *gin.Context) {
    studentID := c.Param("id")
    var req addCourseReq
    if err := c.BindJSON(&req); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "invalid JSON"})
        return
    }

    ctx := c.Request.Context()
    _, err := s.pool.Exec(ctx, `INSERT INTO enrollments (student_id, course_id) VALUES ($1, $2)`, studentID, req.CourseID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    c.JSON(http.StatusOK, gin.H{"student_id": studentID, "course_id": req.CourseID})
}

func (s *Server) removeCourseForStudent(c *gin.Context) {
    studentID := c.Param("id")
    courseID := c.Param("courseId")

    ctx := c.Request.Context()
    cmd, err := s.pool.Exec(ctx, `DELETE FROM enrollments WHERE student_id=$1 AND course_id=$2`, studentID, courseID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
        return
    }

    if cmd.RowsAffected() == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "enrollment not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"student_id": studentID, "course_id": courseID, "status": "removed"})
}

