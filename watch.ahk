#Requires AutoHotkey v2.0

notRecording := false

F8:: {
    global notRecording

    fmediaPath := A_WorkingDir "\bin\fmedia-1.31-windows-x64\fmedia\fmedia.exe"
    whisperPath := A_WorkingDir "\bin\whisper-autohotkey\whisper-autohotkey.exe"

    notRecording := !notRecording
    if notRecording {
        Run(Chr(34) . fmediaPath . Chr(34) . " --record --overwrite --mpeg-quality=16 --rate=12000 --out=rec.mp3 --globcmd=listen", , "Hide")
    } else {
        Run(Chr(34) . fmediaPath . Chr(34) . " --globcmd=stop", , "Hide")
        Sleep(200)
        Run(Chr(34) . whisperPath . Chr(34) . " zh", , "Hide")
    }
}
