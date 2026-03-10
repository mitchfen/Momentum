## About

Momentum is a simple application designed for tracking tasks and habit stacks that need to be done each day.  
Task state is persisted in a backend JSON file and resets each night.  

### Example 
<img src="./screenshot.png" width="600">

## Implementation details
*   **Blazor Server:** Built with .NET 10 for a modern, server-side interactive experience.
*   **Persistent State:** Task completion is stored in a JSON file in the `/data` directory, which can be mounted as a volume.
*   **Containerized:** Multi-stage builds using official .NET 10 runtime images.
*   **Dynamic Configuration:** Tasks are loaded from environment variables (`DAILY_TASKS`) or `appsettings.json`, and the reset timezone can be configured via `TIMEZONE` (e.g., `America/New_York`). This allows for flexible configuration without rebuilding the image.
*   **Kubernetes Ready:** Includes manifests for Deployment, Service, Ingress, PersistentVolumeClaim, and ConfigMap.

## Build and Deploy

### 1. Run Locally (dotnet)
You can run the application directly using the dotnet CLI:
```bash
dotnet run --project src/Momentum.csproj
```

### 2. Build Docker Image
Build the multi-stage image (which handles restoration and publishing internally):
```bash
docker build -t momentum:latest .
```

### 3. Run Locally (Docker)
Run the container to test locally. Access at `http://localhost:5000`
```bash
docker run -p 5000:80 --rm -it momentum:latest
```
*Note: You can pass environment variables to test the config injection:*
```bash
docker run -p 5000:80 \
  -e "DAILY_TASKS=Test Task 1,Test Task 2" \
  -e "TIMEZONE=America/Los_Angeles" \
  --rm -it momentum:latest
```

### 4. Deploy to Kubernetes
Apply the manifests to your cluster.
```bash
kubectl apply -f kubernetes-manifests
```
