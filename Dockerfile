# Build Stage
FROM mcr.microsoft.com/dotnet/sdk:10.0-alpine AS build
WORKDIR /source

# Copy csproj and restore dependencies
COPY src/*.csproj src/
RUN dotnet restore src/Momentum.csproj

# Copy remaining source code and publish
COPY . .
WORKDIR /source/src
RUN dotnet publish Momentum.csproj -c Release -o /app/publish

# Runtime Stage
FROM mcr.microsoft.com/dotnet/aspnet:10.0-alpine
WORKDIR /app

# Copy published output from build stage
COPY --from=build /app/publish .

# The app expects port 80
EXPOSE 80
ENV ASPNETCORE_URLS=http://+:80

ENTRYPOINT ["dotnet", "Momentum.dll"]

LABEL org.opencontainers.image.description="Momentum - a daily habit stacking, task tracking application"
