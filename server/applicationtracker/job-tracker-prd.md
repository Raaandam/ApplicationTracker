# Product Requirements Document
## AI-Powered Job Application Tracker

**Version:** 1.0  
**Status:** Draft  
**Stack:** Go (Gin) · React · MongoDB · LangChain Go · OpenAI API · JWT

---

## 1. Overview

A full-stack web application that helps job seekers track their applications, get AI-powered resume-to-JD fit analysis, and surface patterns in which types of roles they perform best in — all in one place.

---

## 2. Problem Statement

Job seekers applying to multiple roles face three core problems:

- **No structured tracking** — applications are scattered across spreadsheets, notes, and memory
- **No feedback loop** — candidates don't know why they get interviews for some roles and not others
- **Wasted effort** — they apply to roles where their resume is a poor fit without knowing it

---

## 3. Goals

- Allow users to log and track job applications in one place
- Use an LLM to analyze resume-JD fit and return a structured score + actionable feedback
- Surface patterns over time (which industries/roles score highest for the user)
- Keep the stack identical to the reference project: Go + Gin, React, MongoDB, LangChain Go, JWT

---

## 4. Non-Goals (for v1)

- Mobile app
- Email integrations (auto-import from Gmail etc.)
- Team/collaborative features
- Browser extension for one-click job saving

---

## 5. Users

Single-user, self-hosted or deployed. Auth exists to support multiple accounts on the same instance (e.g., a couple sharing a server), not for enterprise multi-tenancy.

---

## 6. Tech Stack

| Layer | Technology |
|---|---|
| Backend | Go with Gin-Gonic |
| Frontend | React |
| Database | MongoDB |
| AI | LangChain Go + OpenAI API |
| Auth | JWT — access token in HTTP-only cookie |
| Deployment | Single binary + MongoDB Atlas (or local) |

---

## 7. Data Models

### 7.1 User
```
{
  _id:          ObjectId
  email:        string (unique)
  password:     string (bcrypt hashed)
  name:         string
  created_at:   timestamp
}
```

### 7.2 Application
```
{
  _id:              ObjectId
  user_id:          ObjectId (ref → User)
  company:          string
  role:             string
  job_description:  string (raw JD text pasted by user)
  resume_version:   string (raw resume text pasted by user)
  status:           enum → ["Saved", "Applied", "Interview", "Offer", "Rejected"]
  applied_date:     date (optional)
  notes:            string (optional free text)
  ai_analysis:      Analysis (see 7.3, nullable until run)
  created_at:       timestamp
  updated_at:       timestamp
}
```

### 7.3 Analysis (embedded in Application)
```
{
  score:          number (0–100)
  fit_label:      enum → ["Strong", "Moderate", "Weak"]
  strengths:      []string
  gaps:           []string
  tip:            string (one actionable improvement suggestion)
  analyzed_at:    timestamp
}
```

---

## 8. API Endpoints

### Auth
| Method | Route | Description |
|---|---|---|
| POST | /api/auth/register | Create account |
| POST | /api/auth/login | Login, set HTTP-only cookie |
| POST | /api/auth/logout | Clear cookie |
| GET | /api/auth/me | Return current user (validates token) |

### Applications
| Method | Route | Description |
|---|---|---|
| GET | /api/applications | List all applications for current user |
| POST | /api/applications | Create new application |
| GET | /api/applications/:id | Get single application |
| PATCH | /api/applications/:id | Update fields (status, notes, etc.) |
| DELETE | /api/applications/:id | Delete application |

### AI Analysis
| Method | Route | Description |
|---|---|---|
| POST | /api/applications/:id/analyze | Run AI analysis on an application |

### Insights (stretch goal)
| Method | Route | Description |
|---|---|---|
| GET | /api/insights | Aggregated stats — avg score by industry, status breakdown, etc. |

---

## 9. Features

### 9.1 Auth (Phase 1)
- Register with email + password
- Login returns JWT stored in HTTP-only cookie (mirrors reference project exactly)
- Protected routes require valid token
- Logout clears cookie

### 9.2 Application CRUD (Phase 1)
- Create an application with: company, role, JD text, resume text, status, notes
- View all applications in a dashboard (filterable by status)
- Edit any field
- Delete an application
- Status pipeline: Saved → Applied → Interview → Offer / Rejected

### 9.3 AI Analysis (Phase 2)
- User triggers analysis manually via a button on the application detail page
- Backend sends JD + resume text to OpenAI via LangChain Go
- Prompt instructs the model to return **only JSON** with: score (0–100), fit_label, strengths (array), gaps (array), tip (string)
- Response is parsed and stored in the `ai_analysis` embedded field
- Frontend displays the analysis results on the detail page

**LLM Prompt structure (for reference):**
```
You are a professional recruiter. Given the job description and resume below,
evaluate how well the candidate's resume fits the role.

Return ONLY a JSON object with these fields:
- score: integer 0-100
- fit_label: one of "Strong", "Moderate", or "Weak"
- strengths: array of 2-3 strings
- gaps: array of 2-3 strings
- tip: one actionable string the candidate can do to improve their resume for this role

Job Description:
{job_description}

Resume:
{resume_text}
```

### 9.4 Dashboard (Phase 1 + 2)
- Cards or table view of all applications
- Filter by status
- Each card shows: company, role, status badge, AI score (if analyzed)
- Click to open detail page

### 9.5 Insights Panel (Phase 3 — stretch)
- Total applications, breakdown by status (pie/bar chart)
- Average AI score across all applications
- Which roles/companies scored highest
- Timeline of applications over time

---

## 10. Frontend Pages

| Page | Route | Description |
|---|---|---|
| Login | /login | Email + password login |
| Register | /register | Create account |
| Dashboard | / | All applications, filterable |
| New Application | /applications/new | Form to create application |
| Application Detail | /applications/:id | Full detail + AI analysis results |
| Edit Application | /applications/:id/edit | Edit form |

---

## 11. Security Considerations

- Passwords hashed with bcrypt
- JWT stored in HTTP-only cookie (not localStorage) — prevents XSS theft
- All `/api/applications` routes are protected by auth middleware
- User can only access their own applications (user_id check on every query)
- OpenAI API key stored as environment variable, never exposed to frontend

---

## 12. Environment Variables

```
PORT=8080
MONGO_URI=mongodb://localhost:27017
MONGO_DB=job_tracker
JWT_SECRET=your_secret_here
OPENAI_API_KEY=sk-...
```

---

## 13. Build Phases

### Phase 1 — Core (follow reference video)
- Auth (register, login, logout, JWT cookie)
- Application CRUD
- Basic dashboard UI in React

### Phase 2 — AI Layer
- POST /analyze endpoint
- LangChain Go + OpenAI integration
- Structured JSON extraction from LLM response
- Analysis display on frontend

### Phase 3 — Insights (stretch)
- MongoDB aggregation queries
- Insights page with charts

---

## 14. Success Criteria

- A user can register, log in, and have their session persist across page refreshes
- A user can create, view, edit, and delete job applications
- A user can trigger AI analysis on any application and see a score, strengths, gaps, and a tip
- All data is user-scoped (one user cannot see another user's applications)
- The app can be run locally with a single `go run` + `npm start`
