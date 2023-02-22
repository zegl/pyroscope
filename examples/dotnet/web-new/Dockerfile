FROM mcr.microsoft.com/dotnet/sdk:6.0 as builder

WORKDIR /sln

# Copy project file and restore
COPY "example/Example.csproj" "."
RUN dotnet restore "Example.csproj" 


# Copy the actual source code
COPY "example" "."



# Build and publish the app
RUN dotnet publish "Example.csproj" -c Release -o /app/publish --self-contained false --no-restore




#FINAL image
FROM mcr.microsoft.com/dotnet/aspnet:6.0
WORKDIR /app
COPY --from=builder /app/publish .

# Set up Pyroscope profiling

ADD https://github.com/pyroscope-io/pyroscope-dotnet/releases/download/v0.7.2-pyroscope/Pyroscope.Linux.ApiWrapper.x64.so /lib/Pyroscope.Linux.ApiWrapper.x64.so
ADD https://github.com/pyroscope-io/pyroscope-dotnet/releases/download/v0.7.2-pyroscope/Pyroscope.Profiler.Native.so /lib/Pyroscope.Profiler.Native.so
ENV PYROSCOPE_APPLICATION_NAME=test.dotnet.app
ENV PYROSCOPE_SERVER_ADDRESS=http://pyroscope-route:80
ENV PYROSCOPE_PROFILING_ENABLED=1
ENV CORECLR_ENABLE_PROFILING=1
ENV CORECLR_PROFILER={BD1A650D-AC5D-4896-B64F-D6FA25D6B26A}
ENV CORECLR_PROFILER_PATH=/lib/Pyroscope.Profiler.Native.so
ENV LD_PRELOAD=/lib/Pyroscope.Linux.ApiWrapper.x64.so

ENTRYPOINT ["dotnet", "example.dll"]