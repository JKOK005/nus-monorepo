go run main.go -stderrthreshold=INFO 1 8001 &
go run main.go -stderrthreshold=INFO 1 8002 &
go run main.go -stderrthreshold=INFO 2 8003 &
go run main.go -stderrthreshold=INFO 2 8004 &
go run main.go -stderrthreshold=INFO 4 8005
