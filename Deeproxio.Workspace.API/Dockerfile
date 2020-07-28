FROM mcr.microsoft.com/dotnet/core/sdk:3.1-alpine AS build-env
WORKDIR /app

# Copy csproj and restore as distinct layers
COPY . .
RUN dotnet restore

WORKDIR /app/Deeproxio.Workspace.API

RUN dotnet publish -c Release -o out

# Build runtime image
FROM mcr.microsoft.com/dotnet/core/aspnet:3.1-alpine
WORKDIR /app
COPY --from=build-env /app/Deeproxio.Workspace.API/out .

EXPOSE 80

ENTRYPOINT ["dotnet", "Deeproxio.Workspace.API.dll"]