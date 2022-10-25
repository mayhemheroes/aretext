package fuzz

import "strconv"
import "github.com/aretext/aretext/text"
import "github.com/aretext/aretext/file"


func mayhemit(bytes []byte) int {

    var num int
    if len(bytes) > 1 {
        num, _ = strconv.Atoi(string(bytes[0]))

        switch num {
    
        case 0:
            content := string(bytes)
            text.Reverse(content)
            return 0

        case 1:
            content := string(bytes)
            file.GlobMatch("mayhem", content)
            return 0

        default:
            content := string(bytes)
            file.RelativePathCwd(content)
            return 0

        }
    }
    return 0
}

func Fuzz(data []byte) int {
    _ = mayhemit(data)
    return 0
}