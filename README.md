## About

Momentum is a simple application for tracking tasks I want/need to do every morning. It allows you to keep tasks individual or group them together in habit stacks. I deploy this in my homelab and make it available on the network so any device in my home can be used to update the status. 

### Example 
<img src="./screenshot.png" width="600">

## Key Features

*   **Task & Habit Stacks:** Define simple tasks or break down habits into discrete steps (e.g., "Meditate + Exercise + Read")
*   **Persistent State:** Daily completion state stored in JSON format, mounted as a volume for data durability
*   **Timezone-Aware:** Configured reset timezone (`TIMEZONE` env var) ensures the daily reset happens at the right time for your location
*   **Fast & Lightweight:** Built in Go with zero external dependencies (uses only stdlib), resulting in a tiny, efficient container
*   **Dynamic Configuration:** Load tasks from `config.json` (file-based) or `DAILY_TASKS` environment variable (container-friendly)
*   **Container-First:** Multi-stage Docker build produces a minimal scratch-based image (~5MB)
*   **Kubernetes Ready:** Includes manifests for Deployment, Service, PersistentVolumeClaim, and ConfigMap

## Task Configuration

Tasks can be configured in two ways:

### 1. File-based (`config.json`)
```json
{
  "DailyTasks": [
    "Habitstack task 1 + habitstack task 2 + habitstack task 3",
    "Single task",
    "Named Stack: Step 1 + Step 2"
  ],
  "TimeZone": "America/New_York"
}
```

### 2. Environment variables (for containers)
```bash
DAILY_TASKS="Task 1,Task 2,Morning: Coffee + Meditate"
TIMEZONE="America/Los_Angeles"
```

#### Task Format
- **Simple task:** `"Task name"` → displays as a single checkbox
- **Unnamed habit stack:** `"Step 1 + Step 2 + Step 3"` → displays as a progress bar with individual steps
- **Named habit stack:** `"My Morning: Step 1 + Step 2"` → named stack with steps
- **Comma-separated:** Multiple tasks in `DAILY_TASKS` are split by comma

## Build and Deploy

### 1. Run Locally (native Go)
Build and run the binary directly:
```bash
cd src
go build -o momentum
./momentum
```
Access at `http://localhost:80` (or set `PORT` env var for a different port)

### 2. Build Docker Image
Multi-stage build produces a tiny, scratch-based image:
```bash
docker build -t momentum:latest .
```

### 3. Run Locally (Docker)
Run the container and access at `http://localhost:8080`
```bash
docker run -p 8080:80 \
  -v momentum-data:/app/data \
  --rm -it momentum:latest
```

With custom configuration:
```bash
docker run -p 8080:80 \
  -e "DAILY_TASKS=Morning: Coffee + Meditate,Workout,Read" \
  -e "TIMEZONE=America/Chicago" \
  -v momentum-data:/app/data \
  --rm -it momentum:latest
```

### 4. Deploy to Kubernetes
Apply the manifests to your cluster:
```bash
kubectl apply -f kubernetes-manifest.yaml
```

This creates:
- `momentum` namespace
- ConfigMap with `DAILY_TASKS` and `TIMEZONE`
- PersistentVolumeClaim for `/app/data`
- Deployment with resource limits
- ClusterIP Service for internal access

Configure your own tasks by editing the ConfigMap in the manifest before applying.
