package main
/*
Безопасное копирование файлов
пример использования defer для менеджмента ресурсов
*/
import (
  "os"
  "io"
  "fmt"
)

func CopyFile(srcName, dstName string) (written int64, err error) {
    src, err := os.Open(srcName)
    if err != nil {
        fmt.Print("Open src error: %s\n",err)
        return
    }
    defer src.Close()

    dst, err := os.Create(dstName)
    if err != nil {
        fmt.Print("Open dst error: %s\n",err)
        return
    }
    defer dst.Close()

    return io.Copy(dst, src)
}

func main () {
  argList := os.Args[1:]
  if len(argList) != 2 {
    fmt.Printf("Expected format: dst src, get %d params\n",len(argList))
    return
  }
  fmt.Printf("Copy: %v\n",argList)
  CopyFile(argList[0],argList[1])
}
