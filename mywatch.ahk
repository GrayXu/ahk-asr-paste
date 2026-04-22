#Requires AutoHotkey v2.0

notRecording := false

CapsLock::Esc

#s:: {
    Send("^s")
}

#v:: {
    Send("^v")
}

#c:: {
    Send("^c")
}

#z:: {
    Send("^z")
}

F8:: {
    global notRecording

    fmediaPath := A_WorkingDir "\bin\fmedia-1.31-windows-x64\fmedia\fmedia.exe"
    whisperPath := A_WorkingDir "\bin\whisper-autohotkey\whisper-autohotkey.exe"

    notRecording := !notRecording
    if notRecording {
        Run(Chr(34) . fmediaPath . Chr(34) . " --record --overwrite --mpeg-quality=16 --rate=12000 --out=rec.mp3 --globcmd=listen", , "Hide")
    } else {
        Run(Chr(34) . fmediaPath . Chr(34) . " --globcmd=stop", , "Hide")
        Sleep(100)
        Run(Chr(34) . whisperPath . Chr(34) . " zh", , "Hide")
    }
}
