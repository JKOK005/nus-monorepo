go run main.go -stderrthreshold=INFO -hash=1 -port=8001 &
go run main.go -stderrthreshold=INFO -hash=1 -port=8003 &
go run main.go -stderrthreshold=INFO -hash=2 -port=8005 &
go run main.go -stderrthreshold=INFO -hash=2 -port=8007 &
go run main.go -stderrthreshold=INFO -hash=4 -port=8009
