# startGo=`date +%s.%N`
echo "Go:"
# echo "Compile-"
go build -o main-go main.go 
# echo
echo "Run-"
time ./main-go
echo
echo
# endGo=`date +%s.%N`

# startC=`date +%s.%N`
echo "C:"
# echo "Compile-"
gcc main.c -o main
# echo
echo "Run-"
time ./main
echo
echo
# endC=`date +%s.%N`

# startPy=`date +%s.%N`
echo "Python:"
echo "Interpretted-"
time python3 main.py
echo
# endPy=`date +%s.%N`

# runtimeGo=$( echo "$endGo - $startGo" | bc -l )
# runtimeC=$( echo "$endC - $startC" | bc -l )
# runtimePy=$( echo "$endC - $startC" | bc -l )
