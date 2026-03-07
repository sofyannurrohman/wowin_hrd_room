# HRD Room - Product Documentation

Welcome to the HRD Room Product Documentation! This guide will help you understand how to deploy, configure, and use the HRD Room application, both from the perspective of an HR Administrator and an Exam Participant.

---

## 1. Introduction & Setup

HRD Room is a comprehensive assessment platform designed to streamline the recruitment process. It allows HR teams to manage job positions, create exam modules, conduct live monitored exam sessions, and evaluate candidates effectively.

### 1.1. Local Setup Requirements
To run the HRD Room application locally, you will need:
- Docker
- Docker Compose

### 1.2. Running the Application
The application is pre-configured with a `docker-compose.yml` file that orchestrates the PostgreSQL database, the Go backend, the Vue frontend, and an Nginx reverse proxy.

1. Clone the repository and navigate to the project root directory:
   ```bash
   cd /path/to/hrd_room
   ```
2. Start the application using Docker Compose:
   ```bash
   docker-compose up -d --build
   ```
3. Once the containers are up and running, you can access:
   - **Frontend (Web App):** `http://localhost:3000`
   - **Backend API:** `http://localhost:8080`
   - **MailHog (Email Testing UI):** `http://localhost:8025`

*(Note: The database runs on port `5433` on the host machine and `5432` internally).*

---

## 2. Admin & HR Dashboard Guide

The HR Dashboard provides administrators with full control over the examination process, from managing questions to monitoring live sessions.

### 2.1. Authentication
- Access the HR Dashboard by navigating to `/login`.
- You must authenticate with an HR-privileged account to gain access to the dashboard (`/dashboard` and related routes).

### 2.2. Master Data Management
Before securing candidates, you need to set up the foundations of the recruitment process:
- **Job Positions (`/job-positions`):** Create and manage the roles you are hiring for. You can define the title and description for each position.
- **Questions Bank (`/questions`):** Create individual questions. The platform supports various question types:
  - Multiple Choice
  - True/False
  - Short Answer (requires manual evaluation)
  - Psychological (custom evaluation formats)
- **Exam Modules (`/modules`):** Group questions logically into modules (e.g., "Basic Programming", "Logical Reasoning"). These modules will later be linked to exam sessions.

### 2.3. Participants & Applications
- **Applicants (`/apply`):** Candidates first register their interest by applying through the `/apply` page, submitting their resume/CV and personal data.
- **Participant Management (`/participants`):** View all registered participants. You can inspect an individual's details securely through the participant detail page (`/participants/:id`).

### 2.4. Exam Sessions Management
- **Creating Sessions (`/sessions/create`):** An exam session is the actual test event. When creating a session, you will link a "Job Position", assign "Modules", define the start and end schedule, and set passing criteria.
- **Session Tokens:** Each created session generates a unique joining token. This token is what participants use to enter the exam.
- **Live Monitoring (`/sessions/:id/monitor`):** While a session is active, HR can monitor it in real-time. This provides a live view of connected participants.
- **Access Control:** The system enforces time constraints and prevents early joins (participants cannot access the exam before the designated start time).

### 2.5. Evaluation & Results
- **Session Results (`/sessions/:id/results`):** After an exam, view the aggregated results and scores of all participants.
- **Manual Evaluation:** Some question types (e.g., Short Answer) require human intervention. HR can review specific participant answers (`/sessions/:id/results/:participantId/answers`) and assign scores manually.
- **Session Violations (`/sessions/:id/violations`):** View any recorded violations (like switching tabs or potential cheating behaviors detected by the system).
- **Analytics (`/sessions/:id/analytics`):** View performance graphs and statistical details of the completed session.

### 2.6. Payroll & Export
- **CSV Export:** For final processing or HR disbursement needs, participant data and examination results can be exported efficiently via CSV.

---

## 3. Participant Guide

This section is for job candidates taking an exam on the platform.

### 3.1. Applying for a Job
- Navigate to the **Application Page (`/apply`)**.
- Fill in your personal details (Name, Email, Phone, etc.) and upload your CV/Resume.
- Successfully submitting your application registers you as a participant in the system.

### 3.2. Joining an Exam
- Once provided an "Exam Token" by HR, navigate to the **Join Page (`/join`)**.
- Enter your token. The system verifies if you are allowed to enter. *Note: You will not be able to join before the scheduled start time.*
- **Camera Check (`/camera-check`):** Before the exam begins, you may be required to pass a system check ensuring your webcam and browser environment are correctly set up.

### 3.3. Taking the Exam
- **The Exam Interface (`/exam/:sessionId`):** You will be presented with questions module by module. Simply follow the on-screen instructions.
- **Session Persistence:** If you accidentally reload the page or briefly lose connection, don't worry! Your session progress, timer, and answers are saved automatically. You will reconnect exactly where you left off.
- **Violations:** The system monitors behavior such as leaving the exam tab. Repeated violations may be recorded and flagged to HR.
- **Submission:** Once all questions are answered or the timer runs out, the exam will automatically submit and redirect you to the Finished Page (`/exam-finished`).

---

## 4. End-to-End Workflow

This section outlines the complete journey from initial HR setup to the final grading of a participant's exam.

**Step 1: HR Setup**
- HR logs into the Dashboard.
- HR defines a **Job Position**, creates **Questions**, and groups them into **Modules**.

**Step 2: Participant Application**
- A prospective candidate navigates to the `/apply` page.
- The candidate submits their personal details and CV/Resume, registering them as a candidate in the system.

**Step 3: Session Creation & Invitation**
- HR creates an **Exam Session** linked to the target Job Position, assigning the relevant Modules and scheduling the start/end times.
- HR retrieves the generated **Exam Token** from the new session and shares it with the candidate.

**Step 4: Participant Takes the Exam**
- At the scheduled time, the participant goes to the `/join` page and enters the Exam Token.
- After passing the preliminary camera check, the participant completes the exam modules.
- The system automatically handles submission when all questions are answered or when the timer expires.

**Step 5: Automated & Manual Grading by HR**
- Objective questions (Multiple Choice, True/False) are graded automatically by the system.
- For subjective questions (Short Answer, Psychological), HR navigates to the session's **Results** page.
- HR selects the candidate to review their specific answers (`Participant Answers`) and inputs manual scores.

**Step 6: Final Review & Export**
- HR reviews the final accumulated scores and checks the **Violations** tab for any suspicious behavior during the session.
- Once finalized, HR can use the **CSV Export** feature to download the results for final decision-making.